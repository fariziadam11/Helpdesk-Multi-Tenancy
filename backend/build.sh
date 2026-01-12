#!/bin/bash

# Werk Ticketing - Backend Build Script
# This script builds the Go backend for production deployment

set -e

echo "🔨 Building Werk Ticketing Backend..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

# Build binary
echo "📦 Building binary..."
go build -ldflags="-s -w" -o werk-ticketing-backend main.go

# Make executable
chmod +x werk-ticketing-backend

echo "✅ Build complete!"
echo "📍 Binary location: ./werk-ticketing-backend"
echo ""
echo "To run: ./werk-ticketing-backend"
