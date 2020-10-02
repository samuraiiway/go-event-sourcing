package processor

import (
	"github.com/samuraiiway/go-event-sourcing/repository"
)

const (
	PROJECTION_CHANNEL_SIZE = 100000
)

var (
	rootProjectionChannel = map[string]chan map[string]interface{}{}
)

func getProjectionChannel(domain string) chan map[string]interface{} {
	channel, ok := rootProjectionChannel[domain]
	if !ok {
		channel = make(chan map[string]interface{}, PROJECTION_CHANNEL_SIZE)
		rootProjectionChannel[domain] = channel
		go registerProjectionTask(channel)
	}

	return channel
}

func registerProjectionTask(channel chan map[string]interface{}) {
	for {
		event, ok := <-channel
		if !ok {
			return
		}
		DoProjection(event)
	}
}

func DoProjection(event map[string]interface{}) {
	domain := event[DOMAIN_NAME].(string)
	key := event[DOMAIN_ID].(string)
	projection := repository.GetProjection(domain, key)
	changed := map[string]interface{}{}

	for key, value := range event {
		if projection[key] != value {
			if key == TIMESTAMP || key == EVENT_ID {
				changed[key] = projection[key]
			} else {
				changed[key+OLD_VALUE] = projection[key]
				changed[key+NEW_VALUE] = value
			}

			projection[key] = value
		}
	}

	repository.SaveProjection(domain, key, projection)
	ProduceProjection(domain, projection)
	if len(changed) > 2 {
		ProduceChanged(domain, changed)
	}
}

func SendEventToProjection(domain string, event map[string]interface{}) {
	channel := getProjectionChannel(domain)
	channel <- event
}
