#!/bin/bash
# setup-docs.sh - Setup and run Docusaurus documentation site

set -e

DOCS_DIR="docs/github-pages"

echo "🚀 Setting up tf2report documentation..."

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
  echo "❌ Error: Must run from tf2report root directory"
  exit 1
fi

# Check Node.js installation
if ! command -v node &> /dev/null; then
  echo "❌ Error: Node.js is not installed"
  echo "Please install Node.js 18+ from https://nodejs.org/"
  exit 1
fi

NODE_VERSION=$(node -v | cut -d'v' -f2 | cut -d'.' -f1)
if [ "$NODE_VERSION" -lt 18 ]; then
  echo "❌ Error: Node.js 18+ is required (current: $(node -v))"
  exit 1
fi

echo "✅ Node.js $(node -v) detected"

# Navigate to docs directory
cd "$DOCS_DIR"

# Install dependencies
echo "📦 Installing dependencies..."
npm install

echo "✅ Setup complete!"
echo ""
echo "Available commands:"
echo "  npm start        - Start development server"
echo "  npm run build    - Build for production"
echo "  npm run serve    - Serve production build"
echo ""
echo "To start the documentation site:"
echo "  cd $DOCS_DIR && npm start"
