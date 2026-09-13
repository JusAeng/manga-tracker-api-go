package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"

	"cloud.google.com/go/pubsub/v2"
	"github.com/JusAeng/manga-tracker-api-go/config"
)

var (
	pubsubClientOnce sync.Once
	pubsubClient     *pubsub.Client
	pubsubClientErr  error
)

func getPubSubClient(ctx context.Context) (*pubsub.Client, error) {
	pubsubClientOnce.Do(func() {
		projectId, _ := config.GetEnv("GCPProjectId")
		if projectId == "" {
			pubsubClientErr = errors.New("GCPPROJECTID is not set")
			return
		}
		pubsubClient, pubsubClientErr = pubsub.NewClient(ctx, projectId)
	})
	return pubsubClient, pubsubClientErr
}

// PublishRatingUpdated tells the recommender service (Python, via Pub/Sub
// -> Cloud Run) that this user's ratings changed, so it can recompute
// their recommendations. Never blocks or fails the rating request itself
// on error — recommendations refreshing is a side effect, not something a
// user should see fail just because Pub/Sub had a bad moment. Callers
// should log the error and move on, same as the LINE webhook's reply
// failures.
func PublishRatingUpdated(userId string) error {
	topicName, _ := config.GetEnv("PubSubRatingTopic")
	if topicName == "" {
		log.Println("PUBSUBRATINGTOPIC not set — skipping recommendation refresh publish")
		return nil
	}

	ctx := context.Background()
	client, err := getPubSubClient(ctx)
	if err != nil {
		return err
	}

	data, err := json.Marshal(map[string]string{"user_id": userId})
	if err != nil {
		return err
	}

	publisher := client.Publisher(topicName)
	defer publisher.Stop()

	result := publisher.Publish(ctx, &pubsub.Message{Data: data})
	_, err = result.Get(ctx)
	return err
}
