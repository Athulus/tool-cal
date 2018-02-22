package main

import (
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
)

var s *miniredis.Miniredis
var startTimeSeed time.Time
var endTimeSeed time.Time

//this is some trickery so i can get my miniredis server setup before
// the init function in cal.go runs
var _ struct{} = testInit()

func testInit() (x struct{}) {
	var err error
	s, err = miniredis.Run()
	if err != nil {
		panic(err)
	}
	log.Println(s.Addr())
	os.Setenv("redisAddress", s.Addr())

	// set up the redis server for testing
	startTimeSeed, err = time.Parse(time.RFC3339, "2018-02-21T01:00:00-00:00")
	endTimeSeed, err = time.Parse(time.RFC3339, "2018-02-21T01:59:00-00:00")
	//toplevel collection of calendars
	s.Lpush("cals", "laser")
	// adding event keys to a 'laser' calendar
	s.Lpush("laser", startTimeSeed.Format(time.RFC3339))
	//adding events
	s.HSet(startTimeSeed.Format(time.RFC3339), "startTime", startTimeSeed.Format(time.RFC3339))
	s.HSet(startTimeSeed.Format(time.RFC3339), "endTime", endTimeSeed.Format(time.RFC3339))
	s.HSet(startTimeSeed.Format(time.RFC3339), "description", "testDescription")
	s.HSet(startTimeSeed.Format(time.RFC3339), "owner", "testOwner")

	return
}

// func tearDown() {
// 	s.Close()
// }
func Test_calendar_fetchEvents(t *testing.T) {
	tests := []struct {
		name string
		cal  calendar
		want []Event
	}{
		// TODO: Add test cases.
		{"test1", "laser", []Event{{startTimeSeed.UTC(), endTimeSeed.UTC(), "testDescription", "testOwner"}}},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cal.fetchEvents(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calendar.fetchEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calendar_addEvent(t *testing.T) {
	type args struct {
		e Event
	}
	tests := []struct {
		name    string
		cal     calendar
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cal.addEvent(tt.args.e); (err != nil) != tt.wantErr {
				t.Errorf("calendar.addEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_calendar_eventFits(t *testing.T) {
	type args struct {
		event Event
	}
	tests := []struct {
		name string
		cal  calendar
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cal.eventFits(tt.args.event); got != tt.want {
				t.Errorf("calendar.eventFits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvent_getDuration(t *testing.T) {
	type fields struct {
		StartTime   time.Time
		EndTime     time.Time
		Description string
		Owner       string
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Event{
				StartTime:   tt.fields.StartTime,
				EndTime:     tt.fields.EndTime,
				Description: tt.fields.Description,
				Owner:       tt.fields.Owner,
			}
			if got := e.getDuration(); got != tt.want {
				t.Errorf("Event.getDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvent_New(t *testing.T) {
	type fields struct {
		StartTime   time.Time
		EndTime     time.Time
		Description string
		Owner       string
	}
	type args struct {
		eventMap map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Event{
				StartTime:   tt.fields.StartTime,
				EndTime:     tt.fields.EndTime,
				Description: tt.fields.Description,
				Owner:       tt.fields.Owner,
			}
			e.New(tt.args.eventMap)
		})
	}
}
