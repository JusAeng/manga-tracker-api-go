package handlers

import (
	"fmt"
	"log"

	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/JusAeng/manga-tracker-api-go/repo"
	"github.com/gofiber/fiber/v2"
)

func FooLogin(c *fiber.Ctx) error {
	sub := "5f563a9da793b25a0952923"
	hashSub := sub+"4"
	
	fmt.Println(hashSub)
	userProfile := repo.GetUserProfile(hashSub)
	return c.JSON(userProfile)
}

func FooAddUser(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		log.Println("Error parsing request body:", err)
    	log.Println("Request Body:", c.Body()) // Print the request body for debugging
		return c.Status(fiber.StatusBadRequest).SendString("nani")
	}
	newUser, err := repo.AddUser(user)
	if err != nil {
		return c.Status(fiber.StatusAccepted).SendString("nnai")
	}
	log.Print(newUser)
	return c.JSON(newUser)
}

func FooDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("title")
	
	err := repo.DeleteUserById(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.Status(fiber.StatusAccepted).SendString(userId)
}