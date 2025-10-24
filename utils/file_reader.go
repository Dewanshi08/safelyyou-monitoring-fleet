package utils

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"safelyyou-monitoring-fleet/models"
	"sync"
)

var (
	// DeviceList keeps all registered devices in memory
	DeviceList = make(map[string]*models.Device)

	// Store error while loading the devices
	DeviceLoadErr error

	// DeviceListMu to protect DeviceList for concurrent reads and writes
	DeviceListMu sync.RWMutex
)

// LoadDevices - reads devices from the provided CSV file,
// creates a new device list, and updates the global DeviceList
func LoadDevices(filename string) {
	// Create a new local map
	m := make(map[string]*models.Device)
	file, err := os.Open(filename)

	// Error handling if file not found, and update the global DeviceLoadErr;
	// and clear DeviceList on error to prevent serving stale or partial data
	if err != nil {
		log.Printf("Error %s", err.Error())
		handleLocks(nil, err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header in devices.csv
	reader.Read()
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Clear DeviceList on error to prevent serving stale or partial data
			handleLocks(nil, err)
			return
		}

		// Add a new Device struct reference to local map
		// The key is the device ID, and the value is a pointer to the Device struct
		m[record[0]] = &models.Device{}
	}

	handleLocks(m, nil)
}

func handleLocks(m map[string]*models.Device, err error) {
	// Lock the DeviceList when updating the value with local map
	DeviceListMu.Lock()
	DeviceList = m
	DeviceLoadErr = err
	// Once the DeviceList is updated, remove the lock
	DeviceListMu.Unlock()
}
