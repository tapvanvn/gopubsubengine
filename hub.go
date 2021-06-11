package gopubsubengine

type Hub interface {
	SubscribeOn(topicPath string) (Subscriber, error)
	PublishOn(topicPath string) (Publisher, error)
}
