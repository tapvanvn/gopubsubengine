package gopubsubengine_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/tapvanvn/gopubsubengine"
	"github.com/tapvanvn/gopubsubengine/wspubsub"
)

func TestHub(t *testing.T) {
	var publisher gopubsubengine.Publisher = nil
	hub, err := wspubsub.NewWSPubSubHub("ws://192.168.1.8:80/ws")
	if err != nil {
		t.Fatal(err)
		return
	}
	publisher, err = hub.PublishOn("test")
	if err != nil {
		t.Fatal(err)
		return
	}

	hub2, err := wspubsub.NewWSPubSubHub("ws://192.168.1.8:80/ws")
	if err != nil {
		t.Fatal(err)
		return
	}
	var subscriber gopubsubengine.Subscriber = nil
	subscriber, err = hub2.SubscribeOn("test")
	if err != nil {
		t.Fatal(err)
		return
	}
	subscriber.SetProcessor(func(message interface{}) {
		fmt.Println("receive:", message)
	})

	for {
		publisher.Publish("test")
		fmt.Println("send")
		time.Sleep(10 * time.Second)
	}
}
