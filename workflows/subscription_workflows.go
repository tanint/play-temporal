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

	// Note: We're not starting the recurring billing workflow here
	// Instead, it should be started separately using the RecurringBillingStarter
	// This avoids any issues with parent-child workflow relationships

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
// This workflow is designed to be used with a cron schedule
func RecurringBillingWorkflow(ctx workflow.Context, params RecurringBillingParams) error {
	logger := workflow.GetLogger(ctx)

	// Get information about the current workflow execution
	info := workflow.GetInfo(ctx)

	// Log workflow execution details including cron schedule
	logger.Info("RecurringBillingWorkflow started",
		"subscriptionID", params.SubscriptionID,
		"workflowID", info.WorkflowExecution.ID,
		"runID", info.WorkflowExecution.RunID,
		"cronSchedule", info.CronSchedule)

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

	// For cron workflows, we don't need to wait for the next billing date
	// The cron schedule will automatically trigger the workflow at the right time
	logger.Info("Processing billing cycle for subscription",
		"subscriptionID", params.SubscriptionID,
		"billingDate", workflow.Now(ctx).Format(time.RFC3339))

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
	err = workflow.ExecuteActivity(ctx, activities.CalculateChargesActivity, subscription).Get(ctx, &amount)
	if err != nil {
		logger.Error("Failed to calculate charges", "error", err)
		return err
	}

	// Step 2: Generate invoice
	var invoice activities.InvoiceDetails
	err = workflow.ExecuteActivity(ctx, activities.GenerateInvoiceActivity, subscription, amount).Get(ctx, &invoice)
	if err != nil {
		logger.Error("Failed to generate invoice", "error", err)
		return err
	}

	// Step 3: Process payment
	var payment activities.PaymentDetails
	err = workflow.ExecuteActivity(ctx, activities.ProcessPaymentActivity, invoice, subscription).Get(ctx, &payment)
	if err != nil {
		logger.Error("Failed to process payment", "error", err)
		return err
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

	// Calculate the next billing date (1 month from now)
	nextBillingDate := workflow.Now(ctx).AddDate(0, 1, 0)

	logger.Info("Completed billing cycle",
		"subscriptionID", params.SubscriptionID,
		"nextBillingDate", nextBillingDate,
		"paymentStatus", payment.Status)

	return nil
}
