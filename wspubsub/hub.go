package wspubsub

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tapvanvn/gopubsubengine"
)

//TODO: when disconnect
//TODO: handle close
const (
	CtlReConnect = iota
)

type Hub struct {
	//ctlmx           sync.Mutex
	topics          map[string]*Topic
	url             string
	conn            *websocket.Conn
	publishTopics   map[string]int
	subscribeTopics map[string]int
	messages        chan *Message
	control         chan int
	onTier2Stop     func()
	hasTier2        bool
}

func NewWSPubSubHub(url string) (*Hub, error) {

	hub := &Hub{
		url:             url,
		conn:            nil,
		topics:          map[string]*Topic{},
		publishTopics:   map[string]int{},
		subscribeTopics: map[string]int{},
		messages:        make(chan *Message, 100),
		control:         make(chan int, 1),
		onTier2Stop:     nil,
		hasTier2:        false,
	}
	hub.control <- CtlReConnect

	go hub.run()

	return hub, nil
}
func (hub *Hub) SetTier2OnStop(fn func()) {
	hub.onTier2Stop = fn
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

	if topicHub, ok := hub.topics[topic]; ok {

		pick := rand.Intn(len(topicHub.subscribers))
		fmt.Println("pick", pick, len(topicHub.subscribers))
		subscriber := topicHub.subscribers[pick]
		go subscriber.processor(message)
	}
}
func (hub *Hub) runWriter() {
	for {

		msg := <-hub.messages
		if hub.conn != nil {
			//hub.ctlmx.Lock()
			err := hub.conn.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
			}
			//hub.ctlmx.Unlock()
		} else {
			hub.messages <- msg
		}
	}
}
func (hub *Hub) handleClose(code int, text string) error {

	fmt.Println("PUBSUB Closed", code, text)
	hub.conn = nil
	hub.control <- CtlReConnect
	return nil
}
func (hub *Hub) run() {

	go hub.runWriter()
	for {
		select {
		case ctl := <-hub.control:
			if ctl == CtlReConnect {
				fmt.Println("PUBSUB reconnect")
				c, _, err := websocket.DefaultDialer.Dial(hub.url, nil)
				if err != nil {

					hub.control <- CtlReConnect
					time.Sleep(time.Second * 2)

					break
				}

				hub.conn = c

				hub.conn.SetCloseHandler(hub.handleClose)
				//resend current pubsub topics
				for topic, _ := range hub.publishTopics {
					register := &Register{
						IsUnregister: false,
						IsPublisher:  true,
						Topic:        topic,
					}

					go hub.SendControl(register)
				}
				for topic, _ := range hub.subscribeTopics {
					register := &Register{
						IsUnregister: false,
						IsPublisher:  false,
						Topic:        topic,
					}

					go hub.SendControl(register)
				}
			}
		default:
			if hub.conn != nil {
				raw := &Message{}
				err := hub.conn.ReadJSON(raw)

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
						tier, hasTier := raw.Attributes["tier"]
						if hasTier {
							if tier == "one" {
								code, hasCode := raw.Attributes["raycode"]
								if hasCode {
									//send responsewith code
									hub.SendResponse(code)
								}
							}
						}
						topicString := raw.Topic
						topics := strings.Split(topicString, ",")
						for _, topic := range topics {
							topic = strings.TrimSpace(topic)
							if len(topic) > 0 {
								hub.pickOne(topic, raw.Message)
							}
						}
					}
				} else {

					fmt.Println("receive json fail:", err)
					hub.control <- CtlReConnect
					hub.conn = nil
					if hub.onTier2Stop != nil {
						hub.onTier2Stop()
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

func (hub *Hub) SendResponse(code string) error {

	msg := &Message{
		Topic:   "response",
		Message: "",
		Attributes: map[string]string{
			"raycode": code,
		},
	}
	return hub.Send(msg)
}

func (hub *Hub) Send(message *Message) error {
	hub.messages <- message
	return nil

}
