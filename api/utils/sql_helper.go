package utils

/*
type Placeholder struct {
	idx  int
	Args []interface{}
}

func (pl *Placeholder) Build(length int, cap int) {
	*pl = Placeholder{
		idx:  0,
		Args: make([]interface{}, length, cap),
	}
}

func (pl *Placeholder) Get(arg interface{}) string {
	pl.Args = append(pl.Args, arg)
	pl.idx += 1
	return fmt.Sprintf("$%d", pl.idx)
}

type SqlBuilder interface {
	Build(pl *Placeholder) string
}
type EqualBuilder struct {
	Key   string
	Value any
}

func (builder EqualBuilder) Build(pl *Placeholder) string {
	if reflect.ValueOf(builder.Value).IsZero() {
		return ""
	}
	return fmt.Sprintf("%s = %s", builder.Key, pl.Get(builder.Value))
}

type LikeBuilder struct {
	Key   string
	Value any
}

func (builder LikeBuilder) Build(pl *Placeholder) string {
	if reflect.ValueOf(builder.Value).IsZero() {
		return ""
	}
	return fmt.Sprintf("%s LIKE %s || '%%'", builder.Key, pl.Get(builder.Value))
}

type ILikeBuilder struct {
	Key   string
	Value any
}

func (builder ILikeBuilder) Build(pl *Placeholder) string {
	if reflect.ValueOf(builder.Value).IsZero() {
		return ""
	}
	return fmt.Sprintf("%s ILIKE %s || '%%'", builder.Key, pl.Get(builder.Value))
}

type FullLikeBuilder struct {
	Key   string
	Value any
}

func (builder FullLikeBuilder) Build(pl *Placeholder) string {
	if reflect.ValueOf(builder.Value).IsZero() {
		return ""
	}
	return fmt.Sprintf("%s ILIKE '%%' || %s || '%%'", builder.Key, pl.Get(builder.Value))
}

type InBuilder struct {
	Key   string
	Value []any
}

func (builder InBuilder) Build(pl *Placeholder) string {
	if reflect.ValueOf(builder.Value).IsZero() {
		return ""
	}
	array := make([]string, len(builder.Value))
	for i, value := range builder.Value {
		array[i] = pl.Get(value)
	}
	return fmt.Sprintf("%s IN (%s)", builder.Key, strings.Join(array, ","))
}

type AndBuilder struct {
	Value []SqlBuilder
}

func (builder *AndBuilder) And(sqlb SqlBuilder) {
	builder.Value = append(builder.Value, sqlb)
}
func (builder AndBuilder) Build(pl *Placeholder) string {
	if len(builder.Value) == 0 {
		return ""
	}
	array := make([]string, 0, len(builder.Value))
	for _, value := range builder.Value {
		if sql := value.Build(pl); sql != "" {
			array = append(array, sql)
		}
	}
	return fmt.Sprintf("(%s)", strings.Join(array, " AND "))
}

type OrBuilder struct {
	Value []SqlBuilder
}

func (builder *OrBuilder) Or(sqlb SqlBuilder) {
	builder.Value = append(builder.Value, sqlb)
}
func (builder OrBuilder) Build(pl *Placeholder) string {
	array := make([]string, 0, len(builder.Value))
	for _, value := range builder.Value {
		if sql := value.Build(pl); sql != "" {
			array = append(array, sql)
		}
	}
	return fmt.Sprintf("(%s)", strings.Join(array, " OR "))
}
func BuildSQLInsert[T any](table string, insert_struct T, pl *Placeholder) (string, error) {
	elements := reflect.ValueOf(&insert_struct).Elem()
	fields := make([]string, 0, elements.NumField())
	values := make([]string, 0, elements.NumField())
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
			values = append(values, "NULL")
		} else {
			switch elements.Field(i).Interface().(type) {
			case map[string]interface{}:
				json, err := json.Marshal(elements.Field(i).Interface())
				if err != nil {
					return "", err
				}
				values = append(values, pl.Get(json))
			default:
				values = append(values, pl.Get(elements.Field(i).Interface()))
			}
		}
		fields = append(fields, key)
	}
	if len(fields) == 0 || len(values) == 0 || len(fields) != len(values) {
		return "", fmt.Errorf("Invalid struct")
	}
	sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, table, strings.Join(fields, ","), strings.Join(values, ","))
	return sql, nil
}

func BuildSQLInsertAll[T any](table string, array_struct []T, pl *Placeholder) (string, error) {
	var fields []string
	array_values := make([]string, 0, len(array_struct))
	for _, insert_struct := range array_struct {
		elements := reflect.ValueOf(&insert_struct).Elem()
		fields = make([]string, 0, elements.NumField())
		values := make([]string, 0, elements.NumField())
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
				values = append(values, "NULL")
			} else {
				switch elements.Field(i).Interface().(type) {
				case map[string]interface{}:
					json, err := json.Marshal(elements.Field(i).Interface())
					if err != nil {
						return "", err
					}
					values = append(values, pl.Get(json))
				default:
					values = append(values, pl.Get(elements.Field(i).Interface()))
				}
			}
			fields = append(fields, key)
		}
		if len(fields) == 0 || len(values) == 0 || len(fields) != len(values) {
			return "", fmt.Errorf("Invalid struct")
		}
		array_values = append(array_values, fmt.Sprintf("(%s)", strings.Join(values, ",")))
	}
	sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES %s`, table, strings.Join(fields, ","), strings.Join(array_values, ","))

	return sql, nil
}

func BuildSQLResponse(row pgx.CollectableRow, response_struct any) error {

	elements := reflect.ValueOf(response_struct).Elem()
	order := make(map[string]int)
	for i, field_descriptor := range row.FieldDescriptions() {
		order[field_descriptor.Name] = i
	}
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
*/
