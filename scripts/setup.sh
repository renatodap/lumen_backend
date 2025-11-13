#!/bin/bash

# LUMEN Backend Setup Script

set -e

echo "🚀 LUMEN Backend Setup"
echo "======================"
echo ""

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✅ Go version: $GO_VERSION"
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
    echo "⚠️  Please update .env with your Supabase credentials"
    echo ""
else
    echo "✅ .env file already exists"
    echo ""
fi

# Install dependencies
echo "📦 Installing dependencies..."
go mod download
go mod verify
echo "✅ Dependencies installed"
echo ""

# Build the application
echo "🔨 Building application..."
go build -o bin/server cmd/server/main.go
echo "✅ Build successful"
echo ""

# Run tests
echo "🧪 Running tests..."
if go test -v ./...; then
    echo "✅ All tests passed"
else
    echo "⚠️  Some tests failed"
fi
echo ""

# Check database connection
echo "🔍 Checking database connection..."
if [ -f .env ]; then
    source .env
    if [ -n "$DB_HOST" ] && [ -n "$DB_PASSWORD" ]; then
        echo "✅ Database credentials found in .env"
        echo "   Host: $DB_HOST"
        echo ""
    else
        echo "⚠️  Database credentials not configured in .env"
        echo "   Please update DB_HOST, DB_PASSWORD, and other database variables"
        echo ""
    fi
else
    echo "⚠️  .env file not found"
    echo ""
fi

# Setup complete
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Update .env with your Supabase credentials"
echo "2. Run 'make dev' or 'go run cmd/server/main.go' to start the server"
echo "3. Visit http://localhost:8080/health to verify"
echo ""
echo "For more information, see docs/README.md"
