package processor

import (
	"fmt"

	"github.com/samuraiiway/go-event-sourcing/repository"
)

const (
	AGGREGATION_CHANNEL_SIZE = 100000
)

var (
	rootAggregationChannel = map[string]chan map[string]interface{}{}
)

func getAggregationChannel(domain string) chan map[string]interface{} {
	channel, ok := rootAggregationChannel[domain]
	if !ok {
		channel = make(chan map[string]interface{}, AGGREGATION_CHANNEL_SIZE)
		rootAggregationChannel[domain] = channel
		go registerAggregationTask(channel)
	}

	return channel
}

func registerAggregationTask(channel chan map[string]interface{}) {
	for {
		event, ok := <-channel
		if !ok {
			return
		}
		DoAggregation(event)
	}
}

func DoAggregation(event map[string]interface{}) {
	domain := event[DOMAIN_NAME].(string)

	if domainTasks, ok := rootAggregationFunction[domain]; ok {
		for name, task := range domainTasks {
			key := getKey(task.GroupByKeyID, event)
			keyId := name + key

			aggregation := repository.GetAggregation(domain, keyId)

			aggregation[EVENT_ID] = event[EVENT_ID]
			aggregation[DOMAIN_NAME] = domain
			aggregation[AGGREGATION_NAME] = name
			aggregation[AGGREGATION_KEY] = key
			aggregation[TIMESTAMP] = event[TIMESTAMP]

			for _, funtion := range task.Functions {
				value := funtion.Function.Apply(aggregation[funtion.FieldName], event[funtion.PropertyName])
				aggregation[funtion.FieldName] = value
			}

			repository.SaveAggregation(domain, keyId, aggregation)
			ProduceAggregation(domain, aggregation)
		}
	}
}

func getKey(groupKey []string, event map[string]interface{}) string {
	keys := ""

	for _, key := range groupKey {
		keys += fmt.Sprintf(":%v", event[key])
	}

	return keys
}

func SendEventToAggregation(domain string, event map[string]interface{}) {
	channel := getAggregationChannel(domain)
	channel <- event
}
