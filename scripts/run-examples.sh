#!/bin/bash

# Script to run all examples in sequence

# Default values
TEMPORAL_HOST=${TEMPORAL_HOST:-localhost:7233}
TEMPORAL_NAMESPACE=${TEMPORAL_NAMESPACE:-default}
NAME=${NAME:-"Temporal User"}
DURATION=${DURATION:-5}
WAIT=${WAIT:-10}

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}   Temporal Go Learning Project Examples   ${NC}"
echo -e "${BLUE}=========================================${NC}"
echo

# Check if Temporal server is running
echo -e "${YELLOW}Checking Temporal server...${NC}"
./scripts/check-temporal.sh
if [ $? -ne 0 ]; then
    echo -e "${YELLOW}Starting Temporal server...${NC}"
    make up
fi

# Start worker in the background
echo -e "${YELLOW}Starting worker in the background...${NC}"
TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/worker/main.go > /tmp/temporal-worker.log 2>&1 &
WORKER_PID=$!

# Function to cleanup worker process
cleanup() {
    echo -e "${YELLOW}Stopping worker...${NC}"
    kill $WORKER_PID 2>/dev/null
    exit 0
}

# Register the cleanup function to be called on exit
trap cleanup EXIT

# Wait for worker to start
echo -e "${YELLOW}Waiting for worker to start...${NC}"
sleep 3

# Run basic workflows
echo -e "${GREEN}Running basic workflows...${NC}"

echo -e "${YELLOW}1. Greeting Workflow${NC}"
TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/starter/main.go -workflow greeting -name "$NAME"
echo

echo -e "${YELLOW}2. Sequential Workflow${NC}"
TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/starter/main.go -workflow sequential -name "$NAME"
echo

echo -e "${YELLOW}3. Parallel Workflow${NC}"
TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/starter/main.go -workflow parallel -name "$NAME"
echo

echo -e "${YELLOW}4. Long-Running Workflow${NC}"
TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/starter/main.go -workflow long-running -duration $DURATION
echo

echo -e "${YELLOW}5. Error-Handling Workflow (Success)${NC}"
TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/starter/main.go -workflow error-handling -fail false
echo

echo -e "${YELLOW}6. Error-Handling Workflow (Failure)${NC}"
TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/starter/main.go -workflow error-handling -fail true
echo

# Run advanced workflows
echo -e "${GREEN}Running advanced workflows...${NC}"

echo -e "${YELLOW}7. Parent-Child Workflow${NC}"
TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/starter/main.go -workflow parent -name "$NAME" -duration $DURATION
echo

echo -e "${YELLOW}8. Signal Workflow${NC}"
TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/starter/main.go -workflow signal -wait $WAIT &
SIGNAL_PID=$!
sleep 2

# Get the workflow ID from the log
WORKFLOW_ID=$(grep -o "signal-workflow-[0-9]*" /tmp/temporal-worker.log | tail -1)
if [ -n "$WORKFLOW_ID" ]; then
    echo -e "${YELLOW}   Sending additional signal to $WORKFLOW_ID...${NC}"
    TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/signal/main.go -w "$WORKFLOW_ID" -action signal -message "Hello from run-examples.sh"
    
    echo -e "${YELLOW}   Querying signals from $WORKFLOW_ID...${NC}"
    TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/signal/main.go -w "$WORKFLOW_ID" -action query
fi

# Wait for signal workflow to complete
wait $SIGNAL_PID
echo

echo -e "${YELLOW}9. Continue-as-New Workflow${NC}"
TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/starter/main.go -workflow continue-as-new -count 0 -max 5
echo

# Run update workflows
echo -e "${GREEN}Running update workflows...${NC}"

echo -e "${YELLOW}10. Counter Workflow${NC}"
TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/update/main.go -workflow counter -action start -initial 0 &
sleep 2

# Get the workflow ID from the log
WORKFLOW_ID=$(grep -o "counter-workflow-[0-9]*" /tmp/temporal-worker.log | tail -1)
if [ -n "$WORKFLOW_ID" ]; then
    echo -e "${YELLOW}   Incrementing counter for $WORKFLOW_ID...${NC}"
    TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/update/main.go -workflow counter -action update -w "$WORKFLOW_ID" -update-type increment -value 5
    
    echo -e "${YELLOW}   Decrementing counter for $WORKFLOW_ID...${NC}"
    TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/update/main.go -workflow counter -action update -w "$WORKFLOW_ID" -update-type decrement -value 2
    
    echo -e "${YELLOW}   Setting counter for $WORKFLOW_ID...${NC}"
    TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/update/main.go -workflow counter -action update -w "$WORKFLOW_ID" -update-type set -value 10
    
    echo -e "${YELLOW}   Querying counter for $WORKFLOW_ID...${NC}"
    TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/update/main.go -workflow counter -action query -w "$WORKFLOW_ID" -query-type get_counter
fi
echo

echo -e "${YELLOW}11. Updateable Workflow${NC}"
TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/update/main.go -workflow updateable -action start &
sleep 2

# Get the workflow ID from the log
WORKFLOW_ID=$(grep -o "updateable-workflow-[0-9]*" /tmp/temporal-worker.log | tail -1)
if [ -n "$WORKFLOW_ID" ]; then
    echo -e "${YELLOW}   Updating state for $WORKFLOW_ID...${NC}"
    TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/update/main.go -workflow updateable -action update -w "$WORKFLOW_ID" -update-type update_state -value '{"name":"Updated Name","description":"This is an updated description","tags":["updated","workflow","example"]}'
    
    echo -e "${YELLOW}   Querying state for $WORKFLOW_ID...${NC}"
    TEMPORAL_HOST=$TEMPORAL_HOST TEMPORAL_NAMESPACE=$TEMPORAL_NAMESPACE go run cmd/update/main.go -workflow updateable -action query -w "$WORKFLOW_ID" -query-type get_state
fi
echo

echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}   All examples completed successfully!   ${NC}"
echo -e "${BLUE}=========================================${NC}"
