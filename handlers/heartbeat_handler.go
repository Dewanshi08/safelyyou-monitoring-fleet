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
	// If devices failed to load at startup (e.g. CSV/DB error),
	// it will return a HTTP 500 because no device records exist
	if utils.DeviceLoadErr != nil {
		http.Error(w, "Error response", http.StatusInternalServerError)
		return
	}

	// Extract device_id from the URL path
	vars := mux.Vars(r)
	deviceID := vars["device_id"]
	log.Printf("RegisterHeartbeat for device_id %s", deviceID)

	var heartbeatReq models.HeartbeatRequest
	err := json.NewDecoder(r.Body).Decode(&heartbeatReq)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Read lock to prevent updates on DeviceList
	utils.DeviceListMu.RLock()
	// Get the Device object from DeviceList using DeviceID as key
	device, exists := utils.DeviceList[deviceID]
	utils.DeviceListMu.RUnlock()

	// Return HTTP 404, if key (DeviceID) not found
	if !exists {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}
	// Append heartbeat sent_at to Device object heartbeats
	log.Println("Appending heartbeats")

	// Write lock required when mutating slice on Device struct
	utils.DeviceListMu.Lock()
	device.Heartbeats = append(device.Heartbeats, heartbeatReq.SentAt)
	utils.DeviceListMu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}
