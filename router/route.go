package router

import (
	"github.com/JusAeng/manga-tracker-api-go/handlers"
	"github.com/gofiber/fiber/v2"
)

func Run(app *fiber.App){
	app.Get("/mangas", handlers.GetMangasHandler)
	app.Get("/mangas/:id", handlers.GetMangaByIdHandler)
	app.Get("/mangas/title/:title", handlers.GetMangaByTitleHandler)
	app.Delete("/mangas/:title", handlers.DeleteMangaByTitleHandler)
	app.Post("/mangas", handlers.AddMangaHandler)
}