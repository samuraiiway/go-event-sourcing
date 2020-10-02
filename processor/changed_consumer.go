package processor

import (
	"github.com/samuraiiway/go-event-sourcing/util"
)

const (
	CHANGED_CONSUMER_BUFFERED_SIZE = 100
)

var (
	// <domain, <group, List<Channel<Changed>>>>
	rootChangedConsumer = map[string]map[string][]chan map[string]interface{}{}
)

func getChangedConsumer(domain string) map[string][]chan map[string]interface{} {
	consumer, ok := rootChangedConsumer[domain]

	if !ok {
		consumer = map[string][]chan map[string]interface{}{}
		rootChangedConsumer[domain] = consumer
	}

	return consumer
}

func RegisterChangedConsumer(domain string, groupName string) chan map[string]interface{} {
	consumer := getChangedConsumer(domain)
	group, ok := consumer[groupName]

	if !ok {
		group = []chan map[string]interface{}{}
	}

	channel := make(chan map[string]interface{}, CHANGED_CONSUMER_BUFFERED_SIZE)
	group = append(group, channel)
	consumer[groupName] = group

	return channel
}

func DeregisterChangedConsumer(domain string, groupName string, channel chan map[string]interface{}) {
	consumer := getChangedConsumer(domain)
	group, ok := consumer[groupName]

	idx := 0

	if ok {
		for i, client := range group {
			if client == channel {
				idx = i
				break
			}
		}

		group = append(group[:idx], group[idx+1:]...)

		if len(group) == 0 {
			delete(consumer, groupName)
		} else {
			consumer[groupName] = group
		}
	}

	close(channel)
}

func ProduceChanged(domain string, changed map[string]interface{}) {
	consumer := getChangedConsumer(domain)

	for _, channels := range consumer {
		channel := channels[util.RandomInt(len(channels))]
		channel <- changed
	}
}
