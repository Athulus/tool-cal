package db

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	bolt "go.etcd.io/bbolt"
)

var startTimeSeed time.Time
var endTimeSeed time.Time

func setup() {
	var err error
	err = Init("test.db")
	if err != nil {
		panic(err)
	}
	startTimeSeed, err = time.Parse(time.RFC3339, "2019-05-21T01:00:00-00:00")
	endTimeSeed, err = time.Parse(time.RFC3339, "2019-05-21T01:59:00-00:00")

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("calTest"))
		value, err := json.Marshal(Event{0, startTimeSeed.UTC(), endTimeSeed.UTC(), "testDescription", "testOwner"})
		if err != nil {
			return err
		}
		err = b.Put(itob(5), value)
		return err
	})
	if err != nil {
		panic(err)
	}

	return
}

func teardown() {
	err := db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte("calTest"))
	})
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func Test_calendar_fetchEvents(t *testing.T) {
	tests := []struct {
		name string
		cal  Calendar
		want []Event
	}{
		{"test1", "test", []Event{{1, startTimeSeed.UTC(), endTimeSeed.UTC(), "testDescription", "testOwner"}}},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if got, err := tt.cal.FetchEvents(); !reflect.DeepEqual(got, tt.want) && err != nil {
				t.Errorf("calendar.FetchEvents() = %v, want %v", got, tt.want)
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
		cal     Calendar
		args    args
		wantErr bool
	}{
		{"duplicate event test", "test", args{Event{0, startTimeSeed, endTimeSeed, "duplicate event", "testUser"}}, true},
		{"overlapping event test 1", "test", args{Event{1, startTimeSeed.Add(30 * time.Minute), endTimeSeed.Add(30 * time.Minute), "overlap event", "testUser"}}, true},
		{"overlapping event test 2", "test", args{Event{2, startTimeSeed.Add(-30 * time.Minute), endTimeSeed.Add(-30 * time.Minute), "overalp event", "testUser"}}, true},
		{"good event test", "test", args{Event{3, startTimeSeed.Add(time.Hour), endTimeSeed.Add(time.Hour), "new event", "testUser"}}, false},
		{"backwards event test", "test", args{Event{4, endTimeSeed.Add(-1 * time.Hour), startTimeSeed.Add(-1 * time.Hour), "backwards event", "testUser"}}, true},
	}
	setup()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cal.AddEvent(tt.args.e); (err != nil) != tt.wantErr {
				t.Errorf("calendar.AddEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	teardown()
}

func Test_calendar_eventFits(t *testing.T) {
	type args struct {
		event Event
	}
	tests := []struct {
		name string
		cal  Calendar
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
