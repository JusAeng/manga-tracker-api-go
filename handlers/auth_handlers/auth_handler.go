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
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"fmt"
	"io"
	"net/http"
	"net/url"
)

type LineProfile struct {
	Name 		string 	`json:"name"`
	Picture 	string 	`json:"picture"`
	Sub 		string 	`json:"sub"`
}

type ILineToken struct {
	Token string `json:"token"`
}

func GetProfileFromLineAPI(tokenId string) (*LineProfile,error) {
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
	body, err := io.ReadAll(resp.Body)
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

	lineProfile,err := GetProfileFromLineAPI(req.Token)
	if err != nil{
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user, err := repo.GetOrCreateUserByLineID(lineProfile.Sub, lineProfile.Name, lineProfile.Picture)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	jwttoken := jwt.New(jwt.SigningMethodHS256)
	claim := jwttoken.Claims.(jwt.MapClaims)
	claim["userId"] = user.ID.String()
	claim["role"] = "user"
	claim["exp"] = time.Now().Add(time.Hour * 6).Unix()

	JWTSignedString,err := config.GetEnv("JWT_SIGNED_STRING")
	if err != nil {
		return errors.New("no env for JWT_SIGNED_STRING")
	}
	token,err := jwttoken.SignedString([]byte(JWTSignedString))
	if err != nil {
		return errors.New("token Id invalid")
	}
	c.Cookie(&fiber.Cookie{
		Name: "token",
		Value: token,
		Expires: time.Now().Add(time.Hour * 6),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"token":token,
		"profile": user,
	})
}

type AdminLoginType struct {
	Username	string	`json:"username"`
	Password	string	`json:"password"`
}

func AdminLogin(c *fiber.Ctx) error {
	req := new(AdminLoginType)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Form invalid!")
	}

	admin, err := repo.GetAdminByUsername(req.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	// Same response whether the username doesn't exist or the password is
	// wrong — don't leak which one it was.
	if admin == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Not found this admin!")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Not found this admin!")
	}
	jwttoken := jwt.New(jwt.SigningMethodHS256)
	claim := jwttoken.Claims.(jwt.MapClaims)
	claim["adminId"] = admin.ID.String()
	claim["username"] = admin.Username
	claim["role"] = "admin"
	claim["exp"] = time.Now().Add(time.Hour * 6).Unix()

	JWTSignedString,err := config.GetEnv("JWT_SIGNED_STRING")
	if err != nil {
		return errors.New("no env for JWT_SIGNED_STRING")
	}
	token,err := jwttoken.SignedString([]byte(JWTSignedString))
	c.Cookie(&fiber.Cookie{
		Name: "token",
		Value: token,
		Expires: time.Now().Add(time.Hour * 6),
		HTTPOnly: true,
	})

	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"token":token,
	})
}

func CheckAdmin(c *fiber.Ctx) error {
	if c.Locals("role").(string) != "admin" {
		return c.Status(fiber.StatusUnauthorized).SendString("no permission")
	}
	return c.Next()
}

func JWTMiddleware(c *fiber.Ctx) error {
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
		// Reject anything that isn't HMAC-signed. Without this check a
		// forged token with alg=none, or an alg swapped to one this server
		// never signs with, could otherwise be accepted (classic JWT
		// "alg confusion" issue).
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	// Set the userID,role in the context
	c.Locals("userId", token.Claims.(jwt.MapClaims)["userId"])
	c.Locals("role", token.Claims.(jwt.MapClaims)["role"])

	// Call the next handler
	return c.Next()
}