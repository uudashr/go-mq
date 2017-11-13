package mq

// Publisher publish message.
type Publisher interface {
	Publish(topic string, body []byte) error
	Stop() error
}
