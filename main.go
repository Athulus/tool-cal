package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

	router.HandleFunc("/calendar/{tool}", addCal).Methods("POST")
	router.HandleFunc("/calendar", getCals).Methods("GET")
	router.HandleFunc("/calendar/{tool}", deleteCal).Methods("DELETE")

	router.HandleFunc("/calendar/{tool}/events", getEvents).Methods("GET")
	router.HandleFunc("/calendar/{tool}/events/{id}", getEvent).Methods("GET")
	router.HandleFunc("/calendar/{tool}/events", addEvent).Methods("POST")
	router.HandleFunc("/calendar/{tool}/events/{id}", deleteEvent).Methods("DELETE")

	// fs := http.FileServer(http.Dir("static"))
	// router.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":8080", router))
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("health check, ok"))
}

func getCals(w http.ResponseWriter, r *http.Request) {
	cals, err := db.FetchCalendars()
	body, err := json.Marshal(cals)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func addCal(w http.ResponseWriter, r *http.Request) {
	tool := mux.Vars(r)["tool"]
	err := db.AddCalendar(db.Calendar(tool))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)

}

func deleteCal(w http.ResponseWriter, r *http.Request) {
	tool := mux.Vars(r)["tool"]
	err := db.DeleteCalendar(db.Calendar(tool))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
	body, err := json.Marshal(events)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not fetch events for the " + tool))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func getEvent(w http.ResponseWriter, r *http.Request) {
	tool := mux.Vars(r)["tool"]
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	event, err := db.Calendar(tool).FetchEvent(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	body, err := json.Marshal(event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not fetch this event for the " + tool))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	tool := mux.Vars(r)["tool"]
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	err = db.Calendar(tool).DeleteEvent(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

}
