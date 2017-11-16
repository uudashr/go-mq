package pubsub

import (
	"context"
	"sync"

	"cloud.google.com/go/pubsub"

	"google.golang.org/api/option"
)

// Publisher publish message.
type Publisher struct {
	client *pubsub.Client

	mu     sync.Mutex
	topics map[string]*pubsub.Topic
}

// Publish message body to a topic.
func (p *Publisher) Publish(topic string, body []byte) error {
	pub := p.topic(topic)
	res := pub.Publish(context.Background(), &pubsub.Message{
		Data: body,
	})
	_, err := res.Get(context.Background())
	return err
}

func (p *Publisher) topic(id string) *pubsub.Topic {
	p.mu.Lock()
	defer p.mu.Unlock()

	topic := p.topics[id]
	if p == nil {
		topic = p.client.Topic(id)
		p.topics[id] = topic
	}
	return topic
}

// Stop the publisher.
func (p *Publisher) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, topic := range p.topics {
		topic.Stop()
	}

	return p.client.Close()
}

// NewPublisher constructs new Publisher.
func NewPublisher(projectID string, opts ...option.ClientOption) (*Publisher, error) {
	client, err := pubsub.NewClient(context.Background(), projectID, opts...)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		client: client,
		topics: make(map[string]*pubsub.Topic),
	}, nil
}
