# Temporal Workflow Examples

This repository contains examples of Temporal workflows implemented in Go. It demonstrates various features and patterns of the Temporal workflow engine.

## Setup

### Prerequisites

- Go 1.16 or later
- Docker and Docker Compose

### Running the Project

1. Start the Temporal server and dependencies:

```bash
make up
```

2. Start the worker:

```bash
make worker
```

3. Run example workflows (see below)

4. Access the Temporal UI at http://localhost:8233

5. Stop the services:

```bash
make down
```

## Basic Workflows

### Greeting Workflow

A simple workflow that calls a single activity to generate a greeting message.

```bash
make greeting NAME="Your Name"
```

**Key concepts:**

- Basic workflow structure
- Executing a single activity
- Passing parameters to workflows and activities

### Sequential Workflow

Demonstrates calling multiple activities in sequence, waiting for each to complete before starting the next.

```bash
make sequential NAME="Your Name"
```

**Key concepts:**

- Activity options configuration
- Retry policies
- Sequential activity execution
- Error handling

### Parallel Workflow

Shows how to execute multiple activities in parallel and wait for all to complete.

```bash
make parallel NAME="Your Name"
```

**Key concepts:**

- Parallel activity execution
- Futures and promises
- Waiting for multiple activities

### Long-Running Workflow

Demonstrates a workflow with a long-running activity that uses heartbeating.

```bash
make long-running DURATION=10
```

**Key concepts:**

- Long-running activities
- Heartbeat timeouts
- Activity cancellation handling

### Error Handling Workflow

Shows how to handle errors in activities and implement retry policies.

```bash
make error-handling FAIL=true|false
```

**Key concepts:**

- Error handling patterns
- Retry policies
- Activity failures

## Advanced Workflows

### Parent-Child Workflow

Demonstrates how to create and manage child workflows.

```bash
make parent NAME="Your Name" DURATION=5
```

**Key concepts:**

- Child workflow execution
- Parent-child workflow relationship
- Child workflow options
- Data passing between parent and child

### Signal Workflow

Shows how to use signals to communicate with running workflows.

```bash
make signal WAIT=60
```

To send a signal to a running workflow:

```bash
make send-signal WORKFLOW_ID="id" MESSAGE="Your message"
```

To query signals received by a workflow:

```bash
make query-signals WORKFLOW_ID="id"
```

**Key concepts:**

- Workflow signals
- Signal channels
- Query handlers
- Long-running workflows
- Selectors for handling multiple events

### Continue-as-New Workflow

Demonstrates the continue-as-new feature to handle long-running workflows.

```bash
make continue-as-new COUNT=0 MAX=10
```

**Key concepts:**

- Continue-as-new pattern
- Workflow history size management
- State passing between workflow runs

## Update Workflows

### Counter Workflow

Demonstrates the update feature with a simple counter.

```bash
make start-counter INITIAL=0
```

To update the counter:

```bash
make increment-counter WORKFLOW_ID="id" VALUE=5
make decrement-counter WORKFLOW_ID="id" VALUE=2
make set-counter WORKFLOW_ID="id" VALUE=10
```

To query the counter value:

```bash
make query-counter WORKFLOW_ID="id"
```

**Key concepts:**

- Update handlers
- Query handlers
- Workflow state management
- Safe state updates

### Updateable Workflow

Shows a more complex update scenario with structured data.

```bash
make start-updateable
```

To update the state:

```bash
make update-state WORKFLOW_ID="id" VALUE='{"name":"New Name","description":"Updated description","tags":["updated","workflow"]}'
```

To query the state:

```bash
make query-state WORKFLOW_ID="id"
```

**Key concepts:**

- Complex state management
- Structured data updates
- Dynamic updates with maps
- Multiple update handlers

## Subscription Workflow System

Simulates a subscription billing system similar to Stripe billing with monthly recurring charges.

### Creating a Subscription

```bash
make subscription CUSTOMER="customer123" PLAN="premium-monthly"
```

**Key concepts:**

- Business process orchestration
- Multi-step activity sequences
- State management
- Error handling with continuation

**Workflow steps:**

1. Create subscription
2. Calculate initial charges
3. Generate invoice
4. Process payment
5. Send invoice email
6. Update subscription status

### Recurring Billing

There are two ways to set up recurring billing:

#### 1. Using CronSchedule in workflow options

```bash
make recurring-billing SUBSCRIPTION="sub_123456" CUSTOMER="customer123"
```

This starts a workflow with a cron schedule that will run monthly. However, this won't appear in the Schedules tab of the Temporal UI.

#### 2. Using Temporal Schedules (recommended)

```bash
make create-schedule SUBSCRIPTION="sub_123456" CUSTOMER="customer123"
```

This creates a schedule using the Temporal CLI that will be visible in the Schedules tab of the Temporal UI.

**Key concepts:**

- Cron-scheduled workflows
- Temporal Schedules
- Independent workflow execution
- Time-based scheduling
- Simulated billing cycles

**Workflow steps:**

1. Calculate charges for the billing period
2. Generate invoice
3. Process payment
4. Send invoice email
5. Update subscription status

The recurring billing workflow runs monthly using either Temporal's CronSchedule feature or the Schedules feature, ensuring reliable execution of billing cycles even after system restarts. Using the Schedules feature provides better visibility and management through the Temporal UI.

## Best Practices Demonstrated

1. **Activity Options**: All workflows set appropriate timeouts for activities
2. **Error Handling**: Proper error handling and logging throughout workflows
3. **Retry Policies**: Configuration of retry policies for activities
4. **Workflow Patterns**: Various workflow patterns (sequential, parallel, parent-child)
5. **Signals and Queries**: Communication with running workflows
6. **Updates**: Safe state updates in long-running workflows
7. **Continue-as-New**: Managing workflow history size
8. **Business Process Modeling**: Subscription workflow demonstrates modeling real-world business processes

## Temporal Concepts Covered

- Workflows and Activities
- Task Queues
- Timeouts and Retries
- Child Workflows
- Signals and Queries
- Updates
- Continue-as-New
- Error Handling
- Workflow Scheduling

## Project Structure

- `cmd/worker/main.go`: Worker implementation
- `cmd/starter/main.go`: Workflow starter
- `cmd/signal/main.go`: Signal sender and query handler
- `cmd/update/main.go`: Update sender and query handler
- `cmd/subscription/main.go`: Subscription workflow starter
- `cmd/billing/main.go`: Recurring billing workflow starter
- `workflows/workflows.go`: Basic workflow implementations
- `workflows/advanced_workflows.go`: Advanced workflow implementations
- `workflows/update_workflows.go`: Update workflow implementations
- `workflows/subscription_workflows.go`: Subscription workflow implementations
- `activities/activities.go`: Activity implementations
- `activities/subscription_activities.go`: Subscription activity implementations
- `config/config.go`: Configuration utilities
- `docker-compose.yml`: Docker Compose configuration for Temporal server
