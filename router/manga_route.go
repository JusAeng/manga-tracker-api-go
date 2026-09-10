package router

import (
	"github.com/JusAeng/manga-tracker-api-go/handlers/manga_handlers"
	"github.com/gofiber/fiber/v2"
)

func routingMangaPath(app *fiber.App) {
	mangaRoute := app.Group("/manga")
	mangaRoute.Get("/", manga_handlers.GetMangaHandler)

	mangaRoute.Get("/trending", manga_handlers.GetMangaTrending)
	mangaRoute.Get("/new", manga_handlers.GetMangaNew)
	mangaRoute.Get("/recommend", manga_handlers.GetMangaRecommend)

	mangaRoute.Get("/:id", manga_handlers.GetMangaByIdHandler)
	mangaRoute.Get("/:id/thai-editions", manga_handlers.GetThaiEditionsHandler)
	mangaRoute.Get("/:id/rating", manga_handlers.GetMangaRating)

	app.Get("/thai-editions/:id/volumes", manga_handlers.GetVolumesHandler)
	app.Get("/publishers/:id", manga_handlers.GetPublisherHandler)
}
