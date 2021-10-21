package gpubsub

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
)

func NewPublisher(topic *pubsub.Topic) *Publisher {
	return &Publisher{
		topic: topic,
	}
}

type Publisher struct {
	topic *pubsub.Topic
}

func (publisher *Publisher) PublishRangeInt(n int) error {
	ctx := context.Background()

	var wg sync.WaitGroup
	var totalErrors uint64

	for i := 0; i < n; i++ {
		result := publisher.topic.Publish(ctx, &pubsub.Message{
			Data: []byte("Message " + strconv.Itoa(i)),
		})

		wg.Add(1)
		go func(i int, res *pubsub.PublishResult) {
			defer wg.Done()
			// The Get method blocks until a server-generated ID or
			// an error is returned for the published message.
			id, err := res.Get(ctx)
			if err != nil {
				// Error handling code can be added here.
				fmt.Printf("Failed to publish: %v", err)
				atomic.AddUint64(&totalErrors, 1)
				return
			}
			fmt.Printf("Published message %d; msg ID: %v\n", i, id)
		}(i, result)
	}

	wg.Wait()

	if totalErrors > 0 {
		return fmt.Errorf("%d of %d messages did not publish successfully", totalErrors, n)
	}
	return nil
}
