package processor

import (
	"github.com/samuraiiway/go-event-sourcing/util"
)

const (
	PROJECTION_CONSUMER_BUFFERED_SIZE = 100
)

var (
	// <domain, <group, List<Channel<Projection>>>>
	rootProjectionConsumer = map[string]map[string][]chan map[string]interface{}{}
)

func getProjectionConsumer(domain string) map[string][]chan map[string]interface{} {
	consumer, ok := rootProjectionConsumer[domain]

	if !ok {
		consumer = map[string][]chan map[string]interface{}{}
		rootProjectionConsumer[domain] = consumer
	}

	return consumer
}

func RegisterProjectionConsumer(domain string, groupName string) chan map[string]interface{} {
	consumer := getProjectionConsumer(domain)
	group, ok := consumer[groupName]

	if !ok {
		group = []chan map[string]interface{}{}
	}

	channel := make(chan map[string]interface{}, PROJECTION_CONSUMER_BUFFERED_SIZE)
	group = append(group, channel)
	consumer[groupName] = group

	return channel
}

func DeregisterProjectionConsumer(domain string, groupName string, channel chan map[string]interface{}) {
	consumer := getProjectionConsumer(domain)
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

func ProduceProjection(domain string, projection map[string]interface{}) {
	consumer := getProjectionConsumer(domain)

	for _, channels := range consumer {
		channel := channels[util.RandomInt(len(channels))]
		channel <- projection
	}
}
