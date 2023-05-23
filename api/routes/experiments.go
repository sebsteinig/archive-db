package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func getExperimentByID(id string) error {
	log.Default().Println("exp id :", id)
	return nil
}

func BuildExperimentRoutes(app *fiber.App) {
	app.Get("/exp_id/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return getExperimentByID(id)
	})
}
