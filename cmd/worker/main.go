package main

import (
	"log"

	"github.com/tanint/play-temporal/activities"
	"github.com/tanint/play-temporal/config"
	"github.com/tanint/play-temporal/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create the client object just once per process
	c, err := client.Dial(config.GetTemporalClientOptions())
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create a Worker instance
	w := worker.New(c, "temporal-learning-task-queue", worker.Options{})

	// Register basic workflows
	w.RegisterWorkflow(workflows.GreetingWorkflow)
	w.RegisterWorkflow(workflows.SequentialWorkflow)
	w.RegisterWorkflow(workflows.ParallelWorkflow)
	w.RegisterWorkflow(workflows.LongRunningWorkflow)
	w.RegisterWorkflow(workflows.ErrorHandlingWorkflow)

	// Register advanced workflows
	w.RegisterWorkflow(workflows.ParentWorkflow)
	w.RegisterWorkflow(workflows.ChildWorkflow)
	w.RegisterWorkflow(workflows.SignalWorkflow)
	w.RegisterWorkflow(workflows.ContinueAsNewWorkflow)

	// Register update workflows
	w.RegisterWorkflow(workflows.CounterWorkflow)
	w.RegisterWorkflow(workflows.UpdateableWorkflow)

	// Register subscription workflows
	w.RegisterWorkflow(workflows.SubscriptionWorkflow)
	w.RegisterWorkflow(workflows.RecurringBillingWorkflow)

	// Register activities
	w.RegisterActivity(activities.GreetingActivity)
	w.RegisterActivity(activities.FarewellActivity)
	w.RegisterActivity(activities.LongRunningActivity)
	w.RegisterActivity(activities.ErrorProneActivity)

	// Register subscription activities
	w.RegisterActivity(activities.CreateSubscriptionActivity)
	w.RegisterActivity(activities.CalculateChargesActivity)
	w.RegisterActivity(activities.GenerateInvoiceActivity)
	w.RegisterActivity(activities.ProcessPaymentActivity)
	w.RegisterActivity(activities.SendInvoiceEmailActivity)
	w.RegisterActivity(activities.UpdateSubscriptionStatusActivity)

	// Start listening to the Task Queue
	log.Println("Starting Temporal worker...")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}
}
