#!/bin/bash

# Werk Ticketing - Frontend Build Script
# This script builds the Vue frontend for production deployment

set -e

echo "🔨 Building Werk Ticketing Frontend..."

# Check if Bun is installed
if ! command -v bun &> /dev/null; then
    echo "❌ Bun is not installed. Please install Bun first."
    exit 1
fi

# Install dependencies
echo "📦 Installing dependencies..."
bun install

# Build for production
echo "🏗️  Building for production..."
bun run build

echo "✅ Build complete!"
echo "📍 Build output: ./dist/"
echo ""
echo "To preview: bun run preview"
