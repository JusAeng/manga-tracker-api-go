package user_handlers

import (
	"errors"
	"log"
	"strconv"

	"github.com/JusAeng/manga-tracker-api-go/repo"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// userIDFromContext reads the userId set by JWTMiddleware. Admin tokens
// carry no userId claim, so c.Locals("userId") can be nil here — asserting
// straight to string would panic instead of failing the request cleanly.
func userIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	userIdClaim, ok := c.Locals("userId").(string)
	if !ok {
		return uuid.UUID{}, errors.New("missing userId claim")
	}
	return uuid.Parse(userIdClaim)
}

func GetUserProfile(c *fiber.Ctx) error {
	userId, err := userIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid userId")
	}
	userProfile, err := repo.GetUserProfileById(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(userProfile)
}

func GetSubscribeList(c *fiber.Ctx) error {
	userId, err := userIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid userId")
	}
	result, err := repo.GetMangaFromSubscribeList(userId)
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
	userId, err := userIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid userId")
	}
	if err = c.BodyParser(req); err != nil {
		log.Println("Error parsing request body:", err)
		log.Println("Request Body:", c.Body())
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if req.Key != "name" && req.Key != "image" {
		return c.Status(fiber.StatusBadRequest).SendString("Not Allow")
	}
	err = repo.UpdateUserProfile(userId, req.Key, req.Value)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusAccepted)
}

func UpdateSubscribe(c *fiber.Ctx) error {
	mangaId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid mangaId")
	}
	userId, err := userIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid userId")
	}
	subscribeList, err := repo.SubscribeMangaById(userId, mangaId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(subscribeList)
}

func UpdateOwnerList(c *fiber.Ctx) error {
	mangaId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("ReqError")
	}
	vol, err := strconv.Atoi(c.Params("vol"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("ReqError")
	}
	userId, err := userIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid userId")
	}
	ownerList, err := repo.UpdateOwnerList(userId, mangaId, vol)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(ownerList)
}

func UpdateRating(c *fiber.Ctx) error {
	mangaId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("ReqError")
	}
	score, err := strconv.Atoi(c.Params("score"))
	if err != nil {
		return errors.New("can't convert score")
	}
	if score < 0 || score > 5 {
		return c.Status(fiber.StatusBadRequest).SendString("Rating between 0 - 5")
	}
	userId, err := userIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Check userId")
	}
	err = repo.UpdateRateList(userId, mangaId, score)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(score)
}
