package main

import (
	"log"
	"net/http"
	"safelyyou-monitoring-fleet/router"
	"safelyyou-monitoring-fleet/utils"

	"github.com/gorilla/mux"
)

func main() {
	// Read devices from devices.csv
	utils.LoadDevices("devices.csv")

	// Create a mux router, which will be used to route incoming http requests
	// to appropriate handle functions
	r := mux.NewRouter()

	// All the endpoints are registered in RegisterRoutes function
	router.RegisterRoutes(r)

	// Starts an http server listening on port 6733
	if err := http.ListenAndServe(":6733", r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

	log.Println("Application started listening on port :: 6733")
}
