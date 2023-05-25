package main

import (
	"archive-api/routes"
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	app := fiber.New(fiber.Config{AppName: "Archive API"})
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer pool.Close()

	routes.BuildInsertRoutes(app, pool)
	routes.BuildSelectRoutes(app, pool)
	routes.BuildSearchRoutes(app, pool)

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("Unable to ping the database:", err)
	}

	if err := app.Listen(":3000"); err != nil {
		log.Fatal("Unable to listen :", err)
	}
}
