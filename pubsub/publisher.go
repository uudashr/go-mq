package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"

	"google.golang.org/api/option"
)

// Publisher publish message.
type Publisher struct {
	client *pubsub.Client
	topic  *pubsub.Topic
}

// Publish message body to a topic.
func (p *Publisher) Publish(topic string, body []byte) error {
	res := p.topic.Publish(context.Background(), &pubsub.Message{
		Data: body,
	})
	_, err := res.Get(context.Background())
	return err
}

// Stop the publisher.
func (p *Publisher) Stop() error {
	p.topic.Stop()
	return p.client.Close()
}

// NewPublisher constructs new Publisher.
func NewPublisher(projectID, topicID string, opts ...option.ClientOption) (*Publisher, error) {
	client, err := pubsub.NewClient(context.Background(), projectID, opts...)
	if err != nil {
		return nil, err
	}

	topic := client.Topic(topicID)

	return &Publisher{
		client: client,
		topic:  topic,
	}, nil
}
