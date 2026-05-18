#!/bin/bash
# Post-creation setup script for development container
# Executes once after container is created
# Use for one-time installations and configuration

set -e

echo "🚀 Initializing development container..."

# Update package managers
echo "📦 Updating package managers..."
apt-get update -qq
apt-get upgrade -qq -y

# Install Go development tools
echo "📦 Installing Go development tools..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/cosmtrek/air@latest

# Install additional utilities
echo "📦 Installing additional development utilities..."
apt-get install -qq -y \
	build-essential \
	git-flow \
	jq \
	ripgrep \
	fd-find

# Verify installed tools
echo "✅ Verifying installations..."
echo "Terraform version: $(terraform version -json | jq -r '.terraform_version')"
echo "Go version: $(go version | awk '{print $3}')"
echo "Node.js version: $(node --version)"
echo "Docker version: $(docker --version)"

# Create required directories
mkdir -p /workspace/logs

echo "✅ Development container initialized successfully!"
echo ""
echo "📝 Quick reference:"
echo "  - Go tools location: $(go env GOPATH)/bin"
echo "  - Working directory: $(pwd)"
echo ""
