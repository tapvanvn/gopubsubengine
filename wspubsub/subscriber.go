package wspubsub

import "github.com/tapvanvn/gopubsubengine"

type Subscriber struct {
	topic     *Topic
	processor func(message string)
}

func (sub *Subscriber) Unsubscribe() {
	//TODO: apply this function
}
func (sub *Subscriber) SetProcessor(processor gopubsubengine.MessageProcessor) {

	sub.processor = processor
}
