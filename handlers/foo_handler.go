package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/JusAeng/manga-tracker-api-go/repo"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminLoginType struct {
	Username	string	`json:"username"`
	Password	string	`json:"password"`
}

func FooLogin(c *fiber.Ctx) error {
	admin := new(AdminLoginType)
	if err := c.BodyParser(admin); err != nil{
		return c.Status(fiber.StatusBadRequest).SendString("Form invalid!")
	}
	if admin.Username != "admin1" || admin.Password != "admin"{
		return c.Status(fiber.StatusUnauthorized).SendString("Not found this admin!")
	}
	jwttoken := jwt.New(jwt.SigningMethodHS256)
	claim := jwttoken.Claims.(jwt.MapClaims)
	claim["username"] = "admin"
	claim["role"] = "admin"
	claim["exp"] = time.Now().Add(time.Hour * 2).Unix()

	token,err := jwttoken.SignedString([]byte("secret"))
	c.Cookie(&fiber.Cookie{
		Name: "token",
		Value: token,
		Expires: time.Now().Add(time.Hour * 2),
		HTTPOnly: true,
	})

	if err != nil {
		fmt.Println("not send token")
		return err
	}
	fmt.Println("token",token)
	return c.JSON(fiber.Map{
		"token":token,
	})
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