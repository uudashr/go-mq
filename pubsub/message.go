package pubsub

import (
	"time"

	"cloud.google.com/go/pubsub"
)

type message struct {
	msg *pubsub.Message
}

func (m *message) ID() string {
	return m.msg.ID
}

func (m *message) Body() []byte {
	return m.msg.Data
}

func (m *message) Timestamp() time.Time {
	return m.msg.PublishTime
}

func (m *message) Ack() error {
	m.msg.Ack()
	return nil
}

func (m *message) Nack() error {
	m.msg.Nack()
	return nil
}
