package router

import (
	"github.com/JusAeng/manga-tracker-api-go/handlers/admin_handlers"
	"github.com/JusAeng/manga-tracker-api-go/handlers/auth_handlers"
	"github.com/gofiber/fiber/v2"
)

func routingAdminPath(app *fiber.App) {
	adminRoute := app.Group("/admin", auth_handlers.CheckAdmin)

	adminRoute.Get("/users", admin_handlers.GetAllUsers)
	adminRoute.Delete("/users/:id", admin_handlers.DeleteUserById)

	adminRoute.Post("/manga", admin_handlers.CreateManga)
	adminRoute.Patch("/manga/:id", admin_handlers.UpdateManga)
	adminRoute.Delete("/manga/:id", admin_handlers.DeleteMangaByID)
	adminRoute.Post("/manga/:id/authors", admin_handlers.AttachMangaAuthor)
	adminRoute.Delete("/manga/:id/authors", admin_handlers.DetachMangaAuthor)
	adminRoute.Post("/manga/:id/genres", admin_handlers.AttachMangaGenre)
	adminRoute.Delete("/manga/:id/genres", admin_handlers.DetachMangaGenre)
	adminRoute.Post("/manga/:id/thai-editions", admin_handlers.CreateThaiEdition)

	adminRoute.Patch("/thai-editions/:id", admin_handlers.UpdateThaiEdition)
	adminRoute.Delete("/thai-editions/:id", admin_handlers.DeleteThaiEdition)
	adminRoute.Post("/thai-editions/:id/volumes", admin_handlers.CreateVolume)

	adminRoute.Patch("/volumes/:id", admin_handlers.UpdateVolume)
	adminRoute.Delete("/volumes/:id", admin_handlers.DeleteVolume)

	adminRoute.Post("/publishers", admin_handlers.CreatePublisher)
	adminRoute.Patch("/publishers/:id", admin_handlers.UpdatePublisher)
	adminRoute.Delete("/publishers/:id", admin_handlers.DeletePublisher)

	adminRoute.Post("/authors", admin_handlers.CreateAuthor)
	adminRoute.Patch("/authors/:id", admin_handlers.UpdateAuthor)
	adminRoute.Delete("/authors/:id", admin_handlers.DeleteAuthor)

	adminRoute.Post("/genres", admin_handlers.CreateGenre)
	adminRoute.Patch("/genres/:id", admin_handlers.UpdateGenre)
	adminRoute.Delete("/genres/:id", admin_handlers.DeleteGenre)
}
