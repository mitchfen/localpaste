#!/bin/bash

# QuickPaste Test Runner
# Simple script to run tests with various options

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

show_usage() {
    cat << EOF
QuickPaste Test Runner

Usage: ./test.sh [options]

Options:
    -v, --verbose     Run tests with verbose output
    -c, --coverage    Run tests with coverage report
    -q, --quiet       Run tests quietly (only show failures)
    -h, --help        Show this help message

Examples:
    ./test.sh              # Run all tests
    ./test.sh -v           # Run with verbose output
    ./test.sh -c           # Run with coverage report
    ./test.sh -v -c        # Run with verbose output and coverage

EOF
}

# Default options
VERBOSE=""
COVERAGE=""
QUIET=""

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE="-v"
            shift
            ;;
        -c|--coverage)
            COVERAGE="true"
            shift
            ;;
        -q|--quiet)
            QUIET="-q"
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Navigate to project root
cd "$(dirname "$0")"

echo -e "${YELLOW}Running QuickPaste Tests${NC}"
echo "========================"

# Run tests
if [ "$COVERAGE" = "true" ]; then
    echo -e "\n${YELLOW}Running tests with coverage...${NC}"
    go test ./internal/... -coverprofile=coverage.out $VERBOSE $QUIET
    
    if [ $? -eq 0 ]; then
        echo -e "\n${GREEN}✓ All tests passed!${NC}"
        
        # Show coverage summary
        COVERAGE_PCT=$(go tool cover -func=coverage.out 2>/dev/null | grep total | awk '{print $3}' || echo "N/A")
        echo -e "${GREEN}Coverage: ${COVERAGE_PCT}${NC}"
        
        # Optionally open coverage report in browser (if -v flag)
        if [ -n "$VERBOSE" ]; then
            echo -e "\n${YELLOW}Generating HTML coverage report...${NC}"
            go tool cover -html=coverage.out -o coverage.html
            echo -e "${GREEN}Coverage report saved to coverage.html${NC}"
        fi
    else
        echo -e "\n${RED}✗ Tests failed!${NC}"
        exit 1
    fi
else
    go test ./internal/... $VERBOSE $QUIET
    
    if [ $? -eq 0 ]; then
        echo -e "\n${GREEN}✓ All tests passed!${NC}"
    else
        echo -e "\n${RED}✗ Tests failed!${NC}"
        exit 1
    fi
fi
