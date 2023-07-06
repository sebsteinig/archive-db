package main

import (
	"archive-api/routes"
	"context"
	"log"
	"os"

	"github.com/gofiber/swagger"

	_ "archive-api/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// @title Fiber Example API
// @version 1.0
// @description API for climatearchive
// @host localhost:3000
// @BasePath /
func main() {
	godotenv.Load(".env")
	app := fiber.New(fiber.Config{AppName: "Archive API"})

	app.Get("/doc/*", swagger.HandlerDefault) // default
	file, err := os.OpenFile("./archive_api.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	app.Use(logger.New(logger.Config{
		Output: file,
		Format: "${time} | [${ip}]:${port} ${status} - ${method} ${path} latency : ${latency}\n\tquery parameters : ${queryParams}\n",
	}))
	app.Use(helmet.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Static("/", "./static")
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer pool.Close()

	routes.BuildInsertRoutes(app, pool)
	routes.BuildSelectRoutes(app, pool)
	routes.BuildSearchRoutes(app, pool)

	/*if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("Unable to ping the database:", err)
	}*/

	if err := app.Listen(":3000"); err != nil {
		log.Fatal("Unable to listen :", err)
	}
}
