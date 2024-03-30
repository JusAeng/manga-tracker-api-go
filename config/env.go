package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func GetEnv(key string) (string,error) {
	err := godotenv.Load()
	if err != nil{
		fmt.Println("Can't Load Env")
		return "",err
	}
	

	return os.Getenv(strings.ToUpper(key)),nil
}