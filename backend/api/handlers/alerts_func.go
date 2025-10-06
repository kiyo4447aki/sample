package handlers

import (
	"fmt"
	"strings"
	"time"

	"backend-proto/utils/storage"
)

func ifEmpty(s, def string) string {
	if strings.TrimSpace(s) == "" { return def }
	return s
}

func getObjectKey(deviceID string, eventID string, occurredAt time.Time, mime string) string {
	ext := storage.GetExtFromMime(mime)
	return fmt.Sprintf("alerts/%s/%04d/%02d/%02d/%s.%s",
		deviceID, occurredAt.Year(), int(occurredAt.Month()), occurredAt.Day(), eventID, ext,
	)
}

func validateObjectKey(key, deviceID, eventID, mime string)bool{
	if key == "" || deviceID == "" || eventID == "" || mime == "" {
		return false
	}
	if strings.Contains(key, "..") || strings.Contains(key, `\`) {
		return false
	}
	if strings.Contains(deviceID, "/") || strings.Contains(eventID, "/") {
		return false
	}
	if !strings.HasPrefix(key, "alerts/"+deviceID + "/"){
		return false 
	}
	if !strings.Contains(key, "/"+eventID+"."){ 
		return false 
	}
	switch mime{
	case "image/jpeg":
		if !strings.HasSuffix(key, ".jpg") { return false }
	case "image/png":
		if !strings.HasSuffix(key, ".png") { return false }
	default:
		return false
	}
	return true
}
