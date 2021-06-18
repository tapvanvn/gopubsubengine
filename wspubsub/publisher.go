package wspubsub

type Publisher struct {
	topic *Topic
}

func (publisher *Publisher) Unpublish() {
	//TODO: apply this function
}
func (publisher *Publisher) Publish(message interface{}) {

	publisher.topic.SendMessage(message)
}
func (publisher *Publisher) PublishAttributes(message interface{}, attributes map[string]string) {
	publisher.topic.SendMessageAttributes(message, attributes)
}
