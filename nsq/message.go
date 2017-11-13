package nsq

import (
	"time"

	nsq "github.com/nsqio/go-nsq"
)

type message struct {
	msg *nsq.Message
}

func (m *message) ID() string {
	return string(m.msg.ID[:])
}

func (m *message) Body() []byte {
	return m.msg.Body
}

func (m *message) Timestamp() time.Time {
	return time.Unix(0, m.msg.Timestamp)
}

func (m *message) Ack() error {
	m.msg.Finish()
	return nil
}

func (m *message) Nack() error {
	m.msg.Requeue(-1)
	return nil
}
