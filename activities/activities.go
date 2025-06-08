package activities

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/activity"
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
	// Create a ticker for heartbeats (every 2 seconds)
	heartbeatTicker := time.NewTicker(2 * time.Second)
	defer heartbeatTicker.Stop()

	// Create a timer for the specified duration
	timer := time.NewTimer(time.Duration(durationSeconds) * time.Second)
	defer timer.Stop()

	// Track progress for heartbeats
	progress := 0
	totalSteps := durationSeconds

	// Log the start of the activity
	fmt.Printf("Starting long-running activity for %d seconds\n", durationSeconds)

	for {
		select {
		case <-ctx.Done():
			// The activity was cancelled
			fmt.Println("Activity was cancelled")
			return "", ctx.Err()
		case <-heartbeatTicker.C:
			// Send heartbeat to Temporal
			progress++
			fmt.Printf("Activity progress: %d/%d\n", progress, totalSteps)
			activity.RecordHeartbeat(ctx, progress)
		case <-timer.C:
			// The timer has expired, complete the activity
			fmt.Println("Activity completed successfully")
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
