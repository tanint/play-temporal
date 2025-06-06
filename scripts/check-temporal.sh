#!/bin/bash

# Script to check if Temporal server is running and ready

# Default values
TEMPORAL_HOST=${TEMPORAL_HOST:-localhost:7233}
MAX_RETRIES=${MAX_RETRIES:-30}
RETRY_INTERVAL=${RETRY_INTERVAL:-2}

echo "Checking Temporal server at $TEMPORAL_HOST..."

# Function to check if Temporal server is ready
check_temporal() {
  # Use netcat to check if the port is open
  if command -v nc &> /dev/null; then
    nc -z $(echo $TEMPORAL_HOST | cut -d: -f1) $(echo $TEMPORAL_HOST | cut -d: -f2) &> /dev/null
    return $?
  fi

  # Fallback to using curl if netcat is not available
  if command -v curl &> /dev/null; then
    curl -s -o /dev/null -w "%{http_code}" $TEMPORAL_HOST &> /dev/null
    if [ $? -eq 0 ]; then
      return 0
    else
      return 1
    fi
  fi

  # If neither netcat nor curl is available, use timeout with a connection attempt
  if command -v timeout &> /dev/null; then
    timeout 1 bash -c "</dev/tcp/$(echo $TEMPORAL_HOST | cut -d: -f1)/$(echo $TEMPORAL_HOST | cut -d: -f2)" &> /dev/null
    return $?
  fi

  echo "Error: Neither nc, curl, nor timeout commands are available. Cannot check Temporal server."
  exit 1
}

# Try to connect to Temporal server
retry_count=0
while [ $retry_count -lt $MAX_RETRIES ]; do
  if check_temporal; then
    echo "Temporal server is running and ready!"
    exit 0
  fi

  retry_count=$((retry_count+1))
  echo "Waiting for Temporal server to be ready... ($retry_count/$MAX_RETRIES)"
  sleep $RETRY_INTERVAL
done

echo "Error: Temporal server is not ready after $MAX_RETRIES attempts."
exit 1
