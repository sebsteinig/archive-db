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

// @Description search for an experiment
// @Param for query string false "string for"
// @Success 200 {object} object "label"
// @Router /search/looking [get]
func QueryExperiment(c *fiber.Ctx, pool *pgxpool.Pool) error {
	type Param struct {
		For string `param:"for"`
	}
	param := new(Param)
	query_parameters, err := utils.BuildQueryParameters(c, param)
	if err != nil {
		log.Default().Println("ERROR <QueryExperiment>")
		log.Default().Println("error :", err)
		return err
	}
	query, err := sql.SQLf(`
		SELECT 
			labels
		FROM table_labels
		WHERE %s
		GROUP BY labels
		`,
		sql.ILikeBuilder{
			Key:   "labels",
			Value: query_parameters["For"],
		})
	type SQLResponse struct {
		Labels string `sql:"labels" json:"labels"`
	}
	responses, err := sql.Receive[SQLResponse](context.Background(), &query, pool)
	if err != nil {
		log.Default().Println("ERROR <QueryExperiment>")
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

func searchExperimentWith(defaults_parameters utils.QueryParameters, labels []string, c *fiber.Ctx, pool *pgxpool.Pool) error {

	param_builder := sql.AndBuilder{
		Value:      []sql.SqlBuilder{},
		And_Prefix: true,
	}
	for key, value := range defaults_parameters {
		param_builder.And(sql.EqualBuilder{
			Key:   strings.ToLower(key),
			Value: value,
		})
	}

	label_builder := sql.OrBuilder{
		Value:        []sql.SqlBuilder{},
		Where_Prefix: true,
	}
	for _, label := range labels {
		label_builder.Or(sql.EqualBuilder{
			Key:   "table_labels.labels",
			Value: strings.ToLower(label),
		})
	}

	query, err := sql.SQLf(`
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
	`, label_builder, param_builder)
	if err != nil {
		log.Default().Println("ERROR <searchExperimentWith>")
		return err
	}
	responses, err := sql.Receive[SearchResponse](context.Background(), &query, pool)
	if err != nil {
		log.Default().Println("ERROR <searchExperimentWith>")
		return err
	}
	return c.JSON(responses)
}

// @Description search for an experiment based on the first character(s)
// @Param like query string false "string like"
// @Success 200 {object} object "experiment"
// @Router /search/ [get]
func SearchExperimentLike(c *fiber.Ctx, pool *pgxpool.Pool) error {

	default_param := new(DefaultParameters)
	query_parameters, err := utils.BuildQueryParameters(c, default_param)
	if err != nil {
		log.Default().Println("ERROR <SearchExperimentLike>")
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

	if like, ok := like_parameter["Like"]; ok {
		param_builder.And(sql.EqualBuilder{
			Key:   "exp_id",
			Value: like,
		})
	}

	query, err := sql.SQLf(`
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
	`, param_builder)
	if err != nil {
		log.Default().Println("ERROR <SearchExperimentLike>")
		return err
	}
	responses, err := sql.Receive[SearchResponse](context.Background(), &query, pool)
	if err != nil {
		log.Default().Println("ERROR <SearchExperimentLike>")
		return err
	}
	return c.JSON(responses)
}

// @Description search for a publication by title, author or journal (at least one these parameters has to be specified)
// @Param title query string false "string title"
// @Param authors_short query string false "string author"
// @Param journal query string false "string journal"
// @Param owner_name query string false "string owner name"
// @Param owner_email query string false "string owner email"
// @Param abstract query string false "string abstract"
// @Param brief_desc query string false "string brief desccription"
// @Param authors_full query string false "string all authors"
// @Param year query int false "int year"
// @Success 200 {object} object "experiment"
// @Router /search/publication [get]
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
	param_builder := sql.OrBuilder{
		Value:      []sql.SqlBuilder{},
		And_Prefix: true,
	}
	for key, value := range query_parameters {
		param_builder.Or(sql.FLikeBuilder{
			Key:   strings.ToLower(key),
			Value: value,
		})
	}
	query, err := sql.SQLf(`
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
	`, param_builder)

	if err != nil {
		log.Default().Println("ERROR <SearchExperimentForPublication>")
		return err
	}
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
	responses, err := sql.Receive[Response](context.Background(), &query, pool)
	if err != nil {
		log.Default().Println("ERROR <SearchExperimentForPublication>")
		return err
	}
	return c.JSON(responses)
}
