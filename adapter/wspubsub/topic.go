package wspubsub

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Topic struct {
	topic       string
	hub         *Hub
	subscribers []*Subscriber
	publishers  []*Publisher
}

func (topic *Topic) SendMessage(message interface{}) error {

	sendMsg := Message{
		Topic:   topic.topic,
		Message: message,
	}
	data, err := json.Marshal(sendMsg)
	if err != nil {
		return err
	}
	err = topic.hub.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}
	return nil
}
