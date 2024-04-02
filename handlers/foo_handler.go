package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/JusAeng/manga-tracker-api-go/config"
	"github.com/gofiber/fiber/v2"
)

func FooHello(c *fiber.Ctx) error {
	return c.SendString("Hello")
}

func FooCheckLineProfileWithLineToken(c *fiber.Ctx) error {
	tokenId := c.Params("id")

	// URL of the API endpoint for the POST request
	client_id,err := config.GetEnv("LineClientId")
	if err != nil{
		return errors.New("LineClientId Fail Load")
	}
	profilePayload := url.Values{
		"id_token":  {tokenId},
		"client_id": {client_id},
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
	body, err := io.ReadAll(resp.Body)
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
		return fmt.Errorf("error: %s", string(body))
	}

	// Return the parsed response body as JSON
	return c.JSON(responseBody)
}