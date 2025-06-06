#!/bin/bash

# Script to clean up the project

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}   Temporal Go Learning Project Cleanup   ${NC}"
echo -e "${BLUE}=========================================${NC}"
echo

# Function to confirm action
confirm() {
    read -p "$(echo -e "${YELLOW}$1 (y/n): ${NC}")" -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        return 1
    fi
    return 0
}

# Stop Docker Compose services
if [ -f "docker-compose.yml" ]; then
    echo -e "${YELLOW}Stopping Docker Compose services...${NC}"
    docker-compose down
else
    echo -e "${RED}docker-compose.yml not found. Skipping Docker Compose cleanup.${NC}"
fi

# Remove data directories
if confirm "Do you want to remove all data directories (data/mysql, data/redis)?"; then
    echo -e "${YELLOW}Removing data directories...${NC}"
    rm -rf data/mysql data/redis
    echo -e "${GREEN}Data directories removed.${NC}"
fi

# Remove Docker volumes
if confirm "Do you want to remove all Docker volumes related to this project?"; then
    echo -e "${YELLOW}Removing Docker volumes...${NC}"
    docker volume rm $(docker volume ls -q | grep -E 'play-temporal|temporal-mysql|temporal-redis') 2>/dev/null || true
    echo -e "${GREEN}Docker volumes removed.${NC}"
fi

# Remove Docker containers
if confirm "Do you want to remove all Docker containers related to this project?"; then
    echo -e "${YELLOW}Removing Docker containers...${NC}"
    docker rm -f $(docker ps -a -q --filter "name=temporal") 2>/dev/null || true
    echo -e "${GREEN}Docker containers removed.${NC}"
fi

# Remove Docker images
if confirm "Do you want to remove all Docker images related to Temporal?"; then
    echo -e "${YELLOW}Removing Docker images...${NC}"
    docker rmi $(docker images | grep -E 'temporal|mysql|redis' | awk '{print $3}') 2>/dev/null || true
    echo -e "${GREEN}Docker images removed.${NC}"
fi

# Remove temporary files
if confirm "Do you want to remove temporary files (logs, etc.)?"; then
    echo -e "${YELLOW}Removing temporary files...${NC}"
    rm -f /tmp/temporal-worker.log
    echo -e "${GREEN}Temporary files removed.${NC}"
fi

echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}   Cleanup completed successfully!   ${NC}"
echo -e "${BLUE}=========================================${NC}"
