package gopubsubengine

type MessageProcessor func(message []byte)

type Subscriber interface {
	Unsubscribe()
	SetProcessor(processor MessageProcessor)
}
