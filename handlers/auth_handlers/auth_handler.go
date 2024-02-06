package auth_handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"

	"github.com/JusAeng/manga-tracker-api-go/repo"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

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

type LineProfile struct {
	Name 		string 	`json:"name"`
	Picture 	string 	`json:"picture"`
	Sub 		string 	`json:"sub"`
}

type ILineToken struct {
	Token string `json:"token"`
}

func GetUserIdFromLineToken(tokenId string) (*LineProfile,error) {
	// URL of the API endpoint for the POST request
	profilePayload := url.Values{
		"id_token":  {tokenId},
		"client_id": {"2002829031"},
	}

	// Create a request with the payload
	req, err := http.NewRequest("POST", "https://api.line.me/oauth2/v2.1/verify", bytes.NewBufferString(profilePayload.Encode()))
	if err != nil {
		return nil,err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil,err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}

	var responseBody LineProfile
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return nil,err
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Token Id invalid")
	}

	return &responseBody,nil
}

func Login(c *fiber.Ctx) error {
	req := new(ILineToken)
	if err := c.BodyParser(req); err != nil {
		log.Println("Error parsing request body:", err)
    	log.Println("Request Body:", c.Body()) // Print the request body for debugging
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}


	lineProfile,err := GetUserIdFromLineToken(req.Token)
	hexId := lineProfile.Sub
	if err != nil{
		return c.SendStatus(fiber.StatusBadRequest)
	}
	hashHexId,err := bcrypt.GenerateFromPassword([]byte(hexId), bcrypt.DefaultCost)
	if err != nil{
		fmt.Println(err)
	}
	hashHexIdString := string(hashHexId)

	// userId, err := primitive.ObjectIDFromHex("5f563a9da793b25a09529123")
	userId, err := primitive.ObjectIDFromHex(hashHexIdString[:24])
	if err != nil{
		return nil
	}
	user := repo.GetUserProfileById(userId)
	// if user == nil{
	// 	user,err = repo.RegisterUser(userId,lineProfile.Name,lineProfile.Picture)
	// 	if err != nil{
	// 		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	// 	}
	// }
	return c.JSON(user)
}