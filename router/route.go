package router

import (
	"github.com/JusAeng/manga-tracker-api-go/handlers"
	// "github.com/JusAeng/manga-tracker-api-go/handlers/admin_handlers"
	"github.com/JusAeng/manga-tracker-api-go/handlers/auth_handlers"
	"github.com/JusAeng/manga-tracker-api-go/handlers/line_handlers"

	// "github.com/JusAeng/manga-tracker-api-go/handlers/user_handlers"
	"github.com/gofiber/fiber/v2"
)

func Run(app *fiber.App){
	// Mock
    app.Get("/hello", handlers.FooHello)
	app.Get("/foo/auth/:id",handlers.FooCheckLineProfileWithLineToken)

	// Authentication
	app.Post("/auth",auth_handlers.Login)
	app.Post("/auth/admin",auth_handlers.AdminLogin)

	// LINE Official Account webhook — public, authenticates itself via
	// X-Line-Signature rather than a JWT, so it must stay above the
	// JWTMiddleware app.Use below.
	app.Post("/line/webhook", line_handlers.Webhook)

	// verify JWT token
	app.Use(auth_handlers.JWTMiddleware)

	routingMangaPath(app)
	routingUserPath(app)
	routingAdminPath(app)
}
