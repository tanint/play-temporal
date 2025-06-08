package workflows

import (
	"time"

	"github.com/tanint/play-temporal/activities"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// SubscriptionParams contains parameters for starting a subscription
type SubscriptionParams struct {
	CustomerID string
	PlanID     string
}

// SubscriptionWorkflow handles the initial subscription creation and setup
func SubscriptionWorkflow(ctx workflow.Context, params SubscriptionParams) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("SubscriptionWorkflow started", "customerID", params.CustomerID, "planID", params.PlanID)

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

	// Step 1: Create the subscription
	var subscription activities.SubscriptionDetails
	err := workflow.ExecuteActivity(ctx, activities.CreateSubscriptionActivity, params.CustomerID, params.PlanID).Get(ctx, &subscription)
	if err != nil {
		logger.Error("Failed to create subscription", "error", err)
		return "", err
	}

	// Step 2: Calculate initial charges
	var amount float64
	err = workflow.ExecuteActivity(ctx, activities.CalculateChargesActivity, subscription).Get(ctx, &amount)
	if err != nil {
		logger.Error("Failed to calculate charges", "error", err)
		return "", err
	}

	// Step 3: Generate the first invoice
	var invoice activities.InvoiceDetails
	err = workflow.ExecuteActivity(ctx, activities.GenerateInvoiceActivity, subscription, amount).Get(ctx, &invoice)
	if err != nil {
		logger.Error("Failed to generate invoice", "error", err)
		return "", err
	}

	// Step 4: Process payment
	var payment activities.PaymentDetails
	err = workflow.ExecuteActivity(ctx, activities.ProcessPaymentActivity, invoice, subscription).Get(ctx, &payment)
	if err != nil {
		logger.Error("Failed to process payment", "error", err)
		return "", err
	}

	// Step 5: Send invoice email
	err = workflow.ExecuteActivity(ctx, activities.SendInvoiceEmailActivity, invoice, subscription.CustomerID).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to send invoice email", "error", err)
		// Continue despite email failure
	}

	// Step 6: Update subscription status based on payment
	var status string
	if payment.Status == "succeeded" {
		status = "active"
	} else {
		status = "payment_failed"
	}

	err = workflow.ExecuteActivity(ctx, activities.UpdateSubscriptionStatusActivity, subscription.ID, status).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to update subscription status", "error", err)
		return "", err
	}

	// Schedule the recurring billing workflow
	// In a real implementation, we would use a cron schedule or timer
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		WorkflowID:          "recurring-billing-" + subscription.ID,
		WorkflowRunTimeout:  365 * 24 * time.Hour, // Run for up to a year
		WorkflowTaskTimeout: 10 * time.Minute,
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)

	// Execute the recurring billing workflow as a child workflow
	recurringParams := RecurringBillingParams{
		SubscriptionID:  subscription.ID,
		CustomerID:      subscription.CustomerID,
		NextBillingDate: workflow.Now(ctx).AddDate(0, 1, 0), // 1 month from now
	}

	childWorkflowFuture := workflow.ExecuteChildWorkflow(ctx, RecurringBillingWorkflow, recurringParams)
	var childWorkflowExecution workflow.Execution
	err = childWorkflowFuture.GetChildWorkflowExecution().Get(ctx, &childWorkflowExecution)
	if err != nil {
		logger.Error("Failed to schedule recurring billing", "error", err)
		// Continue despite scheduling failure
	} else {
		logger.Info("Scheduled recurring billing workflow",
			"childWorkflowID", childWorkflowExecution.ID,
			"childRunID", childWorkflowExecution.RunID)
	}

	logger.Info("SubscriptionWorkflow completed", "subscriptionID", subscription.ID, "status", status)
	return subscription.ID, nil
}

// RecurringBillingParams contains parameters for the recurring billing workflow
type RecurringBillingParams struct {
	SubscriptionID  string
	CustomerID      string
	NextBillingDate time.Time
}

// RecurringBillingWorkflow handles the recurring billing for a subscription
func RecurringBillingWorkflow(ctx workflow.Context, params RecurringBillingParams) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("RecurringBillingWorkflow started", "subscriptionID", params.SubscriptionID)

	// Configure activity options with longer timeouts for reliability
	ao := workflow.ActivityOptions{
		StartToCloseTimeout:    30 * time.Second,
		ScheduleToStartTimeout: time.Minute,
		ScheduleToCloseTimeout: 2 * time.Minute,
		HeartbeatTimeout:       10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Set up a query handler to check the next billing date
	err := workflow.SetQueryHandler(ctx, "get_next_billing_date", func() (time.Time, error) {
		return params.NextBillingDate, nil
	})
	if err != nil {
		logger.Error("Failed to register query handler", "error", err)
		return err
	}

	// This workflow will continue running and processing monthly billing
	// until the subscription is cancelled or fails
	for {
		// Wait until the next billing date
		timeToWait := params.NextBillingDate.Sub(workflow.Now(ctx))
		if timeToWait > 0 {
			logger.Info("Waiting for next billing cycle",
				"subscriptionID", params.SubscriptionID,
				"nextBillingDate", params.NextBillingDate,
				"waitDuration", timeToWait.String())

			// Use a selector with a timer to allow for cancellation and queries during the wait
			timerCtx, cancel := workflow.WithCancel(ctx)
			timer := workflow.NewTimer(timerCtx, timeToWait)

			selector := workflow.NewSelector(ctx)
			selector.AddFuture(timer, func(f workflow.Future) {
				logger.Info("Timer completed, proceeding with billing cycle")
			})

			// Wait for the timer to complete
			selector.Select(ctx)
			cancel()
		} else {
			logger.Info("Next billing date is in the past, processing immediately")
		}

		// Mock getting the subscription details
		// In a real implementation, we would fetch the current subscription state
		subscription := activities.SubscriptionDetails{
			ID:              params.SubscriptionID,
			CustomerID:      params.CustomerID,
			PlanID:          "mock-plan",
			PricePerMonth:   49.99,
			Status:          "active",
			PaymentMethodID: "mock-payment-method",
		}

		// Step 1: Calculate charges for this billing period
		var amount float64
		err := workflow.ExecuteActivity(ctx, activities.CalculateChargesActivity, subscription).Get(ctx, &amount)
		if err != nil {
			logger.Error("Failed to calculate charges", "error", err)
			continue // Try again next billing cycle
		}

		// Step 2: Generate invoice
		var invoice activities.InvoiceDetails
		err = workflow.ExecuteActivity(ctx, activities.GenerateInvoiceActivity, subscription, amount).Get(ctx, &invoice)
		if err != nil {
			logger.Error("Failed to generate invoice", "error", err)
			continue // Try again next billing cycle
		}

		// Step 3: Process payment
		var payment activities.PaymentDetails
		err = workflow.ExecuteActivity(ctx, activities.ProcessPaymentActivity, invoice, subscription).Get(ctx, &payment)
		if err != nil {
			logger.Error("Failed to process payment", "error", err)
			continue // Try again next billing cycle
		}

		// Step 4: Send invoice email
		err = workflow.ExecuteActivity(ctx, activities.SendInvoiceEmailActivity, invoice, subscription.CustomerID).Get(ctx, nil)
		if err != nil {
			logger.Error("Failed to send invoice email", "error", err)
			// Continue despite email failure
		}

		// Step 5: Update subscription status based on payment
		var status string
		if payment.Status == "succeeded" {
			status = "active"
		} else {
			status = "payment_failed"
		}

		err = workflow.ExecuteActivity(ctx, activities.UpdateSubscriptionStatusActivity, subscription.ID, status).Get(ctx, nil)
		if err != nil {
			logger.Error("Failed to update subscription status", "error", err)
			// Continue despite status update failure
		}

		// If payment failed, we might want to retry or cancel the subscription
		// For this example, we'll just continue to the next billing cycle

		// Calculate the next billing date (1 month from now)
		params.NextBillingDate = workflow.Now(ctx).AddDate(0, 1, 0)

		logger.Info("Completed billing cycle",
			"subscriptionID", params.SubscriptionID,
			"nextBillingDate", params.NextBillingDate,
			"paymentStatus", payment.Status)
	}
}
