package nsq_test

import (
	"fmt"
	"testing"
	"time"

	mq "github.com/uudashr/go-mq"
	nsq "github.com/uudashr/go-mq/nsq"

	gonsq "github.com/nsqio/go-nsq"
)

func TestReceive(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip on short mode")
	}

	// Producer
	cfg := gonsq.NewConfig()
	prod, err := gonsq.NewProducer(*flagAddr, cfg)
	if err != nil {
		t.Fatal("err:", err)
	}

	defer prod.Stop()

	testMessage := fmt.Sprintf("Hello World [%s]", time.Now().Format(time.RFC3339Nano))
	err = prod.Publish(*flagTopic, []byte(testMessage))
	if err != nil {
		t.Fatal("err:", err)
	}

	// Receiver
	recv, err := nsq.NewReceiver(*flagTopic, *flagChannel, *flagLookupdAddr, cfg)
	if err != nil {
		t.Fatal("err:", err)
	}

	msgCh := make(chan string)
	errCh := make(chan error, 1)
	go func() {
		err := recv.Listen(mq.HandlerFunc(func(msg mq.Message) {
			if ackErr := msg.Ack(); ackErr != nil {
				t.Error("err:", ackErr)
			}
			msgCh <- string(msg.Body())
		}))

		if err != nil {
			errCh <- err
		}

		close(errCh)
	}()

	defer func() {
		if stopErr := recv.Stop(); stopErr != nil {
			t.Error("err:", stopErr)
		}
	}()

	expectMessage(t, msgCh, testMessage, errCh, 5*time.Second)
}
