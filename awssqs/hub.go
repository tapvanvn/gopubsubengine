package awssqs

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/tapvanvn/gopubsubengine"
)

func NewAWSSQSHubFromSession(sess *session.Session) (*Hub, error) {
	return &Hub{
		sess:           sess,
		topicQueueURLs: map[string]string{},
	}, nil
}

func NewAWSSQSHub() (*Hub, error) {
	return nil, errors.New("not support now")
}

type Hub struct {
	sess           *session.Session
	topicQueueURLs map[string]string
}

func (hub *Hub) SetTopicQueueURL(topic string, queueURL string) {

	//TODO: should we protect this with mutex ?
	hub.topicQueueURLs[topic] = queueURL
}

func (hub *Hub) SetTier2OnStop(fn func()) {
	//hub.onTier2Stop = fn
}

func (hub *Hub) SubscribeOn(topic string) (gopubsubengine.Subscriber, error) {

	if hub.sess == nil {
		return nil, errors.New("Invalid session")
	}
	if queueURL, ok := hub.topicQueueURLs[topic]; ok {
		sub := NewSubscriber(hub.sess, topic, queueURL)
		sub.start()
		return sub, nil
	}
	return nil, gopubsubengine.ErrTopicIsNotExisted
}

func (hub *Hub) PublishOn(topic string) (gopubsubengine.Publisher, error) {
	if hub.sess == nil {
		return nil, errors.New("Invalid session")
	}
	if queueURL, ok := hub.topicQueueURLs[topic]; ok {
		return NewPublisher(hub.sess, queueURL, topic), nil
	}
	return nil, gopubsubengine.ErrTopicIsNotExisted
}
