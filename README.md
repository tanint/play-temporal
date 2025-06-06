# Temporal Go Learning Project

This project is a comprehensive Go application that demonstrates how to use Temporal for workflow orchestration. It includes various examples of workflows and activities to help you learn Temporal's features.

## Project Structure

```
.
├── README.md
├── activities/
│   └── activities.go
├── cmd/
│   ├── signal/
│   │   └── main.go
│   ├── starter/
│   │   └── main.go
│   ├── update/
│   │   └── main.go
│   └── worker/
│       └── main.go
├── go.mod
├── go.sum
└── workflows/
    ├── advanced_workflows.go
    ├── update_workflows.go
    └── workflows.go
```

## Getting Started

### Prerequisites

- Go 1.16 or later
- Temporal server running locally or accessible remotely

### Running the Temporal Server Locally

You can run the Temporal server locally using Docker:

```bash
docker run --rm -p 7233:7233 -p 8233:8233 temporalio/temporal:latest-auto-setup
```

Or using the Temporal CLI:

```bash
temporal server start-dev
```

### Running the Worker

Start the worker to process workflows and activities:

```bash
go run cmd/worker/main.go
```

## Basic Workflows

Run basic workflows with the starter command:

```bash
# Run the greeting workflow
go run cmd/starter/main.go -workflow greeting -name "Your Name"

# Run the sequential workflow
go run cmd/starter/main.go -workflow sequential -name "Your Name"

# Run the parallel workflow
go run cmd/starter/main.go -workflow parallel -name "Your Name"

# Run the long-running workflow
go run cmd/starter/main.go -workflow long-running -duration 10

# Run the error-handling workflow (success)
go run cmd/starter/main.go -workflow error-handling -fail false

# Run the error-handling workflow (failure)
go run cmd/starter/main.go -workflow error-handling -fail true
```

## Advanced Workflows

### Parent-Child Workflow

Run a parent workflow that executes a child workflow:

```bash
go run cmd/starter/main.go -workflow parent -name "Your Name" -duration 5
```

### Signal Workflow

Run a workflow that can receive signals:

```bash
# Start the signal workflow (it will automatically send one signal)
go run cmd/starter/main.go -workflow signal -wait 60

# Send additional signals to the running workflow
# (Note the workflow ID from the previous command)
go run cmd/signal/main.go -w "signal-workflow-123456789" -action signal -message "Hello from CLI"

# Query the workflow to see received signals
go run cmd/signal/main.go -w "signal-workflow-123456789" -action query
```

### Continue-as-New Workflow

Run a workflow that demonstrates the continue-as-new feature:

```bash
go run cmd/starter/main.go -workflow continue-as-new -count 0 -max 10
```

## Update Workflows

### Counter Workflow

Run and interact with a counter workflow:

```bash
# Start the counter workflow with initial value 0
go run cmd/update/main.go -workflow counter -action start -initial 0

# Increment the counter (replace with your workflow ID)
go run cmd/update/main.go -workflow counter -action update -w "counter-workflow-123456789" -update-type increment -value 5

# Decrement the counter
go run cmd/update/main.go -workflow counter -action update -w "counter-workflow-123456789" -update-type decrement -value 2

# Set the counter to a specific value
go run cmd/update/main.go -workflow counter -action update -w "counter-workflow-123456789" -update-type set -value 10

# Query the current counter value
go run cmd/update/main.go -workflow counter -action query -w "counter-workflow-123456789" -query-type get_counter
```

### Updateable Workflow

Run and interact with a more complex updateable workflow:

```bash
# Start the updateable workflow
go run cmd/update/main.go -workflow updateable -action start

# Update the workflow state (replace with your workflow ID)
go run cmd/update/main.go -workflow updateable -action update -w "updateable-workflow-123456789" -update-type update_state -value '{"name":"Updated Name","description":"This is an updated description","tags":["updated","workflow","example"]}'

# Query the current workflow state
go run cmd/update/main.go -workflow updateable -action query -w "updateable-workflow-123456789" -query-type get_state
```

## Learning Resources

- [Temporal Documentation](https://docs.temporal.io/)
- [Temporal Go SDK Documentation](https://pkg.go.dev/go.temporal.io/sdk)
- [Temporal Go SDK GitHub Repository](https://github.com/temporalio/sdk-go)
- [Temporal Samples Repository](https://github.com/temporalio/samples-go)
