package mq

import "time"

// Message from server.
type Message interface {
	ID() string
	Body() []byte
	Timestamp() time.Time
	Ack() error
	Nack() error
}
