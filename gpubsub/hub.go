package gpubsub

import (
	"context"
	"log"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/tapvanvn/gopubsubengine"

	"google.golang.org/api/option"
)

type Hub struct {
	client *pubsub.Client
	topics map[string]*pubsub.Topic
}

func NewGPubsub(connectString string) (*Hub, error) {
	pos := strings.Index(connectString, ":")
	if pos > -1 {
		projectID := connectString[:pos]
		credentialPath := connectString[pos+1:]
		pubsubClient, err := pubsub.NewClient(context.TODO(), projectID, option.WithCredentialsFile(credentialPath))
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
			return nil, err
		}

		hub := &Hub{

			client: pubsubClient,
			topics: map[string]*pubsub.Topic{},
		}
		return hub, nil

	} else {

		projectID := connectString

		pubsubClient, err := pubsub.NewClient(context.TODO(), projectID)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
			return nil, err
		}
		hub := &Hub{
			client: pubsubClient,
			topics: map[string]*pubsub.Topic{},
		}
		return hub, nil
	}
}

func (hub *Hub) SetTier2OnStop(fn func()) {
	//hub.onTier2Stop = fn
}

func (hub *Hub) SubscribeOn(topic string) (gopubsubengine.Subscriber, error) {

	sub := hub.client.Subscription(topic)
	return NewSubscriber(sub)
}

func (hub *Hub) PublishOn(topic string) (gopubsubengine.Publisher, error) {

	if pubsubTopic, exist := hub.topics[topic]; exist {
		return NewPublisher(pubsubTopic), nil
	} else {
		pubsubTopic := hub.client.Topic(topic)
		if pubsubTopic == nil {
			return nil, gopubsubengine.ErrTopicIsNotExisted
		}
		hub.topics[topic] = pubsubTopic
		return NewPublisher(pubsubTopic), nil
	}
}

func Test() {

}
