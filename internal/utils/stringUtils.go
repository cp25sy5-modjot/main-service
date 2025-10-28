package utils

// Utility
func FirstOrEmpty(msg []string, fallback string) string {
	if len(msg) > 0 && msg[0] != "" {
		return msg[0]
	}
	return fallback
}
