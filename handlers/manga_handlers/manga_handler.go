package manga_handlers

import (
	"log"
	"net/url"

	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/JusAeng/manga-tracker-api-go/repo"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Read
func GetMangasHandler(c *fiber.Ctx) error {
	loadedMangas, err := repo.GetMangas()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(loadedMangas)
}

func GetMangaByIdHandler(c *fiber.Ctx) error {
	mangaId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	manga, err := repo.GetMangaById(mangaId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(manga)
}

func GetMangaByTitleHandler(c *fiber.Ctx) error {
	manga, err := repo.GetMangaByTitle(c.Params("title"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(manga)
}

// Create
func AddMangaHandler(c *fiber.Ctx) error {
	manga := new(models.Manga)

	if err := c.BodyParser(manga); err != nil {
		log.Println("Error parsing request body:", err)
    	log.Println("Request Body:", c.Body()) // Print the request body for debugging
		return c.Status(fiber.StatusBadRequest).SendString("nani")
	}
	newManga, err := repo.AddManga(manga)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("nnai")
	}
	log.Print(newManga)
	return c.JSON(newManga)
}

func VolumeAdding(c* fiber.Ctx) error {
	mangaId, err := primitive.ObjectIDFromHex(c.Params("id"))
	vol := new(models.Vol)
	if err := c.BodyParser(vol); err != nil {
		log.Println("Error parsing request body:", err)
    	log.Println("Request Body:", c.Body()) // Print the request body for debugging
		return c.Status(fiber.StatusBadRequest).SendString("requestBody")
	}
	newVol, err := repo.AddMangaVol(mangaId,*vol)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(newVol)
}

// Delete
func DeleteMangaByTitleHandler(c *fiber.Ctx) error{
	encodedTitle := c.Params("title")

	mangaTitle, err := url.PathUnescape(encodedTitle)
	if err != nil {
		log.Printf("Error decoding title: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Error decoding title")
	}
	
	err = repo.DeleteMangaByTitle(mangaTitle)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.Status(fiber.StatusAccepted).SendString(mangaTitle)
}

// Put
