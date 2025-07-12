#!/bin/bash

# Script to run tests with Docker and build application only if tests pass
# Usage: ./run_docker_tests.sh

set -e  # Exit on any error

echo "🧪 Starting Docker tests for Go Users API..."
echo "=============================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

# Stop any existing containers
print_status "Stopping existing containers..."
docker compose -f docker-compose.test.yml down --remove-orphans 2>/dev/null || true

# Clean up any existing images
print_status "Cleaning up existing images..."
docker rmi $(docker images -q go-users-api) 2>/dev/null || true

# Start MongoDB for tests
print_status "Starting MongoDB for tests..."
docker compose -f docker-compose.test.yml up -d mongodb

# Wait for MongoDB to be ready
print_status "Waiting for MongoDB to be ready..."
sleep 10

# Run tests
print_status "Running tests..."
if docker compose -f docker-compose.test.yml run --rm test; then
    print_success "✅ All tests passed!"
    
    # Build the application
    print_status "Building application..."
    if docker compose -f docker-compose.test.yml build api; then
        print_success "✅ Application built successfully!"
        
        # Start the application
        print_status "Starting application..."
        docker compose -f docker-compose.test.yml up -d api mongo-express
        
        print_success "🎉 Application is now running!"
        echo ""
        echo "📋 Service URLs:"
        echo "   - API: http://localhost:8080"
        echo "   - Swagger Docs: http://localhost:8080/swagger/index.html"
        echo "   - MongoDB Express: http://localhost:8081"
        echo ""
        echo "📊 Test data has been loaded with 10 users in MongoDB"
        echo "🔧 To view logs: docker compose -f docker-compose.test.yml logs -f"
        echo "🛑 To stop: docker compose -f docker-compose.test.yml down"
        
    else
        print_error "❌ Failed to build application"
        exit 1
    fi
else
    print_error "❌ Tests failed! Application will not be built."
    print_status "Stopping test containers..."
    docker compose -f docker-compose.test.yml down
    exit 1
fi 