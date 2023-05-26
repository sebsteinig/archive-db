package routes

import (
	"archive-api/utils/services"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildSelectRoutes(app *fiber.App, pool *pgxpool.Pool) {

	select_routes := app.Group("/select")

	select_routes.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return services.GetExperimentByID(id, c, pool)
	})

	select_routes.Get("/collection", func(c *fiber.Ctx) error {
		return services.GetExperimentsByIDs(c, pool)
	})

}
