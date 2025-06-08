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

	// Create workflow options with cron schedule
	workflowOptions := client.StartWorkflowOptions{
		ID:           fmt.Sprintf("recurring-billing-%s", *subscriptionID),
		TaskQueue:    "temporal-learning-task-queue",
		CronSchedule: "@monthly", // Run once a month
	}

	// Create recurring billing parameters
	params := workflows.RecurringBillingParams{
		SubscriptionID:  *subscriptionID,
		CustomerID:      *customerID,
		NextBillingDate: time.Now(), // Start billing immediately
	}

	// Start the recurring billing workflow
	log.Printf("Starting recurring billing workflow for subscription %s\n", *subscriptionID)
	workflowRun, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.RecurringBillingWorkflow, params)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Printf("Recurring billing workflow started with ID: %s and RunID: %s\n", workflowRun.GetID(), workflowRun.GetRunID())
}
