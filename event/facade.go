package event

var Facade = NewManager("default")

func On(name string, listener Listener, priority ...int) {
	Facade.On(name, listener, priority...)
}

func Listen(name string, listener Listener, priority ...int) {
	Facade.Listen(name, listener, priority...)
}

func Subscribe(sbr Subscriber) {
	Facade.Subscribe(sbr)
}

func AddSubscriber(sbr Subscriber) {
	Facade.AddSubscriber(sbr)
}

func AsyncFire(e Event) {
	Facade.AsyncFire(e)
}

func Trigger(name string, params M) (error, Event) {
	return Facade.Fire(name, params)
}

func Fire(name string, params M) (error, Event) {
	return Facade.Fire(name, params)
}

func FireEvent(e Event) error {
	return Facade.FireEvent(e)
}

func TriggerEvent(e Event) error {
	return Facade.FireEvent(e)
}

func MustFire(name string, params M) Event {
	return Facade.MustFire(name, params)
}

func MustTrigger(name string, params M) Event {
	return Facade.MustFire(name, params)
}

func FireBatch(es ...any) []error {
	return Facade.FireBatch(es...)
}

func HasListeners(name string) bool {
	return Facade.HasListeners(name)
}

func Reset() {
	Facade.Clear()
}

func AddEvent(e Event) {
	Facade.AddEvent(e)
}

func GetEvent(name string) (Event, bool) {
	return Facade.GetEvent(name)
}

func HasEvent(name string) bool {
	return Facade.HasEvent(name)
}
