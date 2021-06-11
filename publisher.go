package gopubsubengine

type Publisher interface {
	Unpublish()
	Publish(message interface{})
}
