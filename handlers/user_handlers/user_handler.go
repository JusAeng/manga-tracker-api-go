package user_handlers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/JusAeng/manga-tracker-api-go/repo"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
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

func UpdateSubscribe(c *fiber.Ctx) error {
	mangaId, err := primitive.ObjectIDFromHex(c.Params("id"))
	tempId := "5f563a9da793b25a09529123"
	userId,err := primitive.ObjectIDFromHex(tempId)
	if err != nil {
		fmt.Println("Convert primitiveID from hex error")
		return err
	}
	err = repo.SubscribeMangaById(userId,mangaId)

	return err
}

func UpdateOwnerList(c *fiber.Ctx) error {
	mangaId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("ReqError")
	}
	vol,err := strconv.Atoi(c.Params("vol"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("ReqError")
	}
	tempId := "5f563a9da793b25a09529123"
	userId,err := primitive.ObjectIDFromHex(tempId)
	if err != nil {
		fmt.Println("Convert primitiveID from hex error")
		return err
	}
	err = repo.UpdateOwnerList(userId,mangaId,vol)

	return err
}
func UpdateRating(c *fiber.Ctx) error {return nil}