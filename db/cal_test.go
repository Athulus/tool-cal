package db

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
var _ struct{} = caltestInit()

func caltestInit() (x struct{}) {
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
	s.Lpush("cals", "test")
	// adding event keys to a 'laser' calendar
	s.Lpush("test", startTimeSeed.Format(time.RFC3339))
	//adding events
	s.HSet(startTimeSeed.Format(time.RFC3339), "startTime", startTimeSeed.Format(time.RFC3339))
	s.HSet(startTimeSeed.Format(time.RFC3339), "endTime", endTimeSeed.Format(time.RFC3339))
	s.HSet(startTimeSeed.Format(time.RFC3339), "description", "testDescription")
	s.HSet(startTimeSeed.Format(time.RFC3339), "owner", "testOwner")

	return
}
func Test_calendar_fetchEvents(t *testing.T) {
	tests := []struct {
		name string
		cal  calendar
		want []Event
	}{
		{"test1", "test", []Event{{startTimeSeed.UTC(), endTimeSeed.UTC(), "testDescription", "testOwner"}}},
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
		{"duplicate event test", "test", args{Event{startTimeSeed, endTimeSeed, "duplicate event", "testUser"}}, true},
		{"overlapping event test 1", "test", args{Event{startTimeSeed.Add(30 * time.Minute), endTimeSeed.Add(30 * time.Minute), "overlap event", "testUser"}}, true},
		{"overlapping event test 2", "test", args{Event{startTimeSeed.Add(-30 * time.Minute), endTimeSeed.Add(-30 * time.Minute), "overalp event", "testUser"}}, true},
		{"good event test", "test", args{Event{startTimeSeed.Add(time.Hour), endTimeSeed.Add(time.Hour), "new event", "testUser"}}, false},
		{"backwards event test", "test", args{Event{endTimeSeed.Add(-1 * time.Hour), startTimeSeed.Add(-1 * time.Hour), "backwards event", "testUser"}}, true},
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
		{"regular event", fields{startTimeSeed, endTimeSeed, "", ""}, (59 * time.Minute)},
		{"zero time event", fields{startTimeSeed, startTimeSeed, "", ""}, 0},
		{"backwards event", fields{endTimeSeed, startTimeSeed, "", ""}, (-59 * time.Minute)},
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

func TestEvent_NewEvent(t *testing.T) {
	type args struct {
		eventMap map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		{"event test 1",
			args{map[string]string{
				"startTime":   startTimeSeed.Format(time.RFC3339),
				"endTime":     endTimeSeed.Format(time.RFC3339),
				"description": "testDescription",
				"owner":       "testOwner",
			},
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEvent(tt.args.eventMap)
			if e.StartTime.Format(time.RFC3339) != tt.args.eventMap["startTime"] {
				t.Errorf("NewEvent() startTime is not equal to map")
			}
			if e.EndTime.Format(time.RFC3339) != tt.args.eventMap["endTime"] {
				t.Errorf("NewEvent() endTime is not equal to map")
			}
			if e.Description != tt.args.eventMap["description"] {
				t.Errorf("NewEvent() description is not equal to map")
			}
			if e.Owner != tt.args.eventMap["owner"] {
				t.Errorf("NewEvent() owner is not equal to map")
			}

		})
	}
}
