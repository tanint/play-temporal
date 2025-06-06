package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/temporalio/play-temporal/config"
	"github.com/temporalio/play-temporal/workflows"
	"go.temporal.io/sdk/client"
)

func main() {
	// Define command line flags
	workflowID := flag.String("w", "", "Workflow ID to send signal to")
	runID := flag.String("r", "", "Run ID of the workflow (optional)")
	action := flag.String("action", "signal", "Action to perform: signal, query")
	message := flag.String("message", "Signal from command line", "Message to send in the signal")
	flag.Parse()

	if *workflowID == "" {
		log.Fatalln("Workflow ID is required. Use -w flag to specify it.")
	}

	// Create the client object
	c, err := client.Dial(config.GetTemporalClientOptions())
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Perform the requested action
	switch *action {
	case "signal":
		// Create signal data
		signalData := workflows.SignalWorkflowData{
			Message: *message,
			Time:    time.Now(),
		}

		// Send signal to the workflow
		err = c.SignalWorkflow(context.Background(), *workflowID, *runID, "signal-channel", signalData)
		if err != nil {
			log.Fatalln("Failed to send signal", err)
		}
		log.Println("Signal sent successfully")

	case "query":
		// Query the workflow
		resp, err := c.QueryWorkflow(context.Background(), *workflowID, *runID, "get_signals")
		if err != nil {
			log.Fatalln("Failed to query workflow", err)
		}

		// Decode the query result
		var signals []workflows.SignalWorkflowData
		if err := resp.Get(&signals); err != nil {
			log.Fatalln("Failed to decode query result", err)
		}

		// Print the signals received by the workflow
		log.Printf("Workflow has received %d signals:\n", len(signals))
		for i, signal := range signals {
			log.Printf("  %d: %s (received at %v)\n", i+1, signal.Message, signal.Time.Format(time.RFC3339))
		}

	default:
		log.Fatalf("Unknown action: %s. Use 'signal' or 'query'.", *action)
	}
}
