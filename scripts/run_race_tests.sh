#!/bin/bash

# Race Condition Test Runner for Universal Go Service
# This script starts the Go server and runs comprehensive race condition tests

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color
BOLD='\033[1m'

# Configuration
SERVER_PORT=${SERVER_PORT:-3000}
SERVER_URL="http://localhost:${SERVER_PORT}"
GO_SERVER_PID=""
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Logging functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_title() {
    echo -e "${BOLD}${CYAN}🚀 $1${NC}"
}

# Cleanup function
cleanup() {
    if [ ! -z "$GO_SERVER_PID" ]; then
        log_info "Stopping Go server (PID: $GO_SERVER_PID)..."
        kill $GO_SERVER_PID 2>/dev/null || true
        wait $GO_SERVER_PID 2>/dev/null || true
        log_success "Go server stopped"
    fi
}

# Set trap for cleanup
trap cleanup EXIT INT TERM

# Check dependencies
check_dependencies() {
    log_title "Checking Dependencies"
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed or not in PATH"
        exit 1
    fi
    log_success "Go is available: $(go version)"
    
    # Check if Node.js is installed
    if ! command -v node &> /dev/null; then
        log_error "Node.js is not installed or not in PATH"
        log_info "Please install Node.js to run the race condition tests"
        exit 1
    fi
    log_success "Node.js is available: $(node --version)"
    
    # Check if the race test script exists
    if [ ! -f "$SCRIPT_DIR/race_test.js" ]; then
        log_error "race_test.js not found in $SCRIPT_DIR"
        exit 1
    fi
    log_success "Race test script found"
    
    echo ""
}

# Start the Go server
start_go_server() {
    log_title "Starting Go Server"
    
    cd "$PROJECT_ROOT"
    
    # Check if server is already running
    if curl -s "$SERVER_URL/api/v1/items?page=1&limit=1" > /dev/null 2>&1; then
        log_warning "Server is already running on port $SERVER_PORT"
        log_info "Using existing server instance"
        return 0
    fi
    
    # Build the project
    log_info "Building Go project..."
    if ! go build -o ./tmp/server ./cmd/server/; then
        log_error "Failed to build Go project"
        exit 1
    fi
    log_success "Go project built successfully"
    
    # Start the server in background
    log_info "Starting Go server on port $SERVER_PORT..."
    ./tmp/server > ./tmp/server.log 2>&1 &
    GO_SERVER_PID=$!
    
    # Wait for server to start
    log_info "Waiting for server to start..."
    for i in {1..30}; do
        if curl -s "$SERVER_URL/api/v1/items?page=1&limit=1" > /dev/null 2>&1; then
            log_success "Go server is running (PID: $GO_SERVER_PID)"
            return 0
        fi
        sleep 1
    done
    
    log_error "Failed to start Go server within 30 seconds"
    if [ ! -z "$GO_SERVER_PID" ]; then
        log_info "Server logs:"
        cat ./tmp/server.log 2>/dev/null || echo "No logs available"
    fi
    exit 1
}

# Run race condition tests
run_race_tests() {
    log_title "Running Race Condition Tests"
    
    cd "$SCRIPT_DIR"
    
    # Set environment variables
    export BASE_URL="$SERVER_URL"
    
    # Run the Node.js test script
    log_info "Executing race condition test suite..."
    echo ""
    
    if node race_test.js; then
        log_success "All race condition tests passed! 🎉"
        return 0
    else
        log_error "Some race condition tests failed"
        return 1
    fi
}

# Simple bash-based race test as fallback
run_simple_bash_race_test() {
    log_title "Running Simple Bash Race Test"
    
    local test_name="BashRaceTest_$(date +%s)"
    local success_count=0
    local conflict_count=0
    
    log_info "Testing with item name: $test_name"
    
    # Create two parallel requests
    (
        response=$(curl -s -w "%{http_code}" -o /dev/null \
            -X POST "$SERVER_URL/api/v1/items" \
            -H "Content-Type: application/json" \
            -d "{\"name\":\"$test_name\",\"amount\":10}")
        echo "Request1:$response" > ./tmp/race_result1.txt
    ) &
    
    (
        response=$(curl -s -w "%{http_code}" -o /dev/null \
            -X POST "$SERVER_URL/api/v1/items" \
            -H "Content-Type: application/json" \
            -d "{\"name\":\"$test_name\",\"amount\":20}")
        echo "Request2:$response" > ./tmp/race_result2.txt
    ) &
    
    wait
    
    # Check results
    if [ -f ./tmp/race_result1.txt ]; then
        result1=$(cat ./tmp/race_result1.txt | cut -d: -f2)
        if [ "$result1" = "201" ]; then
            ((success_count++))
            log_success "Request 1: Created (201)"
        elif [ "$result1" = "409" ] || [ "$result1" = "400" ]; then
            ((conflict_count++))
            log_success "Request 1: Conflict prevented ($result1)"
        else
            log_warning "Request 1: Unexpected status ($result1)"
        fi
    fi
    
    if [ -f ./tmp/race_result2.txt ]; then
        result2=$(cat ./tmp/race_result2.txt | cut -d: -f2)
        if [ "$result2" = "201" ]; then
            ((success_count++))
            log_success "Request 2: Created (201)"
        elif [ "$result2" = "409" ] || [ "$result2" = "400" ]; then
            ((conflict_count++))
            log_success "Request 2: Conflict prevented ($result2)"
        else
            log_warning "Request 2: Unexpected status ($result2)"
        fi
    fi
    
    # Clean up temporary files
    rm -f ./tmp/race_result1.txt ./tmp/race_result2.txt
    
    # Evaluate results
    if [ $success_count -eq 1 ] && [ $conflict_count -eq 1 ]; then
        log_success "✅ RACE CONDITION TEST PASSED: Exactly one item created"
        return 0
    elif [ $success_count -eq 2 ]; then
        log_error "❌ RACE CONDITION TEST FAILED: Both requests succeeded (race condition exists!)"
        return 1
    else
        log_warning "⚠️  Unexpected result: $success_count successes, $conflict_count conflicts"
        return 1
    fi
}

# Main execution
main() {
    log_title "🧪 Universal Go Service - Race Condition Test Runner"
    echo ""
    
    # Create tmp directory for temporary files
    mkdir -p "$PROJECT_ROOT/tmp"
    
    # Check dependencies
    check_dependencies
    
    # Start Go server
    start_go_server
    echo ""
    
    # Run tests based on available tools
    local test_passed=false
    
    if command -v node &> /dev/null && [ -f "$SCRIPT_DIR/race_test.js" ]; then
        # Run comprehensive Node.js tests
        if run_race_tests; then
            test_passed=true
        fi
    else
        log_warning "Node.js not available, running simple bash test"
        if run_simple_bash_race_test; then
            test_passed=true
        fi
    fi
    
    echo ""
    
    if [ "$test_passed" = true ]; then
        log_title "🎉 All Tests Completed Successfully!"
        log_success "Your atomic transaction implementation prevents race conditions"
        exit 0
    else
        log_title "⚠️  Some Tests Failed"
        log_error "Race conditions may still exist in your implementation"
        exit 1
    fi
}

# Show help
show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --port PORT     Set server port (default: 8080)"
    echo "  --url URL       Set server URL (default: http://localhost:8080)"
    echo "  --help          Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                    # Run tests with defaults"
    echo "  $0 --port 3000       # Use port 3000"
    echo "  $0 --url http://localhost:3000  # Use custom URL"
    echo ""
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --port)
            SERVER_PORT="$2"
            SERVER_URL="http://localhost:${SERVER_PORT}"
            shift 2
            ;;
        --url)
            SERVER_URL="$2"
            shift 2
            ;;
        --help)
            show_help
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Make script executable if run directly
if [ "${BASH_SOURCE[0]}" == "${0}" ]; then
    main "$@"
fi