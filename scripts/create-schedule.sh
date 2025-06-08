#!/bin/bash

# Check if required arguments are provided
if [ "$#" -lt 2 ]; then
    echo "Usage: $0 <subscription_id> <customer_id>"
    exit 1
fi

SUBSCRIPTION_ID=$1
CUSTOMER_ID=$2
CURRENT_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "Creating schedule for subscription $SUBSCRIPTION_ID for customer $CUSTOMER_ID"

# Create the schedule using Temporal CLI
temporal schedule create \
    --schedule-id "recurring-billing-schedule-$SUBSCRIPTION_ID" \
    --cron "0 0 1 * *" \
    --workflow-id "recurring-billing-$SUBSCRIPTION_ID" \
    --task-queue "temporal-learning-task-queue" \
    --type "RecurringBillingWorkflow" \
    --input "{\"SubscriptionID\":\"$SUBSCRIPTION_ID\",\"CustomerID\":\"$CUSTOMER_ID\",\"NextBillingDate\":\"$CURRENT_TIME\"}"

echo "Schedule created successfully!"
echo "You can now see it in the Schedules tab of the Temporal UI."
