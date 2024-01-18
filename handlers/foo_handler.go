package handlers

import (
	"fmt"
	"log"

	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/JusAeng/manga-tracker-api-go/repo"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FooLogin(c *fiber.Ctx) error {
	sub := "5f563a9da793b25a0952923"
	hashSub := sub+"4"
	
	fmt.Println(hashSub)
	objectID, err := primitive.ObjectIDFromHex(hashSub)
	if err != nil{
		return nil
	}
	userProfile := repo.GetUserProfileById(objectID)
	
	return c.JSON(userProfile)
}

func FooAddUser(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		log.Println("Error parsing request body:", err)
    	log.Println("Request Body:", c.Body()) // Print the request body for debugging
		return c.Status(fiber.StatusBadRequest).SendString("nani")
	}
	newUser, err := repo.CreateUser(user)
	if err != nil {
		return c.Status(fiber.StatusAccepted).SendString("nnai")
	}
	log.Print(newUser)
	return c.JSON(newUser)
}

func FooDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	err = repo.DeleteUserById(objectID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.Status(fiber.StatusAccepted).SendString(userId)
}