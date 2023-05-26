package routes

import (
	"archive-api/utils/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildSearchRoutes(app *fiber.App, pool *pgxpool.Pool) {

	search_routes := app.Group("/search")
	search_routes.Use(cache.New(cache.Config{
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}))

	search_routes.Get("/", func(c *fiber.Ctx) error {
		return services.SearchExperimentLike(c, pool)
	})

}
