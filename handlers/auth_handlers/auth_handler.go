package auth_handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/JusAeng/manga-tracker-api-go/config"
	"github.com/JusAeng/manga-tracker-api-go/repo"
	"github.com/JusAeng/manga-tracker-api-go/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

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
		return fmt.Errorf("error: %s", string(body))
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

func GetUserFromLineToken(tokenId string) (*LineProfile,error) {
	// URL of the API endpoint for the POST request
	client_id,err := config.GetEnv("LineClientId")
	if err != nil{
		return nil,errors.New("LineClientId Fail Load")
	}
	profilePayload := url.Values{
		"id_token":  {tokenId},
		"client_id": {client_id},
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

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("token Id invalid")
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}

	var responseBody LineProfile
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return nil,err
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

	lineProfile,err := GetUserFromLineToken(req.Token)
	if err != nil{
		return c.SendStatus(fiber.StatusBadRequest)
	}
	hexId := lineProfile.Sub
	encryptHexId,err := service.EncryptHexId(hexId)
	if err != nil{
		fmt.Println(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	userId, err := primitive.ObjectIDFromHex(encryptHexId)
	if err != nil{
		return nil
	}
	user := repo.GetUserProfileById(userId)
	if user == nil{
		user,err = repo.RegisterUser(userId,lineProfile.Name,lineProfile.Picture)
		if err != nil{
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
	}

	jwttoken := jwt.New(jwt.SigningMethodHS256)
	claim := jwttoken.Claims.(jwt.MapClaims)
	claim["userId"] = user.ID.Hex()
	claim["role"] = "user"
	claim["exp"] = time.Now().Add(time.Hour * 72).Unix()

	JWTSignedString,err := config.GetEnv("JWT_SIGNED_STRING")
	if err != nil {
		return errors.New("no env for JWT_SIGNED_STRING")
	}
	token,err := jwttoken.SignedString([]byte(JWTSignedString))
	if err != nil {
		return errors.New("token Id invalid")
	}

	return c.JSON(fiber.Map{
		"token":token,
	})
}

func AuthMiddleware(c *fiber.Ctx) error {
	JWTSignedString,err := config.GetEnv("JWT_SIGNED_STRING")
	if err != nil {
		return errors.New("no env for JWT_SIGNED_STRING")
	}
	var jwtKey = []byte(JWTSignedString)

	// Extract the JWT token from the request header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	// Check if the Authorization header is formatted correctly
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	tokenString := parts[1]
	// Parse and validate the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	// Set the userID in the context
	c.Locals("userId", token.Claims.(jwt.MapClaims)["userId"])

	// Call the next handler
	return c.Next()
}