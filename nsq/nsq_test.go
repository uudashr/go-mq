package nsq_test

import (
	"flag"
	"testing"
	"time"
)

var (
	flagAddr        = flag.String("nsq.nsqd-addr", "127.0.0.1:4150", "nsqd address")
	flagLookupdAddr = flag.String("nsq.lookupd-addr", "127.0.0.1:4161", "nsqlookupd address")
	flagTopic       = flag.String("nsq.topic", "greet", "topic")
	flagChannel     = flag.String("nsq.channel", "public", "channel")
)

func expectMessage(t *testing.T, msgCh <-chan string, wantText string, errCh <-chan error, timeout time.Duration) {
	deadline := time.After(timeout)
	for {
		select {
		case msg := <-msgCh:
			if msg == wantText {
				return
			}
		case err := <-errCh:
			if err != nil {
				t.Error("err:", err)
			}
			t.Error("no message")
			return
		case <-deadline:
			t.Error("timeout")
			return
		}
	}
}
