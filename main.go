package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kr/pretty"
)

var tools = make(map[string]*calendar)

func main() {
	pretty.Println(time.Now().Format(time.RFC3339))
	tools["printer"] = new(calendar)
	tools["laser"] = new(calendar)
	tools["cnc"] = new(calendar)

	router := mux.NewRouter()
	router.Use(toolValidation)
	router.HandleFunc("/health", health)
	router.HandleFunc("/calendar/{tool}/events", getEvents).Methods("GET")
	router.HandleFunc("/calendar/{tool}/events", addEvent).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("health check, ok"))
}

func auth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not implimented!"))
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
	err = tools[tool].addEvent(event)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("could not add the event: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(body)

}

func getEvents(w http.ResponseWriter, r *http.Request) {
	tool := mux.Vars(r)["tool"]
	log.Println(tool)
	log.Println(tools[tool].events)
	events, err := json.Marshal(tools[tool].events)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not fetch events for the " + tool))
	}
	w.Write(events)
}
