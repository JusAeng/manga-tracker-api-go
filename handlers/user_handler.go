package handlers

import (
	"github.com/JusAeng/manga-tracker-api-go/repo"
	"github.com/gofiber/fiber/v2"
)

func SubscribeByIdHandler(c *fiber.Ctx) error {
	mangaId := c.Params("id");
	err := repo.SubscribeMangaById(mangaId)

	return err
}