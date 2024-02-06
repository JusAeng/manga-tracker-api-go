package service

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func EncryptHexId(hexId string) (string,error){
	err := godotenv.Load()
	if err != nil{
		fmt.Println("Can't Load Env")
		return "",err
	}
	encryptHexId := os.Getenv("EncryptionKey")
	fmt.Println(encryptHexId)

	return encryptHexId,nil
}