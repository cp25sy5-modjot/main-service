package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"
)

// Utility
func FirstOrEmpty(msg []string, fallback string) string {
	if len(msg) > 0 && msg[0] != "" {
		return msg[0]
	}
	return fallback
}

// random color generator
func GenerateRandomColor() string {
	bytes := make([]byte, 3)
	_, err := rand.Read(bytes)
	if err != nil {
		return "#000000" // return black if error occurs
	}
	return "#" + hex.EncodeToString(bytes)
}

func ConvertStringToTime(s string) *time.Time {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		log.Printf("Failed to parse date: %v", err)
		return nil
	}
	return &t
}
