package main

import (
	"sync"

	"github.com/rs/xid"
)

type value = []byte

type broker struct {
	store     storer
	consumers map[string][]consumer
	sync.RWMutex
}

func newBroker(store storer) *broker {
	return &broker{
		store:     store,
		consumers: map[string][]consumer{},
	}
}

// Publish a message to a topic.
func (b *broker) Publish(topic string, val value) error {
	if err := b.store.Insert(topic, val); err != nil {
		return err
	}

	b.NotifyConsumer(topic, eventTypePublish)

	return nil
}

// Subscribe to a topic and return a consumer for the topic.
func (b *broker) Subscribe(topic string) *consumer {
	b.Lock()
	defer b.Unlock()

	cons := consumer{
		id:        xid.New().String(),
		topic:     topic,
		store:     b.store,
		eventChan: make(chan eventType),
		notifier:  b,
	}

	b.consumers[topic] = append(b.consumers[topic], cons)

	return &cons
}

// Shutdown the broker.
func (b *broker) Shutdown() error {
	return b.store.Close()
}

// NotifyConsumers notifies a waiting consumer of a topic that an event has
// occurred.
func (b *broker) NotifyConsumer(topic string, ev eventType) {
	b.RLock()
	defer b.RUnlock()

	for _, c := range b.consumers[topic] {
		select {
		case c.eventChan <- ev:
			return
		default: // If there is noone listening noop
		}
	}
}
