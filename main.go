package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	"github.com/samuraiiway/go-event-sourcing/handler"
)

func main() {
	router := mux.NewRouter()
	handler.RegisterPaths(router)
	go monitoring()
	http.ListenAndServe(":8080", router)
}

func monitoring() {
	for {
		time.Sleep(4 * time.Second)
		fmt.Printf("========== Monitoring : %v ==========\n", time.Now())
		fmt.Printf("Thread : %v\n", runtime.NumGoroutine())
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		fmt.Printf("Alloc = %v MiB", bToMb(mem.Alloc))
		fmt.Printf("\tTotalAlloc = %v MiB", bToMb(mem.TotalAlloc))
		fmt.Printf("\tSys = %v MiB", bToMb(mem.Sys))
		fmt.Printf("\tNumGC = %v\n", mem.NumGC)
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
