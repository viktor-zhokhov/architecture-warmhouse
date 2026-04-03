package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"
)

type TemperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

func getLocationBySensorID(sensorID string) string {
	switch sensorID {
	case "1":
		return "Living Room"
	case "2":
		return "Bedroom"
	case "3":
		return "Kitchen"
	default:
		return "Unknown"
	}
}

func getSensorIDByLocation(location string) string {
	switch location {
	case "Living Room":
		return "1"
	case "Bedroom":
		return "2"
	case "Kitchen":
		return "3"
	default:
		return "0"
	}
}

func generateTemperature(location, sensorID string) TemperatureResponse {
	value := 15.0 + rand.Float64()*15.0 // 15.0 - 30.0

	return TemperatureResponse{
		Value:       math.Round(value*100) / 100,
		Unit:        "°C",
		Timestamp:   time.Now(),
		Location:    location,
		Status:      "active",
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: fmt.Sprintf("Temperature in %s: %.1f°C", location, value),
	}
}

// GET /temperature?location=Living+Room
func handleTemperature(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	sensorID := r.URL.Query().Get("sensorId")

	if location == "" && sensorID == "" {
		http.Error(w, `{"error": "location or sensorId is required"}`, http.StatusBadRequest)
		return
	}

	if location == "" {
		location = getLocationBySensorID(sensorID)
	}

	if sensorID == "" {
		sensorID = getSensorIDByLocation(location)
	}

	resp := generateTemperature(location, sensorID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GET /temperature/{sensorId}
func handleTemperatureByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	sensorID := path[len("/temperature/"):]

	if sensorID == "" {
		http.Error(w, `{"error": "sensorId is required"}`, http.StatusBadRequest)
		return
	}

	location := getLocationBySensorID(sensorID)
	resp := generateTemperature(location, sensorID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/temperature", handleTemperature)
	mux.HandleFunc("/temperature/", handleTemperatureByID)

	fmt.Println("Temperature API starting on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
