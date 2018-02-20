package main

import (
	"errors"
	"log"
	"time"

	"github.com/mediocregopher/radix.v2/pool"
)

var db *pool.Pool

func init() {
	var err error

	//set up redis
	db, err = pool.New("tcp", "localhost:6379", 10)
	if err != nil {
		log.Panic(err)
	}

}

type calendar string

func (cal calendar) fetchEvents() []Event {
	var events []Event
	conn, err := db.Get()
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Put(conn)

	keys, err := db.Cmd("lrange", cal, 0, -1).Array()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(keys)

	for _, key := range keys {
		var event Event
		eventMap, err := db.Cmd("HGETALL", key).Map()
		if err != nil {
			log.Fatalln(err.Error())
		}
		event.New(eventMap)
		events = append(events, event)

	}
	return events
}

func (cal calendar) addEvent(e Event) error {
	if cal.eventFits(e) {
		conn, err := db.Get()
		if err != nil {
			return err
		}
		defer db.Put(conn)
		err = db.Cmd("HSET", e.StartTime.Format(time.RFC3339), "startTime", e.StartTime.Format(time.RFC3339),
			"endTime", e.EndTime.Format(time.RFC3339), "Description", e.Description, "Owner", e.Owner).Err
		if err != nil {
			return err
		}
		err = db.Cmd("LPUSH", cal, e.StartTime.Format(time.RFC3339)).Err
		if err != nil {
			return err
		}
	} else {
		return errors.New("this event conflicts with another event in the calendar")
	}
	return nil
}

func (cal calendar) eventFits(event Event) bool {

	for _, e := range cal.fetchEvents() {
		//this conditional is gross looking
		// basiccally if there is an overlapping event return false
		if (e.StartTime.Before(event.StartTime) && e.EndTime.After(event.StartTime)) ||
			(e.StartTime.Before(event.EndTime) && e.EndTime.After(event.EndTime)) ||
			e.StartTime.Equal(event.StartTime) || e.StartTime.Equal(event.EndTime) ||
			e.EndTime.Equal(event.StartTime) || e.EndTime.Equal(event.EndTime) {
			return false
		}

	}
	return true
}

//An Event is an allocation of time held by a user. Events go on calendars
type Event struct {
	StartTime   time.Time `json:"start_time,omitempty"`
	EndTime     time.Time `json:"end_time,omitempty"`
	Description string    `json:"description,omitempty"`
	Owner       string    `json:"owner,omitempty"`
}

func (e Event) getDuration() time.Duration {
	return e.EndTime.Sub(e.StartTime)
}

//New creates a new Event from the redis map
func (e *Event) New(eventMap map[string]string) {
	var err error
	var ok bool
	e.StartTime, err = time.Parse(time.RFC3339, eventMap["startTime"])
	if err != nil {
		log.Fatalln(err.Error())
	}
	e.EndTime, err = time.Parse(time.RFC3339, eventMap["endTime"])
	if err != nil {
		log.Fatalln(err.Error())
	}
	e.Description, ok = eventMap["Description"]
	if !ok {
		log.Fatalln("problem creating event: description")
	}
	e.Owner, ok = eventMap["Owner"]
	if !ok {
		log.Fatalln("problem creating event: owner")
	}
}
