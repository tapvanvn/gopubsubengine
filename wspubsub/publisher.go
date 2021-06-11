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
