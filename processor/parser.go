package processor

import (
	"time"

	"github.com/samuraiiway/go-event-sourcing/util"
)

const (
	EVENT_ID         = "_event_id"
	DOMAIN_NAME      = "_domain_name"
	DOMAIN_ID        = "_domain_id"
	TIMESTAMP        = "_timestamp"
	OLD_VALUE        = "_old"
	NEW_VALUE        = "_new"
	AGGREGATION_NAME = "_aggregation_name"
	AGGREGATION_KEY  = "_aggregation_key"
)

func ParseEventInternalProperties(domain string, event map[string]interface{}) {
	event[EVENT_ID] = util.GetUUID()
	event[DOMAIN_NAME] = domain
	event[TIMESTAMP] = time.Now().UnixNano()

	_, ok := event[DOMAIN_ID]
	if !ok {
		event[DOMAIN_ID] = util.GetUUID()
	}
}
