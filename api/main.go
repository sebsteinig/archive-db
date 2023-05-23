package main

import (
	"archive-api/routes"
	"context"
	"log"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Variable struct {
	Id                 int       `json:"id"`
	Name               string    `json:"name"`
	Exp_id             string    `json:"exp_id"`
	Paths              []string  `json:"paths"`
	Created_at         time.Time `json:"created_at"`
	Config_name        string    `json:"config_name"`
	Levels             int       `json:"levels"`
	Timesteps          int       `json:"timesteps"`
	Xsize              int       `json:"xsize"`
	Xfirst             float32   `json:"xfirst"`
	Xinc               float32   `json:"xinc"`
	Ysize              int       `json:"ysize"`
	Yfirst             float32   `json:"yfirst"`
	Yinc               float32   `json:"yinc"`
	Extension          string    `json:"extension"`
	Lossless           bool      `json:"lossless"`
	Nan_value_encoding int       `json:"nan_value_encoding"`
	Threshold          float32   `json:"threshold"`
	Chunks             int       `json:"chunks"`
	Metadata           string    `json:"metadata"`
}

func main() {
	godotenv.Load(".env")
	app := fiber.New(fiber.Config{AppName: "test"})
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
		os.Exit(1)
	}
	defer pool.Close()

	routes.BuildExperimentRoutes(app)

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("Unable to ping the database:", err)
		os.Exit(1)
	}
	app.Get("/", func(c *fiber.Ctx) error {
		variables_query := squirrel.Select("*").From("table_variables")
		sql, _, _ := variables_query.ToSql()
		rows, err := pool.Query(context.Background(), sql)
		if err != nil {
			log.Fatal("Unable to query:", sql, "error :", err)
		}
		defer rows.Close()
		variables, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Variable, error) {
			var variable Variable
			err := row.Scan(
				&variable.Id,
				&variable.Name,
				&variable.Exp_id,
				&variable.Paths,
				&variable.Created_at,
				&variable.Config_name,
				&variable.Levels,
				&variable.Timesteps,
				&variable.Xsize,
				&variable.Xfirst,
				&variable.Xinc,
				&variable.Ysize,
				&variable.Yfirst,
				&variable.Yinc,
				&variable.Extension,
				&variable.Lossless,
				&variable.Nan_value_encoding,
				&variable.Threshold,
				&variable.Chunks,
				&variable.Metadata,
			)
			if err != nil {
				log.Fatal(err)
			}
			return variable, err
		})

		return c.JSON(variables)
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatal("Unable to listen :", err)
		os.Exit(1)
	}
}
