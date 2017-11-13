package mq_test

import (
	mq "github.com/uudashr/go-mq"
)

func ExampleReceiver() {
	var recv mq.Receiver
	// ...

	handler := mq.HandlerFunc(func(msg mq.Message) {
		// TODO: handler message
		_ = msg.Ack()
	})

	err := recv.Listen(handler)
	if err != nil {
		// TODO: handle error
	}

	// Stop the receiver on other routine
	// recv.Stop()
}
