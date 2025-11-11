package logger

import "go.uber.org/zap"

var Log *zap.Logger

// Init - to initialize the Log variable based on environment
func Init(env string) {
	if env == "local" || env == "dev" {
		Log, _ = zap.NewDevelopment()
	} else {
		Log, _ = zap.NewProduction()
	}
}

// MaskDeviceID - to mask the deviceId when logging
func MaskDeviceID(id string) zap.Field {
	if len(id) < 3 {
		return zap.String("deviceID", id)
	}
	return zap.String("deviceID", "*****-"+id[len(id)-2:])
}
