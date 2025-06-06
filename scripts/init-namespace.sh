#!/bin/bash

# Script to initialize Temporal namespace

# Default values
TEMPORAL_HOST=${TEMPORAL_HOST:-localhost:7233}
TEMPORAL_NAMESPACE=${TEMPORAL_NAMESPACE:-default}
RETENTION_DAYS=${RETENTION_DAYS:-7}

echo "Initializing Temporal namespace '$TEMPORAL_NAMESPACE' at $TEMPORAL_HOST..."

# Check if tctl is installed
if ! command -v tctl &> /dev/null; then
    echo "Error: tctl is not installed. Please install it first."
    echo "Installation instructions: https://docs.temporal.io/tctl/install"
    exit 1
fi

# Check if namespace exists
if tctl --address $TEMPORAL_HOST namespace describe $TEMPORAL_NAMESPACE &> /dev/null; then
    echo "Namespace '$TEMPORAL_NAMESPACE' already exists."
else
    # Create namespace
    echo "Creating namespace '$TEMPORAL_NAMESPACE'..."
    tctl --address $TEMPORAL_HOST namespace register \
        --retention $RETENTION_DAYS \
        --description "Namespace for Temporal Go Learning Project" \
        $TEMPORAL_NAMESPACE

    if [ $? -eq 0 ]; then
        echo "Namespace '$TEMPORAL_NAMESPACE' created successfully."
    else
        echo "Error: Failed to create namespace '$TEMPORAL_NAMESPACE'."
        exit 1
    fi
fi

# Check if namespace is active
status=$(tctl --address $TEMPORAL_HOST namespace describe $TEMPORAL_NAMESPACE | grep "State:" | awk '{print $2}')
if [ "$status" == "REGISTERED" ]; then
    echo "Namespace '$TEMPORAL_NAMESPACE' is active."
else
    echo "Warning: Namespace '$TEMPORAL_NAMESPACE' is not active. Status: $status"
fi

# Display namespace details
echo "Namespace details:"
tctl --address $TEMPORAL_HOST namespace describe $TEMPORAL_NAMESPACE

echo "Namespace initialization complete."
exit 0
