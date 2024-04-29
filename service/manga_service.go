package service

import (
	"math/rand"

	"github.com/JusAeng/manga-tracker-api-go/models"
)

func Shuffle(arr []*models.Manga) {
    for i := range arr {
        j := rand.Intn(i + 1)
        arr[i], arr[j] = arr[j], arr[i]
    }
}