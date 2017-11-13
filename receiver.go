package mq

// Receiver receives message.
type Receiver interface {
	Listen(Handler) error
	Stop() error
}

// Handler handles the message.
type Handler interface {
	Handle(Message)
}

// HandlerFunc is the function adapter for Handler.
type HandlerFunc func(Message)

// Handle the message.
func (h HandlerFunc) Handle(m Message) {
	h(m)
}
