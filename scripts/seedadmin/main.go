// Command seedadmin creates a row in the admins table. There's no HTTP
// route for this on purpose — creating the first admin has no
// authenticated caller yet, so it has to happen out-of-band against the
// real database (DATABASE_URL from your .env).
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/repo"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	username := flag.String("username", "", "admin username (required)")
	password := flag.String("password", "", "admin password (required) — only the bcrypt hash is stored")
	flag.Parse()

	if *username == "" || *password == "" {
		log.Fatal("-username and -password are both required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	db.Connect()
	admin, err := repo.CreateAdmin(*username, string(hash))
	if err != nil {
		log.Fatalf("failed to create admin (already exists?): %v", err)
	}

	fmt.Printf("created admin %q (id: %s)\n", admin.Username, admin.ID)
}
