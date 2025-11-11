package utils_test

import (
	"safelyyou-monitoring-fleet/logger"
	"safelyyou-monitoring-fleet/utils"
	"testing"
)

// TestLoadDevices_Success - verifies that LoadDevices successfully loads data
// from a valid CSV file and populates the DeviceList map
func TestLoadDevices_Success(t *testing.T) {
	logger.Init("")
	utils.LoadDevices("../devices.csv")

	// Assertion
	if utils.DeviceLoadErr != nil && len(utils.DeviceList) == 0 {
		t.Errorf("Assertion Failed: %s", utils.DeviceLoadErr.Error())
	}
}

// TestLoadDevices_Error - verifies that LoadDevices handles an invalid or missing CSV file path
// This test expects an error scenario where the file does not exist
func TestLoadDevices_Error(t *testing.T) {
	logger.Init("dev")
	// Loading CSV file which does not exists
	utils.LoadDevices("file.csv")

	// Assertion
	if utils.DeviceLoadErr == nil && len(utils.DeviceList) > 0 {
		t.Errorf("Assertion Failed: %s", utils.DeviceLoadErr.Error())
	}
}
