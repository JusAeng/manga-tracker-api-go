package router

import (
	"github.com/JusAeng/manga-tracker-api-go/handlers/manga_handlers"
	"github.com/gofiber/fiber/v2"
)
func routingMangaPath(app *fiber.App) {
	mangaRoute := app.Group("/manga")
	mangaRoute.Get("/",manga_handlers.GetMangaHandler)

	mangaRoute.Get("/highlight", manga_handlers.GetMangaHighlight)
	mangaRoute.Get("/trending",manga_handlers.GetMangaTrending)
	mangaRoute.Get("/new",manga_handlers.GetMangaNew)
	mangaRoute.Get("/recommend",manga_handlers.GetMangaRecommend)
	
	mangaRoute.Get("/:id", manga_handlers.GetMangaByIdHandler)
	// app.Get("/mangas/title/:title", manga_handlers.GetMangaByTitleHandler)
}