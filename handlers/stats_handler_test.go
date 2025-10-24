package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"safelyyou-monitoring-fleet/utils"
	"testing"
)

// TestRegisterStats_Success - verifies that the valid stat is registered
// for an existing device and returns HTTP 204
func TestRegisterStats_Success(t *testing.T) {
	// Load Devices from devices.csv
	utils.LoadDevices("../devices.csv")

	// Define a valid device_id that exists in CSV file and a valid JSON payload
	device_id := "18-b8-87-e7-1f-06"
	body := []byte(`{"sent_at": "2025-10-22T12:00:00Z", "upload_time":273561280315}`)

	// Execute the request and display the response
	rr, err := performStatsRequest(device_id, body, http.MethodPost)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// Assertion
	if rr.Code != http.StatusNoContent {
		t.Errorf("Assertion Failed: Expected - 204, Actual - %d", rr.Code)
	}
}

// TestRegisterStats_NotFound - verifies that an invalid device_id is used for
// registering the stats and returns HTTP 404
func TestRegisterStats_NotFound(t *testing.T) {
	// Load Devices from devices.csv
	utils.LoadDevices("../devices.csv")

	// Define an invalid device_id that does not exists in CSV file and a valid JSON payload
	device_id := "device123"
	body := []byte(`{"sent_at": "2025-10-22T12:00:00Z", "upload_time":273561280315}`)

	// Execute the request and display the response
	rr, err := performStatsRequest(device_id, body, http.MethodPost)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// Assertion
	if rr.Code != http.StatusNotFound {
		t.Errorf("Assertion Failed: Expected - 404, Actual - %d", rr.Code)
	}
}

// TestRegisterStats_BadRequest - verifies that an invalid JSON payload is used for
// registering the stats and returns HTTP 400
func TestRegisterStats_BadRequest(t *testing.T) {
	// Load Devices from devices.csv
	utils.LoadDevices("../devices.csv")

	// Define a valid device_id that exists in CSV file and an invalid JSON payload
	device_id := "device123"
	body := []byte(`{sent_at: 2025-10-22T12:00:00Z, "upload_time":273561280315}}`)

	// Execute the request and display the response
	rr, err := performStatsRequest(device_id, body, http.MethodPost)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// Assertion
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Assertion Failed: Expected - 204, Actual - %d", rr.Code)
	}
}

// TestRegisterStats_ServerError - verifies that when an error occurs while
// loading the DeviceList, then endpoint should return HTTP 500
func TestRegisterStats_ServerError(t *testing.T) {
	// Using non-existing csv file to generate error while loading the DeviceList
	utils.LoadDevices("file.csv")
	// Define a valid device_id that exists in CSV file and an invalid JSON payload
	device_id := "18-b8-87-e7-1f-06"
	body := []byte(`{"sent_at": "2025-10-22T12:00:00Z"}`)

	// Execute the request and display the response
	rr, err := performStatsRequest(device_id, body, http.MethodPost)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// Assertion
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Assertion Failed: Expected - 500, Actual - %d", rr.Code)
	}
}

// TestGetStats_Success - verifies that a valid device_id
// should fetch the stats and returns HTTP 200
func TestGetStats_Success(t *testing.T) {
	// Load Devices from devices.csv
	utils.LoadDevices("../devices.csv")

	// Define a valid device_id that exists in CSV file and a valid JSON payload
	device_id := "18-b8-87-e7-1f-06"
	bodies := [][]byte{
		[]byte(`{"sent_at": "2025-10-22T12:00:10Z"}`),
		[]byte(`{"sent_at": "2025-10-22T12:00:00Z"}`),
	}

	for _, body := range bodies {
		// Register the heartbeats and stats to device object
		_, err := performHeartbeatRequest(device_id, body)
		if err != nil {
			t.Fatalf("Failed to perform request: %v", err)
		}

		_, err = performStatsRequest(device_id, body, http.MethodPost)
		if err != nil {
			t.Fatalf("Failed to perform request: %v", err)
		}
	}

	// Execute the request and display the response
	rr, err := performStatsRequest(device_id, nil, http.MethodGet)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// Assertion
	if rr.Code != http.StatusOK {
		t.Errorf("Assertion Failed: Expected - 200, Actual - %d", rr.Code)
	}
}

// TestGetStats_NotFound - verifies that an invalid device_id
// should throw an error while trying to fetch the stats and returns HTTP 404
func TestGetStats_NotFound(t *testing.T) {
	// Load Devices from devices.csv
	utils.LoadDevices("../devices.csv")

	// Define an invalid device_id that does not exists in CSV file
	device_id := "device_123"

	// Execute the request and display the response
	rr, err := performStatsRequest(device_id, nil, http.MethodGet)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// Assertion
	if rr.Code != http.StatusNotFound {
		t.Errorf("Assertion Failed: Expected - 404, Actual - %d", rr.Code)
	}
}

// TestGetStats_Exception - verifies that a valid device_id with empty details
// should return invalid data while trying to fetch the stats for it and returns HTTP 400
func TestGetStats_BadRequest1(t *testing.T) {
	// Load Devices from devices.csv
	utils.LoadDevices("../devices.csv")

	// Define a valid device_id for which heartbeats and upload times are empty (invalid data)
	device_id := "18-b8-87-e7-1f-06"

	// Execute the request and display the response
	rr, err := performStatsRequest(device_id, nil, http.MethodGet)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// Assertion
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Assertion Failed: Expected - 204, Actual - %d", rr.Code)
	}
}

// TestGetStats_Exception - verifies that a valid device_id with empty upload times
// should return no upload time found while trying to fetch the stats for it and returns HTTP 400
func TestGetStats_BadRequest2(t *testing.T) {
	utils.LoadDevices("../devices.csv")

	// Define a valid device_id that exists in CSV file and a valid JSON payload
	device_id := "18-b8-87-e7-1f-06"
	body := []byte(`{"sent_at": "2025-10-22T12:00:00Z"}`)

	// Register the heartbeat
	_, err := performHeartbeatRequest(device_id, body)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// Define a valid device_id for which upload times are empty (invalid data)
	device_id = "60-6b-44-84-dc-64"

	// Execute the request and display the response
	rr, err := performStatsRequest(device_id, nil, http.MethodGet)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// Assertion
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Assertion Failed: Expected - 204, Actual - %d", rr.Code)
	}
}

// performStatsRequest - a helper function that sends an HTTP GET/POST request
// to the `/api/v1/devices/{device_id}/stats` endpoint
func performStatsRequest(deviceID string, body []byte, method string) (*httptest.ResponseRecorder, error) {
	// Create a <method> request with JSON body for the given device_id
	req, err := http.NewRequest(method, "/api/v1/devices/"+deviceID+"/stats", bytes.NewBuffer(body))

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
