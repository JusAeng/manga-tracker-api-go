package user_handlers

import (
	"errors"

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
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if req.Key != "displayName" && req.Key != "pictureUrl" {
		return c.Status(fiber.StatusBadRequest).SendString("Not Allow")
	}
	err = repo.UpdateUserProfile(userId, req.Key, req.Value)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusAccepted)
}

func GetFollowedManga(c *fiber.Ctx) error {
	userId, err := userIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid userId")
	}
	result, err := repo.GetUserFollowedManga(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(result)
}

func UpdateFollow(c *fiber.Ctx) error {
	mangaId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid mangaId")
	}
	userId, err := userIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid userId")
	}
	followedIds, err := repo.ToggleFollow(userId, mangaId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(followedIds)
}
