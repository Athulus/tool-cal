package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/alicebob/miniredis"
)

var _ struct{} = maintestInit()

func maintestInit() (x struct{}) {
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
func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_health(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"health test",
			args{
				httptest.NewRecorder(),
				httptest.NewRequest(http.MethodGet, "http://test/health", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			health(tt.args.w, tt.args.r)
			resp := tt.args.w.Result()
			body, _ := ioutil.ReadAll(resp.Body)

			if string(body) != "health check, ok" {
				t.Errorf("expected response: 'health check, ok'. got: '%v'", string(body))
			}
		})
	}
}

func Test_addEvent(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"add event test",
			args{
				httptest.NewRecorder(),
				httptest.NewRequest(http.MethodPost, "http://test/calendar/test/events", strings.NewReader("[{\"start_time\":\"2018-03-21T01:00:00Z\",\"end_time\":\"2018-03-21T01:59:00Z\",\"description\":\"testDescription\",\"owner\":\"testOwner\"}]")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addEvent(tt.args.w, tt.args.r)
		})
	}
}

func Test_getEvents(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"get event test",
			args{
				httptest.NewRecorder(),
				httptest.NewRequest(http.MethodGet, "http://test/calendar/test/events", nil),
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.r = mux.SetURLVars(tt.args.r, map[string]string{"tool": "test"})
			getEvents(tt.args.w, tt.args.r)
			resp := tt.args.w.Result()
			body, _ := ioutil.ReadAll(resp.Body)
			if string(body) != "[{\"start_time\":\"2018-02-21T01:00:00Z\",\"end_time\":\"2018-02-21T01:59:00Z\",\"description\":\"testDescription\",\"owner\":\"testOwner\"}]" {
				t.Errorf("incorrect response body: %v", string(body))
			}

		})
	}
}
