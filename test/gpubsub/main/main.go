package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/tapvanvn/gopubsubengine/test/gpubsub"
	"github.com/tapvanvn/goutil"
)

func main() {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := goutil.MustGetEnv("PROJECT_ID")
	// Creates a client.
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Sets the id for the new topic.
	topicID := "test-topic"

	// Creates the new topic.
	/*topic, err := client.CreateTopic(ctx, topicID)
	if err != nil {
		log.Fatalf("Failed to create topic: %v", err)
	}*/

	topic := client.Topic(topicID)
	//topic.Publish(context.Background(), )
	publisher := gpubsub.NewPublisher(topic)

	publisher.PublishRangeInt(10)

	subID := fmt.Sprintf("%s.1", topicID)
	subID2 := fmt.Sprintf("%s.2", topicID)

	/*sub2, err := client.CreateSubscription(ctx, subID2, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 20 * time.Second,
	})
	if err != nil {
		log.Fatal("err :", err)
	}*/

	sub := client.Subscription(subID)
	sub2 := client.Subscription(subID)
	_ = sub
	//subscriber := gpubsub.NewSubscriber(sub, subID)

	subscriber2 := gpubsub.NewSubscriber(sub2, subID2)

	go pullMessage(subscriber2)
	//go pullMessage(subscriber)
	/*if err != nil {
		log.Fatal("err :", err)
	}*/

	//go pullMessage(subscriber2)
	/*if err != nil {
		log.Fatal("err :", err)
	}*/
	//fmt.Printf("Topic %v created.\n", topic)

	time.Sleep(15 * time.Second)
}
func pullMessage(subscriber *gpubsub.Subscriber) {
	err := subscriber.PullMessage()
	if err != nil {
		log.Fatal("err :", err)
	}
}
