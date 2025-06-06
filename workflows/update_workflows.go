package workflows

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

// Counter is a simple struct to hold a counter value
type Counter struct {
	Value int
}

// CounterWorkflow demonstrates the update feature
func CounterWorkflow(ctx workflow.Context, initialValue int) (int, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("CounterWorkflow started", "initialValue", initialValue)

	// Create a counter with the initial value
	counter := Counter{Value: initialValue}

	// Register update handler for incrementing the counter
	err := workflow.SetUpdateHandler(ctx, "increment", func(ctx workflow.Context, amount int) (int, error) {
		logger.Info("Incrementing counter", "currentValue", counter.Value, "amount", amount)
		counter.Value += amount
		return counter.Value, nil
	})
	if err != nil {
		logger.Error("Failed to register increment update handler", "error", err)
		return 0, err
	}

	// Register update handler for decrementing the counter
	err = workflow.SetUpdateHandler(ctx, "decrement", func(ctx workflow.Context, amount int) (int, error) {
		logger.Info("Decrementing counter", "currentValue", counter.Value, "amount", amount)
		counter.Value -= amount
		return counter.Value, nil
	})
	if err != nil {
		logger.Error("Failed to register decrement update handler", "error", err)
		return 0, err
	}

	// Register update handler for setting the counter to a specific value
	err = workflow.SetUpdateHandler(ctx, "set", func(ctx workflow.Context, value int) (int, error) {
		logger.Info("Setting counter", "oldValue", counter.Value, "newValue", value)
		counter.Value = value
		return counter.Value, nil
	})
	if err != nil {
		logger.Error("Failed to register set update handler", "error", err)
		return 0, err
	}

	// Register query handler to get the current counter value
	err = workflow.SetQueryHandler(ctx, "get_counter", func() (int, error) {
		return counter.Value, nil
	})
	if err != nil {
		logger.Error("Failed to register query handler", "error", err)
		return 0, err
	}

	// Keep the workflow running for a specified time
	// In a real application, you might want to keep it running indefinitely
	// or until a specific condition is met
	if err := workflow.Sleep(ctx, 24*time.Hour); err != nil {
		return counter.Value, err
	}

	logger.Info("CounterWorkflow completed", "finalValue", counter.Value)
	return counter.Value, nil
}

// UpdateableWorkflow demonstrates a more complex update scenario
func UpdateableWorkflow(ctx workflow.Context) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("UpdateableWorkflow started")

	// State that can be updated
	state := struct {
		Name        string
		Description string
		Tags        []string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}{
		Name:        "Initial State",
		Description: "This is the initial state of the workflow",
		Tags:        []string{"initial", "workflow"},
		CreatedAt:   workflow.Now(ctx),
		UpdatedAt:   workflow.Now(ctx),
	}

	// Register update handler for updating the state
	err := workflow.SetUpdateHandler(ctx, "update_state", func(ctx workflow.Context, updates map[string]interface{}) (map[string]interface{}, error) {
		logger.Info("Updating state", "updates", updates)

		// Apply updates
		for key, value := range updates {
			switch key {
			case "name":
				if name, ok := value.(string); ok {
					state.Name = name
				}
			case "description":
				if desc, ok := value.(string); ok {
					state.Description = desc
				}
			case "tags":
				if tags, ok := value.([]string); ok {
					state.Tags = tags
				}
			}
		}

		// Update the timestamp
		state.UpdatedAt = workflow.Now(ctx)

		// Return the current state
		return map[string]interface{}{
			"name":        state.Name,
			"description": state.Description,
			"tags":        state.Tags,
			"createdAt":   state.CreatedAt,
			"updatedAt":   state.UpdatedAt,
		}, nil
	})
	if err != nil {
		logger.Error("Failed to register update handler", "error", err)
		return "", err
	}

	// Register query handler to get the current state
	err = workflow.SetQueryHandler(ctx, "get_state", func() (map[string]interface{}, error) {
		return map[string]interface{}{
			"name":        state.Name,
			"description": state.Description,
			"tags":        state.Tags,
			"createdAt":   state.CreatedAt,
			"updatedAt":   state.UpdatedAt,
		}, nil
	})
	if err != nil {
		logger.Error("Failed to register query handler", "error", err)
		return "", err
	}

	// Keep the workflow running for a specified time
	if err := workflow.Sleep(ctx, 24*time.Hour); err != nil {
		return fmt.Sprintf("Workflow interrupted: %s", state.Name), err
	}

	logger.Info("UpdateableWorkflow completed", "finalState", state)
	return fmt.Sprintf("Workflow completed: %s", state.Name), nil
}
