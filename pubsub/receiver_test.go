package pubsub_test

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/uudashr/go-mq"
	"github.com/uudashr/go-mq/pubsub"
	tilde "gopkg.in/mattes/go-expand-tilde.v1"

	gpubsub "cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func TestReceive(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip non short mode")
	}

	if *flagProjectID == "" || *flagTopicID == "" || *flagSubscriptinID == "" || *flagCredentialsFile == "" {
		flag.Usage()
		t.Fatal("Test flag required")
	}

	credsFile, err := tilde.Expand(*flagCredentialsFile)
	if err != nil {
		t.Fatal("err:", err)
	}

	// Publisher
	client, err := gpubsub.NewClient(context.Background(), *flagProjectID, option.WithCredentialsFile(credsFile))
	if err != nil {
		t.Fatal("err:", err)
	}

	defer func() {
		if err = client.Close(); err != nil {
			t.Error("err:", err)
		}
	}()

	topic := client.Topic(*flagTopicID)
	exists, err := topic.Exists(context.Background())
	if err != nil {
		t.Fatal("err:", err)
	}

	if !exists {
		topic, err = client.CreateTopic(context.Background(), *flagTopicID)
		if err != nil {
			t.Fatal("err:", err)
		}
	}

	defer func() {
		if err = topic.Delete(context.Background()); err != nil {
			t.Error("err:", err)
		}
	}()

	subs := client.Subscription(*flagSubscriptinID)
	exists, err = subs.Exists(context.Background())
	if err != nil {
		t.Fatal("err:", err)
	}

	if !exists {
		subs, err = client.CreateSubscription(context.Background(), *flagSubscriptinID, gpubsub.SubscriptionConfig{Topic: topic})
		if err != nil {
			t.Fatal("err:", err)
		}
	}

	defer func() {
		if err = subs.Delete(context.Background()); err != nil {
			t.Error("err:", err)
		}
	}()

	testMessage := fmt.Sprintf("Hello World [%s]", time.Now().Format(time.RFC3339Nano))
	pubRes := topic.Publish(context.Background(), &gpubsub.Message{Data: []byte(testMessage)})
	_, err = pubRes.Get(context.Background())
	if err != nil {
		t.Fatal("err:", err)
	}

	topic.Stop()

	recv, err := pubsub.NewReceiver(*flagProjectID, *flagSubscriptinID, option.WithCredentialsFile(credsFile))
	if err != nil {
		t.Fatal("err:", err)
	}

	msgCh := make(chan string)
	errCh := make(chan error, 1)
	go func() {
		err = recv.Listen(mq.HandlerFunc(func(msg mq.Message) {
			if ackErr := msg.Ack(); ackErr != nil {
				t.Error("err:", err)
			}
			msgCh <- string(msg.Body())
		}))

		if err != nil {
			errCh <- err
		}

		close(errCh)
	}()

	expectMessage(t, msgCh, testMessage, errCh, 5*time.Second)
}
