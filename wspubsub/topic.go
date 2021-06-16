package wspubsub

import (
	"encoding/json"
)

type Topic struct {
	topic       string
	hub         *Hub
	subscribers []*Subscriber
	publishers  []*Publisher
}

func (topic *Topic) SendMessage(message interface{}) error {

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	sendMsg := Message{
		Topic:   topic.topic,
		Message: string(data),
	}
	return topic.hub.Send(&sendMsg)

}
