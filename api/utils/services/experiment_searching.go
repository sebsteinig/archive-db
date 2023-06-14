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
		//for _, column := range row.RawValues() {
		//	fmt.Println(string(column))
		//}
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
			Value: strings.ToLower(label),
		})
	}
	labels_sql := ""
	if len(param_builder.Value) > 0 {
		labels_sql += "WHERE " + param_builder.Build(pl)
	}
	sql := fmt.Sprintf(`
		WITH valid_exp AS (
			SELECT exp_id
			FROM table_labels
			%s
			GROUP BY exp_id
		)
		
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
		
		INNER JOIN valid_exp
			ON table_nimbus_execution.exp_id = valid_exp.exp_id
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

func SearchExperimentForPublication(c *fiber.Ctx, pool *pgxpool.Pool) error {
	type PublicationParam struct {
		Title         string `json:"title" sql:"title" param:"title"`
		Authors_short string `json:"authors_short" sql:"authors_short" param:"authors_short"`
		//Authors_full  string `json:"authors_full" sql:"authors_full" param:"authors"`
		Journal      string `json:"journal" sql:"journal" param:"journal"`
		Owner_name   string `json:"owner_name" sql:"owner_name"`
		Owner_email  string `json:"owner_email" sql:"owner_email"`
		Abstract     string `json:"abstract" sql:"abstract"`
		Brief_desc   string `json:"brief_desc" sql:"brief_desc"`
		Authors_full string `json:"authors_full" sql:"authors_full"`
		Year         int    `json:"year" sql:"year"`
	}
	publication_param := new(PublicationParam)
	query_parameters, err := utils.BuildQueryParameters(c, publication_param)
	if err != nil {
		log.Default().Println("error :", err)
		return err
	}
	if len(query_parameters) == 0 {
		return fmt.Errorf("some parameters must be specified")
	}
	param_builder := utils.OrBuilder{
		Value: []utils.SqlBuilder{},
	}
	for key, value := range query_parameters {
		param_builder.Or(utils.FullLikeBuilder{
			Key:   strings.ToLower(key),
			Value: value,
		})
	}

	pl := new(utils.Placeholder)
	pl.Build(0, len(query_parameters))

	filters := param_builder.Build(pl)
	if len(query_parameters) > 0 {
		filters = " AND " + filters
	}
	sql := fmt.Sprintf(`
		SELECT 
			ARRAY_AGG(join_publication_exp.exp_id) as exps,
			table_publication.title,
			table_publication.journal,
			table_publication.owner_name,
			table_publication.owner_email,
			--table_publication.abstract,
			--table_publication.brief_desc,
			table_publication.year,
			table_publication.authors_full,
			table_publication.authors_short
		FROM table_publication
		INNER JOIN (
		select 
			join_publication_exp.publication_id
		from 
			join_publication_exp
		except
		select
			join_publication_exp.publication_id
		from 
			join_publication_exp
		where
			join_publication_exp.requested_exp_id is not NULL or join_publication_exp.exp_id is NULL
		) as res
		ON res.publication_id = table_publication.id 
		INNER JOIN join_publication_exp
		ON join_publication_exp.publication_id = table_publication.id %s
		GROUP BY table_publication.id,table_publication.title,join_publication_exp.publication_id
	`, filters)
	rows, err := pool.Query(context.Background(), sql, pl.Args...)
	if err != nil {
		log.Default().Println("Unable to query:", sql, "error :", err)
		return err
	}
	defer rows.Close()
	type Response struct {
		Exps          []string `sql:"exps"`
		Title         string   `sql:"title"`
		Journal       string   `sql:"journal"`
		Owner_name    string   `sql:"owner_name"`
		Owner_email   string   `sql:"owner_email"`
		Abstract      string   `sql:"abstract"`
		Brief_desc    string   `sql:"brief_desc"`
		Year          int      `sql:"year"`
		Authors_full  string   `sql:"authors_full"`
		Authors_short string   `sql:"authors_short"`
	}
	responses, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Response, error) {
		var res Response
		err := utils.BuildSQLResponse(row, &res)
		return res, err
	})
	if err != nil {
		log.Default().Println(err)
		return err
	}
	return c.JSON(responses)
}
