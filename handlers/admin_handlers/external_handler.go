package admin_handlers

import (
	"strconv"

	"github.com/JusAeng/manga-tracker-api-go/service"
	"github.com/gofiber/fiber/v2"
)

func SearchExternalManga(c *fiber.Ctx) error {
	results, err := service.SearchExternalManga(c.Query("q"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(results)
}

func GetExternalMangaDraft(c *fiber.Ctx) error {
	anilistID, err := strconv.Atoi(c.Params("anilistId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid anilistId")
	}
	draft, err := service.GetExternalMangaDraft(anilistID)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).SendString(err.Error())
	}
	return c.JSON(draft)
}
