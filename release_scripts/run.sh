#!/bin/bash
# CyberLab Server — Startup Script (Linux / macOS)
#
# Usage: ./run.sh
# The script reads from .env if it exists, otherwise prompts for required values.

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

# Load .env if it exists
if [ -f .env ]; then
  echo "📁 Loading configuration from .env..."
  set -a
  source .env
  set +a
fi

# Check required environment variables
missing=()
[ -z "$JWT_SECRET" ] && missing+=("JWT_SECRET")
[ -z "$DB_PASSWORD" ] && missing+=("DB_PASSWORD")
[ -z "$FLAG_HASH_SALT" ] && missing+=("FLAG_HASH_SALT")

if [ ${#missing[@]} -gt 0 ]; then
  echo "⚠️  Missing required environment variables: ${missing[*]}"
  echo ""
  echo "Please create a .env file with the following:"
  echo ""
  echo "  JWT_SECRET=your-random-secret-key"
  echo "  DB_PASSWORD=your-database-password"
  echo "  FLAG_HASH_SALT=your-random-64-char-salt"
  echo ""
  echo "Optionally:"
  echo "  AI_API_KEY=sk-your-openai-api-key"
  echo ""
  
  # Generate a random .env if one doesn't exist
  if [ ! -f .env ]; then
    echo "🔧 Generating .env with random values..."
    RANDOM_SECRET=$(openssl rand -hex 32 2>/dev/null || cat /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c64)
    RANDOM_SALT=$(openssl rand -hex 32 2>/dev/null || cat /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c64)
    cat > .env << EOF
JWT_SECRET=${RANDOM_SECRET:-change-me-jwt-secret}
DB_PASSWORD=change-me-db-password
FLAG_HASH_SALT=${RANDOM_SALT:-change-me-flag-salt}
# AI_API_KEY=sk-your-openai-api-key
EOF
    echo "✅ .env file created. Edit it to set your database password, then re-run."
    echo "   nano .env"
    exit 1
  else
    echo "   Edit .env and fill in the missing values, then re-run."
    exit 1
  fi
fi

# Detect platform and choose binary
UNAME_S=$(uname -s)
UNAME_M=$(uname -m)

if [ "$UNAME_S" = "Linux" ]; then
  BINARY="./cyberlab-server-linux"
elif [ "$UNAME_S" = "Darwin" ]; then
  if [ "$UNAME_M" = "arm64" ]; then
    BINARY="./cyberlab-server-darwin-arm64"
  else
    BINARY="./cyberlab-server-darwin"
  fi
else
  echo "❌ Unsupported platform: $UNAME_S"
  exit 1
fi

if [ ! -f "$BINARY" ]; then
  echo "❌ Binary not found: $BINARY"
  echo "   Did you extract the full release package?"
  exit 1
fi

chmod +x "$BINARY"

echo ""
echo "🚀 CyberLab Server v1.0.0"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "   API:       http://localhost:17420"
echo "   Swagger:   http://localhost:17420/swagger/index.html"
echo "   Frontend:  http://localhost:17420"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

exec "$BINARY"
