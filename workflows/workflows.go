package workflows

import (
	"time"

	"github.com/tanint/play-temporal/activities"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// GreetingWorkflow is a simple workflow that calls the greeting activity
func GreetingWorkflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("GreetingWorkflow started", "name", name)

	var result string
	err := workflow.ExecuteActivity(ctx, activities.GreetingActivity, name).Get(ctx, &result)
	if err != nil {
		logger.Error("GreetingActivity failed", "error", err)
		return "", err
	}

	logger.Info("GreetingWorkflow completed", "result", result)
	return result, nil
}

// SequentialWorkflow demonstrates calling multiple activities in sequence
func SequentialWorkflow(ctx workflow.Context, name string) ([]string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("SequentialWorkflow started", "name", name)

	// Configure activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var greetingResult string
	err := workflow.ExecuteActivity(ctx, activities.GreetingActivity, name).Get(ctx, &greetingResult)
	if err != nil {
		logger.Error("GreetingActivity failed", "error", err)
		return nil, err
	}

	var farewellResult string
	err = workflow.ExecuteActivity(ctx, activities.FarewellActivity, name).Get(ctx, &farewellResult)
	if err != nil {
		logger.Error("FarewellActivity failed", "error", err)
		return nil, err
	}

	results := []string{greetingResult, farewellResult}
	logger.Info("SequentialWorkflow completed", "results", results)
	return results, nil
}

// ParallelWorkflow demonstrates calling multiple activities in parallel
func ParallelWorkflow(ctx workflow.Context, name string) ([]string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("ParallelWorkflow started", "name", name)

	// Configure activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Start both activities in parallel
	greetingFuture := workflow.ExecuteActivity(ctx, activities.GreetingActivity, name)
	farewellFuture := workflow.ExecuteActivity(ctx, activities.FarewellActivity, name)

	// Wait for both activities to complete
	var greetingResult string
	if err := greetingFuture.Get(ctx, &greetingResult); err != nil {
		logger.Error("GreetingActivity failed", "error", err)
		return nil, err
	}

	var farewellResult string
	if err := farewellFuture.Get(ctx, &farewellResult); err != nil {
		logger.Error("FarewellActivity failed", "error", err)
		return nil, err
	}

	results := []string{greetingResult, farewellResult}
	logger.Info("ParallelWorkflow completed", "results", results)
	return results, nil
}

// LongRunningWorkflow demonstrates a workflow with a long-running activity
func LongRunningWorkflow(ctx workflow.Context, durationSeconds int) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("LongRunningWorkflow started", "duration", durationSeconds)

	// Configure activity options with a longer timeout
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Duration(durationSeconds+5) * time.Second,
		HeartbeatTimeout:    5 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	err := workflow.ExecuteActivity(ctx, activities.LongRunningActivity, durationSeconds).Get(ctx, &result)
	if err != nil {
		logger.Error("LongRunningActivity failed", "error", err)
		return "", err
	}

	logger.Info("LongRunningWorkflow completed", "result", result)
	return result, nil
}

// ErrorHandlingWorkflow demonstrates how to handle activity errors
func ErrorHandlingWorkflow(ctx workflow.Context, shouldFail bool) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("ErrorHandlingWorkflow started", "shouldFail", shouldFail)

	// Configure activity options with retry policy
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	err := workflow.ExecuteActivity(ctx, activities.ErrorProneActivity, shouldFail).Get(ctx, &result)
	if err != nil {
		logger.Error("ErrorProneActivity failed", "error", err)
		return "", err
	}

	logger.Info("ErrorHandlingWorkflow completed", "result", result)
	return result, nil
}
