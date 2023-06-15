package services

import (
	"archive-api/utils"
	"archive-api/utils/sql"
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
	default_param := new(DefaultParameters)
	query_parameters, err := utils.BuildQueryParameters(c, default_param)
	if err != nil {
		log.Default().Println("ERROR <GetExperimentByID>")
		log.Default().Println("error :", err)
		return err
	}
	param_builder := sql.AndBuilder{
		Value:      []sql.SqlBuilder{},
		And_Prefix: true,
	}
	for key, value := range query_parameters {
		param_builder.And(sql.EqualBuilder{
			Key:   strings.ToLower(key),
			Value: value,
		})
	}

	type VariablesParams struct {
		Variables []string `param:"vars"`
	}
	variable_params := new(VariablesParams)
	_, err_v := utils.BuildQueryParameters(c, variable_params)
	if err_v != nil {
		log.Default().Println("ERROR <GetExperimentByID>")
		log.Default().Println("error :", err)
		return err
	}

	params_vars_builder := sql.OrBuilder{
		Value:      []sql.SqlBuilder{},
		And_Prefix: true,
	}
	for _, value := range variable_params.Variables {
		params_vars_builder.Or(sql.EqualBuilder{
			Key:   "variable_name",
			Value: value,
		})
	}

	query, err := sql.SQLf(`WITH nimbus_run AS 
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
	ON table_variable.id = joined.variable_id %s`, id, param_builder, params_vars_builder)

	if err != nil {
		log.Default().Println("ERROR <GetExperimentByID>")
		return err
	}
	responses, err := sql.Receive[Response](context.Background(), &query, pool)
	if err != nil {
		log.Default().Println("ERROR <GetExperimentByID>")
		return err
	}
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
	default_param := new(DefaultParameters)
	query_parameters, err := utils.BuildQueryParameters(c, default_param)
	if err != nil {
		log.Default().Println("ERROR <GetExperimentsByIDs>")
		log.Default().Println("error :", err)
		return err
	}

	type Ids struct {
		Ids []string `param:"ids,required"`
	}
	ids_param := new(Ids)
	_, err = utils.BuildQueryParameters(c, ids_param)
	if err != nil {
		log.Default().Println("ERROR <GetExperimentsByIDs>")
		log.Default().Println("error :", err)
		return err
	}

	in_builder := sql.InBuilder{
		Key:   "exp_id",
		Value: toAnyList(ids_param.Ids),
	}

	param_builder := sql.AndBuilder{
		Value:      []sql.SqlBuilder{},
		And_Prefix: true,
	}
	for key, value := range query_parameters {
		param_builder.And(sql.EqualBuilder{
			Key:   strings.ToLower(key),
			Value: value,
		})
	}

	type VariablesParams struct {
		Variables []string `param:"vars"`
	}
	variable_params := new(VariablesParams)
	_, err_v := utils.BuildQueryParameters(c, variable_params)
	if err_v != nil {
		log.Default().Println("ERROR <GetExperimentsByIDs>")
		log.Default().Println("error :", err)
		return err
	}

	params_vars_builder := sql.OrBuilder{
		Value:      []sql.SqlBuilder{},
		And_Prefix: true,
	}
	for _, value := range variable_params.Variables {
		params_vars_builder.Or(sql.EqualBuilder{
			Key:   "variable_name",
			Value: value,
		})
	}

	query, err := sql.SQLf(`WITH nimbus_run AS 
	(
		SELECT *
		FROM table_nimbus_execution 
		WHERE %s %s
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
	ON table_variable.id = joined.variable_id %s`, in_builder, param_builder, params_vars_builder)
	if err != nil {
		log.Default().Println("ERROR <GetExperimentsByIDs>")
		return err
	}

	rows, err := sql.ReceiveRows(context.Background(), &query, pool)
	if err != nil {
		log.Default().Println("ERROR <GetExperimentsByIDs>")
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
		log.Default().Println("ERROR <GetExperimentsByIDs>")
		log.Default().Println("map failed, error :", err_map)
		return err_map
	}
	return c.JSON(map_exp)
}
