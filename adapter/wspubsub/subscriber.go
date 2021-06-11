package wspubsub

type Subscriber struct {
	topic     *Topic
	processor func(message []byte)
}

func (sub *Subscriber) Unsubscribe() {
	//TODO: apply this function
}
func (sub *Subscriber) SetProcessor(processor func(message []byte)) {

	sub.processor = processor
}
