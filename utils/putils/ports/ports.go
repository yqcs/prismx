package ports

import (
	"strconv"
)

// IsValid checks if a port is valid
func IsValid(v interface{}) bool {
	switch p := v.(type) {
	case string:
		return IsValidWithString(p)
	case int:
		return IsValidWithInt(p)
	}
	return false
}

// IsValidWithString checks if a string port is valid
func IsValidWithString(p string) bool {
	port, err := strconv.Atoi(p)
	return err == nil && IsValidWithInt(port)
}

// IsValidWithInt checks if an int port is valid
func IsValidWithInt(port int) bool {
	return port >= 1 && port <= 65535
}
