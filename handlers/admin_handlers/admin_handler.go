package admin_handlers

import (
	"fmt"

	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/JusAeng/manga-tracker-api-go/repo"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User Handlers
func GetAllUsers(c *fiber.Ctx) error {
	allUsers,err := repo.GetAllUsers()
	if err != nil{
		return nil
	}
	return c.JSON(allUsers)
}
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

// Manga Hanlers
func CreateManga(c *fiber.Ctx) error { 
	manga := new(models.Manga)

	if err := c.BodyParser(manga); err != nil {
		fmt.Print(err)
		return c.Status(fiber.StatusBadRequest).SendString("Error BodyParser")
	}
	newManga, err := repo.AddManga(manga)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusForbidden).SendString("")
	}
	return c.JSON(newManga) 
}

func UpdateMangaByID(c *fiber.Ctx) error {return nil}
func DeleteMangaByID(c *fiber.Ctx) error {
	objectID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil{
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	err = repo.DeleteMangaById(objectID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.Status(fiber.StatusAccepted).SendString(objectID.String())
}

// Vol Handlers
func CreateVol(c *fiber.Ctx) error {return nil}
func DeleteVol(c *fiber.Ctx) error {return nil}
func UpdateVol(c *fiber.Ctx) error {return nil}