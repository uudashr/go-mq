package pubsub

import (
	"context"
	"errors"
	"sync/atomic"

	"cloud.google.com/go/pubsub"

	mq "github.com/uudashr/go-mq"
	"google.golang.org/api/option"
)

const (
	stateCreated int32 = iota
	stateConnecting
	stateListening
	stateStopping
	stateStopped
)

// Receiver receives message.
type Receiver struct {
	projectID      string
	subscriptionID string
	opts           []option.ClientOption

	recvCtx  context.Context
	stopRecv context.CancelFunc

	state int32
}

// Listen to the incoming message.
func (r *Receiver) Listen(h mq.Handler) (retErr error) {
	if !atomic.CompareAndSwapInt32(&r.state, stateCreated, stateConnecting) {
		return errors.New("pubsub: not in created state")
	}

	defer atomic.StoreInt32(&r.state, stateStopped)

	client, err := pubsub.NewClient(context.Background(), r.projectID, r.opts...)
	if err != nil {
		return err
	}

	defer func() {
		if err := client.Close(); err != nil && retErr != nil {
			retErr = err
		}
	}()

	atomic.StoreInt32(&r.state, stateListening)
	subscription := client.Subscription(r.subscriptionID)
	defer r.stopRecv()
	err = subscription.Receive(r.recvCtx, func(ctx context.Context, msg *pubsub.Message) {
		wrap := &message{msg: msg}
		h.Handle(wrap)
	})

	if err != context.Canceled {
		return err
	}

	return nil
}

// Stop the receiver.
func (r *Receiver) Stop() error {
	if !atomic.CompareAndSwapInt32(&r.state, stateListening, stateStopping) {
		return errors.New("pubsub: cannot stop non-listening receiver")
	}

	r.stopRecv()
	return nil
}

// NewReceiver construct new Receiver.
func NewReceiver(projectID, subscriptionID string, opts ...option.ClientOption) (*Receiver, error) {
	recvCtx, stopRecv := context.WithCancel(context.Background())
	return &Receiver{
		projectID:      projectID,
		subscriptionID: subscriptionID,
		opts:           opts,
		recvCtx:        recvCtx,
		stopRecv:       stopRecv,
	}, nil
}
