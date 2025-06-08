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
	customerID := flag.String("customer", "cust123", "Customer ID for the subscription")
	planID := flag.String("plan", "basic-monthly", "Plan ID for the subscription")
	flag.Parse()

	// Create the client object
	c, err := client.Dial(config.GetTemporalClientOptions())
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create workflow options
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("subscription-%s-%v", *customerID, time.Now().Unix()),
		TaskQueue: "temporal-learning-task-queue",
	}

	// Create subscription parameters
	params := workflows.SubscriptionParams{
		CustomerID: *customerID,
		PlanID:     *planID,
	}

	// Start the subscription workflow
	log.Printf("Starting subscription workflow for customer %s with plan %s\n", *customerID, *planID)
	workflowRun, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.SubscriptionWorkflow, params)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Printf("Subscription workflow started with ID: %s and RunID: %s\n", workflowRun.GetID(), workflowRun.GetRunID())

	// Wait for workflow completion
	var subscriptionID string
	err = workflowRun.Get(context.Background(), &subscriptionID)
	if err != nil {
		log.Fatalln("Workflow failed", err)
	}

	log.Printf("Subscription created successfully with ID: %s\n", subscriptionID)
	log.Println("The recurring billing workflow has been scheduled and will run monthly.")
}
