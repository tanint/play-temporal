# Makefile for Temporal Go Learning Project

# Variables
TEMPORAL_HOST ?= localhost:7233
TEMPORAL_NAMESPACE ?= default
TASK_QUEUE ?= temporal-learning-task-queue

# Docker Compose commands
.PHONY: up
up:
	docker-compose up -d
	@echo "Waiting for services to start..."
	@sleep 5
	@./scripts/check-temporal.sh

.PHONY: down
down:
	docker-compose down

.PHONY: down-v
down-v:
	docker-compose down -v

.PHONY: logs
logs:
	docker-compose logs -f

.PHONY: status
status:
	./scripts/status.sh

.PHONY: check-temporal
check-temporal:
	./scripts/check-temporal.sh

.PHONY: init-namespace
init-namespace:
	./scripts/init-namespace.sh

.PHONY: run-examples
run-examples:
	./scripts/run-examples.sh

.PHONY: cleanup
cleanup:
	./scripts/cleanup.sh

# Worker commands
.PHONY: worker
worker:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/worker/main.go

# Workflow commands
.PHONY: greeting
greeting:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/starter/main.go -workflow greeting -name "$(NAME)"

.PHONY: sequential
sequential:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/starter/main.go -workflow sequential -name "$(NAME)"

.PHONY: parallel
parallel:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/starter/main.go -workflow parallel -name "$(NAME)"

.PHONY: long-running
long-running:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/starter/main.go -workflow long-running -duration $(DURATION)

.PHONY: error-handling
error-handling:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/starter/main.go -workflow error-handling -fail $(FAIL)

.PHONY: parent
parent:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/starter/main.go -workflow parent -name "$(NAME)" -duration $(DURATION)

.PHONY: signal
signal:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/starter/main.go -workflow signal -wait $(WAIT)

.PHONY: continue-as-new
continue-as-new:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/starter/main.go -workflow continue-as-new -count $(COUNT) -max $(MAX)

# Signal commands
.PHONY: send-signal
send-signal:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/signal/main.go -w "$(WORKFLOW_ID)" -action signal -message "$(MESSAGE)"

.PHONY: query-signals
query-signals:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/signal/main.go -w "$(WORKFLOW_ID)" -action query

# Update commands
.PHONY: start-counter
start-counter:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/update/main.go -workflow counter -action start -initial $(INITIAL)

.PHONY: increment-counter
increment-counter:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/update/main.go -workflow counter -action update -w "$(WORKFLOW_ID)" -update-type increment -value $(VALUE)

.PHONY: decrement-counter
decrement-counter:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/update/main.go -workflow counter -action update -w "$(WORKFLOW_ID)" -update-type decrement -value $(VALUE)

.PHONY: set-counter
set-counter:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/update/main.go -workflow counter -action update -w "$(WORKFLOW_ID)" -update-type set -value $(VALUE)

.PHONY: query-counter
query-counter:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/update/main.go -workflow counter -action query -w "$(WORKFLOW_ID)" -query-type get_counter

.PHONY: start-updateable
start-updateable:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/update/main.go -workflow updateable -action start

.PHONY: update-state
update-state:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/update/main.go -workflow updateable -action update -w "$(WORKFLOW_ID)" -update-type update_state -value '$(VALUE)'

.PHONY: query-state
query-state:
	TEMPORAL_HOST=$(TEMPORAL_HOST) TEMPORAL_NAMESPACE=$(TEMPORAL_NAMESPACE) go run cmd/update/main.go -workflow updateable -action query -w "$(WORKFLOW_ID)" -query-type get_state

# Help
.PHONY: help
help:
	@echo "Temporal Go Learning Project Makefile"
	@echo ""
	@echo "Docker Compose Commands:"
	@echo "  make up              Start all services"
	@echo "  make down            Stop all services"
	@echo "  make down-v          Stop all services and remove volumes"
	@echo "  make logs            Show logs from all services"
	@echo "  make status          Show status of all services"
	@echo "  make check-temporal  Check if Temporal server is running"
	@echo "  make init-namespace  Initialize Temporal namespace"
	@echo "  make run-examples    Run all example workflows in sequence"
	@echo "  make cleanup         Clean up the project (remove data, containers, etc.)"
	@echo ""
	@echo "Worker Commands:"
	@echo "  make worker          Start the worker"
	@echo ""
	@echo "Workflow Commands:"
	@echo "  make greeting NAME=\"Your Name\"                  Run greeting workflow"
	@echo "  make sequential NAME=\"Your Name\"                Run sequential workflow"
	@echo "  make parallel NAME=\"Your Name\"                  Run parallel workflow"
	@echo "  make long-running DURATION=10                     Run long-running workflow"
	@echo "  make error-handling FAIL=true|false               Run error-handling workflow"
	@echo "  make parent NAME=\"Your Name\" DURATION=5         Run parent-child workflow"
	@echo "  make signal WAIT=60                               Run signal workflow"
	@echo "  make continue-as-new COUNT=0 MAX=10               Run continue-as-new workflow"
	@echo ""
	@echo "Signal Commands:"
	@echo "  make send-signal WORKFLOW_ID=\"id\" MESSAGE=\"msg\"  Send signal to workflow"
	@echo "  make query-signals WORKFLOW_ID=\"id\"               Query signals from workflow"
	@echo ""
	@echo "Update Commands:"
	@echo "  make start-counter INITIAL=0                      Start counter workflow"
	@echo "  make increment-counter WORKFLOW_ID=\"id\" VALUE=5   Increment counter"
	@echo "  make decrement-counter WORKFLOW_ID=\"id\" VALUE=2   Decrement counter"
	@echo "  make set-counter WORKFLOW_ID=\"id\" VALUE=10        Set counter value"
	@echo "  make query-counter WORKFLOW_ID=\"id\"               Query counter value"
	@echo "  make start-updateable                             Start updateable workflow"
	@echo "  make update-state WORKFLOW_ID=\"id\" VALUE='{...}'  Update workflow state"
	@echo "  make query-state WORKFLOW_ID=\"id\"                 Query workflow state"
	@echo ""
	@echo "Environment Variables:"
	@echo "  TEMPORAL_HOST         Temporal server host (default: localhost:7233)"
	@echo "  TEMPORAL_NAMESPACE    Temporal namespace (default: default)"
	@echo "  TASK_QUEUE           Task queue name (default: temporal-learning-task-queue)"
