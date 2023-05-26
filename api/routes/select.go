package routes

import (
	"archive-api/utils/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildSelectRoutes(app *fiber.App, pool *pgxpool.Pool) {

	select_routes := app.Group("/select")
	//we add a limiter
	select_routes.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 30 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

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
