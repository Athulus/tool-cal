package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Athulus/tool-cal/db"
	"github.com/gorilla/mux"
)

var tools = make(map[string]db.Calendar)

func main() {
	log.Println(time.Now().Format(time.RFC3339))

	db.Init("cal.db")

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
	var event db.Event
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
	err = db.Calendar(tool).AddEvent(event)
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
	events, err := db.Calendar(tool).FetchEvents()
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("could not fetch the events: " + err.Error()))
		return
	}
	log.Println(events)
	e, err := json.Marshal(events)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not fetch events for the " + tool))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(e)
}
