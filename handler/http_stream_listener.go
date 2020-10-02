package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/samuraiiway/go-event-sourcing/processor"
)

func getProjectionStreamListener(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domain := vars["domain"]
	group := vars["group"]

	flusher, flusherOk := w.(http.Flusher)

	if !flusherOk {
		fmt.Println("Unsupported steaming")
		http.Error(w, "Unsupported steaming", http.StatusInternalServerError)
		return
	}

	endSignal := r.Context().Done()
	ch := processor.RegisterProjectionConsumer(domain, group)

	go func(ch chan map[string]interface{}) {
		<-endSignal
		processor.DeregisterProjectionConsumer(domain, group, ch)
	}(ch)

	w.Header().Set("Content-Type", "text/event-stream")

	for {
		data, ok := <-ch
		if !ok {
			return
		}
		response, _ := json.Marshal(data)
		fmt.Fprintf(w, "%s\n\n", response)
		flusher.Flush()
	}
}

func getChangedStreamListener(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domain := vars["domain"]
	group := vars["group"]

	flusher, flusherOk := w.(http.Flusher)

	if !flusherOk {
		fmt.Println("Unsupported steaming")
		http.Error(w, "Unsupported steaming", http.StatusInternalServerError)
		return
	}

	endSignal := r.Context().Done()
	ch := processor.RegisterChangedConsumer(domain, group)

	go func(ch chan map[string]interface{}) {
		<-endSignal
		processor.DeregisterChangedConsumer(domain, group, ch)
	}(ch)

	w.Header().Set("Content-Type", "text/event-stream")

	for {
		data, ok := <-ch
		if !ok {
			return
		}
		response, _ := json.Marshal(data)
		fmt.Fprintf(w, "%s\n\n", response)
		flusher.Flush()
	}
}

func getAggregationStreamListener(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domain := vars["domain"]
	group := vars["group"]

	flusher, flusherOk := w.(http.Flusher)

	if !flusherOk {
		fmt.Println("Unsupported steaming")
		http.Error(w, "Unsupported steaming", http.StatusInternalServerError)
		return
	}

	endSignal := r.Context().Done()
	ch := processor.RegisterAggregationConsumer(domain, group)

	go func(ch chan map[string]interface{}) {
		<-endSignal
		processor.DeregisterAggregationConsumer(domain, group, ch)
	}(ch)

	w.Header().Set("Content-Type", "text/event-stream")

	for {
		data, ok := <-ch
		if !ok {
			return
		}
		response, _ := json.Marshal(data)
		fmt.Fprintf(w, "%s\n\n", response)
		flusher.Flush()
	}
}
