package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

/* ForecastResponse structure
getWeathers JSON response structure
*/
type ForecastResponse struct {
	Dt  string `json:"dt"`
	Min string `json:"min"`
	Max string `json:"max"`
	Day string `json:"day"`
}

/* WeatherAPIResponse
Top level structure representing the response from
https://api.openweathermap.org/data/2.5/onecall
*/
type WeatherAPIResponse struct {
	Lat            float64                   `json:"lat"`
	Lon            float64                   `json:"lon"`
	Timezone       string                    `json:"timezone"`
	TimezoneOffset int                       `json:"timezone_offset"`
	Daily          []DailyWeatherAPIResponse `json:"daily"`
}

/* DailyWeatherAPIResponse
Daily structure from
https://api.openweathermap.org/data/2.5/onecall
*/
type DailyWeatherAPIResponse struct {
	Dt         int64     `json:"dt"`
	Sunrise    int       `json:"sunrise"`
	Sunset     int       `json:"sunset"`
	Temp       Temp      `json:"temp"`
	FeelsLike  FeelsLike `json:"feels_like"`
	Pressure   int       `json:"pressure"`
	Humidity   int       `json:"humidity"`
	DewPoint   float32   `json:"dew_point"`
	WindSpeed  float32   `json:"wind_speed"`
	WindDegree int64     `json:"wind_deg"`
	Weather    []Weather `json:"weather"`
	Clouds     int       `json:"clouds"`
	Pop        float32   `json:"pop"`
	Rain       float32   `json:"rain"`
	Uvi        float32   `json:"uvi"`
}

/* Temp
Temp structure underneath the Daily structure from
https://api.openweathermap.org/data/2.5/onecall
*/
type Temp struct {
	Day   float32 `json:"day"`
	Min   float32 `json:"min"`
	Max   float32 `json:"max"`
	Night float32 `json:"night"`
	Eve   float32 `json:"eve"`
	Morn  float32 `json:"morn"`
}

/* FeelsLike
FeelsLike structure underneath the Daily structure from
https://api.openweathermap.org/data/2.5/onecall
*/
type FeelsLike struct {
	Day   float32 `json:"day"`
	Night float32 `json:"night"`
	Eve   float32 `json:"eve"`
	Morn  float32 `json:"morn"`
}

/* "Weather"
weather structure underneath the Daily structure from
https://api.openweathermap.org/data/2.5/onecall
*/
type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

/* Gets the weather from https://api.openweathermap.org/data/2.5/onecall
excludes all data except for daily data.
Has the following pathParams:
	lat float
	lon float
	units string

Example URL:
	localhost:8080/weatherapi/v1/latitude/33.441792/longitude/-94.037689/units/metric
*/
func getWeather(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	// For production apiKey would move into an environment file
	apiKey := "2fa6857ac0dea1e59286b2426654ca4c"
	exclude := "current,minutely,hourly,alerts"
	var lat float64
	var lon float64
	var units string
	var unitAbbr string
	var err error

	// check path params for expected values and formats
	if val, ok := pathParams["lat"]; ok {
		lat, err = strconv.ParseFloat(val, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "accepts a float/int for latitude"}`))
			return
		}
	}

	if val, ok := pathParams["lon"]; ok {
		lon, err = strconv.ParseFloat(val, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "accepts a float/int for longitude"}`))
			return
		}
	}

	if val, ok := pathParams["units"]; ok {
		if val != "metric" && val != "imperial" && val != "standard" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "accepts 'units' as 'metric', 'imperial' or 'standard'"}`))
			return
		}
		switch val {
		case "metric":
			unitAbbr = "C"
		case "imperial":
			unitAbbr = "F"
		case "standard":
			unitAbbr = "K"
		}
		units = val
	}

	requestURL := fmt.Sprintf(`https://api.openweathermap.org/data/2.5/onecall?lat=%f&lon=%f&units=%s&exclude=%s&appid=%s`, lat, lon, units, exclude, apiKey)

	// Make a call to openweathermap and map it to our structure WeatherAPIResponse
	response, err := http.Get(requestURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Unable to get response from requestURL"}`))
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Unable to read response body"`))
		return
	}

	var responseObject WeatherAPIResponse
	json.Unmarshal(responseData, &responseObject)

	// Format incoming reponse to our response structure
	var forecastResponseArray []ForecastResponse
	for _, day := range responseObject.Daily {
		var forecastResponse ForecastResponse
		formattedTime, weekday, err := formatUnixTime(day.Dt, responseObject.Timezone)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Something went wrong formatting time"`))
			return
		}
		if weekday == "Saturday" || weekday == "Sunday" {
			continue
		}
		forecastResponse.Max = fmt.Sprintf(`%.2f%s`, day.Temp.Max, unitAbbr)
		forecastResponse.Min = fmt.Sprintf(`%.2f%s`, day.Temp.Min, unitAbbr)
		forecastResponse.Dt = formattedTime
		forecastResponse.Day = weekday
		forecastResponseArray = append(forecastResponseArray, forecastResponse)
	}
	jsonBytes, err := json.Marshal(forecastResponseArray)

	if err != nil {
		log.Fatalf("Unable to encode")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Unable to encode json"`))
		return
	}

	w.WriteHeader(http.StatusOK)
	res := fmt.Sprintf(`{"forecast": %s}`, jsonBytes)
	w.Write([]byte(res))
}

// Takes a unix timestamp and returns the date in the form of mm/dd/yyyy and the Weekday
func formatUnixTime(dt int64, timezone string) (string, string, error) {
	t := time.Unix(dt, 0)

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Fatalf("Unable to load location")
		return "", "", errors.New("Unable to load timezone location")
	}

	t = t.In(loc)

	return t.Format("01/02/2006"), t.Weekday().String(), nil
}

// CORS for development
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set headers
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// next
		next.ServeHTTP(w, r)
		return
	})
}

// Catch all
func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "Not found"}`))
}

func main() {
	r := mux.NewRouter()

	r.Use(CORS)

	api := r.PathPrefix("/weatherapi/v1").Subrouter()
	api.HandleFunc("/latitude/{lat}/longitude/{lon}/units/{units}", getWeather).Methods(http.MethodGet, http.MethodOptions)

	// This catches anything that doesn't conform to the above route(s)
	api.HandleFunc("", notFound)

	log.Fatal(http.ListenAndServe(":8080", r))
}
