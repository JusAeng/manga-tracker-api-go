package admin_handlers

import (
	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/JusAeng/manga-tracker-api-go/repo"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Users

func GetAllUsers(c *fiber.Ctx) error {
	allUsers, err := repo.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(allUsers)
}

func DeleteUserById(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := repo.DeleteUserById(id); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusAccepted)
}

// Manga

func CreateManga(c *fiber.Ctx) error {
	manga := new(models.Manga)
	if err := c.BodyParser(manga); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	newManga, err := repo.AddManga(manga)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(newManga)
}

func UpdateManga(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	manga := new(models.Manga)
	if err := c.BodyParser(manga); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	manga.ID = id
	updated, err := repo.UpdateManga(manga)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(updated)
}

func DeleteMangaByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := repo.DeleteMangaById(id); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.Status(fiber.StatusAccepted).SendString(id.String())
}

// Manga <-> Author / Genre attachments

type AttachAuthorRequest struct {
	AuthorID uuid.UUID `json:"authorId"`
	Role     string    `json:"role"`
}

func AttachMangaAuthor(c *fiber.Ctx) error {
	mangaId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	req := new(AttachAuthorRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := repo.AddMangaAuthor(mangaId, req.AuthorID, req.Role); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusAccepted)
}

func DetachMangaAuthor(c *fiber.Ctx) error {
	mangaId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	req := new(AttachAuthorRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := repo.RemoveMangaAuthor(mangaId, req.AuthorID, req.Role); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusAccepted)
}

type AttachGenreRequest struct {
	GenreID uuid.UUID `json:"genreId"`
}

func AttachMangaGenre(c *fiber.Ctx) error {
	mangaId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	req := new(AttachGenreRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := repo.AddMangaGenre(mangaId, req.GenreID); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusAccepted)
}

func DetachMangaGenre(c *fiber.Ctx) error {
	mangaId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	req := new(AttachGenreRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := repo.RemoveMangaGenre(mangaId, req.GenreID); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusAccepted)
}

// Publishers

func CreatePublisher(c *fiber.Ctx) error {
	publisher := new(models.Publisher)
	if err := c.BodyParser(publisher); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	newPublisher, err := repo.AddPublisher(publisher)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(newPublisher)
}

func UpdatePublisher(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	publisher := new(models.Publisher)
	if err := c.BodyParser(publisher); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	publisher.ID = id
	updated, err := repo.UpdatePublisher(publisher)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(updated)
}

func DeletePublisher(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := repo.DeletePublisherById(id); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusAccepted)
}

// Thai editions — :id is the manga's id on create (an edition is always
// created under a manga), and the edition's own id on update/delete.

func CreateThaiEdition(c *fiber.Ctx) error {
	mangaId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	edition := new(models.ThaiEdition)
	if err := c.BodyParser(edition); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	edition.MangaID = mangaId
	newEdition, err := repo.AddThaiEdition(edition)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(newEdition)
}

func UpdateThaiEdition(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	edition := new(models.ThaiEdition)
	if err := c.BodyParser(edition); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	edition.ID = id
	updated, err := repo.UpdateThaiEdition(edition)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(updated)
}

func DeleteThaiEdition(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := repo.DeleteThaiEditionById(id); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusAccepted)
}

// Volumes — :id is the thai edition's id on create, and the volume's own
// id on update/delete (volumes have a real primary key, so there's no
// need to route through the parent for a single-volume operation).

func CreateVolume(c *fiber.Ctx) error {
	thaiEditionId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	vol := new(models.Volume)
	if err := c.BodyParser(vol); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	vol.ThaiEditionID = thaiEditionId
	newVol, err := repo.AddVolume(*vol)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(newVol)
}

func UpdateVolume(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	vol := new(models.Volume)
	if err := c.BodyParser(vol); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	updated, err := repo.UpdateVolume(id, *vol)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(updated)
}

func DeleteVolume(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := repo.DeleteVolumeById(id); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusAccepted)
}

// Authors

func CreateAuthor(c *fiber.Ctx) error {
	author := new(models.Author)
	if err := c.BodyParser(author); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	newAuthor, err := repo.AddAuthor(author)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(newAuthor)
}

func UpdateAuthor(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	author := new(models.Author)
	if err := c.BodyParser(author); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	author.ID = id
	updated, err := repo.UpdateAuthor(author)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(updated)
}

func DeleteAuthor(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := repo.DeleteAuthorById(id); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusAccepted)
}

// Genres

func CreateGenre(c *fiber.Ctx) error {
	genre := new(models.Genre)
	if err := c.BodyParser(genre); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	newGenre, err := repo.AddGenre(genre)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(newGenre)
}

func UpdateGenre(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	genre := new(models.Genre)
	if err := c.BodyParser(genre); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	genre.ID = id
	updated, err := repo.UpdateGenre(genre)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(updated)
}

func DeleteGenre(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := repo.DeleteGenreById(id); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusAccepted)
}
