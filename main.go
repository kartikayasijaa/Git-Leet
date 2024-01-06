package main

import (
	"fmt"
	"gitleet/config"
	"gitleet/routes"
	"gitleet/services"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	logger := config.InitLogger()
	db, ctx := config.InitDB()

	dbServices := services.DBServicesHandler(db, ctx)
	h := routes.NewHandlers(logger, db, ctx, dbServices)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "Running",
		})
	})

	api := app.Group("/api")
	routes.SubmissionRoute(api, h)
	routes.AuthRoutes(api, h)

	PORT := os.Getenv("PORT")
	err := app.Listen(fmt.Sprintf(":%s", PORT))
	if err != nil {
		logger.Println(err)
	}
}
