package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"safelyyou-monitoring-fleet/models"
	"safelyyou-monitoring-fleet/utils"
	"time"

	"github.com/gorilla/mux"
)

// Post Method - RegisterStats records device statistics such as upload durations
func RegisterStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceID := vars["device_id"]
	log.Printf("RegisterStats for device_id: %s", deviceID)

	var st models.StatsRequest
	json.NewDecoder(r.Body).Decode(&st)
	log.Println("Stats for int64: ", st.SentAt, " upload time: ", st.UploadTime)

	// Get the Device object from DeviceList using DeviceID as key
	device, exists := utils.DeviceList[deviceID]
	// If key (DeviceID) not found return StatusNotFound
	if !exists {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	} else {
		// Append stats upload_time to Device object uploadTimes
		log.Println("Appending upload times")
		device.UploadTimes = append(device.UploadTimes, st.UploadTime)
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

// Get Method - GetStats retrieves aggregated statistics like uptime and average upload time
func GetStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceID := vars["device_id"]
	log.Printf("Getting the stats for device %s", deviceID)

	// TODO: Device is considered 'OFFLINE' if there is no heartbeat for that minute

	// Get the Device object from DeviceList using DeviceID as key
	device, exists := utils.DeviceList[deviceID]
	// If key (DeviceID) not found return StatusNotFound
	if !exists {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	// Calculating the uptime
	var uptime float64
	heartbeats := device.Heartbeats
	// If number of heartbeats are less than 2, uptime cannot be calculated
	if len(heartbeats) < 2 {
		uptime = 0.0
	} else {
		// Get the duration for last and first heartbeat
		duration := heartbeats[len(heartbeats)-1].Sub(heartbeats[0]).Minutes()
		uptime = (float64(len(heartbeats)) / duration) * 100
	}

	// Calculating average upload time
	var sum int64
	for _, t := range device.UploadTimes {
		sum += t
	}
	avg := sum / int64(len(device.UploadTimes))

	// Creating the response for uptime with avgUploadTime in duration format
	response := models.StatsResponse{
		Uptime:        uptime,
		AvgUploadTime: time.Duration(avg).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
