package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/temporalio/play-temporal/workflows"
	"go.temporal.io/sdk/client"
)

func main() {
	// Define command line flags
	workflowType := flag.String("workflow", "counter", "Workflow type to run (counter, updateable)")
	action := flag.String("action", "start", "Action to perform: start, update, query")
	workflowID := flag.String("w", "", "Workflow ID (required for update and query)")
	runID := flag.String("r", "", "Run ID (optional for update and query)")
	initialValue := flag.Int("initial", 0, "Initial value for counter workflow")
	updateType := flag.String("update-type", "", "Update type (increment, decrement, set for counter; update_state for updateable)")
	updateValue := flag.String("value", "", "Update value (amount for counter; JSON for updateable)")
	queryType := flag.String("query-type", "", "Query type (get_counter for counter; get_state for updateable)")
	flag.Parse()

	// Create the client object
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Perform the requested action
	switch *action {
	case "start":
		startWorkflow(c, *workflowType, *initialValue)
	case "update":
		if *workflowID == "" {
			log.Fatalln("Workflow ID is required for update. Use -w flag.")
		}
		if *updateType == "" {
			log.Fatalln("Update type is required. Use -update-type flag.")
		}
		updateWorkflow(c, *workflowType, *workflowID, *runID, *updateType, *updateValue)
	case "query":
		if *workflowID == "" {
			log.Fatalln("Workflow ID is required for query. Use -w flag.")
		}
		if *queryType == "" {
			log.Fatalln("Query type is required. Use -query-type flag.")
		}
		queryWorkflow(c, *workflowType, *workflowID, *runID, *queryType)
	default:
		log.Fatalf("Unknown action: %s. Use 'start', 'update', or 'query'.", *action)
	}
}

func startWorkflow(c client.Client, workflowType string, initialValue int) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("%s-workflow-%v", workflowType, time.Now().Unix()),
		TaskQueue: "temporal-learning-task-queue",
	}

	var workflowRun client.WorkflowRun
	var err error

	switch workflowType {
	case "counter":
		workflowRun, err = c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.CounterWorkflow, initialValue)
	case "updateable":
		workflowRun, err = c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.UpdateableWorkflow)
	default:
		log.Fatalf("Unknown workflow type: %s", workflowType)
	}

	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Printf("Workflow started with ID: %s and RunID: %s\n", workflowRun.GetID(), workflowRun.GetRunID())
}

func updateWorkflow(c client.Client, workflowType, workflowID, runID, updateType, updateValue string) {
	switch workflowType {
	case "counter":
		updateCounterWorkflow(c, workflowID, runID, updateType, updateValue)
	case "updateable":
		updateUpdateableWorkflow(c, workflowID, runID, updateValue)
	default:
		log.Fatalf("Unknown workflow type: %s", workflowType)
	}
}

func updateCounterWorkflow(c client.Client, workflowID, runID, updateType, updateValue string) {
	// Parse the update value as an integer
	amount, err := strconv.Atoi(updateValue)
	if err != nil {
		log.Fatalf("Invalid update value: %s. Must be an integer.", updateValue)
	}

	// Create update options
	updateOptions := client.UpdateWorkflowOptions{
		WorkflowID: workflowID,
		RunID:      runID,
		UpdateName: updateType,
		Args:       []interface{}{amount},
	}

	// Send the update
	resp, err := c.UpdateWorkflow(context.Background(), updateOptions)
	if err != nil {
		log.Fatalln("Failed to update workflow", err)
	}

	// Get the update result
	var result int
	if err := resp.Get(context.Background(), &result); err != nil {
		log.Fatalln("Failed to get update result", err)
	}

	log.Printf("Update successful. New counter value: %d\n", result)
}

func updateUpdateableWorkflow(c client.Client, workflowID, runID, updateValue string) {
	// Parse the update value as a JSON object
	var updates map[string]interface{}
	if err := json.Unmarshal([]byte(updateValue), &updates); err != nil {
		log.Fatalf("Invalid update value: %s. Must be a valid JSON object.", updateValue)
	}

	// Create update options
	updateOptions := client.UpdateWorkflowOptions{
		WorkflowID: workflowID,
		RunID:      runID,
		UpdateName: "update_state",
		Args:       []interface{}{updates},
	}

	// Send the update
	resp, err := c.UpdateWorkflow(context.Background(), updateOptions)
	if err != nil {
		log.Fatalln("Failed to update workflow", err)
	}

	// Get the update result
	var result map[string]interface{}
	if err := resp.Get(context.Background(), &result); err != nil {
		log.Fatalln("Failed to get update result", err)
	}

	// Print the update result
	log.Println("Update successful. New state:")
	printMap(result, "  ")
}

func queryWorkflow(c client.Client, workflowType, workflowID, runID, queryType string) {
	// Send the query
	resp, err := c.QueryWorkflow(context.Background(), workflowID, runID, queryType)
	if err != nil {
		log.Fatalln("Failed to query workflow", err)
	}

	// Process the query result based on workflow type
	switch workflowType {
	case "counter":
		var result int
		if err := resp.Get(&result); err != nil {
			log.Fatalln("Failed to decode query result", err)
		}
		log.Printf("Current counter value: %d\n", result)

	case "updateable":
		var result map[string]interface{}
		if err := resp.Get(&result); err != nil {
			log.Fatalln("Failed to decode query result", err)
		}
		log.Println("Current state:")
		printMap(result, "  ")

	default:
		log.Fatalf("Unknown workflow type: %s", workflowType)
	}
}

// Helper function to print a map with indentation
func printMap(m map[string]interface{}, indent string) {
	for k, v := range m {
		switch val := v.(type) {
		case map[string]interface{}:
			log.Printf("%s%s:\n", indent, k)
			printMap(val, indent+"  ")
		case []interface{}:
			log.Printf("%s%s: [%s]\n", indent, k, strings.Join(interfaceSliceToStringSlice(val), ", "))
		default:
			log.Printf("%s%s: %v\n", indent, k, v)
		}
	}
}

// Helper function to convert a slice of interface{} to a slice of strings
func interfaceSliceToStringSlice(slice []interface{}) []string {
	result := make([]string, len(slice))
	for i, v := range slice {
		result[i] = fmt.Sprintf("%v", v)
	}
	return result
}
