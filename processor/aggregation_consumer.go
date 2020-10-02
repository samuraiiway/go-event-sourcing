package processor

import (
	"github.com/samuraiiway/go-event-sourcing/util"
)

const (
	AGGREGATION_CONSUMER_BUFFERED_SIZE = 100
)

var (
	// <domain, <group, List<Channel<Aggregation>>>>
	rootAggregationConsumer = map[string]map[string][]chan map[string]interface{}{}
)

func getAggregationConsumer(domain string) map[string][]chan map[string]interface{} {
	consumer, ok := rootAggregationConsumer[domain]

	if !ok {
		consumer = map[string][]chan map[string]interface{}{}
		rootAggregationConsumer[domain] = consumer
	}

	return consumer
}

func RegisterAggregationConsumer(domain string, groupName string) chan map[string]interface{} {
	consumer := getAggregationConsumer(domain)
	group, ok := consumer[groupName]

	if !ok {
		group = []chan map[string]interface{}{}
	}

	channel := make(chan map[string]interface{}, AGGREGATION_CONSUMER_BUFFERED_SIZE)
	group = append(group, channel)
	consumer[groupName] = group

	return channel
}

func DeregisterAggregationConsumer(domain string, groupName string, channel chan map[string]interface{}) {
	consumer := getAggregationConsumer(domain)
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

func ProduceAggregation(domain string, aggregation map[string]interface{}) {
	consumer := getAggregationConsumer(domain)

	for _, channels := range consumer {
		channel := channels[util.RandomInt(len(channels))]
		channel <- aggregation
	}
}
