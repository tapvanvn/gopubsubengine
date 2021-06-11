package wspubsub

type Subscriber struct {
	topic     *Topic
	processor func(message interface{})
}

func (sub *Subscriber) Unsubscribe() {
	//TODO: apply this function
}
func (sub *Subscriber) SetProcessor(processor func(message interface{})) {

	sub.processor = processor
}
