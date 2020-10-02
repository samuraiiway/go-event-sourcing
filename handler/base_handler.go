package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/samuraiiway/go-event-sourcing/processor"
)

const (
	POST_DOMAIN_EVENT_PATH           = "/api/event/{domain}"
	PROJECTION_STREAM_LISTENER_PATH  = "/stream/projection/{domain}/{group}"
	CHANGED_STREAM_LISTENER_PATH     = "/stream/changed/{domain}/{group}"
	AGGREGATION_STREAM_LISTENER_PATH = "/stream/aggregation/{domain}/{group}"
	LOAD_TEST_PATH                   = "/test/{domain}/{number}"
)

func RegisterPaths(router *mux.Router) {
	router.HandleFunc(POST_DOMAIN_EVENT_PATH, postDomainEventHandler).Methods("POST")
	router.HandleFunc(PROJECTION_STREAM_LISTENER_PATH, getProjectionStreamListener).Methods("GET")
	router.HandleFunc(CHANGED_STREAM_LISTENER_PATH, getChangedStreamListener).Methods("GET")
	router.HandleFunc(AGGREGATION_STREAM_LISTENER_PATH, getAggregationStreamListener).Methods("GET")
	router.HandleFunc(LOAD_TEST_PATH, loadTestHandler).Methods("POST")
}

func parseBodyToMap(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, "bad_request", http.StatusBadRequest)
		return nil, readErr
	}

	result := map[string]interface{}{}
	jsonErr := json.Unmarshal(body, &result)
	if jsonErr != nil {
		http.Error(w, "bad_request", http.StatusBadRequest)
		return nil, jsonErr
	}

	return result, nil
}

func loadTestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domain := vars["domain"]
	number, _ := strconv.Atoi(vars["number"])

	for i := 0; i < number; i++ {
		event := map[string]interface{}{}
		event["service_id"] = "transfer"
		event["payer_id"] = "12345"
		event["payee_id"] = "98765"
		event["amount"] = float64(i)
		event["status"] = "success"

		processor.ParseEventInternalProperties(domain, event)
		processor.SaveEvent(event)
		processor.SendEventToProjection(domain, event)
		processor.SendEventToAggregation(domain, event)
	}

	w.Header().Set("Content-Type", "application/json")
}
