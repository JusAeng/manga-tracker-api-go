package router

import (
	"github.com/JusAeng/manga-tracker-api-go/handlers"
	// "github.com/JusAeng/manga-tracker-api-go/handlers/admin_handlers"
	"github.com/JusAeng/manga-tracker-api-go/handlers/auth_handlers"
	"github.com/JusAeng/manga-tracker-api-go/handlers/manga_handlers"

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

	// verify JWT token
	app.Use(auth_handlers.JWTMiddleware)

	// General
	app.Get("/manga", manga_handlers.GetMangaHandler)
	app.Get("/manga/:id", manga_handlers.GetMangaByIdHandler)
	// app.Get("/mangas/title/:title", manga_handlers.GetMangaByTitleHandler)

	routingUserPath(app)
	routingAdminPath(app)
}
