package utils

import (
	"fmt"
	"time"
)

type RetryFunc func() (stop bool)

// Retry run function until true
func Retry(interval time.Duration, maxAttempts int, doFunc RetryFunc) error {
	ticker := time.NewTicker(interval)
	currentAttempt := 0
	for range ticker.C {
		currentAttempt++
		if currentAttempt > maxAttempts {
			return fmt.Errorf("number of attempts (%d) to perform test exceeded", maxAttempts)
		}
		stop := doFunc()
		if stop {
			break
		}
	}
	ticker.Stop()
	return nil
}
