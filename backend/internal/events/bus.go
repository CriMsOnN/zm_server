package events

import "sync"

type Event struct {
	Name string
	Data any
}

type Handler func(Event)

type Bus interface {
	Clear()
	Publish(event Event)
	Subscribe(eventName string, handler Handler) (unsubscribe func())
}

type InMemoryBus struct {
	mu       sync.RWMutex
	handlers map[string]map[uint64]Handler
	nextID   uint64
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		handlers: make(map[string]map[uint64]Handler),
		nextID:   0,
	}
}

func (b *InMemoryBus) Subscribe(eventName string, handler Handler) func() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.nextID++
	id := b.nextID

	if _, ok := b.handlers[eventName]; !ok {
		b.handlers[eventName] = make(map[uint64]Handler)
	}

	b.handlers[eventName][id] = handler

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()

		if _, ok := b.handlers[eventName][id]; ok {
			return
		}

		delete(b.handlers[eventName], id)
		if len(b.handlers[eventName]) == 0 {
			delete(b.handlers, eventName)
		}
	}
}

func (b *InMemoryBus) Publish(event Event) {
	b.mu.RLock()
	eventHandlers, ok := b.handlers[event.Name]
	if !ok {
		b.mu.RUnlock()
		return
	}

	if len(eventHandlers) == 0 {
		b.mu.RUnlock()
		return
	}

	ch := make([]Handler, 0, len(eventHandlers))
	for _, h := range eventHandlers {
		ch = append(ch, h)
	}

	b.mu.RUnlock()

	for _, h := range ch {
		h(event)
	}
}

func (b *InMemoryBus) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers = make(map[string]map[uint64]Handler)
	b.nextID = 0
}
