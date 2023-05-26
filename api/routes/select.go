package routes

import (
	"archive-api/utils/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildSelectRoutes(app *fiber.App, pool *pgxpool.Pool) {

	select_routes := app.Group("/select")
	select_routes.Use(cache.New(cache.Config{
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}))

	select_routes.Get("/collection", func(c *fiber.Ctx) error {
		return services.GetExperimentsByIDs(c, pool)
	})

	select_routes.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return services.GetExperimentByID(id, c, pool)
	})

}
