package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// Tests to make sure we get 200 when we have a valid request
func TestGetWeatherHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/weatherapi/v1/latitude/33.34913/longitude/-94.39194/units/imperial", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r := mux.NewRouter()

	api := r.PathPrefix("/weatherapi/v1").Subrouter()

	api.HandleFunc("/latitude/{lat}/longitude/{lon}/units/{units}", getWeather).Methods(http.MethodGet, http.MethodOptions)

	api.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// Tests to make sure we get a 404 when we have missing path params in request URL
func TestGetWeatherHandlerWithMissingPathParam(t *testing.T) {
	req, err := http.NewRequest("GET", "/weatherapi/v1/latitude/33.34913/longitude/-94.39194/units/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r := mux.NewRouter()

	api := r.PathPrefix("/weatherapi/v1").Subrouter()

	api.HandleFunc("/latitude/{lat}/longitude/{lon}/units/{units}", getWeather).Methods(http.MethodGet, http.MethodOptions)

	api.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

// Tests to make sure we get 500 when we have incorrect formats for path params in request URL
func TestGetWeatherHandlerWithInvalidPathParam(t *testing.T) {
	req, err := http.NewRequest("GET", "/weatherapi/v1/latitude/f/longitude/-94.39194/units/metric", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r := mux.NewRouter()

	api := r.PathPrefix("/weatherapi/v1").Subrouter()

	api.HandleFunc("/latitude/{lat}/longitude/{lon}/units/{units}", getWeather).Methods(http.MethodGet, http.MethodOptions)

	api.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	req, err = http.NewRequest("GET", "/weatherapi/v1/latitude/33.34913/longitude/x/units/metric", nil)
	if err != nil {
		t.Fatal(err)
	}

	api.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	req, err = http.NewRequest("GET", "/weatherapi/v1/latitude/33.34913/longitude/-94.39194/units/cats", nil)
	if err != nil {
		t.Fatal(err)
	}

	api.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

// Tests to make sure Unix Time formatter is returning the correct Date and Weekday
func TestFormatUnixTime(t *testing.T) {
	var unixTimestamp int64
	unixTimestamp = 712904400
	timezone := "America/Chicago"

	testFormatDate, testFormatWeekday, err := formatUnixTime(unixTimestamp, timezone)
	if err != nil {
		t.Errorf("Something when wrong: %v", err)
	}

	expectedDate := "08/04/1992"

	if testFormatDate != expectedDate {
		t.Errorf("formatter returned wrong date: got %v want %v", testFormatDate, expectedDate)
	}

	expectedWeekday := "Tuesday"

	if testFormatWeekday != expectedWeekday {
		t.Errorf("formatter returned wrong weekday: got %v want %v", testFormatWeekday, expectedWeekday)
	}
}
