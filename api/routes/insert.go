package routes

import (
	"archive-api/utils"
	"archive-api/utils/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildInsertRoutes(app *fiber.App, pool *pgxpool.Pool) {

	insert_routes := app.Group("/insert")

	insert_routes.Get("/clean", func(c *fiber.Ctx) error {
		return services.Clean(pool)
	})
	insert_routes.Post("/:id", func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		request := new(utils.Request)
		id := c.Params("id")
		if err := c.BodyParser(request); err != nil {
			log.Default().Println(err)
			return err
		}
		return services.AddVariablesWithExp(id, request, pool)
	})
}
