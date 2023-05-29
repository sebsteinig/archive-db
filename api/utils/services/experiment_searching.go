package services

import (
	"archive-api/utils"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SearchExperimentLike(c *fiber.Ctx, pool *pgxpool.Pool) error {
	pl := new(utils.Placeholder)
	pl.Build(0, 9)
	params := make(utils.Params)
	params.ParseParams(c, "like", "config_name", "extension", "lossless", "threshold", "rx", "ry", "chunks")

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
