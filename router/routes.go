package router

import (
	"safelyyou-monitoring-fleet/handlers"

	"github.com/gorilla/mux"
)

// RegisterRoutes defines all endpoints, and mapping them to its corresponding controllers
func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/devices/{device_id}/heartbeat", handlers.RegisterHeartbeat).Methods("POST")
	router.HandleFunc("/api/v1/devices/{device_id}/stats", handlers.RegisterStats).Methods("POST")
	router.HandleFunc("/api/v1/devices/{device_id}/stats", handlers.GetStats).Methods("GET")
}
