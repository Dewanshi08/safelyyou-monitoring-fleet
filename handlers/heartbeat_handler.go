package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"safelyyou-monitoring-fleet/models"
	"safelyyou-monitoring-fleet/utils"

	"github.com/gorilla/mux"
)

// Post Method - RegisterHeartbeat handles heartbeat signals sent by devices
func RegisterHeartbeat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceID := vars["device_id"]
	log.Printf("RegisterHeartbeat for device_id %s", deviceID)

	var hb models.HeartbeatRequest
	json.NewDecoder(r.Body).Decode(&hb)

	// Get the Device object from DeviceList using DeviceID as key
	device, exists := utils.DeviceList[deviceID]
	// If key (DeviceID) not found return StatusNotFound
	if !exists {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	} else {
		// Append heartbeat sent_at to Device object heartbeats
		log.Println("Appending heartbeats")
		device.Heartbeats = append(device.Heartbeats, hb.SentAt)
		w.WriteHeader(http.StatusNoContent)
		return
	}
}
