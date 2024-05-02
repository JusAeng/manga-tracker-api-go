package manga_handlers

import (
	"log"
	"math/rand"
	"net/url"

	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/JusAeng/manga-tracker-api-go/repo"
	"github.com/JusAeng/manga-tracker-api-go/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Read
func GetMangaHandler(c *fiber.Ctx) error {
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

func GetMangaHighlight(c *fiber.Ctx) error {
	mangaId, err := primitive.ObjectIDFromHex("662d5f00d657e10679478b83")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	manga, err := repo.GetMangaById(mangaId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(manga)
}

func GetMangaTrending(c *fiber.Ctx) error {
	allManga,err := repo.GetMangas()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	service.Shuffle(allManga)
	selected := make([]*models.Manga, 3)
	runner := 0
    for i := 0; i < len(allManga); i++ {
        pick := rand.Intn(2) == 0
		if (pick){
			selected[runner] = allManga[i]
			runner +=1
		}
		if (runner >= 3 || len(selected) >= 3) {
			break
		}
		if (i + 3 == len(allManga)+len(selected)){
			break
		}
    }
	return c.JSON(selected)
}

func GetMangaRecommend(c *fiber.Ctx) error {
	allManga,err := repo.GetMangas()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	service.Shuffle(allManga)
	selected := make([]*models.Manga, 8)
	runner := 0
    for i := 0; i < len(allManga); i++ {
        pick := rand.Intn(2) == 0
		if (pick){
			selected[runner] = allManga[i]
			runner +=1
		}
		if (runner >= 8 || len(selected) >= 8) {
			break
		}
		if (i + 8 == len(allManga)+len(selected)){
			break
		}
    }
	return c.JSON(selected)
}

func GetMangaNew(c *fiber.Ctx) error {
	allManga,err := repo.GetMangas()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	service.Shuffle(allManga)
	selected := make([]*models.Manga, 5)
	runner := 0
    for i := 0; i < len(allManga); i++ {
        pick := rand.Intn(2) == 0
		if (pick){
			selected[runner] = allManga[i]
			runner +=1
		}
		if (runner >= 5 || len(selected) >= 5) {
			break
		}
		if (i + 5 == len(allManga)+len(selected)){
			break
		}
    }
	return c.JSON(selected)
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
	if err != nil {
		return err
	}
	vol := new(models.Vol)
	if err := c.BodyParser(vol); err != nil {
		log.Println("Error parsing request body:", err)
    	log.Println("Request Body:", c.Body()) // Print the request body for debugging
		return c.Status(fiber.StatusBadRequest).SendString("requestBody")
	}
	vol.MangaID = mangaId
	newVol, err := repo.AddMangaVol(*vol)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(newVol)
}

func VolumeDeleteAll(c* fiber.Ctx) error {
	mangaId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil{
		return err
	}
	err = repo.DeleteAllVolsByMangaId(mangaId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return nil
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
