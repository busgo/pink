package bus

import (
	"errors"
	"sync"
)

// event bus
type EventBus struct {
	subNodes map[string]*SubNode
	rw       sync.RWMutex
}

// new a event bus
func NewEventBus() *EventBus {

	return &EventBus{
		subNodes: make(map[string]*SubNode),
		rw:       sync.RWMutex{},
	}
}

func (b *EventBus) Subscribe(topic string) *Sub {
	b.rw.Lock()
	defer b.rw.Unlock()
	sub := NewSub()
	if node, ok := b.subNodes[topic]; ok {
		node.rw.Lock()
		node.subs = append(node.subs, sub)
		node.rw.Unlock()
		return sub
	}
	node := NewSubNode()
	node.subs = append(node.subs, sub)
	b.subNodes[topic] = node
	return sub
}

func (b *EventBus) Unsubscribe(topic string, sub *Sub) {

	b.rw.Lock()
	defer b.rw.Unlock()
	n := b.subNodes[topic]
	if n == nil {
		return
	}
	n.delete(sub)
}

func (b *EventBus) Publish(topic string, msg interface{}) error {
	b.rw.RLock()
	defer b.rw.RUnlock()

	if n, ok := b.subNodes[topic]; ok {

		go func(subs []*Sub, msg interface{}) {

			for _, sub := range subs {
				sub.receive(msg)
			}
		}(n.subs, msg)

		return nil
	}
	return errors.New("not found the topic")
}
