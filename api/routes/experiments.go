package routes

import (
	"archive-api/utils"
	"archive-api/utils/services"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildExperimentRoutes(app *fiber.App, pool *pgxpool.Pool, psql *squirrel.StatementBuilderType) {

	experiments_routes := app.Group("/experiment")

	experiments_routes.Get("/clean", func(c *fiber.Ctx) error {
		return services.Clean(pool)
	})

	experiments_routes.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return services.GetExperimentByID(id, c, pool, psql)
	})
	experiments_routes.Post("/:id/add", func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		request := new(utils.Request)
		id := c.Params("id")
		if err := c.BodyParser(request); err != nil {
			log.Default().Println(err)
			return err
		}
		return services.AddVariablesWithExp(id, request, pool, psql)
	})
}
