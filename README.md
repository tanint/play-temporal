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

This project includes a Docker Compose configuration that sets up a complete Temporal development environment with:

- MySQL (for persistence)
- Redis (for visibility storage)
- Redis Commander (Redis UI, available at http://localhost:8081)
- Temporal Server
- Temporal Web UI (available at http://localhost:8233)

To start the Temporal server and all its dependencies:

```bash
# Using Docker Compose directly
docker-compose up -d

# Or using the Makefile (recommended)
make up
```

The `make up` command will:

1. Start all services
2. Wait for them to initialize
3. Check if the Temporal server is ready to accept connections

To stop all services:

```bash
# Using Docker Compose directly
docker-compose down

# Or using the Makefile
make down
```

If you want to remove all data volumes when stopping:

```bash
# Using Docker Compose directly
docker-compose down -v

# Or using the Makefile
make down-v
```

Alternatively, you can use the Temporal CLI for a simpler setup (without persistence):

```bash
temporal server start-dev
```

### Initializing the Temporal Namespace

After starting the Temporal server, you may need to initialize a namespace:

```bash
# Initialize the default namespace
make init-namespace

# Or specify a custom namespace
TEMPORAL_NAMESPACE=your-namespace make init-namespace

# Or specify both a custom host and namespace
TEMPORAL_HOST=your-host:7233 TEMPORAL_NAMESPACE=your-namespace make init-namespace
```

This will:

1. Check if the namespace exists
2. Create it if it doesn't exist
3. Display the namespace details

### Using the Makefile

This project includes a Makefile to simplify common operations:

```bash
# Show available commands
make help

# Start all Docker Compose services
make up

# Stop all Docker Compose services
make down

# Start the worker
make worker

# Run a workflow (examples)
make greeting NAME="Your Name"
make sequential NAME="Your Name"
make long-running DURATION=10
make error-handling FAIL=false

# Run all examples in sequence
make run-examples
```

### Running All Examples

To run all examples in sequence and see Temporal in action:

```bash
# Using the Makefile
make run-examples

# Or directly with the script
./scripts/run-examples.sh

# You can customize the execution with environment variables
NAME="John Doe" DURATION=3 WAIT=15 ./scripts/run-examples.sh
```

This will:

1. Check if the Temporal server is running (and start it if needed)
2. Start a worker in the background
3. Run all basic workflows
4. Run all advanced workflows
5. Run all update workflows
6. Clean up the worker process when done

This is a great way to see all the features of Temporal in action and understand how they work together.

### Checking Status

You can check the status of all services:

```bash
# Using the Makefile
make status

# Or directly with the script
./scripts/status.sh
```

The status script will check:

1. Docker status
2. Docker Compose services
3. Temporal server
4. Temporal UI
5. Redis
6. Redis Commander
7. MySQL

This is useful for troubleshooting and ensuring all services are running correctly.

### Cleaning Up

When you're done with the project, you can clean up all resources:

```bash
# Using the Makefile
make cleanup

# Or directly with the script
./scripts/cleanup.sh
```

The cleanup script will:

1. Stop all Docker Compose services
2. Prompt you to remove data directories
3. Prompt you to remove Docker volumes
4. Prompt you to remove Docker containers
5. Prompt you to remove Docker images
6. Prompt you to remove temporary files

This ensures that no resources are left behind when you're done with the project.

### Running the Worker

Start the worker to process workflows and activities:

```bash
# Using the Makefile
make worker

# Or directly with Go
go run cmd/worker/main.go

# Or specify a custom Temporal server host
TEMPORAL_HOST=your-temporal-host:7233 go run cmd/worker/main.go

# Or specify both a custom host and namespace
TEMPORAL_HOST=your-temporal-host:7233 TEMPORAL_NAMESPACE=your-namespace go run cmd/worker/main.go
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

## Docker Compose Services

The included Docker Compose configuration provides the following services:

### MySQL

- **Purpose**: Persistence store for Temporal
- **Port**: 3306
- **Credentials**: Username: `temporal`, Password: `temporal`

### Redis

- **Purpose**: Used for visibility storage
- **Port**: 6379

### Redis Commander

- **Purpose**: Web UI for Redis
- **Port**: 8081
- **URL**: http://localhost:8081

### Temporal Server

- **Purpose**: The main Temporal service
- **Port**: 7233
- **Connection**: Use `localhost:7233` in your Temporal client

### Temporal UI

- **Purpose**: Web interface for Temporal
- **Port**: 8233
- **URL**: http://localhost:8233

## Environment Variables

The following environment variables can be used to configure the Temporal client:

- `TEMPORAL_HOST`: The host and port of the Temporal server (default: `localhost:7233`)
- `TEMPORAL_NAMESPACE`: The namespace to use (default: `default`)

## Learning Resources

- [Temporal Documentation](https://docs.temporal.io/)
- [Temporal Go SDK Documentation](https://pkg.go.dev/go.temporal.io/sdk)
- [Temporal Go SDK GitHub Repository](https://github.com/temporalio/sdk-go)
- [Temporal Samples Repository](https://github.com/temporalio/samples-go)
- [Temporal Docker Compose Repository](https://github.com/temporalio/docker-compose)
