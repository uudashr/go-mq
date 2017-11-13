package pubsub_test

import (
	"flag"
	"testing"
	"time"
)

var (
	flagProjectID       = flag.String("gcp.project-id", "", "Google Cloud Project ID")
	flagTopicID         = flag.String("gcp.topic-id", "", "Google Cloud Topic ID")
	flagSubscriptinID   = flag.String("gcp.subscription-id", "", "Google Cloud Subscription ID")
	flagCredentialsFile = flag.String("gcp.credentials-file", "", "Google Cloud Credentials File")
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
