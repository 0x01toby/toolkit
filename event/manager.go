package event

import (
	"strings"
	"sync"
)

type Manager struct {
	sync.Mutex
	EnableLock    bool
	name          string
	sample        *BasicEvent
	events        map[string]Event
	listeners     map[string]*ListenerQueue
	listenedNames map[string]int
}

func NewManager(name string) *Manager {
	em := &Manager{
		name:          name,
		events:        make(map[string]Event),
		listeners:     make(map[string]*ListenerQueue),
		listenedNames: make(map[string]int),
	}
	return em
}

func (m *Manager) AddEvent(e Event) {
	name := goodName(e.Name())
	m.events[name] = e
}

func (m *Manager) On(name string, listener Listener, priority ...int) {
	pv := Normal
	if len(priority) > 0 {
		pv = priority[0]
	}
	m.addListenerItem(name, &ListenerItem{pv, listener})
}

func (m *Manager) Fire(name string, params M) (err error, e Event) {
	name = goodName(name)

	if m.HasListeners(name) == false && m.HasListeners(Wildcard) == false {
		pos := strings.LastIndexByte(name, '.')
		if pos < 0 || pos == len(name)-1 {
			return
		}

		groupName := name[:pos+1] + Wildcard
		if m.HasListeners(groupName) == false {
			return
		}
	}

	if e, ok := m.events[name]; ok {
		if params != nil {
			e.SetData(params)
		}
		err = m.FireEvent(e)
		return
	}
	e = m.newBasicEvent(name, params)
	err = m.FireEvent(e)
	return
}

func (m *Manager) GetEvent(name string) (e Event, ok bool) {
	e, ok = m.events[name]
	return
}

func (m *Manager) HasEvent(name string) bool {
	_, ok := m.events[name]
	return ok
}

func (m *Manager) RemoveEvent(name string) {
	if _, ok := m.events[name]; ok {
		delete(m.events, name)
	}
}

func (m *Manager) RemoveEvents() {
	m.events = map[string]Event{}
}

func (m *Manager) AddListener(name string, listener Listener, priority ...int) {
	m.On(name, listener, priority...)
}

func (m *Manager) Listen(name string, listener Listener, priority ...int) {
	m.On(name, listener, priority...)
}

func (m *Manager) HasListener(name string) bool {
	_, ok := m.listeners[name]
	return ok
}

func (m *Manager) Listeners() map[string]*ListenerQueue {
	return m.listeners
}

func (m *Manager) ListenersByName(name string) *ListenerQueue {
	return m.listeners[name]
}

func (m *Manager) ListenersCount(name string) int {
	if lq, ok := m.listeners[name]; ok {
		return lq.Len()
	}
	return 0
}

func (m *Manager) ListenedNames() map[string]int {
	return m.listenedNames
}

func (m *Manager) RemoveListener(name string, listener Listener) {
	if name != "" {
		if lq, ok := m.listeners[name]; ok {
			lq.Remove(listener)
			if lq.IsEmpty() {
				delete(m.listeners, name)
				delete(m.listenedNames, name)
			}
			return
		}
	}

	for name, lq := range m.listeners {
		lq.Remove(listener)
		if lq.IsEmpty() {
			delete(m.listeners, name)
			delete(m.listenedNames, name)
		}
	}
}

func (m *Manager) RemoveListeners(name string) {
	_, ok := m.listenedNames[name]
	if !ok {
		return
	}
	m.listeners[name].Clear()
	delete(m.listeners, name)
	delete(m.listenedNames, name)
}

func (m *Manager) Clear() {
	m.Reset()
}

func (m *Manager) Reset() {
	for _, lq := range m.listeners {
		lq.Clear()
	}
	m.name = ""
	m.events = make(map[string]Event)
	m.listeners = make(map[string]*ListenerQueue)
	m.listenedNames = make(map[string]int)
}

func (m *Manager) Subscribe(sbr Subscriber) {
	m.AddSubscriber(sbr)
}

func (m *Manager) AddSubscriber(sbr Subscriber) {
	for name, listener := range sbr.SubscribedEvents() {
		switch lt := listener.(type) {
		case Listener:
			m.On(name, lt)
		case ListenerItem:
			m.addListenerItem(name, &lt)
		default:
			panic("event: the value must be an listener or ListenerItem instance")
		}
	}
}

func (m *Manager) addListenerItem(name string, li *ListenerItem) {
	if name != Wildcard {
		name = goodName(name)
	}

	if li.Listener == nil {
		panic("event: the event '" + name + "' listener can not be empty")
	}

	if lq, ok := m.listeners[name]; ok {
		lq.Push(li)
	} else {
		m.listenedNames[name] = 1
		m.listeners[name] = (&ListenerQueue{}).Push(li)
	}
}

func (m *Manager) MustTrigger(name string, params M) Event {
	return m.MustFire(name, params)
}

func (m *Manager) Trigger(name string, params M) (error, Event) {
	return m.Fire(name, params)
}

func (m *Manager) AsyncFire(e Event) {
	go func(e Event) {
		_ = m.FireEvent(e)
	}(e)
}

func (m *Manager) AwaitFire(e Event) (err error) {
	ch := make(chan error)

	go func(e Event) {
		err = m.FireEvent(e)
		ch <- err
	}(e)

	err = <-ch
	close(ch)
	return
}

func (m *Manager) FireBatch(es ...any) (ers []error) {
	var err error
	for _, e := range es {
		if name, ok := e.(string); ok {
			err, _ = m.Fire(name, nil)
		} else if evt, ok := e.(Event); ok {
			err = m.FireEvent(evt)
		}
		if err != nil {
			ers = append(ers, err)
		}
	}
	return
}

func (m *Manager) MustFire(name string, params M) Event {
	err, e := m.Fire(name, params)
	if err != nil {
		panic(err)
	}
	return e
}

func (m *Manager) FireEvent(e Event) (err error) {
	if m.EnableLock {
		m.Lock()
		defer m.Unlock()
	}
	e.Abort(false)
	name := e.Name()
	lq, ok := m.listeners[name]
	if ok {
		for _, li := range lq.Sort().Items() {
			err = li.Listener.Handle(e)
			if err != nil || e.IsAborted() {
				return
			}
		}
	}
	pos := strings.LastIndexByte(name, '.')
	if pos > 0 && pos < len(name) {
		groupName := name[:pos+1] + Wildcard
		if lq, ok = m.listeners[groupName]; ok {
			for _, li := range lq.Sort().Items() {
				if li.Listener.Handle(e) != nil || e.IsAborted() {
					return
				}
			}
		}
	}
	if lq, ok = m.listeners[Wildcard]; ok {
		for _, li := range lq.Sort().Items() {
			if li.Listener.Handle(e) != nil || e.IsAborted() {
				break
			}
		}
	}
	return
}

func (m *Manager) HasListeners(name string) bool {
	_, ok := m.listenedNames[name]
	return ok
}

func (m *Manager) newBasicEvent(name string, data M) *BasicEvent {
	var cp = *m.sample
	cp.SetName(name)
	cp.SetData(data)
	return &cp
}
