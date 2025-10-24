package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"safelyyou-monitoring-fleet/router"
	"safelyyou-monitoring-fleet/utils"
	"testing"

	"github.com/gorilla/mux"
)

// Create a global mux router instance which will be reused across tests
var r = setupRouter()

// setupRouter - initializes a new mux router
func setupRouter() *mux.Router {
	// Create a new router and register API endpoints
	r := mux.NewRouter()
	router.RegisterRoutes(r)
	return r
}

// TestRegisterHeartbeat_Success - verifies that a valid heartbeat is registered
// for an existing device and returns HTTP 204
func TestRegisterHeartbeat_Success(t *testing.T) {
	// Load Devices from devices.csv
	utils.LoadDevices("../devices.csv")

	// Define a valid device_id that exists in CSV file and a valid JSON payload
	device_id := "18-b8-87-e7-1f-06"
	body := []byte(`{"sent_at": "2025-10-22T12:00:00Z"}`)

	// Execute the request and display the response
	rr, err := performHeartbeatRequest(device_id, body)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// Assertion
	if rr.Code != http.StatusNoContent {
		t.Errorf("Assertion Failed: Expected - 204, Actual - %d", rr.Code)
	}
}

// TestRegisterHeartbeat_NotFound - verifies that an invalid device_id is used for
// registering the heartbeat and returns HTTP 404
func TestRegisterHeartbeat_NotFound(t *testing.T) {
	// Load Devices from devices.csv
	utils.LoadDevices("../devices.csv")

	// Define an invalid device_id that does not exists in CSV file and a valid JSON payload
	device_id := "device123"
	body := []byte(`{"sent_at": "2025-10-22T12:00:00Z"}`)

	// Execute the request and display the response
	rr, err := performHeartbeatRequest(device_id, body)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// Assertion
	if rr.Code != http.StatusNotFound {
		t.Errorf("Assertion Failed: Expected - 404, Actual - %d", rr.Code)
	}
}

// TestRegisterHeartbeat_BadRequest - verifies that an invalid JSON payload is used for
// registering the heartbeat and returns HTTP 400
func TestRegisterHeartbeat_BadRequest(t *testing.T) {
	// Load Devices from devices.csv
	utils.LoadDevices("../devices.csv")

	// Define a valid device_id that exists in CSV file and an invalid JSON payload
	device_id := "18-b8-87-e7-1f-06"
	body := []byte(`{sent_at: 2025-10-22T12:00:00Z}`)

	// Execute the request and display the response
	rr, err := performHeartbeatRequest(device_id, body)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// Assertion
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Assertion Failed: Expected - 400, Actual - %d", rr.Code)
	}
}

// TestRegisterHeartbeat_ServerError - verifies that when an error occurs while
// loading the DeviceList, then endpoint should return HTTP 500
func TestRegisterHeartbeat_ServerError(t *testing.T) {
	// Load Devices from devices.csv
	utils.LoadDevices("../devices.csv")

	// Using non-existing csv file to generate error while loading the DeviceList
	utils.LoadDevices("file.csv")
	// Define a valid device_id that exists in CSV file and an invalid JSON payload
	device_id := "18-b8-87-e7-1f-06"
	body := []byte(`{"sent_at": "2025-10-22T12:00:00Z"}`)

	// Execute the request and display the response
	rr, err := performHeartbeatRequest(device_id, body)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// Assertion
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Assertion Failed: Expected - 500, Actual - %d", rr.Code)
	}
}

// performHeartbeatRequest - a helper function that sends an HTTP POST request
// to the `/api/v1/devices/{device_id}/heartbeat` endpoint
func performHeartbeatRequest(device_id string, body []byte) (*httptest.ResponseRecorder, error) {
	// Create a POST request with JSON body for the given device_id
	req, err := http.NewRequest(http.MethodPost, "/api/v1/devices/"+device_id+"/heartbeat", bytes.NewBuffer(body))

	// Incase error occurs while implementing request, stop the executiong and return the error
	if err != nil {
		return nil, err
	}

	// Set request content type to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to capture handler's reponse and serve the request via registered router
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	return rr, nil
}
