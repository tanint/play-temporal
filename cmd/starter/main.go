package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/temporalio/play-temporal/config"
	"github.com/temporalio/play-temporal/workflows"
	"go.temporal.io/sdk/client"
)

func main() {
	// Define command line flags
	workflowType := flag.String("workflow", "greeting",
		"Workflow type to run (greeting, sequential, parallel, long-running, error-handling, parent, signal, continue-as-new)")
	name := flag.String("name", "World", "Name to use in greeting workflows")
	duration := flag.Int("duration", 5, "Duration in seconds for long-running workflow")
	shouldFail := flag.Bool("fail", false, "Whether the error-handling workflow should fail")
	waitTime := flag.Int("wait", 30, "Wait time in seconds for signal workflow")
	count := flag.Int("count", 0, "Starting count for continue-as-new workflow")
	maxCount := flag.Int("max", 10, "Maximum count for continue-as-new workflow")
	flag.Parse()

	// Create the client object just once per process
	c, err := client.Dial(config.GetTemporalClientOptions())
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("%s-workflow-%v", *workflowType, time.Now().Unix()),
		TaskQueue: "temporal-learning-task-queue",
	}

	var workflowRun client.WorkflowRun
	var workflowID string

	// Execute the selected workflow
	switch *workflowType {
	case "greeting":
		workflowRun, err = c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.GreetingWorkflow, *name)
		workflowID = "Greeting"
	case "sequential":
		workflowRun, err = c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.SequentialWorkflow, *name)
		workflowID = "Sequential"
	case "parallel":
		workflowRun, err = c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.ParallelWorkflow, *name)
		workflowID = "Parallel"
	case "long-running":
		workflowRun, err = c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.LongRunningWorkflow, *duration)
		workflowID = "LongRunning"
	case "error-handling":
		workflowRun, err = c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.ErrorHandlingWorkflow, *shouldFail)
		workflowID = "ErrorHandling"
	case "parent":
		params := workflows.ChildWorkflowParams{
			Name:     *name,
			Duration: time.Duration(*duration) * time.Second,
		}
		workflowRun, err = c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.ParentWorkflow, params)
		workflowID = "Parent"
	case "signal":
		waitDuration := time.Duration(*waitTime) * time.Second
		workflowRun, err = c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.SignalWorkflow, waitDuration)
		workflowID = "Signal"

		// If workflow started successfully, send a signal to it
		if err == nil {
			log.Println("Sending signal to workflow...")
			signalData := workflows.SignalWorkflowData{
				Message: "Hello from signal sender!",
				Time:    time.Now(),
			}
			err = c.SignalWorkflow(context.Background(), workflowRun.GetID(), workflowRun.GetRunID(), "signal-channel", signalData)
			if err != nil {
				log.Fatalln("Failed to send signal", err)
			}
		}
	case "continue-as-new":
		workflowRun, err = c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.ContinueAsNewWorkflow, *count, *maxCount)
		workflowID = "ContinueAsNew"
	default:
		log.Fatalf("Unknown workflow type: %s", *workflowType)
	}

	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Printf("%s Workflow started with ID: %s and RunID: %s\n", workflowID, workflowRun.GetID(), workflowRun.GetRunID())

	// Wait for workflow completion (optional)
	var result interface{}
	err = workflowRun.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Workflow failed", err)
	}

	// Print the result based on workflow type
	switch *workflowType {
	case "greeting", "long-running", "error-handling", "parent":
		log.Printf("Workflow result: %s\n", result)
	case "sequential", "parallel", "signal":
		results, ok := result.([]string)
		if !ok {
			log.Println("Unable to parse workflow result")
			os.Exit(1)
		}
		log.Println("Workflow results:")
		for i, res := range results {
			log.Printf("  %d: %s\n", i+1, res)
		}
	case "continue-as-new":
		finalCount, ok := result.(int)
		if !ok {
			log.Println("Unable to parse workflow result")
			os.Exit(1)
		}
		log.Printf("Final count: %d\n", finalCount)
	}
}
