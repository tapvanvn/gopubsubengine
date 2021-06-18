package wspubsub

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/tapvanvn/gopubsubengine"
)

//TODO: when disconnect
//TODO: handle close

type Hub struct {
	topics          map[string]*Topic
	url             string
	conn            *websocket.Conn
	publishTopics   map[string]int
	subscribeTopics map[string]int
	messages        chan *Message
}

func NewWSPubSubHub(url string) (*Hub, error) {
	hub := &Hub{
		url:             url,
		topics:          map[string]*Topic{},
		publishTopics:   map[string]int{},
		subscribeTopics: map[string]int{},
		messages:        make(chan *Message),
	}

	go hub.run()

	return hub, nil
}

func (hub *Hub) getTopic(topic string) *Topic {
	if topicHub, ok := hub.topics[topic]; ok {
		return topicHub
	}
	topicHub := &Topic{
		hub:         hub,
		topic:       topic,
		subscribers: make([]*Subscriber, 0),
		publishers:  make([]*Publisher, 0),
	}
	hub.topics[topic] = topicHub
	return topicHub
}

func (hub *Hub) broadcast(topic string, message string) {

	if topicHub, ok := hub.topics[topic]; ok {

		for _, subscriber := range topicHub.subscribers {
			if subscriber.processor != nil {
				go subscriber.processor(message)
			}
		}
	}
}
func (hub *Hub) pickOne(topic string, message string) {
	fmt.Println("pick one", message)
	if topicHub, ok := hub.topics[topic]; ok {

		pick := rand.Intn(len(topicHub.subscribers))
		fmt.Println("choice", pick)
		subscriber := topicHub.subscribers[pick]
		go subscriber.processor(message)
	}
}
func (hub *Hub) runWriter() {
	for {
		msg := <-hub.messages
		err := hub.conn.WriteJSON(msg)
		if err != nil {
			fmt.Println(err)
		}
	}
}
func (hub *Hub) run() {

	c, _, err := websocket.DefaultDialer.Dial(hub.url, nil)
	if err != nil {

		fmt.Println(err)
		return
	}

	hub.conn = c
	go hub.runWriter()
	for {
		_, message, err := hub.conn.ReadMessage()
		//fmt.Println("receive:", message)
		if err != nil {
			// handle error
			fmt.Println("error:", err)
		}
		raw := &Message{}

		err = json.Unmarshal(message, &raw)
		if err == nil {
			msgType, ok := raw.Attributes["type"]
			if !ok || msgType != "pick_one" {
				topicString := raw.Topic
				topics := strings.Split(topicString, ",")
				for _, topic := range topics {
					topic = strings.TrimSpace(topic)
					if len(topic) > 0 {
						hub.broadcast(topic, raw.Message)
					}
				}
			} else {
				topicString := raw.Topic
				topics := strings.Split(topicString, ",")
				for _, topic := range topics {
					topic = strings.TrimSpace(topic)
					if len(topic) > 0 {
						hub.pickOne(topic, raw.Message)
					}
				}
			}
		}
	}
}

func (hub *Hub) SubscribeOn(topic string) (gopubsubengine.Subscriber, error) {

	topicHub := hub.getTopic(topic)

	subscriber := &Subscriber{
		topic: topicHub,
	}
	needSubscribe := false
	if numSubscribe, ok := hub.subscribeTopics[topic]; ok {
		if numSubscribe <= 0 {
			needSubscribe = true
		}
	} else {
		needSubscribe = true
		hub.subscribeTopics[topic] = 0
	}
	hub.subscribeTopics[topic] = hub.subscribeTopics[topic] + 1

	if needSubscribe {

		register := &Register{
			IsUnregister: false,
			IsPublisher:  false,
			Topic:        topic,
		}
		go hub.SendControl(register)
		fmt.Println("send subscribe on topic:", topic)
	}

	topicHub.subscribers = append(topicHub.subscribers, subscriber)

	return subscriber, nil
}

func (hub *Hub) PublishOn(topic string) (gopubsubengine.Publisher, error) {
	topicHub := hub.getTopic(topic)

	publisher := &Publisher{
		topic: topicHub,
	}
	topicHub.publishers = append(topicHub.publishers, publisher)

	needPublish := false
	if numPublisher, ok := hub.publishTopics[topic]; ok {
		if numPublisher <= 0 {
			needPublish = true
		}
	} else {
		needPublish = true
		hub.publishTopics[topic] = 0
	}
	hub.publishTopics[topic] = hub.publishTopics[topic] + 1

	if needPublish {

		register := &Register{
			IsUnregister: false,
			IsPublisher:  true,
			Topic:        topic,
		}

		go hub.SendControl(register)

		fmt.Println("send publish on topic:", topic)
	}

	topicHub.publishers = append(topicHub.publishers, publisher)
	return publisher, nil
}

func (hub *Hub) SendControl(register *Register) error {
	data, err := json.Marshal(register)
	if err != nil {
		return err
	}
	msg := &Message{
		Topic:   "control",
		Message: string(data),
	}
	return hub.Send(msg)
}

func (hub *Hub) Send(message *Message) error {
	hub.messages <- message
	return nil

	//return hub.conn.WriteJSON(message)
}
