package gpubsub

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
)

func NewPublisher(topic *pubsub.Topic) *Publisher {
	return &Publisher{
		topic: topic,
	}
}

type Publisher struct {
	topic *pubsub.Topic
}

func (publisher *Publisher) Unpublish() {
	//TODO: apply this function
}
func (publisher *Publisher) Publish(message interface{}) {

	data, err := json.Marshal(message)
	if err != nil {
		//TODO: handle error
		return
	}

	result := publisher.topic.Publish(context.Background(), &pubsub.Message{
		Data: data,
	})

	//TODO: handle result
	_ = result
}
func (publisher *Publisher) PublishAttributes(message interface{}, attributes map[string]string) {
	//publisher.topic.SendMessageAttributes(message, attributes)
	data, err := json.Marshal(message)
	if err != nil {
		//TODO: handle error
		return
	}
	msg := &pubsub.Message{
		Data: data,
	}
	for key, value := range attributes {

		msg.Attributes[key] = value
	}

	result := publisher.topic.Publish(context.Background(), msg)

	//TODO: handle result
	_ = result
}
