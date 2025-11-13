#!/bin/bash

# LUMEN Backend Test Script

set -e

echo "🧪 Running LUMEN Backend Tests"
echo "=============================="
echo ""

# Run tests with coverage
echo "📊 Running tests with coverage..."
go test -v -race -coverprofile=coverage.out ./...

# Generate coverage report
echo ""
echo "📈 Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

# Display coverage summary
echo ""
echo "Coverage summary:"
go tool cover -func=coverage.out | tail -1

echo ""
echo "✅ Tests complete!"
echo "📄 Coverage report: coverage.html"
