package gopubsubengine

import "errors"

type Hub interface {
	SubscribeOn(topic string) (Subscriber, error)
	PublishOn(topic string) (Publisher, error)
	SetTier2OnStop(fn func())
}

var ErrTopicIsNotExisted = errors.New("Topic is not existed")
