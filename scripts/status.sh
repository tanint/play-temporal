#!/bin/bash

# Script to check the status of all Docker containers related to the project

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}   Temporal Go Learning Project Status   ${NC}"
echo -e "${BLUE}=========================================${NC}"
echo

# Check if Docker is running
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker is not installed or not in PATH.${NC}"
    exit 1
fi

# Check Docker status
echo -e "${YELLOW}Docker Status:${NC}"
docker info --format '{{.ServerVersion}}' > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo -e "${RED}Docker is not running.${NC}"
    exit 1
else
    echo -e "${GREEN}Docker is running.${NC}"
fi
echo

# Check Docker Compose
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}Error: Docker Compose is not installed or not in PATH.${NC}"
    exit 1
fi

# Check if docker-compose.yml exists
if [ ! -f "docker-compose.yml" ]; then
    echo -e "${RED}Error: docker-compose.yml not found in the current directory.${NC}"
    exit 1
fi

# Check Docker Compose services
echo -e "${YELLOW}Docker Compose Services:${NC}"
docker-compose ps
echo

# Check Temporal server
echo -e "${YELLOW}Temporal Server:${NC}"
TEMPORAL_HOST=${TEMPORAL_HOST:-localhost:7233}
nc -z $(echo $TEMPORAL_HOST | cut -d: -f1) $(echo $TEMPORAL_HOST | cut -d: -f2) &> /dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Temporal server is running at $TEMPORAL_HOST.${NC}"
else
    echo -e "${RED}Temporal server is not running at $TEMPORAL_HOST.${NC}"
fi
echo

# Check Temporal UI
echo -e "${YELLOW}Temporal UI:${NC}"
TEMPORAL_UI_PORT=${TEMPORAL_UI_PORT:-8233}
nc -z localhost $TEMPORAL_UI_PORT &> /dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Temporal UI is running at http://localhost:$TEMPORAL_UI_PORT.${NC}"
else
    echo -e "${RED}Temporal UI is not running at http://localhost:$TEMPORAL_UI_PORT.${NC}"
fi
echo

# Check Redis
echo -e "${YELLOW}Redis:${NC}"
REDIS_PORT=${REDIS_PORT:-6379}
nc -z localhost $REDIS_PORT &> /dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Redis is running at localhost:$REDIS_PORT.${NC}"
else
    echo -e "${RED}Redis is not running at localhost:$REDIS_PORT.${NC}"
fi
echo

# Check Redis Commander
echo -e "${YELLOW}Redis Commander:${NC}"
REDIS_UI_PORT=${REDIS_UI_PORT:-8081}
nc -z localhost $REDIS_UI_PORT &> /dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Redis Commander is running at http://localhost:$REDIS_UI_PORT.${NC}"
else
    echo -e "${RED}Redis Commander is not running at http://localhost:$REDIS_UI_PORT.${NC}"
fi
echo

# Check MySQL
echo -e "${YELLOW}MySQL:${NC}"
nc -z localhost 3306 &> /dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}MySQL is running at localhost:3306.${NC}"
else
    echo -e "${RED}MySQL is not running at localhost:3306.${NC}"
fi
echo

echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}   Status check completed!   ${NC}"
echo -e "${BLUE}=========================================${NC}"
