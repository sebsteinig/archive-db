package routes

import (
	"archive-api/utils/services"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildSearchRoutes(app *fiber.App, pool *pgxpool.Pool) {

	search_routes := app.Group("/search")

	search_routes.Get("/", func(c *fiber.Ctx) error {
		return services.SearchExperimentLike(c, pool)
	})

}
