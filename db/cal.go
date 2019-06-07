package db

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/kr/pretty"
	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB
var eventBucket = []byte("events")

//Init opens or creates the database
func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	return err
}

// FetchCalendars shows a list of the current calendars
func FetchCalendars() ([]Calendar, error) {
	var calendars []Calendar
	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			cal := Calendar(k)
			calendars = append(calendars, cal)
		}
		return nil
	})
	return calendars, err
}

//AddCalendar creates a new bucket for a tool if it doesn't exist
func AddCalendar(cal Calendar) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(cal))
		return err
	})
}

//DeleteCalendar removes the bucket for a tool if it exists
func DeleteCalendar(cal Calendar) error {
	return db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(cal))
	})
}

// A Calendar is how we will split events by tools.
// there could be a laser calendar or a 3d printer calendar
type Calendar string

// FetchEvents returns all events for a calendar in a slice
func (cal Calendar) FetchEvents() ([]Event, error) {
	var events []Event
	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(cal)).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var event Event
			err := json.Unmarshal(v, &event)
			if err != nil {
				return err
			}
			events = append(events, event)
		}
		return nil
	})
	return events, err
}

// AddEvent adds an event to a calendar
func (cal Calendar) AddEvent(e Event) error {
	var err error
	if e.StartTime.After(e.EndTime) {
		return errors.New("the event must start before it ends")
	}
	if cal.eventFits(e) {
		var id int
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(cal))
			id64, _ := b.NextSequence()
			id = int(id64)
			e.ID = id
			key := itob(id)
			j, err := json.Marshal(e)
			if err != nil {
				return err
			}
			return b.Put(key, j)
		})
	} else {
		return errors.New("this event conflicts with another event in the calendar")
	}
	return err
}

//DeleteEvent removes an evnet from the calendar bucket
func (cal Calendar) DeleteEvent(id int) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(cal))
		err := b.Delete(itob(id))
		return err
	})
	return err
}

//FetchEvent gets on event form a calendar by id
func (cal Calendar) FetchEvent(id int) (Event, error) {
	var e Event
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(cal))
		value := b.Get(itob(id))
		err := json.Unmarshal(value, &e)
		return err
	})
	return e, err
}

func (cal Calendar) eventFits(event Event) bool {

	events, err := cal.FetchEvents()
	if err != nil {
		return false
	}
	for _, e := range events {
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
	ID          int       `json:"id,omitempty"`
	StartTime   time.Time `json:"start_time,omitempty"`
	EndTime     time.Time `json:"end_time,omitempty"`
	Description string    `json:"description,omitempty"`
	Owner       string    `json:"owner,omitempty"`
}

func (e Event) getDuration() time.Duration {
	return e.EndTime.Sub(e.StartTime)
}

//NewEvent creates a new Event from the redis map
//TODO is this still needed?
func NewEvent(eventMap map[string]string) Event {
	var err error
	var ok bool
	var e Event
	e.StartTime, err = time.Parse(time.RFC3339, eventMap["startTime"])
	if err != nil {
		log.Fatalln(err.Error())
	}
	e.EndTime, err = time.Parse(time.RFC3339, eventMap["endTime"])
	if err != nil {
		log.Fatalln(err.Error())
	}
	e.Description, ok = eventMap["description"]
	if !ok {
		pretty.Println(eventMap)
		log.Fatalln("problem creating event: description")
	}
	e.Owner, ok = eventMap["owner"]
	if !ok {
		log.Fatalln("problem creating event: owner")
	}
	return e
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
