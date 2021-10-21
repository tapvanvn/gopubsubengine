package gpubsub

import (
	"github.com/tapvanvn/gopubsubengine"
	_ "github.com/tapvanvn/gopubsubengine"
	//cloud.google.com/go/pubsub
)

type Hub struct {
}

func NewGPubsub(connectString string) (*Hub, error) {

	return &Hub{}, nil
}
func (hub *Hub) SetTier2OnStop(fn func()) {
	//hub.onTier2Stop = fn
}

func (hub *Hub) SubscribeOn(topic string) (gopubsubengine.Subscriber, error) {
	return nil, nil
}

func (hub *Hub) PublishOn(topic string) (gopubsubengine.Publisher, error) {
	return nil, nil
}

func Test() {

}
