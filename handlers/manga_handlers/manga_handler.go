package manga_handlers

import (
	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/JusAeng/manga-tracker-api-go/repo"
	"github.com/JusAeng/manga-tracker-api-go/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// MangaDetail bundles a manga with its credited authors and genres — the
// three things "view manga information" means in practice — so the
// frontend doesn't need three round-trips for one screen.
type MangaDetail struct {
	*models.Manga
	Authors []*models.MangaAuthor `json:"authors"`
	Genres  []*models.Genre       `json:"genres"`
}

// GetMangaHandler lists manga, optionally filtered by ?q= against either
// title.
func GetMangaHandler(c *fiber.Ctx) error {
	search := c.Query("q")
	loadedMangas, err := repo.GetMangas(search)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(loadedMangas)
}

func GetMangaByIdHandler(c *fiber.Ctx) error {
	mangaId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	manga, err := repo.GetMangaById(mangaId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if manga == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	authors, err := repo.GetAuthorsByMangaId(mangaId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	genres, err := repo.GetGenresByMangaId(mangaId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(MangaDetail{Manga: manga, Authors: authors, Genres: genres})
}

// GetThaiEditionsHandler lists the official Thai edition(s) of a manga.
func GetThaiEditionsHandler(c *fiber.Ctx) error {
	mangaId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	editions, err := repo.GetThaiEditionsByMangaId(mangaId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(editions)
}

// GetVolumesHandler lists the Thai volumes (with release dates) of a
// specific Thai edition.
func GetVolumesHandler(c *fiber.Ctx) error {
	thaiEditionId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	volumes, err := repo.GetVolumesByThaiEditionId(thaiEditionId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(volumes)
}

func GetPublisherHandler(c *fiber.Ctx) error {
	publisherId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	publisher, err := repo.GetPublisherById(publisherId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if publisher == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.JSON(publisher)
}

func GetMangaTrending(c *fiber.Ctx) error {
	allManga, err := repo.GetMangas("")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	selected := service.RandomManga(allManga, 3)
	return c.JSON(selected)
}

func GetMangaRecommend(c *fiber.Ctx) error {
	allManga, err := repo.GetMangas("")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	selected := service.RandomManga(allManga, 8)
	return c.JSON(selected)
}

func GetMangaNew(c *fiber.Ctx) error {
	allManga, err := repo.GetMangas("")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	selected := service.RandomManga(allManga, 5)
	return c.JSON(selected)
}
