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

func QueryExperiment(c *fiber.Ctx, pool *pgxpool.Pool) error {
	params := make(utils.Params)
	params.ParseParams(c, "for")

	query, ok := params["query"]
	if !ok {
		return fmt.Errorf("for clause must be specified when looking")
	}
	pl := new(utils.Placeholder)
	pl.Build(0, 9)
	queries := query.Value.([]string)

	labels := make([]string, len(queries))
	for i, q := range queries {
		labels[i] = fmt.Sprintf("labels ILIKE %s || '%%'", pl.Get(q))
	}
	labels_sql := strings.Join(labels, " OR ")
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
	responses, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (string, error) {
		var res string
		err := row.Scan(&res)
		if err != nil {
			log.Default().Println(err)
		}
		return res, err
	})
	if err != nil {
		return err
	}
	return c.JSON(responses)
}

func searchExperimentWith(params *utils.Params, labels []string, c *fiber.Ctx, pool *pgxpool.Pool) error {
	pl := new(utils.Placeholder)
	pl.Build(0, 9+len(labels))
	params_sql := params.ParamToSql(pl)
	if params_sql != "" {
		params_sql = " AND " + params_sql
	}
	labels_str_array := make([]string, len(labels))
	for i, q := range labels {
		labels_str_array[i] = fmt.Sprintf("table_labels.labels = %s", pl.Get(q))
	}
	labels_sql := ""
	if len(labels_str_array) > 0 {
		labels_sql = "WHERE " + strings.Join(labels_str_array, " AND ")
	}
	sql := fmt.Sprintf(`
		WITH valid_exp AS (
			SELECT exp_id
			FROM table_labels
			%s
		)
		
		SELECT 
		table_nimbus_execution.exp_id,
		created_at,
		config_name,
		ARRAY_AGG(join_nimbus_execution_variables.variable_name) as available_variables
		
		FROM table_nimbus_execution 
			INNER JOIN valid_exp
			ON table_nimbus_execution.exp_id = valid_exp.exp_id
			INNER JOIN join_nimbus_execution_variables
			ON table_nimbus_execution.id = join_nimbus_execution_variables.id_nimbus_execution
			%s
			GROUP BY id,table_nimbus_execution.exp_id
		
		ORDER BY created_at DESC;
	`, labels_sql, params_sql)
	rows, err := pool.Query(context.Background(), sql, pl.Args...)
	if err != nil {
		log.Default().Println("Unable to query:", sql, "error :", err)
		return err
	}
	defer rows.Close()

	type Response struct {
		Created_at          time.Time `json:"created_at"`
		Config_name         string    `json:"config_name"`
		Exp_id              string    `json:"exp_id"`
		Available_variables []string  `json:"available_variables"`
	}
	responses, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Response, error) {
		var res Response
		err := row.Scan(
			&res.Exp_id,
			&res.Created_at,
			&res.Config_name,
			&res.Available_variables,
		)
		if err != nil {
			log.Default().Println(err)
		}
		return res, err
	})
	if err != nil {
		return err
	}
	return c.JSON(responses)
}

func SearchExperimentLike(c *fiber.Ctx, pool *pgxpool.Pool) error {
	params := make(utils.Params)
	params.ParseParams(c, "like", "with", "config_name", "extension", "lossless", "threshold", "rx", "ry", "chunks")
	labels, ok := params["labels"]
	if ok {
		return searchExperimentWith(&params, labels.Value.([]string), c, pool)
	}
	pl := new(utils.Placeholder)
	pl.Build(0, 9)

	params_sql := params.ParamToSql(pl)
	sql := fmt.Sprintf(`
		SELECT 
		
		exp_id,
		created_at,
		config_name,
		ARRAY_AGG(join_nimbus_execution_variables.variable_name) as available_variables

		FROM table_nimbus_execution 
		INNER JOIN join_nimbus_execution_variables
		ON table_nimbus_execution.id = join_nimbus_execution_variables.id_nimbus_execution
		 AND %s
		GROUP BY id,exp_id
		ORDER BY created_at DESC;
	`, params_sql)
	rows, err := pool.Query(context.Background(), sql, pl.Args...)
	if err != nil {
		log.Default().Println("Unable to query:", sql, "error :", err)
		return err
	}
	defer rows.Close()

	type Response struct {
		Created_at          time.Time `json:"created_at"`
		Config_name         string    `json:"config_name"`
		Exp_id              string    `json:"exp_id"`
		Available_variables []string  `json:"available_variables"`
	}
	responses, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Response, error) {
		var res Response
		err := row.Scan(
			&res.Exp_id,
			&res.Created_at,
			&res.Config_name,
			&res.Available_variables,
		)
		if err != nil {
			log.Default().Println(err)
		}
		return res, err
	})
	if err != nil {
		return err
	}
	return c.JSON(responses)
}
