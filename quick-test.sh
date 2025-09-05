#!/bin/bash

echo "🚀 Quick Test: Railway Backend Deployment"
echo "========================================"

# Test the deployed Railway API
RAILWAY_URL="https://swipe-sports-backend-production-d5d9f688.up.railway.app"

echo "Testing your deployed Railway API..."
echo "URL: $RAILWAY_URL"
echo ""

# Test health endpoint
echo "1. Testing /health endpoint:"
HEALTH=$(curl -s $RAILWAY_URL/health)
echo "   Response: $HEALTH"
echo ""

# Test auth signup endpoint
echo "2. Testing /auth/signup endpoint:"
SIGNUP=$(curl -s -X POST $RAILWAY_URL/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"provider": "auth0", "token": "test_production_token"}')
echo "   Response: $SIGNUP"
echo ""

if [[ $HEALTH == *"healthy"* ]]; then
    echo "✅ SUCCESS: Your Railway API is working!"
    echo ""
    echo "🎯 Next steps:"
    echo "1. Add this URL to your Vercel environment variables:"
    echo "   VITE_API_BASE_URL=$RAILWAY_URL"
    echo ""
    echo "2. Redeploy your frontend in Vercel"
    echo ""
    echo "3. Test signup/login on swipesports.co"
else
    echo "❌ API not responding. Check Railway deployment."
fi
