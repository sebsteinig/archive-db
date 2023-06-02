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

type DefaultParameters struct {
	Config_name        string  `param:"config_name"`
	Extension          string  `param:"extension" `
	Lossless           bool    `param:"lossless" `
	Nan_value_encoding int     `param:"nan_value_encoding" `
	Threshold          float64 `param:"threshold" `
	Chunks             int     `param:"chunks"`
	Rx                 float64 `param:"rx"`
	Ry                 float64 `param:"ry"`
}

func QueryExperiment(c *fiber.Ctx, pool *pgxpool.Pool) error {
	type Param struct {
		For string `param:"for"`
	}
	param := new(Param)
	query_parameters, err := utils.BuildQueryParameters(c, param)
	if err != nil {
		log.Default().Println("error :", err)
		return err
	}
	pl := new(utils.Placeholder)
	pl.Build(0, 1)
	labels_sql := utils.ILikeBuilder{
		Key:   "labels",
		Value: query_parameters["For"],
	}.Build(pl)

	sql := fmt.Sprintf(`
		SELECT 
			labels
		FROM table_labels
		WHERE %s
		GROUP BY labels
		`, labels_sql)
	rows, err := pool.Query(context.Background(), sql, pl.Args...)
	if err != nil {
		log.Default().Println("Unable to query:", sql, "error :", err)
		return err
	}
	defer rows.Close()

	type SQLResponse struct {
		Labels string `sql:"labels" json:"labels"`
	}

	responses, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (SQLResponse, error) {
		var res SQLResponse
		err := utils.BuildSQLResponse(row, &res)
		return res, err
	})
	if err != nil {
		return err
	}
	return c.JSON(responses)
}

type SearchResponse struct {
	Created_at          time.Time `sql:"created_at" json:"created_at"`
	Config_name         string    `sql:"config_name" json:"config_name"`
	Exp_id              string    `sql:"exp_id" json:"exp_id"`
	Available_variables []string  `sql:"available_variables" json:"available_variables"`
}

func retrieveQueryResponse(rows pgx.Rows) ([]SearchResponse, error) {
	responses, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (SearchResponse, error) {
		var res SearchResponse
		err := utils.BuildSQLResponse(row, &res)
		return res, err
	})
	return responses, err
}

func searchExperimentWith(defaults_parameters utils.QueryParameters, labels []string, c *fiber.Ctx, pool *pgxpool.Pool) error {
	pl := new(utils.Placeholder)
	pl.Build(0, 9+len(labels))
	param_sql := ""

	if len(defaults_parameters) > 0 {
		param_builder := utils.AndBuilder{
			Value: []utils.SqlBuilder{},
		}
		for key, value := range defaults_parameters {
			param_builder.And(utils.EqualBuilder{
				Key:   strings.ToLower(key),
				Value: value,
			})
		}
		param_sql += " AND " + param_builder.Build(pl)
	}
	param_builder := utils.OrBuilder{
		Value: []utils.SqlBuilder{},
	}
	for _, label := range labels {
		param_builder.Or(utils.EqualBuilder{
			Key:   "table_labels.labels",
			Value: label,
		})
	}
	labels_sql := ""
	if len(param_builder.Value) > 0 {
		labels_sql = "WHERE " + param_builder.Build(pl)
	}
	sql := fmt.Sprintf(`
		WITH valid_exp AS (
			SELECT exp_id
			FROM table_labels
			%s
		)
		
		SELECT 
			table_exp.exp_id,
			created_at,
			config_name,
			ARRAY_AGG(join_nimbus_execution_variables.variable_name) as available_variables
		
		FROM table_nimbus_execution 
		
		INNER JOIN valid_exp
			ON table_nimbus_execution.exp_id = valid_exp.exp_id
		INNER JOIN join_nimbus_execution_variables
			ON table_nimbus_execution.id = join_nimbus_execution_variables.id_nimbus_execution
			%s
		INNER JOIN table_exp
			ON table_nimbus_execution.exp_id = table_exp.exp_id
		GROUP BY id,table_exp.exp_id
		ORDER BY created_at DESC;
	`, labels_sql, param_sql)
	rows, err := pool.Query(context.Background(), sql, pl.Args...)
	if err != nil {
		log.Default().Println("Unable to query:", sql, "error :", err)
		return err
	}
	defer rows.Close()
	responses, err := retrieveQueryResponse(rows)
	if err != nil {
		return err
	}
	return c.JSON(responses)
}

func SearchExperimentLike(c *fiber.Ctx, pool *pgxpool.Pool) error {

	default_param := new(DefaultParameters)
	query_parameters, err := utils.BuildQueryParameters(c, default_param)
	if err != nil {
		log.Default().Println("error :", err)
		return err
	}

	type LabelsParam struct {
		Labels []string `param:"with"`
	}
	labelsParam := new(LabelsParam)
	labels_parameters, err := utils.BuildQueryParameters(c, labelsParam)
	if labels, ok := labels_parameters["Labels"]; ok {
		return searchExperimentWith(query_parameters, labels.([]string), c, pool)
	}
	type LikeParam struct {
		Like []string `param:"like"`
	}
	likeParam := new(LikeParam)
	like_parameter, err := utils.BuildQueryParameters(c, likeParam)
	param_sql := ""

	param_builder := utils.AndBuilder{
		Value: []utils.SqlBuilder{},
	}
	if len(query_parameters) > 0 {
		for key, value := range query_parameters {
			param_builder.And(utils.EqualBuilder{
				Key:   strings.ToLower(key),
				Value: value,
			})
		}
	}
	if like, ok := like_parameter["Like"]; ok {
		param_builder.And(utils.EqualBuilder{
			Key:   "exp_id",
			Value: like,
		})
	}
	pl := new(utils.Placeholder)
	pl.Build(0, 9)
	if len(param_builder.Value) > 0 {
		param_sql += " AND " + param_builder.Build(pl)
	}
	sql := fmt.Sprintf(`
		SELECT 
		
			table_exp.exp_id,
			created_at,
			config_name,
			ARRAY_AGG(join_nimbus_execution_variables.variable_name) as available_variables

		FROM table_nimbus_execution

		INNER JOIN join_nimbus_execution_variables
			ON table_nimbus_execution.id = join_nimbus_execution_variables.id_nimbus_execution
			%s
		INNER JOIN table_exp
			ON table_nimbus_execution.exp_id = table_exp.exp_id

		GROUP BY id,table_exp.exp_id
		ORDER BY created_at DESC;
	`, param_sql)
	rows, err := pool.Query(context.Background(), sql, pl.Args...)
	if err != nil {
		log.Default().Println("Unable to query:", sql, "error :", err)
		return err
	}
	defer rows.Close()
	responses, err := retrieveQueryResponse(rows)
	if err != nil {
		return err
	}
	return c.JSON(responses)
}
