package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tapvanvn/gopubsubengine/gpubsub"
)

func onMessage(message string) {
	fmt.Printf("receive:%s\n", message)
}

func main() {
	hub, err := gpubsub.NewGPubSubHub("newcontinent-platform:../../credential.json")
	if err != nil {
		log.Fatal(err)
	}
	publisher, err := hub.PublishOn("test-topic")
	if err != nil {
		log.Fatal(err)
	}
	subscriber, err := hub.SubscribeOn("test-topic")
	if err != nil {
		log.Fatal(err)
	}
	subscriber.SetProcessor(onMessage)

	publisher.Publish("test message")

	time.Sleep(15 * time.Second)
}
