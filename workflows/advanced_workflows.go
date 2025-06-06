package workflows

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// ChildWorkflowParams contains parameters for the child workflow
type ChildWorkflowParams struct {
	Name     string
	Duration time.Duration
}

// ChildWorkflowResult contains the result of the child workflow
type ChildWorkflowResult struct {
	Message string
	Elapsed time.Duration
}

// ParentWorkflow demonstrates how to use child workflows
func ParentWorkflow(ctx workflow.Context, params ChildWorkflowParams) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("ParentWorkflow started", "params", params)

	// Set up child workflow options
	childOptions := workflow.ChildWorkflowOptions{
		WorkflowID:         "child-workflow",
		WorkflowRunTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithChildOptions(ctx, childOptions)

	// Execute child workflow
	var childResult ChildWorkflowResult
	err := workflow.ExecuteChildWorkflow(ctx, ChildWorkflow, params).Get(ctx, &childResult)
	if err != nil {
		logger.Error("Child workflow failed", "error", err)
		return "", err
	}

	result := fmt.Sprintf("Parent workflow completed. Child result: %s (took %v)",
		childResult.Message, childResult.Elapsed)
	logger.Info("ParentWorkflow completed", "result", result)
	return result, nil
}

// ChildWorkflow is executed as a child workflow
func ChildWorkflow(ctx workflow.Context, params ChildWorkflowParams) (ChildWorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("ChildWorkflow started", "params", params)

	// Configure activity options with timeout
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	startTime := workflow.Now(ctx)

	// Simulate some work
	if err := workflow.Sleep(ctx, params.Duration); err != nil {
		return ChildWorkflowResult{}, err
	}

	elapsed := workflow.Now(ctx).Sub(startTime)
	result := ChildWorkflowResult{
		Message: fmt.Sprintf("Hello, %s from child workflow!", params.Name),
		Elapsed: elapsed,
	}

	logger.Info("ChildWorkflow completed", "result", result)
	return result, nil
}

// SignalWorkflowData contains the data for the signal
type SignalWorkflowData struct {
	Message string
	Time    time.Time
}

// SignalWorkflow demonstrates how to use signals in a workflow
func SignalWorkflow(ctx workflow.Context, waitTime time.Duration) ([]string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("SignalWorkflow started", "waitTime", waitTime)

	// Configure activity options with timeout
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Define a channel for the signal
	signalChan := workflow.GetSignalChannel(ctx, "signal-channel")

	// Define a channel for the query result
	var receivedSignals []SignalWorkflowData

	// Set up query handler
	err := workflow.SetQueryHandler(ctx, "get_signals", func() ([]SignalWorkflowData, error) {
		return receivedSignals, nil
	})
	if err != nil {
		logger.Error("Failed to register query handler", "error", err)
		return nil, err
	}

	// Create a timer
	timer := workflow.NewTimer(ctx, waitTime)

	// Keep the workflow open until the timer fires or is cancelled
	timerFired := false
	for !timerFired {
		selector := workflow.NewSelector(ctx)

		// Add signal channel to selector
		selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
			var signal SignalWorkflowData
			c.Receive(ctx, &signal)
			receivedSignals = append(receivedSignals, signal)
			logger.Info("Received signal", "signal", signal)
		})

		// Add timer to selector
		selector.AddFuture(timer, func(f workflow.Future) {
			logger.Info("Timer fired")
			timerFired = true
		})

		// Wait for one of the conditions
		selector.Select(ctx)
	}

	// Process received signals
	var messages []string
	for _, signal := range receivedSignals {
		messages = append(messages, fmt.Sprintf("Signal received at %v: %s",
			signal.Time.Format(time.RFC3339), signal.Message))
	}

	logger.Info("SignalWorkflow completed", "signals", len(receivedSignals))
	return messages, nil
}

// ContinueAsNewWorkflow demonstrates the continue-as-new feature
func ContinueAsNewWorkflow(ctx workflow.Context, count int, maxCount int) (int, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("ContinueAsNewWorkflow execution", "count", count, "maxCount", maxCount)

	if count >= maxCount {
		logger.Info("ContinueAsNewWorkflow completed", "finalCount", count)
		return count, nil
	}

	// Increment the counter
	count++

	// Continue as new
	return count, workflow.NewContinueAsNewError(ctx, ContinueAsNewWorkflow, count, maxCount)
}
