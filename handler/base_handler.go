package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	POST_DOMAIN_EVENT_PATH           = "/api/event/{domain}"
	PROJECTION_STREAM_LISTENER_PATH  = "/stream/projection/{domain}/{group}"
	CHANGED_STREAM_LISTENER_PATH     = "/stream/changed/{domain}/{group}"
	AGGREGATION_STREAM_LISTENER_PATH = "/stream/aggregation/{domain}/{group}"
)

func RegisterPaths(router *mux.Router) {
	router.HandleFunc(POST_DOMAIN_EVENT_PATH, postDomainEventHandler).Methods("POST")
	router.HandleFunc(PROJECTION_STREAM_LISTENER_PATH, getProjectionStreamListener).Methods("GET")
	router.HandleFunc(CHANGED_STREAM_LISTENER_PATH, getChangedStreamListener).Methods("GET")
	router.HandleFunc(AGGREGATION_STREAM_LISTENER_PATH, getAggregationStreamListener).Methods("GET")
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
