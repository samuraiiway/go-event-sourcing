package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/samuraiiway/go-event-sourcing/processor"
)

func postDomainEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domain := vars["domain"]
	event, err := parseBodyToMap(w, r)
	if err != nil {
		return
	}

	processor.ParseEventInternalProperties(domain, event)
	processor.SaveEvent(event)
	processor.SendEventToProjection(domain, event)
	processor.SendEventToAggregation(domain, event)

	response, _ := json.Marshal(event)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
