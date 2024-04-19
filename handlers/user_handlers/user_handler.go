package user_handlers

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/JusAeng/manga-tracker-api-go/repo"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserProfile(c *fiber.Ctx) error {
	userId,err := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	if err != nil {
		fmt.Println("Convert primitiveID from hex error")
		return err
	}
	userProfile := repo.GetUserProfileById(userId)
	return c.JSON(userProfile)
}

func GetSubscribeList(c *fiber.Ctx) error {
	userId,err := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	if err != nil {
		fmt.Println("Convert primitiveID from hex error")
		return err
	}
	result,err := repo.GetMangaFromSubscribeList(userId)
	if err != nil {
		return errors.New("get manga from subscribe list error")
	}
	return c.JSON(result)
}

type UpdateUserProfileRequest struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

func UpdateUserProfile(c *fiber.Ctx) error {
	req := new(UpdateUserProfileRequest)
	userId,err := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	if err != nil{
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err = c.BodyParser(req); err != nil {
		log.Println("Error parsing request body:", err)
    	log.Println("Request Body:", c.Body())
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	fmt.Println("handler",req.Key,req.Value)
	if (req.Key != "name" && req.Key != "image") {
		return c.Status(fiber.StatusBadRequest).SendString("Not Allow")
	}
	err = repo.UpdateUserProfile(userId,req.Key,req.Value)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return  c.SendStatus(fiber.StatusAccepted)
}

func UpdateSubscribe(c *fiber.Ctx) error {
	mangaId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		fmt.Println("Convert primitiveID from hex error")
		return err
	}
	userId,err := primitive.ObjectIDFromHex(c.Locals("userId").(string))
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
	userId,err := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	if err != nil {
		fmt.Println("Convert primitiveID from hex error")
		return err
	}
	ownerList,err := repo.UpdateOwnerList(userId,mangaId,vol)
	if err != nil{
		return err
	}

	return c.JSON(ownerList)
}
func UpdateRating(c *fiber.Ctx) error {
	mangaId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("ReqError")
	}
	score,err := strconv.Atoi(c.Params("score"))
	if err != nil{
		return errors.New("can't convert score")
	}
	if score < 1 || score > 10 {
		return c.Status(fiber.StatusBadRequest).SendString("Rating between 1 - 10")
	}
	userId,err := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	if err != nil{
		return c.Status(fiber.StatusBadRequest).SendString("Check userId")
	}
	err = repo.UpdateRateList(userId,mangaId,score)

	return err
}