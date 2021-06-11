package wspubsub

type Message struct {
	Topic   string      `json:"topic"`
	Message interface{} `json:"message"`
}
