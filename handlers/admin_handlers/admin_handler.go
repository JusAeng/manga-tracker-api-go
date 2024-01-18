package admin_handlers

import (
	"github.com/JusAeng/manga-tracker-api-go/repo"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteUserById(c *fiber.Ctx) error {
	objectID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil{
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	err = repo.DeleteUserById(objectID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusAccepted)
}

func CreateManga(c *fiber.Ctx) error { return nil }
func DeleteMangaByID(c *fiber.Ctx) error {return nil}
func DeleteMangaByTitle(c *fiber.Ctx) error {return nil}
func UpdateMangaByID(c *fiber.Ctx) error {return nil}

func CreateVol(c *fiber.Ctx) error {return nil}
func DeleteVol(c *fiber.Ctx) error {return nil}
func UpdateVol(c *fiber.Ctx) error {return nil}