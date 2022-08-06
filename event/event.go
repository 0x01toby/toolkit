package event

type BasicEvent struct {
	name    string
	data    map[string]any
	target  any
	aborted bool
}

func NewBasic(name string, data M) *BasicEvent {
	if data == nil {
		data = make(M)
	}
	return &BasicEvent{
		name: name,
		data: data,
	}
}

func (b *BasicEvent) Name() string {
	return b.name
}

func (b *BasicEvent) Get(key string) any {
	if v, ok := b.data[key]; ok {
		return v
	}
	return nil
}

func (b *BasicEvent) Set(key string, val any) {
	if b.data == nil {
		b.data = make(M)
	}
	b.data[key] = val
}

func (b *BasicEvent) Add(key string, val any) {
	if _, ok := b.data[key]; !ok {
		b.Set(key, val)
	}
}

func (b *BasicEvent) Data() M {
	return b.data
}

func (b *BasicEvent) Abort(abort bool) {
	b.aborted = abort
}

func (b *BasicEvent) IsAborted() bool {
	return b.aborted
}

func (b *BasicEvent) Fill(target any, data M) *BasicEvent {
	if data != nil {
		b.data = data
	}
	b.target = target
	return b
}

func (b *BasicEvent) Target() any {
	return b.target
}

func (b *BasicEvent) SetName(name string) *BasicEvent {
	b.name = name
	return b
}

func (b *BasicEvent) SetData(data M) Event {
	if data != nil {
		b.data = data
	}
	return b
}

func (b *BasicEvent) SetTarget(target any) *BasicEvent {
	b.target = target
	return b
}
