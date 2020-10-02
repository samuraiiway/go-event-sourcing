package processor

import "github.com/samuraiiway/go-event-sourcing/repository"

func SaveEvent(event map[string]interface{}) {
	domain := event[DOMAIN_NAME].(string)
	key := event[EVENT_ID].(string)

	repository.SaveEvent(domain, key, event)
}
