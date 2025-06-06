package activities

import (
	"context"
	"fmt"
	"time"
)

// GreetingActivity is a simple activity that returns a greeting message
func GreetingActivity(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("Hello, %s!", name), nil
}

// FarewellActivity is a simple activity that returns a farewell message
func FarewellActivity(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("Goodbye, %s!", name), nil
}

// LongRunningActivity simulates a long-running process
func LongRunningActivity(ctx context.Context, durationSeconds int) (string, error) {
	// Create a ticker to check for context cancellation
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Create a timer for the specified duration
	timer := time.NewTimer(time.Duration(durationSeconds) * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			// The activity was cancelled
			return "", ctx.Err()
		case <-ticker.C:
			// Just a tick, continue waiting
			continue
		case <-timer.C:
			// The timer has expired, complete the activity
			return fmt.Sprintf("Completed long-running activity after %d seconds", durationSeconds), nil
		}
	}
}

// ErrorProneActivity demonstrates how to handle errors in activities
func ErrorProneActivity(ctx context.Context, shouldFail bool) (string, error) {
	if shouldFail {
		return "", fmt.Errorf("activity failed as requested")
	}
	return "Activity completed successfully", nil
}
