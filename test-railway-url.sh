#!/bin/bash

echo "🚂 Test Your Railway API URL"
echo "============================"

# Replace this with your actual Railway URL
RAILWAY_URL="YOUR_RAILWAY_URL_HERE"

echo "Instructions:"
echo "1. Copy your Railway URL from the dashboard"
echo "2. Replace YOUR_RAILWAY_URL_HERE in this script"
echo "3. Run: ./test-railway-url.sh"
echo ""

if [[ $RAILWAY_URL == "YOUR_RAILWAY_URL_HERE" ]]; then
    echo "❌ Please replace YOUR_RAILWAY_URL_HERE with your actual Railway URL"
    echo ""
    echo "Example:"
    echo "RAILWAY_URL=\"https://swipe-sports-backend-production-abc123.up.railway.app\""
    exit 1
fi

echo "Testing Railway API: $RAILWAY_URL"
echo ""

# Test health endpoint
echo "1. Testing /health:"
curl -s $RAILWAY_URL/health | jq '.' || curl -s $RAILWAY_URL/health
echo ""

# Test auth endpoint  
echo "2. Testing /auth/signup:"
curl -s -X POST $RAILWAY_URL/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"provider": "auth0", "token": "test"}' | jq '.' || curl -s -X POST $RAILWAY_URL/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"provider": "auth0", "token": "test"}'
echo ""

echo "✅ If you see {\"status\":\"healthy\"}, your API is working!"
echo ""
echo "🎯 Next: Add this URL to Vercel environment variables:"
echo "   VITE_API_BASE_URL=$RAILWAY_URL"
