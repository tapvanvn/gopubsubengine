package gopubsubengine

type Publisher interface {
	Unpublish()
	Publish(message interface{})
	PublishAttributes(message interface{}, attributes map[string]string)
}
