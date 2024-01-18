package handlers

import (
	"fmt"
	"log"

	"github.com/JusAeng/manga-tracker-api-go/repo"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SubscribeByIdHandler(c *fiber.Ctx) error {
	mangaId := c.Params("id");
	err := repo.SubscribeMangaById(mangaId)

	return err
}

type UpdateUserRequest struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

func UpdateUserProfile(c *fiber.Ctx) error {
	req := new(UpdateUserRequest)
	userId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil{
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err = c.BodyParser(req); err != nil {
		log.Println("Error parsing request body:", err)
    	log.Println("Request Body:", c.Body())
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	fmt.Println("handler",req.Key,req.Value)
	err = repo.UpdateUserProfile(userId,req.Key,req.Value)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return  c.SendStatus(fiber.StatusAccepted)
}