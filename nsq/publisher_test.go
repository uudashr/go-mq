package nsq_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	gonsq "github.com/nsqio/go-nsq"
	nsq "github.com/uudashr/go-mq/nsq"
)

func TestPublish(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip on short mode")
	}

	// Publisher
	pub, err := nsq.NewPublisher(*flagAddr)
	if err != nil {
		t.Fatal("err:", err)
	}

	defer func() {
		if stopErr := pub.Stop(); stopErr != nil {
			t.Error("err:", err)
		}
	}()

	testMessage := fmt.Sprintf("Hello World [%s]", time.Now().Format(time.RFC3339Nano))
	err = pub.Publish(*flagTopic, []byte(testMessage))
	if err != nil {
		t.Fatal("err:", err)
	}

	// Consumer
	cfg := gonsq.NewConfig()
	cons, err := gonsq.NewConsumer(*flagTopic, *flagChannel, cfg)
	if err != nil {
		t.Fatal("err:", err)
	}

	defer func() {
		cons.Stop()
		<-cons.StopChan
	}()

	msgCh := make(chan string)
	errCh := make(chan error, 1)
	cons.AddHandler(gonsq.HandlerFunc(func(msg *gonsq.Message) error {
		msg.DisableAutoResponse()
		msg.Finish()

		msgCh <- string(msg.Body)
		return nil
	}))

	log.Println("Connecting...")
	err = cons.ConnectToNSQLookupd(*flagLookupdAddr)
	if err != nil {
		t.Fatal("err:", err)
	}

	log.Println("Connected")

	expectMessage(t, msgCh, testMessage, errCh, 5*time.Second)
}
