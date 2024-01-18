package auth_handlers

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/JusAeng/manga-tracker-api-go/repo"
	"github.com/gofiber/fiber/v2"

	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetLineProfileByTokenIdHandler(c *fiber.Ctx) error {
	tokenId := c.Params("id")

	// manga, err := repo.GetMangaById(mangaId)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	// }

	// URL of the API endpoint for the POST request
	profilePayload := url.Values{
		"id_token":  {tokenId},
		"client_id": {"2002829031"},
	}

	// Create a request with the payload
	req, err := http.NewRequest("POST", "https://api.line.me/oauth2/v2.1/verify", bytes.NewBufferString(profilePayload.Encode()))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parse the response body into a map[string]interface{}
	var responseBody map[string]interface{}
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return err
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error: %s", string(body))
	}

	// Return the parsed response body as JSON
	return c.JSON(responseBody)
}

func Login(c *fiber.Ctx) error {

	if false {
		Register(c)
	}

	return nil
}

func Register(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		log.Println("Error parsing request body:", err)
    	log.Println("Request Body:", c.Body()) // Print the request body for debugging
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	newUser, err := repo.CreateUser(user)
	if err != nil {
		return c.Status(fiber.StatusAccepted).SendString(err.Error())
	}

	return c.JSON(newUser)
}