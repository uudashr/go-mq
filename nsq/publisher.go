package nsq

import (
	"github.com/nsqio/go-nsq"
)

// Publisher publish message.
type Publisher struct {
	prod *nsq.Producer
}

// Publish message.
func (p *Publisher) Publish(topic string, body []byte) error {
	return p.prod.Publish(topic, body)
}

// Stop the publisher.
func (p *Publisher) Stop() error {
	p.prod.Stop()
	return nil
}

// NewPublisher constructs new Publisher.
func NewPublisher(addr string) (*Publisher, error) {
	cfg := nsq.NewConfig()
	prod, err := nsq.NewProducer(addr, cfg)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		prod: prod,
	}, nil
}
