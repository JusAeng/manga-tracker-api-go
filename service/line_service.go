package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/JusAeng/manga-tracker-api-go/config"
	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/JusAeng/manga-tracker-api-go/repo"
)

const noUpdatesText = "ยังไม่มีมังงะที่คุณติดตามออกเล่มใหม่"

// HandleUpdateRequest answers a "my manga updates" request from the LINE
// webhook: look up the sender by their LINE sub, find any volumes
// published for manga they follow in the last 2 months, and reply. An
// unrecognized sender (no matching users.line_user_id — they've never
// logged into the app) gets the same "no updates" reply as a known user
// with nothing new; the spec only defines those two outcomes.
func HandleUpdateRequest(lineUserId, replyToken string) error {
	user, err := repo.GetUserByLineID(lineUserId)
	if err != nil {
		return err
	}
	if user == nil {
		return ReplyMessage(replyToken, noUpdatesText)
	}

	since := time.Now().AddDate(0, -2, 0)
	updates, err := repo.GetRecentVolumesForFollowedManga(user.ID, since)
	if err != nil {
		return err
	}

	return ReplyMessage(replyToken, BuildReplyText(updates))
}

// BuildReplyText is pure formatting, kept separate from the repo/network
// calls in HandleUpdateRequest so the message shape can be tested/reasoned
// about on its own.
func BuildReplyText(updates []*models.MangaVolumeUpdate) string {
	if len(updates) == 0 {
		return noUpdatesText
	}

	var b strings.Builder
	b.WriteString("มังงะที่คุณติดตามมีเล่มใหม่:\n")
	for _, u := range updates {
		title := u.MangaTitleEN
		if title == "" {
			title = u.MangaTitleOriginal
		}
		fmt.Fprintf(&b, "- %s เล่ม %d (%s)\n", title, u.VolumeNumber, u.PublishDate.Format("2006-01-02"))
	}
	return strings.TrimRight(b.String(), "\n")
}

type lineReplyRequest struct {
	ReplyToken string            `json:"replyToken"`
	Messages   []lineTextMessage `json:"messages"`
}

type lineTextMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ReplyMessage calls LINE's Reply API. Same stdlib net/http style as
// handlers/foo_handler.go's existing LINE API call, just a JSON body
// instead of form-encoded since the Reply API requires JSON.
func ReplyMessage(replyToken, text string) error {
	accessToken, _ := config.GetEnv("LineChannelAccessToken")

	payload, err := json.Marshal(lineReplyRequest{
		ReplyToken: replyToken,
		Messages:   []lineTextMessage{{Type: "text", Text: text}},
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.line.me/v2/bot/message/reply", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("line reply api error: %s", string(body))
	}
	return nil
}
