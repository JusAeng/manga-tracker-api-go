package line_handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"

	"github.com/JusAeng/manga-tracker-api-go/config"
	"github.com/JusAeng/manga-tracker-api-go/service"
	"github.com/gofiber/fiber/v2"
)

const updateCommandText = "อัพเดทมังงะของฉัน"

type webhookPayload struct {
	Events []webhookEvent `json:"events"`
}

type webhookEvent struct {
	Type       string         `json:"type"`
	ReplyToken string         `json:"replyToken"`
	Source     webhookSource  `json:"source"`
	Message    webhookMessage `json:"message"`
}

type webhookSource struct {
	Type   string `json:"type"`
	UserId string `json:"userId"`
}

type webhookMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Webhook receives inbound events from the LINE Official Account. Unlike
// every other route in this codebase it's public (registered before
// auth_handlers.JWTMiddleware in router/route.go) — LINE authenticates
// itself via the X-Line-Signature header instead of a JWT.
func Webhook(c *fiber.Ctx) error {
	body := c.Body() // fasthttp exposes the raw body directly, no capturing middleware needed
	secret, _ := config.GetEnv("LineChannelSecret")
	if !verifySignature(secret, body, c.Get("X-Line-Signature")) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var payload webhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	for _, event := range payload.Events {
		if event.Type != "message" || event.Message.Type != "text" {
			continue
		}
		if strings.TrimSpace(event.Message.Text) != updateCommandText {
			continue
		}
		// A failed reply (e.g. LINE's API briefly down) must not fail the
		// webhook response itself, or LINE will retry the whole payload —
		// including events that already succeeded. Logged since this
		// codebase has no structured logger to route this through instead.
		if err := service.HandleUpdateRequest(event.Source.UserId, event.ReplyToken); err != nil {
			log.Println("line webhook reply failed:", err)
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

// verifySignature is the first signature-verification code in this repo —
// HMAC-SHA256 over the raw body using the channel secret, base64-encoded,
// compared with hmac.Equal (constant-time) rather than == to avoid a
// timing side-channel.
func verifySignature(secret string, body []byte, signature string) bool {
	if secret == "" || signature == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
