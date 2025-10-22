package utils

import (
	"encoding/csv"
	"io"
	"os"
	"safelyyou-monitoring-fleet/models"
)

// DeviceList keeps all registered devices in memory
var DeviceList = make(map[string]*models.Device)

func LoadDevices(filename string) {
	// TODO: Add error handling if file not found
	file, _ := os.Open(filename)
	reader := csv.NewReader(file)

	// Skip header in device.csv
	reader.Read()
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		// Add a new Device struct reference to DeviceList
		// The key is the device ID, and the value is a pointer to the Device struct
		DeviceList[record[0]] = &models.Device{}
	}
}
