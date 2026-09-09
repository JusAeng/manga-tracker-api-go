package main

import (
	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New()

	app.Use(recover.New())

	// Apply CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Adjust this to be more restrictive if needed
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	db.Connect();

	router.Run(app)
  
	app.Listen(":8080")
  }