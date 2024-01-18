package router

import (
	"github.com/JusAeng/manga-tracker-api-go/handlers"
	"github.com/JusAeng/manga-tracker-api-go/handlers/admin_handlers"
	"github.com/JusAeng/manga-tracker-api-go/handlers/auth_handlers"
	"github.com/JusAeng/manga-tracker-api-go/handlers/manga_handlers"
	"github.com/JusAeng/manga-tracker-api-go/handlers/user_handlers"
	"github.com/gofiber/fiber/v2"
)

func Run(app *fiber.App){
	//foo
	app.Post("/foo/user",handlers.FooAddUser)
	app.Get("/foo/user/:id",handlers.FooLogin)
	app.Delete("/foo/user/:id",handlers.FooDeleteUser)

	// manga
	app.Get("/mangas", manga_handlers.GetMangasHandler)
	app.Get("/mangas/:id", manga_handlers.GetMangaByIdHandler)
	app.Get("/mangas/title/:title", manga_handlers.GetMangaByTitleHandler)
	app.Put("/manga/vol",manga_handlers.VolumeAdding)
	app.Put("/manga/voting",manga_handlers.VotingAdding)
	app.Put("/manga/subscribe",manga_handlers.SubscribeAdding)

	// user
	app.Patch("/user/:id",user_handlers.UpdateUserProfile)
	app.Put("/user/subscribe/:id",user_handlers.UpdateSubscribe)
	app.Put("/user/ownerlist",user_handlers.UpdateOwnerList)
	app.Put("/user/rating",user_handlers.UpdateRating)

	// Auth
	app.Get("/auth/:id",auth_handlers.Login)

	// Admin
	app.Delete("/admin/user/:id",admin_handlers.DeleteUserById)

	app.Post("/admin/manga",admin_handlers.CreateManga)
	app.Patch("/admin/manga/:id",admin_handlers.UpdateMangaByID)
	app.Delete("/admin/manga/id/:id",admin_handlers.DeleteMangaByID)
	app.Delete("/admin/manga/title/:title",admin_handlers.DeleteMangaByTitle)

	app.Post("/admin/vol",admin_handlers.CreateVol)
	app.Patch("/admin/vol",admin_handlers.UpdateVol)
	app.Delete("/admin/vol",admin_handlers.DeleteVol)
}