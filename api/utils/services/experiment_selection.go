package services

import (
	"archive-api/utils"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Response struct {
	VariableName       string                 `sql:"variable_name" json:"variable_name"`
	Path_ts            []string               `sql:"paths_ts" json:"paths_ts"`
	Path_mean          []string               `sql:"paths_mean" json:"paths_mean"`
	Levels             int                    `sql:"levels" json:"levels"`
	Timesteps          int                    `sql:"timesteps" json:"timesteps"`
	Xsize              int                    `sql:"xsize" json:"xsize"`
	Xfirst             float32                `sql:"xfirst" json:"xfirst"`
	Yinc               float32                `sql:"xinc" json:"xinc"`
	Ysize              int                    `sql:"ysize" json:"ysize"`
	Yfirst             float32                `sql:"yfirst" json:"yfirst"`
	Xinc               float32                `sql:"yinc" json:"yinc"`
	Metadata           map[string]interface{} `sql:"metadata" json:"metadata"`
	Created_at         time.Time              `sql:"created_at" json:"created_at"`
	Config_name        string                 `sql:"config_name" json:"config_name"`
	Extension          string                 `sql:"extension" json:"extension"`
	Lossless           bool                   `sql:"lossless" json:"lossless"`
	Nan_value_encoding int                    `sql:"nan_value_encoding" json:"nan_value_encoding"`
	Chunks             int                    `sql:"chunks" json:"chunks"`
	Rx                 float64                `sql:"rx" json:"rx"`
	Ry                 float64                `sql:"ry" json:"ry"`
	Exp_id             string                 `sql:"exp_id" json:"exp_id"`
	Threshold          float32                `sql:"threshold" json:"threshold"`
}

type SelectDefaultParameters struct {
	Config_name        string  `param:"config_name"`
	Extension          string  `param:"extension" `
	Lossless           bool    `param:"lossless" `
	Nan_value_encoding int     `param:"nan_value_encoding" `
	Threshold          float64 `param:"threshold" `
	Chunks             int     `param:"chunks"`
	Rx                 float64 `param:"rx"`
	Ry                 float64 `param:"ry"`
}

func GetExperimentByID(id string, c *fiber.Ctx, pool *pgxpool.Pool) error {
	pl := new(utils.Placeholder)
	pl.Build(0, 8)

	default_param := new(DefaultParameters)
	query_parameters, err := utils.BuildQueryParameters(c, default_param)
	if err != nil {
		log.Default().Println("error :", err)
		return err
	}
	params_sql := ""
	if len(query_parameters) > 0 {
		param_builder := utils.AndBuilder{
			Value: []utils.SqlBuilder{},
		}
		for key, value := range query_parameters {
			param_builder.And(utils.EqualBuilder{
				Key:   strings.ToLower(key),
				Value: value,
			})
		}
		params_sql += " AND " + param_builder.Build(pl)
	}

	sql := fmt.Sprintf(`WITH nimbus_run AS 
	(
		SELECT *
		FROM table_nimbus_execution 
		WHERE exp_id = %s %s
		ORDER BY created_at desc
		LIMIT 1
	)
	SELECT 
		name AS variable_name,
		paths_ts,
		paths_mean,levels,
		timesteps,
		xsize,
		xfirst,
		xinc,
		ysize,
		yfirst,
		yinc,
		metadata,
		created_at,
		config_name,
		extension,
		lossless,
		nan_value_encoding,
		chunks,
		rx,
		ry,
		exp_id,
		threshold
	FROM table_variable
	INNER JOIN 
		( 
			SELECT * 
			FROM join_nimbus_execution_variables
			INNER JOIN nimbus_run 
			ON join_nimbus_execution_variables.id_nimbus_execution = nimbus_run.id
		) AS joined
	ON table_variable.id = joined.variable_id`, pl.Get(id), params_sql)
	rows, err := pool.Query(context.Background(), sql, pl.Args...)
	if err != nil {
		log.Default().Println("Unable to query:", sql, "error :", err)
		return err
	}
	defer rows.Close()
	responses, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Response, error) {
		var res Response
		err := utils.BuildSQLResponse(row, &res)
		return res, err
	})
	return c.JSON(responses)
}
func toAnyList[T any](input []T) []any {
	list := make([]any, len(input))
	for i, v := range input {
		list[i] = v
	}
	return list
}
func GetExperimentsByIDs(c *fiber.Ctx, pool *pgxpool.Pool) error {

	pl := new(utils.Placeholder)
	pl.Build(0, 8)

	default_param := new(DefaultParameters)
	query_parameters, err := utils.BuildQueryParameters(c, default_param)
	if err != nil {
		log.Default().Println("error :", err)
		return err
	}
	params_sql := ""

	type Ids struct {
		Ids []string `param:"ids,required"`
	}
	ids_param := new(Ids)
	_, err = utils.BuildQueryParameters(c, ids_param)
	if err != nil {
		log.Default().Println("error :", err)
		return err
	}

	params_sql += utils.InBuilder{
		Key:   "exp_id",
		Value: toAnyList(ids_param.Ids),
	}.Build(pl)

	if len(query_parameters) > 0 {
		param_builder := utils.AndBuilder{
			Value: []utils.SqlBuilder{},
		}
		for key, value := range query_parameters {
			param_builder.And(utils.EqualBuilder{
				Key:   strings.ToLower(key),
				Value: value,
			})
		}
		params_sql += " AND " + param_builder.Build(pl)
	}
	fmt.Println(params_sql)
	fmt.Println(pl.Args...)
	sql := fmt.Sprintf(`WITH nimbus_run AS 
	(
		SELECT *
		FROM table_nimbus_execution 
		WHERE ` + params_sql + `
		ORDER BY created_at desc
	)
	SELECT 
		name AS variable_name,
		paths_ts,
		paths_mean,levels,
		timesteps,
		xsize,
		xfirst,
		xinc,
		ysize,
		yfirst,
		yinc,
		metadata,
		created_at,
		config_name,
		extension,
		lossless,
		nan_value_encoding,
		chunks,
		rx,
		ry,
		exp_id,
		threshold
	FROM table_variable
	INNER JOIN 
		( 
			SELECT * 
			FROM join_nimbus_execution_variables
			INNER JOIN nimbus_run 
			ON join_nimbus_execution_variables.id_nimbus_execution = nimbus_run.id
		) AS joined
	ON table_variable.id = joined.variable_id`)
	rows, err := pool.Query(context.Background(), sql, pl.Args...)
	if err != nil {
		fmt.Println(pl.Args...)
		log.Default().Println("Unable to query:", sql, "error :", err)
		return err
	}
	defer rows.Close()
	var map_exp map[string][]Response = make(map[string][]Response)
	var res Response
	_, err_map := pgx.ForEachRow(rows, []any{
		&res.VariableName,
		&res.Path_ts,
		&res.Path_mean,
		&res.Levels,
		&res.Timesteps,
		&res.Xsize,
		&res.Xfirst,
		&res.Yinc,
		&res.Ysize,
		&res.Yfirst,
		&res.Xinc,
		&res.Metadata,
		&res.Created_at,
		&res.Config_name,
		&res.Extension,
		&res.Lossless,
		&res.Nan_value_encoding,
		&res.Chunks,
		&res.Rx,
		&res.Ry,
		&res.Exp_id,
		&res.Threshold,
	}, func() error {
		fmt.Println(res.Exp_id)
		map_exp[res.Exp_id] = append(map_exp[res.Exp_id], res)
		return nil
	})
	if err_map != nil {
		log.Default().Println("map failed, error :", err_map)
		return err_map
	}
	return c.JSON(map_exp)
}
