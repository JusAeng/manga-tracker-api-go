// Command devtoken mints a locally-signed JWT so you can hit protected
// endpoints without going through the real LINE login flow. It only needs
// JWT_SIGNED_STRING from your .env — it never talks to LINE or the
// database.
//
// DEV ONLY. Never run this against production secrets: whoever holds
// JWT_SIGNED_STRING can mint a token that's indistinguishable from a real
// login, for any userId/role they like.
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/JusAeng/manga-tracker-api-go/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func main() {
	role := flag.String("role", "user", `"user" or "admin"`)
	userID := flag.String("user-id", "", "existing user's id (uuid). If empty, a fresh one is generated — useful for hitting routes structurally, but repo.GetUserProfileById will find no matching row unless a user with this id actually exists")
	hours := flag.Float64("hours", 6, "token lifetime in hours (login flow uses 6)")
	flag.Parse()

	if *role != "user" && *role != "admin" {
		log.Fatalf("-role must be \"user\" or \"admin\", got %q", *role)
	}

	jwtSignedString, err := config.GetEnv("JWT_SIGNED_STRING")
	if err != nil || jwtSignedString == "" {
		log.Fatal("JWT_SIGNED_STRING is not set — check your .env")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Duration(*hours * float64(time.Hour))).Unix()

	if *role == "admin" {
		claims["username"] = "admin"
		claims["role"] = "admin"
	} else {
		id := *userID
		if id == "" {
			id = uuid.NewString()
			fmt.Printf("# no -user-id given, generated: %s\n", id)
			fmt.Println("# GetUserProfileById will return null for this id unless a matching user row actually exists.")
		}
		claims["userId"] = id
		claims["role"] = "user"
	}

	signed, err := token.SignedString([]byte(jwtSignedString))
	if err != nil {
		log.Fatalf("failed to sign token: %v", err)
	}

	fmt.Println(signed)
	fmt.Printf("\n# try it:\ncurl -H \"Authorization: Bearer %s\" http://localhost:8080/user/profile\n", signed)
}
