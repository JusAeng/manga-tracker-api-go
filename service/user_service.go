package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/JusAeng/manga-tracker-api-go/config"
)

func EncryptHexId(subIdLine string) (string,error){
	encryptionKey,err := config.GetEnv("EncryptionKey")
	if err != nil {
		return "",errors.New("EncryptionKey Failed")
	}

	insertedHexId := ""
	subIdLine = subIdLine[1:]

	// fill 40 char by added ? every 3 char
	ce := string(encryptionKey[0])
	added := 0
	if len(subIdLine) < 40{
		for i:=0; i<len(subIdLine); i++{
			if i%3 == 0{
				insertedHexId += ce
				added += 1
				ce = string(encryptionKey[added])
			}
			if added >= 40 - len(subIdLine){
				insertedHexId += subIdLine[i:]
				break
			}
			insertedHexId += string(subIdLine[i])
		}
	}

	// now we get hexId as 40 chars -> to fill 48 chars by rotated ? step forward
	fullHexId := insertedHexId
	if len(insertedHexId) < 48 {
		for i:=0; i< 48 - len(insertedHexId); i++ {
			s, err := shieftHex(string(insertedHexId[i]),string(encryptionKey[added]))
			if err != nil{
				return "",errors.New("Not Heximal form")
			}
			fullHexId += s
			added += 1
			added %= len(encryptionKey)
		}
	}

	if len(fullHexId) != 48 {
		return "",errors.New("Hex Id is invalid length")
	}

	// now we going to process by half it !
	halfHexId := ""
	for i:=0; i<24; i++{
		s,err := meanHex(string(fullHexId[i]),string(fullHexId[24+i]))
		if err != nil{
			return "",errors.New("Can't process this heximal")
		}
		halfHexId += s
	}

	return halfHexId,nil
}

func shieftHex(s string,step string) (string,error){
	s_decimal, err := strconv.ParseInt(s, 16, 64)
	if err != nil{
		return "",err
	}
	step_decimal, err := strconv.ParseInt(step, 16, 64)
	if err != nil{
		return "",err
	}
	result := fmt.Sprintf("%X",(s_decimal+step_decimal)%16)

	return result,nil
}

func meanHex(s1 string,s2 string) (string,error){
	s1_decimal, err := strconv.ParseInt(s1, 16, 64)
	if err != nil{
		return "",err
	}
	s2_decimal, err := strconv.ParseInt(s2, 16, 64)
	if err != nil{
		return "",err
	}

	result := fmt.Sprintf("%x",(s1_decimal+s2_decimal)/2)

	return result,nil
}