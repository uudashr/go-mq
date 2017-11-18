package pubsub_test

import (
	"context"
	"flag"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
)

var (
	flagProjectID       = flag.String("gcp.project-id", "", "Google Cloud Project ID")
	flagTopicID         = flag.String("gcp.topic-id", "", "Google Cloud Topic ID")
	flagSubscriptinID   = flag.String("gcp.subscription-id", "", "Google Cloud Subscription ID")
	flagCredentialsFile = flag.String("gcp.credentials-file", "", "Google Cloud Credentials File")
)

func ensureTopic(client *pubsub.Client, id string) (*pubsub.Topic, error) {
	topic := client.Topic(id)
	exists, err := topic.Exists(context.Background())
	if err != nil {
		return nil, err
	}

	if !exists {
		topic, err = client.CreateTopic(context.Background(), id)
		if err != nil {
			return nil, err
		}
	}

	return topic, nil
}

func ensureSubscription(client *pubsub.Client, topic *pubsub.Topic, id string) (*pubsub.Subscription, error) {
	subs := client.Subscription(id)
	exists, err := subs.Exists(context.Background())
	if err != nil {
		return nil, err
	}

	if !exists {
		cfg := pubsub.SubscriptionConfig{Topic: topic}
		subs, err = client.CreateSubscription(context.Background(), id, cfg)
		if err != nil {
			return nil, err
		}
	}

	return subs, nil
}

func expectMessage(t *testing.T, msgCh <-chan string, wantText string, errCh <-chan error, timeout time.Duration) {
	deadline := time.After(timeout)
	for {
		select {
		case msg := <-msgCh:
			if msg == wantText {
				return
			}
		case err := <-errCh:
			if err != nil {
				t.Error("err:", err)
			}
			t.Error("no message")
			return
		case <-deadline:
			t.Error("timeout")
			return
		}
	}
}
