package gpubsub

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
)

func NewSubscriber(sub *pubsub.Subscription, subID string) *Subscriber {
	return &Subscriber{
		sub:   sub,
		subID: subID,
	}
}

type Subscriber struct {
	sub   *pubsub.Subscription
	subID string
}

func (subscriber *Subscriber) PullMessage() error {
	var mu sync.Mutex
	received := 0
	ctx := context.Background()

	cctx, cancel := context.WithCancel(ctx)
	err := subscriber.sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		fmt.Printf("%s %d Got message: %q\n", subscriber.subID, received, string(msg.Data))
		msg.Ack()
		received++
		if received == 2 {
			cancel()
		}
	})
	if err != nil {
		return fmt.Errorf("%s Receive: %v", subscriber.subID, err)
	}
	return nil
}
