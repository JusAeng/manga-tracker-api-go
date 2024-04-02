package router

import (
	"github.com/JusAeng/manga-tracker-api-go/handlers/admin_handlers"
	"github.com/gofiber/fiber/v2"
)
func routingAdminPath(app *fiber.App) {
	adminRoute := app.Group("/admin")
	adminRoute.Delete("/user/:id",admin_handlers.DeleteUserById)

	adminRoute.Post("/manga",admin_handlers.CreateManga)
	// adminRoute.Patch("/manga/:id",admin_handlers.UpdateMangaByID)
	adminRoute.Delete("/manga/:id",admin_handlers.DeleteMangaByID)

	// adminRoute.Post("/vol",admin_handlers.CreateVol)
	// adminRoute.Patch("/vol",admin_handlers.UpdateVol)
	// adminRoute.Delete("/vol",admin_handlers.DeleteVol)
}