package gopubsubengine

type MessageProcessor func(message string)

type Subscriber interface {
	Unsubscribe()
	SetProcessor(processor MessageProcessor)
}
