package main

import (
	"errors"
	"time"
)

type calendar struct {
	events []Event
}

func (cal *calendar) addEvent(e Event) error {
	if cal.eventFits(e) {
		cal.events = append(cal.events, e)
	} else {
		return errors.New("this event conflicts with another event in the calendar")
	}
	return nil
}

func (cal *calendar) eventFits(event Event) bool {

	for _, e := range cal.events {
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
