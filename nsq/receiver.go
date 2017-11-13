package nsq

import (
	"errors"
	"sync/atomic"

	nsq "github.com/nsqio/go-nsq"
	mq "github.com/uudashr/go-mq"
)

const (
	stateCreated int32 = iota
	stateConnecting
	stateListening
	stateStopping
	stateStopped
)

// Receiver receive message.
type Receiver struct {
	lookupdAddr string
	cons        *nsq.Consumer
	state       int32
}

// Listen to incoming message. This will block until error found or the listening stoppped.
func (r *Receiver) Listen(h mq.Handler) error {
	if !atomic.CompareAndSwapInt32(&r.state, stateCreated, stateConnecting) {
		return errors.New("nsq: not in created state")
	}

	defer func() {
		atomic.StoreInt32(&r.state, stateStopped)
	}()

	r.cons.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
		msg.DisableAutoResponse()
		wrap := &message{msg: msg}
		h.Handle(wrap)
		return nil
	}))

	if err := r.cons.ConnectToNSQLookupd(r.lookupdAddr); err != nil {
		return err
	}

	atomic.StoreInt32(&r.state, stateListening)
	<-r.cons.StopChan
	return nil
}

// Stop the receiver.
func (r *Receiver) Stop() error {
	if !atomic.CompareAndSwapInt32(&r.state, stateListening, stateStopping) {
		return errors.New("nsq: cannot stop non-listening receiver")
	}

	r.cons.Stop()
	return nil
}

// NewReceiver constructs new Receiver.
func NewReceiver(topic, channel, lookupdAddr string, cfg *nsq.Config) (*Receiver, error) {
	cons, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		return nil, err
	}

	return &Receiver{
		lookupdAddr: lookupdAddr,
		cons:        cons,
	}, nil
}
