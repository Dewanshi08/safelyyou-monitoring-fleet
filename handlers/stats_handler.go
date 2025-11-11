package handlers

import (
	"encoding/json"
	"net/http"
	"safelyyou-monitoring-fleet/logger"
	"safelyyou-monitoring-fleet/models"
	"safelyyou-monitoring-fleet/utils"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Post Method - RegisterStats records device statistics such as upload durations
func RegisterStats(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Inside RegisterStats of a device")
	// If devices failed to load at startup (e.g. CSV/DB error),
	// it will return a HTTP 500 because no device records exist
	if utils.DeviceLoadErr != nil {
		logger.Log.Error("Unable to load devices")
		http.Error(w, "Error response", http.StatusInternalServerError)
		return
	}

	// Extract device_id from the URL path
	vars := mux.Vars(r)
	deviceID := vars["device_id"]
	logger.Log.Info("RegisterStats for", logger.MaskDeviceID(deviceID))

	var uploadStatReq models.UploadStatsRequest
	err := json.NewDecoder(r.Body).Decode(&uploadStatReq)
	if err != nil {
		logger.Log.Error("Invalid JSON body")
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	logger.Log.Debug("Stats for upload", zap.Time("sentAt", uploadStatReq.SentAt), zap.Int64("uploadTime", uploadStatReq.UploadTime))

	// Read lock to prevent updates on DeviceList
	utils.DeviceListMu.RLock()
	// Get the Device object from DeviceList using DeviceID as key
	device, exists := utils.DeviceList[deviceID]
	utils.DeviceListMu.RUnlock()

	// Return HTTP 404, if key (DeviceID) not found
	if !exists {
		logger.Log.Error("Device not found")
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}
	// Append stats upload_time to Device object uploadTimes
	logger.Log.Info("Appending upload times")

	// Write lock required when mutating slice on Device struct
	utils.DeviceListMu.Lock()
	device.UploadTimes = append(device.UploadTimes, uploadStatReq.UploadTime)
	utils.DeviceListMu.Unlock()

	logger.Log.Info("Successfully registered device statistics")
	w.WriteHeader(http.StatusNoContent)
}

// Get Method - GetStats retrieves aggregated statistics like uptime and average upload time
func GetStats(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Inside GetStats of a device")
	// If devices failed to load at startup (e.g. CSV/DB error),
	// it will return a HTTP 500 because no device records exist
	if utils.DeviceLoadErr != nil {
		logger.Log.Error("Unable to load devices")
		http.Error(w, "Error response", http.StatusInternalServerError)
		return
	}

	// Extract device_id from the URL path
	vars := mux.Vars(r)
	deviceID := vars["device_id"]
	logger.Log.Info("GetStats for", logger.MaskDeviceID(deviceID))

	// Read lock to prevent updates on DeviceList
	utils.DeviceListMu.RLock()
	// Get the Device object from DeviceList using DeviceID as key
	device, exists := utils.DeviceList[deviceID]
	utils.DeviceListMu.RUnlock()

	// Return HTTP 404, if key (DeviceID) not found
	if !exists {
		logger.Log.Error("Device not found")
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}
	logger.Log.Debug("Stats", logger.MaskDeviceID(deviceID), zap.Int("#RegisterHeartbeats", len(device.Heartbeats)), zap.Int("#UploadTimes", len(device.UploadTimes)))

	var uptime float64
	heartbeats := device.Heartbeats

	// Uptime will be -1.0 if no heartbeats are registered or device is offline
	if len(heartbeats) == 0 || (len(heartbeats) >= 2 && isDeviceOffline(heartbeats)) {
		logger.Log.Warn("Device is Offline", logger.MaskDeviceID(deviceID))
		uptime = -1.0
	} else if len(heartbeats) < 2 {
		// If only one heartbeat is registered then return uptime 0.0
		uptime = 0.0
	} else {
		// duration - is difference between last and first heartbeat
		duration := heartbeats[len(heartbeats)-1].Sub(heartbeats[0]).Minutes()
		if duration == 0 {
			logger.Log.Error("Invalid data: first and last heartbeat has same timestamp")
			http.Error(w, "Invalid data: first and last heartbeat has same timestamp", http.StatusBadRequest)
			return
		}
		// Calculating the uptime; get number of heartbeats and divide it with duration
		uptime = (float64(len(heartbeats)) / duration) * 100
	}

	// If number of upload time is 0, i.e., return expection as division by zero is not possible
	if len(device.UploadTimes) == 0 {
		logger.Log.Error("No upload times found")
		http.Error(w, "No upload times found", http.StatusBadRequest)
		return
	}

	// Calculating average upload time
	var sum int64
	for _, t := range device.UploadTimes {
		sum += t
	}

	avg := sum / int64(len(device.UploadTimes))

	logger.Log.Info("Generating a reponse object with uptime and avgUploadTime")
	// Creating the response for uptime with avgUploadTime in duration format
	response := models.GetDeviceStatsResponse{
		Uptime:        uptime,
		AvgUploadTime: time.Duration(avg).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(response)

	// Return HTTP 500, if an error occurs while implementing a JSON response
	if err != nil {
		logger.Log.Error("Unable to encode reponse to JSON format")
		http.Error(w, "Error response: ", http.StatusInternalServerError)
	}
}

// isDeviceOffline - checks if the device is considered offline
// A device is assumed offline when the gap between the two most recent
// heartbeats exceeds one minute
//
// *NOTE: For real-time device status, lastHeartBeat should be compared with
// current time (time.Now().UTC().Sub(lastHeartbeat))
//
// Here, using comparison between last two heartbeats timestamps, because
// the simulator registers heartbeats in old time
func isDeviceOffline(heartbeats []time.Time) bool {
	// Get the last two heartbeat timestamps
	lastHeartbeat := heartbeats[len(heartbeats)-1]
	previousHeartbeat := heartbeats[len(heartbeats)-2]

	// Compute time difference between them
	diff := previousHeartbeat.Sub(lastHeartbeat)
	logger.Log.Info("", zap.Time("previousHeartbeat", previousHeartbeat), zap.Time("lastHeartbeat", lastHeartbeat), zap.Duration("diff", diff))

	// Device is consider offline, if  difference is more than 1 minute
	return diff > time.Minute
}
