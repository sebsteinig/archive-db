package sql

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type SQL_Runner interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	// Close()
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
}

type SQL_Query struct {
	query  string
	holder Placeholder
	args   []any
}
type SQL_Value struct {
	Name     string
	Value    any
	Nullable bool
	Is_null  bool
}

func SQLf(query string, args ...any) (SQL_Query, error) {
	sql_query := SQL_Query{
		query:  query,
		holder: BuildPlaceholder(len(args)),
		args:   args,
	}
	return sql_query, nil
}

func (sql *SQL_Query) format() {
	var pushed_args []any
	sql.holder = BuildPlaceholder(len(sql.args))
	for _, arg := range sql.args {
		if builder, ok := arg.(SqlBuilder); ok {
			arg_str := builder.Build(&sql.holder)
			pushed_args = append(pushed_args, arg_str)
		} else if value, ok := arg.(SQL_Value); ok {
			if value.Is_null {
				pushed_args = append(pushed_args, "NULL")
			} else {
				pushed_args = append(pushed_args, sql.holder.Push(value.Value))
			}
		} else if values, ok := arg.([]SQL_Value); ok {
			var values_str []string
			for _, value := range values {
				if value.Is_null {
					values_str = append(values_str, "NULL")
				} else {
					values_str = append(values_str, sql.holder.Push(value.Value))
				}
			}
			pushed_args = append(pushed_args, to_any_list(values_str)...)
		} else {
			pushed_args = append(pushed_args, sql.holder.Push(arg))
		}
	}
	sql.query = fmt.Sprintf(sql.query, pushed_args...)
}
func (sql *SQL_Query) Suffixe(query string) *SQL_Query {
	sql.query += query
	return sql
}

func (sql *SQL_Query) Append(others ...*SQL_Query) *SQL_Query {
	for _, other := range others {
		sql.query += other.query
		sql.args = append(sql.args, other.args...)
	}
	return sql
}

func Insert[T any](table string, values ...T) (SQL_Query, error) {
	sql_values := make([][]SQL_Value, 0, len(values))
	for _, obj := range values {
		if sql_value, err_value := sql_value_from_struct[T](obj); err_value == nil {
			sql_values = append(sql_values, sql_value)
		} else {
			return SQL_Query{}, err_value
		}
	}
	if len(sql_values) == 0 {
		return SQL_Query{}, fmt.Errorf("Zero Value for insert")
	}
	var column_names []string
	for _, value := range sql_values[0] {
		column_names = append(column_names, value.Name)
	}
	insert_sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES `,
		table, strings.Join(column_names, ","))

	var values_str []string
	for _, sql_value := range sql_values {
		var value_str []string
		for range sql_value {
			value_str = append(value_str, "%s")
		}
		values_str = append(values_str, fmt.Sprintf("(%s)", strings.Join(value_str, ",")))
	}
	insert_sql += strings.Join(values_str, ",")
	sql_query, sql_error := SQLf(insert_sql, to_any_list[[]SQL_Value](sql_values)...)

	return sql_query, sql_error
}

func Exec(ctx context.Context, sql *SQL_Query, runner SQL_Runner) error {
	sql.format()
	_, err := runner.Exec(ctx, sql.query, sql.holder.Args...)
	// defer runner.Close()
	if err != nil {
		log.Default().Println("ERROR ::", err, "\nON SQL :", sql.query)
	}
	return err
}

func Receive[T any](ctx context.Context, sql *SQL_Query, runner SQL_Runner) ([]T, error) {
	sql.format()
	rows, err := runner.Query(ctx, sql.query, sql.holder.Args...)
	defer rows.Close()
	// defer runner.Close()
	if err != nil {
		log.Default().Println("ERROR ::", err, "\nON SQL :", sql.query)
	}
	order := make(map[string]int)
	for i, field_descriptor := range rows.FieldDescriptions() {
		order[field_descriptor.Name] = i
	}
	res, err_rows := pgx.CollectRows(rows, func(row pgx.CollectableRow) (T, error) {
		var res T
		err := BuildSQLResponse(row, &res, order)
		return res, err
	})
	return res, err_rows
}
func ReceiveRows(ctx context.Context, sql *SQL_Query, runner SQL_Runner) (pgx.Rows, error) {
	sql.format()
	rows, err := runner.Query(ctx, sql.query, sql.holder.Args...)
	if err != nil {
		log.Default().Println("ERROR ::", err, "\nON SQL :", sql.query)
	}
	return rows, err
}
func BuildSQLResponse(row pgx.CollectableRow, response_struct any, order map[string]int) error {

	elements := reflect.ValueOf(response_struct).Elem()
	pointers := make([]interface{}, len(row.FieldDescriptions()))
	for i := 0; i < elements.NumField(); i++ {
		tag := elements.Type().Field(i).Tag
		key, ok := tag.Lookup("sql")
		if !ok {
			continue
		}
		if index, ok := order[key]; ok {
			pointers[index] = elements.Field(i).Addr().Interface()
		}
	}
	err := row.Scan(pointers...)
	return err
}

func sql_value_from_struct[T any](obj T) ([]SQL_Value, error) {
	elements := reflect.ValueOf(&obj).Elem()
	values := make([]SQL_Value, 0, elements.NumField())
	for i := 0; i < elements.NumField(); i++ {
		tag := elements.Type().Field(i).Tag
		key, ok := tag.Lookup("sql")
		if !ok {
			continue
		}
		nullable := false
		if strings.Contains(key, "nullable") {
			key = strings.Replace(key, "nullable", "", -1)
			nullable = true
		}
		key = strings.Replace(key, ",", "", -1)

		if nullable && elements.Field(i).IsZero() {
			values = append(values, SQL_Value{
				Name:     key,
				Nullable: nullable,
				Is_null:  true,
			})
		} else {
			switch elements.Field(i).Interface().(type) {
			case map[string]interface{}:
				json, err := json.Marshal(elements.Field(i).Interface())
				if err != nil {
					return []SQL_Value{}, err
				}
				values = append(values,
					SQL_Value{
						Name:     key,
						Value:    string(json),
						Nullable: nullable,
						Is_null:  false,
					})
			default:
				values = append(values,
					SQL_Value{
						Name:     key,
						Value:    elements.Field(i).Interface(),
						Nullable: nullable,
						Is_null:  false,
					})
			}
		}
	}
	if len(values) == 0 {
		return []SQL_Value{}, fmt.Errorf("Invalid struct")
	}
	return values, nil
}

func to_any_list[T any](arr []T) []any {
	any_arr := make([]any, len(arr))
	for i, e := range arr {
		any_arr[i] = e
	}
	return any_arr
}
