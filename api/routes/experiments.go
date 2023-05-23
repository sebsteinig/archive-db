package routes

import (
	"archive-api/utils"
	"context"
	"database/sql"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Variables struct {
	Variables [12]utils.Variable `json:"variables"`
}

func getExperimentByID(id string) error {
	log.Default().Println("exp id :", id)
	return nil
}

func addVariablesWithExp(exp_id string, variables []utils.Variable, pool *pgxpool.Pool, psql *squirrel.StatementBuilderType) error {
	if err := pgx.BeginFunc(context.Background(), pool,
		func(tx pgx.Tx) error {
			upsert_exp, arg, err_upsert_exp := psql.Insert("table_experiments").Columns("exp_id").Values(exp_id).ToSql()
			if err_upsert_exp != nil {
				log.Default().Println("pqsl sql error :", err_upsert_exp)
				return err_upsert_exp
			}
			_, err_upsert_exp_cmd := tx.Exec(context.Background(), upsert_exp+" ON CONFLICT (exp_id) DO NOTHING", arg...)
			if err_upsert_exp_cmd != nil {
				log.Default().Println("tx exec error :", upsert_exp+"ON CONFLICT (exp_id) DO NOTHING", err_upsert_exp_cmd)
				return err_upsert_exp_cmd
			}
			_, err := tx.CopyFrom(
				context.Background(),
				pgx.Identifier{"table_variables"},
				[]string{
					"name",
					"exp_id",
					"paths_ts",
					"paths_mean",
					"config_name",
					"levels",
					"timesteps",
					"xsize",
					"xfirst",
					"xinc",
					"ysize",
					"yfirst",
					"yinc",
					"extension",
					"lossless",
					"nan_value_encoding",
					"threshold",
					"chunks",
					"rx",
					"ry",
					"metadata"},
				pgx.CopyFromSlice(len(variables), func(i int) ([]any, error) {
					rx := sql.NullFloat64{
						Float64: variables[i].Rx,
						Valid:   variables[i].Rx != 0,
					}
					ry := sql.NullFloat64{
						Float64: variables[i].Rx,
						Valid:   variables[i].Rx != 0,
					}
					return []any{
						variables[i].Name,
						variables[i].Exp_id,
						variables[i].Paths_ts,
						variables[i].Paths_mean,
						variables[i].Config_name,
						variables[i].Levels,
						variables[i].Timesteps,
						variables[i].Xsize,
						variables[i].Xfirst,
						variables[i].Xinc,
						variables[i].Ysize,
						variables[i].Yfirst,
						variables[i].Yinc,
						variables[i].Extension,
						variables[i].Lossless,
						variables[i].Nan_value_encoding,
						variables[i].Threshold,
						variables[i].Chunks,
						rx,
						ry,
						variables[i].Metadata,
					}, nil
				}),
			)
			if err != nil {
				log.Default().Println("copy from error :", err)
				return err
			}
			return nil
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

		variables := new(Variables)
		id := c.Params("id")
		if err := c.BodyParser(variables); err != nil {
			log.Default().Println(err)
			return err
		}
		return addVariablesWithExp(id, variables.Variables[:], pool, psql)
	})
}
