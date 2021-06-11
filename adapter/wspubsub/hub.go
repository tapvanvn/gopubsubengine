package wspubsub

import (
	"encoding/json"
	"strings"

	"github.com/gorilla/websocket"
)

//TODO: when disconnect
//TODO: handle close

type Hub struct {
	topics map[string]*Topic
	url    string
	conn   *websocket.Conn
}

func (hub *Hub) getTopic(topic string) *Topic {
	if topicHub, ok := hub.topics[topic]; ok {
		return topicHub
	}
	topicHub := &Topic{
		hub:         hub,
		subscribers: make([]*Subscriber, 0),
		publishers:  make([]*Publisher, 0),
	}
	hub.topics[topic] = topicHub
	return topicHub
}

func (hub *Hub) broadcast(topic string, message interface{}) {

	if topicHub, ok := hub.topics[topic]; ok {

		for _, subscriber := range topicHub.subscribers {
			if subscriber.processor != nil {
				subscriber.processor(message)
			}
		}
	}
}
func (hub *Hub) run() {
	for {
		_, message, err := hub.conn.ReadMessage()
		if err != nil {
			// handle error
		}
		raw := &Message{}

		err = json.Unmarshal(message, &raw)
		if err == nil {

			topicString := raw.Topic
			topics := strings.Split(topicString, ",")
			for _, topic := range topics {
				topic = strings.TrimSpace(topic)
				if len(topic) > 0 {
					hub.broadcast(topic, raw.Message)
				}
			}

		}
	}
}

func (hub *Hub) SubscribeOn(topic string) (*Subscriber, error) {

	topicHub := hub.getTopic(topic)

	subscriber := &Subscriber{
		topic: topicHub,
	}
	topicHub.subscribers = append(topicHub.subscribers, subscriber)

	return subscriber, nil
}

func (hub *Hub) PublishOn(topic string) (*Publisher, error) {
	topicHub := hub.getTopic(topic)

	publisher := &Publisher{
		topic: topicHub,
	}
	topicHub.publishers = append(topicHub.publishers, publisher)

	return publisher, nil
}

func NewWSPubSubHub(url string) (*Hub, error) {
	hub := &Hub{
		url:    url,
		topics: map[string]*Topic{},
	}

	c, _, err := websocket.DefaultDialer.Dial(hub.url, nil)
	if err != nil {
		return nil, err
	}

	hub.conn = c

	// send message
	/*err := c.WriteMessage(websocket.TextMessage, {message})
	  if err != nil {
	      // handle error
	  }
	      .......
	  // receive message
	*/

	return hub, nil
}
