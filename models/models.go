package models

import "time"

// Device stores the list of timestamp of heartbeats and uploadTimes duration
type Device struct {
	Heartbeats  []time.Time
	UploadTimes []int64
}

// HeartbeatRequest represents the payload sent to /devices/{device_id}/heartbeat
type HeartbeatRequest struct {
	SentAt time.Time `json:"sent_at"`
}

// StatsRequest represents the payload sent to /devices/{device_id}/stats
type StatsRequest struct {
	SentAt     time.Time `json:"sent_at"`
	UploadTime int64     `json:"upload_time"`
}

// StatsResponse represents the API response for GET /devices/{device_id}/stats
type StatsResponse struct {
	Uptime        float64 `json:"uptime"`
	AvgUploadTime string  `json:"avg_upload_time"`
}
