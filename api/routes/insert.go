package routes

import (
	"archive-api/utils"
	"archive-api/utils/services"
	"crypto/sha256"
	"crypto/subtle"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/jackc/pgx/v5/pgxpool"
)

func validateAPIKEY_for(valid_key string) func(c *fiber.Ctx, key string) (bool, error) {
	return func(c *fiber.Ctx, key string) (bool, error) {
		hashedAPIKey := sha256.Sum256([]byte(valid_key))
		hashedKey := sha256.Sum256([]byte(key))

		if subtle.ConstantTimeCompare(hashedAPIKey[:], hashedKey[:]) == 1 {
			return true, nil
		}
		return false, keyauth.ErrMissingOrMalformedAPIKey
	}
}

func BuildInsertRoutes(app *fiber.App, pool *pgxpool.Pool) {

	insert_routes := app.Group("/insert")

	insert_routes.Use(keyauth.New(keyauth.Config{
		KeyLookup: "cookie:access_token",
		Validator: validateAPIKEY_for(os.Getenv("API_KEY")),
	}))

	insert_routes.Get("/clean", func(c *fiber.Ctx) error {
		return services.Clean(pool)
	})
	insert_routes.Post("/labels/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		type RequestLabels struct {
			Labels []string `json:"labels"`
		}
		labels := new(RequestLabels)
		if err := c.BodyParser(labels); err != nil {
			log.Default().Println(err)
			return err
		}
		return services.AddLabelsForId(id, labels.Labels, pool)
	})
	insert_routes.Post("/:id", func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		request := new(utils.Request)
		id := c.Params("id")
		if err := c.BodyParser(request); err != nil {
			log.Default().Println("error : ", err)
			return err
		}
		table_experiment := request.Request.ExperimentJSON.ToTable()
		request.Request.Table_experiment = table_experiment
		return services.InsertAll(id, request, pool)
	})
}
