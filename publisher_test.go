package mq_test

import "github.com/uudashr/go-mq"

func ExamplePublisher() {
	const topic = "greetings"

	var pub mq.Publisher
	err := pub.Publish(topic, []byte("Hello World"))
	if err != nil {
		// TODO: handle error
	}

	err = pub.Stop()
	if err != nil {
		// TODO: handle error
	}
}
