package gopubsubengine

type Hub interface {
	SubscribeOn(topic string) (Subscriber, error)
	PublishOn(topic string) (Publisher, error)
	SetTier2OnStop(fn func())
}
