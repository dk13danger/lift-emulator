package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dk13danger/lift-emulator/service"
)

// Run runs web server for processing external calls.
func Run(eventsCh chan<- *service.Event) error {
	http.HandleFunc("/external", callHandler(service.ExternalCall, eventsCh))
	http.HandleFunc("/internal", callHandler(service.InternalCall, eventsCh))

	log.Println("Web server is running on port: 9090")
	return http.ListenAndServe(":9090", nil)
}

func callHandler(callType int, eventsCh chan<- *service.Event) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		// simple validating incoming floor param
		floor, err := strconv.Atoi(r.Form.Get("floor"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Incorrect 'floor' param: %v", err)))
			return
		}
		if floor < 1 || floor > 10 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("floor must be in range: from 1 to 10")))
			return
		}

		eventsCh <- &service.Event{
			Type:  callType,
			Floor: floor,
		}
	}
}
