package main

import (
	"log"
	"net/http"
	"os"
	"safelyyou-monitoring-fleet/logger"
	"safelyyou-monitoring-fleet/router"
	"safelyyou-monitoring-fleet/utils"

	"github.com/gorilla/mux"
)

func main() {
	env := os.Getenv("APP_ENV")
	logger.Init(env)
	logger.Log.Info("Load devices")
	// Read devices from devices.csv
	utils.LoadDevices("devices.csv")

	// Create a mux router, which will be used to route incoming http requests
	// to appropriate handle functions
	r := mux.NewRouter()

	// All the endpoints are registered in RegisterRoutes function
	router.RegisterRoutes(r)

	logger.Log.Info("Application started listening on port :: 6733")
	// Starts an http server listening on port 6733
	if err := http.ListenAndServe(":6733", r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
