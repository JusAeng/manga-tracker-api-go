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

func RandomManga(allManga []*models.Manga,numb int) ([]*models.Manga){
	var selected []*models.Manga
    Shuffle(allManga)

    for i := 0; i < len(allManga); i++ {
        pick := rand.Intn(2) == 0
		if (i + numb >= len(allManga)+len(selected)){
			pick = true
		}
		if (pick){
			selected = append(selected, allManga[i])
		}
		if (len(selected) >= numb) {
			break
		}
    }
    return selected
}