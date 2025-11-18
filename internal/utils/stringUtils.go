package utils

import (
	"crypto/rand"
	"encoding/hex"
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
