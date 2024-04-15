package router

import (
	"github.com/JusAeng/manga-tracker-api-go/handlers/user_handlers"
	"github.com/gofiber/fiber/v2"
)

func routingUserPath(app *fiber.App) {
	userRoute := app.Group("/user")
	
	userRoute.Get("/profile",user_handlers.GetUserProfile)
	userRoute.Patch("/profile",user_handlers.UpdateUserProfile)

	userRoute.Get("/subscribelist",user_handlers.GetSubscribeList)

	userRoute.Put("/subscribe/:id",user_handlers.UpdateSubscribe)
	userRoute.Put("/ownerlist/:id/:vol",user_handlers.UpdateOwnerList)
	userRoute.Put("/rating/:id/:score",user_handlers.UpdateRating)
}