package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var tools = make(map[string]calendar)

func main() {
	log.Println(time.Now().Format(time.RFC3339))

	router := mux.NewRouter()
	router.HandleFunc("/health", health)
	router.HandleFunc("/calendar/{tool}/events", getEvents).Methods("GET")
	router.HandleFunc("/calendar/{tool}/events", addEvent).Methods("POST")

	// fs := http.FileServer(http.Dir("static"))
	// router.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":8080", router))
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("health check, ok"))
}

func addEvent(w http.ResponseWriter, r *http.Request) {
	tool := mux.Vars(r)["tool"]
	var event Event
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = json.Unmarshal(body, &event)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(event)
	err = calendar(tool).addEvent(event)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("could not add the event: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)

}

func getEvents(w http.ResponseWriter, r *http.Request) {
	tool := mux.Vars(r)["tool"]
	log.Println(tool)
	log.Println(calendar(tool).fetchEvents())
	events, err := json.Marshal(calendar(tool).fetchEvents())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not fetch events for the " + tool))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(events)
}
