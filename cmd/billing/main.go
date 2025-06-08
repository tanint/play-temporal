package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/tanint/play-temporal/config"
	"github.com/tanint/play-temporal/workflows"
	"go.temporal.io/sdk/client"
)

func main() {
	// Define command line flags
	subscriptionID := flag.String("subscription", "", "Subscription ID for recurring billing")
	customerID := flag.String("customer", "", "Customer ID for recurring billing")
	flag.Parse()

	if *subscriptionID == "" || *customerID == "" {
		log.Fatalln("Subscription ID and Customer ID are required")
	}

	// Create the client object
	c, err := client.Dial(config.GetTemporalClientOptions())
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create recurring billing parameters
	params := workflows.RecurringBillingParams{
		SubscriptionID:  *subscriptionID,
		CustomerID:      *customerID,
		NextBillingDate: time.Now(), // Start billing immediately
	}

	// Create workflow options
	workflowOptions := client.StartWorkflowOptions{
		ID:                  fmt.Sprintf("recurring-billing-%s", *subscriptionID),
		TaskQueue:           "temporal-learning-task-queue",
		WorkflowRunTimeout:  24 * time.Hour,
		WorkflowTaskTimeout: 10 * time.Minute,
		CronSchedule:        "0 0 1 * *", // Run at midnight on the 1st day of each month
	}

	log.Printf("Starting recurring billing workflow with cron schedule: %s\n", workflowOptions.CronSchedule)

	// Start the workflow
	workflowRun, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.RecurringBillingWorkflow, params)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Printf("Started recurring billing workflow with ID: %s and RunID: %s\n", workflowRun.GetID(), workflowRun.GetRunID())
	log.Printf("The workflow will run according to cron schedule: %s\n", workflowOptions.CronSchedule)
	log.Printf("NOTE: This will not appear in the Schedules tab of the Temporal UI.")
	log.Printf("To create a visible schedule, use the Temporal CLI:")
	log.Printf("temporal schedule create --cron \"0 0 1 * *\" --workflow-id \"recurring-billing-%s\" --task-queue \"temporal-learning-task-queue\" --workflow-type \"RecurringBillingWorkflow\" --input \"{\\\"SubscriptionID\\\":\\\"%s\\\",\\\"CustomerID\\\":\\\"%s\\\",\\\"NextBillingDate\\\":\\\"%s\\\"}\"",
		*subscriptionID, *subscriptionID, *customerID, time.Now().Format(time.RFC3339))
}
