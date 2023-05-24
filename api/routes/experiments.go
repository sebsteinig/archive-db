package routes

import (
	"archive-api/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RequestBody struct {
	Table_nimbus_execution utils.NimbusExecution `json:"table_nimbus_execution"`
	Table_variable         [12]utils.Variable    `json:"table_variable"`
}

type Request struct {
	Request RequestBody `json:"request"`
}

func getExperimentByID(id string) error {
	log.Default().Println("exp id :", id)
	return nil
}

func arrayToString(arr []string) string {
	res := "{"
	for i, v := range arr {
		res += fmt.Sprintf(`"%s"`, v)
		if i < len(arr)-1 {
			res += ","
		}
	}
	res += "}"
	return res
}
func placeholder_wrap(i *int, arg interface{}, args *[]interface{}) string {
	*args = append(*args, fmt.Sprintf("'%v'", arg))
	*i += 1
	return fmt.Sprintf("$%d", *i)
}
func placeholder(i *int, arg interface{}, args *[]interface{}) string {
	*args = append(*args, arg)
	*i += 1
	return fmt.Sprintf("$%d", *i)
}
func addVariablesWithExp(exp_id string, request *Request, pool *pgxpool.Pool, psql *squirrel.StatementBuilderType) error {

	if err := pgx.BeginFunc(context.Background(), pool,
		func(tx pgx.Tx) error {
			idx := 0
			args := make([]interface{}, 0, 144)

			insert_into_table_nimbus := "INSERT INTO table_nimbus_execution (exp_id,config_name,extension,lossless,nan_value_encoding,threshold,chunks,rx,ry) VALUES "

			insert_into_table_nimbus += fmt.Sprintf("(%s,%s,%s,%s,%s,%s,%s,%s,%s)",
				placeholder_wrap(&idx, request.Request.Table_nimbus_execution.Exp_id, &args),
				placeholder_wrap(&idx, request.Request.Table_nimbus_execution.Config_name, &args),
				placeholder_wrap(&idx, request.Request.Table_nimbus_execution.Extension, &args),
				placeholder(&idx, request.Request.Table_nimbus_execution.Lossless, &args),
				placeholder(&idx, request.Request.Table_nimbus_execution.Nan_value_encoding, &args),
				placeholder(&idx, request.Request.Table_nimbus_execution.Threshold, &args),
				placeholder(&idx, request.Request.Table_nimbus_execution.Chunks, &args),
				placeholder(&idx, request.Request.Table_nimbus_execution.Rx, &args),
				placeholder(&idx, request.Request.Table_nimbus_execution.Ry, &args),
			)
			insert_into_table_variable := "INSERT INTO table_variable (name, paths_ts, paths_mean, levels, timesteps, xsize, xfirst, xinc, ysize, yfirst, yinc, metadata) VALUES "
			for i, v := range request.Request.Table_variable {
				metadata, err := json.Marshal(v.Metadata)
				if err != nil {
					return err
				}
				insert_into_table_variable += fmt.Sprintf("(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)",
					placeholder_wrap(&idx, v.Name, &args),
					placeholder(&idx, v.Paths_ts, &args),
					placeholder(&idx, v.Paths_mean, &args),
					placeholder(&idx, v.Levels, &args),
					placeholder(&idx, v.Timesteps, &args),
					placeholder(&idx, v.Xsize, &args),
					placeholder(&idx, v.Xfirst, &args),
					placeholder(&idx, v.Xinc, &args),
					placeholder(&idx, v.Ysize, &args),
					placeholder(&idx, v.Yfirst, &args),
					placeholder(&idx, v.Yinc, &args),
					placeholder(&idx, metadata, &args),
				)
				if i < len(request.Request.Table_variable)-1 {
					insert_into_table_variable += ","
				}
			}
			fmt.Println(psql.Insert("table").Values("test", "t2").ToSql())
			fmt.Println(len(args))
			sql := fmt.Sprintf("WITH nimbus_id AS (%s RETURNING id), var_ids_name AS (%s RETURNING name,id) INSERT INTO join_nimbus_execution_variables SELECT nimbus_id.id AS id_nimbus, var_ids_name.name AS var_name, var_ids_name.id AS var_id FROM var_ids_name CROSS JOIN nimbus_id;", insert_into_table_nimbus, insert_into_table_variable)
			fmt.Println(sql)
			_, err := tx.Exec(context.Background(), sql, args...)
			return err
		},
	); err != nil {
		log.Default().Println("transactions error :", err)
		return err
	}
	return nil
}

func BuildExperimentRoutes(app *fiber.App, pool *pgxpool.Pool, psql *squirrel.StatementBuilderType) {

	experiments_routes := app.Group("/experiment")

	experiments_routes.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return getExperimentByID(id)
	})

	experiments_routes.Post("/:id/add", func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		request := new(Request)
		id := c.Params("id")
		if err := c.BodyParser(request); err != nil {
			log.Default().Println(err)
			return err
		}
		return addVariablesWithExp(id, request, pool, psql)
	})
}
