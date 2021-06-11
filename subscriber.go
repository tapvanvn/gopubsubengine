package gopubsubengine

type MessageProcessor func(message interface{})

type Subscriber interface {
	Unsubscribe()
	SetProcessor(processor MessageProcessor)
}
