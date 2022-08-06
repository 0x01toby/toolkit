package event

type M map[string]any

// Listener interface
type Listener interface {
	Handle(e Event) error
}

type Event interface {
	Name() string
	Get(key string) any
	Set(key string, val any)
	Add(key string, val any)
	Data() M
	SetData(M) Event
	Abort(bool)
	IsAborted() bool
}

type ManagerFace interface {
	AddEvent(event Event)
	On(name string, listener Listener, priority ...int)
	Fire(name string, params M) (error, Event)
}

type Subscriber interface {
	SubscribedEvents() map[string]any
}
