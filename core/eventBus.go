package core

import (
	"github.com/asaskevich/EventBus"
)

var GlobalEventBus EventBus.Bus

func init() {
	GlobalEventBus = EventBus.New()
}

func Subscribe(topic string, fn interface{}) error {
	return GlobalEventBus.Subscribe(topic, fn)
}

func SubscribeAsync(topic string, fn interface{}, transactional bool) error {
	return GlobalEventBus.SubscribeAsync(topic, fn, transactional)
}

func Publish(topic string, args ...interface{}) {
	GlobalEventBus.Publish(topic, args...)
}
