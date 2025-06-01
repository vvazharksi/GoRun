package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano()) // ensures random numbers differ each run
}

// Generates a timestamp-based ID: YYMMDDHHMMSS + 3-digit random number
func GenerateTimestampID() string {
	now := time.Now()
	timestamp := now.Format("060102150405") // YYMMDDHHMMSS
	randomPart := rand.Intn(1000)
	return fmt.Sprintf("%s%03d", timestamp, randomPart)
}
