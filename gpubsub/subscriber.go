package gpubsub

import (
	"context"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/tapvanvn/gopubsubengine"
)

func NewSubscriber(sub *pubsub.Subscription) (*Subscriber, error) {

	subscriber := &Subscriber{
		sub: sub,
	}
	go subscriber.run()
	return subscriber, nil
}

type Subscriber struct {
	sub        *pubsub.Subscription
	processor  func(message string)
	cancelFunc context.CancelFunc
}

func (sub *Subscriber) run() {
	ctx := context.Background()
	cctx, cancel := context.WithCancel(ctx)
	sub.cancelFunc = cancel
	var mu sync.Mutex
	err := sub.sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		msg.Ack()
		defer mu.Unlock()

		if sub.processor != nil {
			go sub.processor(string(msg.Data))
		}
	})
	if err != nil {

	}
}
func (sub *Subscriber) Unsubscribe() {
	//TODO: apply this function
}

func (sub *Subscriber) SetProcessor(processor gopubsubengine.MessageProcessor) {

	sub.processor = processor
}
