package main

import (
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/tapvanvn/gopubsubengine/awssqs"
	"github.com/tapvanvn/goutil"
)

func processor(message string) {

}

func main() {
	accessKey := goutil.MustGetEnv("ACCESS_KEY")
	accessSecret := goutil.MustGetEnv("ACCESS_SECRET")
	topic := "test"
	topicURL := goutil.MustGetEnv("QUEUE_URL")

	sessionConfig := &aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(accessKey, accessSecret, ""),
	}

	sess, err := session.NewSession(sessionConfig)
	if err != nil {
		log.Println("error establishing session")
		os.Exit(1)
	}
	hub, err := awssqs.NewAWSSQSHubFromSession(sess)
	if err != nil {
		log.Println("error creating hub")
	}

	hub.SetTopicQueueURL(topic, topicURL)
	publisher, err := hub.PublishOn(topic)
	if err != nil {
		log.Println("error creating publisher on topic")
	}
	sub, err := hub.SubscribeOn(topic)
	sub.SetProcessor(processor)
	if err != nil {
		log.Println("error creating subscribling on topic")
	}
	publisher.Publish("test")
	time.Sleep(15 * time.Second)
}
