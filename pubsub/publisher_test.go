package pubsub_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/uudashr/go-mq/pubsub"
	"google.golang.org/api/option"
	"gopkg.in/mattes/go-expand-tilde.v1"

	gpubsub "cloud.google.com/go/pubsub"
)

func TestPublish(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip non short mode")
	}

	credsFile, err := tilde.Expand(*flagCredentialsFile)
	if err != nil {
		t.Fatal("err:", err)
	}

	client, err := gpubsub.NewClient(context.Background(), *flagProjectID, option.WithCredentialsFile(credsFile))
	if err != nil {
		t.Fatal("err:", err)
	}

	topic, err := ensureTopic(client, *flagTopicID)
	if err != nil {
		t.Fatal("err:", err)
	}
	defer func() {
		if err := topic.Delete(context.Background()); err != nil {
			t.Error("err:", err)
		}
	}()

	subs, err := ensureSubscription(client, topic, *flagSubscriptinID)
	if err != nil {
		t.Fatal("err:", err)
	}
	defer func() {
		if err := subs.Delete(context.Background()); err != nil {
			t.Error("err:", err)
		}
	}()

	// Publish
	pub, err := pubsub.NewPublisher(context.Background(), *flagProjectID, option.WithCredentialsFile(credsFile))
	if err != nil {
		t.Fatal("err:", err)
	}

	testMessage := fmt.Sprintf("Hello World [%s]", time.Now().Format(time.RFC3339Nano))
	if err = pub.Publish(*flagTopicID, []byte(testMessage)); err != nil {
		t.Fatal("err:", err)
	}

	if err = pub.Stop(); err != nil {
		t.Fatal("err:", err)
	}

	// Receiver
	recvCtx, stopRecv := context.WithCancel(context.Background())
	msgCh := make(chan string)
	errCh := make(chan error, 1)
	go func() {
		err := subs.Receive(recvCtx, func(ctx context.Context, msg *gpubsub.Message) {
			msg.Ack()
			msgCh <- string(msg.Data)
		})

		if err != nil {
			errCh <- err
		}

		close(errCh)
	}()

	defer stopRecv()

	expectMessage(t, msgCh, testMessage, errCh, 5*time.Second)
}
