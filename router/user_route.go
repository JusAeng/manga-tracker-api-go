package router

import (
	"github.com/JusAeng/manga-tracker-api-go/handlers/user_handlers"
	"github.com/gofiber/fiber/v2"
)

func routingUserPath(app *fiber.App) {
	userRoute := app.Group("/user")

	userRoute.Get("/profile", user_handlers.GetUserProfile)
	userRoute.Patch("/profile", user_handlers.UpdateUserProfile)

	userRoute.Get("/following", user_handlers.GetFollowedManga)
	userRoute.Put("/follow/:id", user_handlers.UpdateFollow)

	userRoute.Put("/rating/:id", user_handlers.RateManga)
	userRoute.Delete("/rating/:id", user_handlers.DeleteRating)
}
