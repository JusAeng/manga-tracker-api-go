package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const anilistEndpoint = "https://graphql.anilist.co"

type ExternalMangaSummary struct {
	AnilistID     int    `json:"anilistId"`
	TitleRomaji   string `json:"titleRomaji"`
	TitleEnglish  string `json:"titleEnglish"`
	TitleNative   string `json:"titleNative"`
	Year          *int   `json:"year"`
	CoverImageURL string `json:"coverImageUrl"`
	Format        string `json:"format"`
}

type ExternalAuthorDraft struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type ExternalMangaDraft struct {
	TitleOriginal string                `json:"titleOriginal"`
	TitleEn       string                `json:"titleEn"`
	Introduction  string                `json:"introduction"`
	ImageURL      string                `json:"imageUrl"`
	FirstDateJp   *string               `json:"firstDateJp"`
	Status        string                `json:"status"`
	Authors       []ExternalAuthorDraft `json:"authors"`
	Genres        []string              `json:"genres"`
}

type anilistDate struct {
	Year  *int `json:"year"`
	Month *int `json:"month"`
	Day   *int `json:"day"`
}

func (d anilistDate) toISODate() *string {
	if d.Year == nil {
		return nil
	}
	month, day := 1, 1
	if d.Month != nil {
		month = *d.Month
	}
	if d.Day != nil {
		day = *d.Day
	}
	s := fmt.Sprintf("%04d-%02d-%02d", *d.Year, month, day)
	return &s
}

func anilistPost(query string, variables map[string]any, out any) error {
	body, err := json.Marshal(map[string]any{"query": query, "variables": variables})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, anilistEndpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("anilist request failed: %s", res.Status)
	}
	return json.NewDecoder(res.Body).Decode(out)
}

const searchQuery = `
query ($search: String) {
  Page(page: 1, perPage: 8) {
    media(search: $search, type: MANGA, sort: SEARCH_MATCH) {
      id
      title { romaji english native }
      coverImage { large }
      startDate { year }
      format
    }
  }
}`

func SearchExternalManga(query string) ([]ExternalMangaSummary, error) {
	if strings.TrimSpace(query) == "" {
		return nil, errors.New("search query is required")
	}

	var resp struct {
		Data struct {
			Page struct {
				Media []struct {
					ID    int `json:"id"`
					Title struct {
						Romaji  string `json:"romaji"`
						English string `json:"english"`
						Native  string `json:"native"`
					} `json:"title"`
					CoverImage struct {
						Large string `json:"large"`
					} `json:"coverImage"`
					StartDate anilistDate `json:"startDate"`
					Format    string      `json:"format"`
				} `json:"media"`
			} `json:"Page"`
		} `json:"data"`
	}
	if err := anilistPost(searchQuery, map[string]any{"search": query}, &resp); err != nil {
		return nil, err
	}

	results := make([]ExternalMangaSummary, 0, len(resp.Data.Page.Media))
	for _, m := range resp.Data.Page.Media {
		results = append(results, ExternalMangaSummary{
			AnilistID:     m.ID,
			TitleRomaji:   m.Title.Romaji,
			TitleEnglish:  m.Title.English,
			TitleNative:   m.Title.Native,
			Year:          m.StartDate.Year,
			CoverImageURL: m.CoverImage.Large,
			Format:        m.Format,
		})
	}
	return results, nil
}

const detailQuery = `
query ($id: Int) {
  Media(id: $id, type: MANGA) {
    title { romaji english native }
    description(asHtml: false)
    coverImage { large }
    startDate { year month day }
    status
    genres
    staff {
      edges {
        role
        node { name { full } }
      }
    }
  }
}`

func anilistStatusToOurs(status string) string {
	switch status {
	case "FINISHED":
		return "completed"
	case "RELEASING":
		return "ongoing"
	case "NOT_YET_RELEASED":
		return "upcoming"
	case "CANCELLED":
		return "cancelled"
	case "HIATUS":
		return "hiatus"
	default:
		return ""
	}
}

// anilistRoleToOurs narrows AniList's free-text staff roles (e.g. "Story &
// Art", "Character Design", "Translator") down to the four credit types the
// schema tracks. Roles that aren't clearly a writing/art credit are dropped
// rather than guessed at.
func anilistRoleToOurs(role string) (string, bool) {
	r := strings.ToLower(role)
	hasStory := strings.Contains(r, "story")
	hasArt := strings.Contains(r, "art")
	switch {
	case hasStory && hasArt:
		return "author", true
	case hasStory:
		return "story", true
	case hasArt:
		return "artist", true
	default:
		return "", false
	}
}

func GetExternalMangaDraft(anilistID int) (*ExternalMangaDraft, error) {
	var resp struct {
		Data struct {
			Media struct {
				Title struct {
					Romaji  string `json:"romaji"`
					English string `json:"english"`
					Native  string `json:"native"`
				} `json:"title"`
				Description string `json:"description"`
				CoverImage  struct {
					Large string `json:"large"`
				} `json:"coverImage"`
				StartDate anilistDate `json:"startDate"`
				Status    string      `json:"status"`
				Genres    []string    `json:"genres"`
				Staff     struct {
					Edges []struct {
						Role string `json:"role"`
						Node struct {
							Name struct {
								Full string `json:"full"`
							} `json:"name"`
						} `json:"node"`
					} `json:"edges"`
				} `json:"staff"`
			} `json:"Media"`
		} `json:"data"`
	}
	if err := anilistPost(detailQuery, map[string]any{"id": anilistID}, &resp); err != nil {
		return nil, err
	}

	m := resp.Data.Media
	titleOriginal := m.Title.Native
	if titleOriginal == "" {
		titleOriginal = m.Title.Romaji
	}
	titleEn := m.Title.English
	if titleEn == "" {
		titleEn = m.Title.Romaji
	}

	seen := map[string]bool{}
	authors := make([]ExternalAuthorDraft, 0)
	for _, edge := range m.Staff.Edges {
		role, ok := anilistRoleToOurs(edge.Role)
		if !ok {
			continue
		}
		key := edge.Node.Name.Full + "|" + role
		if seen[key] {
			continue
		}
		seen[key] = true
		authors = append(authors, ExternalAuthorDraft{Name: edge.Node.Name.Full, Role: role})
	}

	return &ExternalMangaDraft{
		TitleOriginal: titleOriginal,
		TitleEn:       titleEn,
		Introduction:  m.Description,
		ImageURL:      m.CoverImage.Large,
		FirstDateJp:   m.StartDate.toISODate(),
		Status:        anilistStatusToOurs(m.Status),
		Authors:       authors,
		Genres:        m.Genres,
	}, nil
}
