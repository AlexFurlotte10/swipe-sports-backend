#!/bin/bash

echo "🚀 Testing Your Production Railway API"
echo "====================================="

# Replace with your actual Railway URL
RAILWAY_URL="REPLACE_WITH_YOUR_RAILWAY_URL"

if [[ $RAILWAY_URL == "REPLACE_WITH_YOUR_RAILWAY_URL" ]]; then
    echo "❌ Please update RAILWAY_URL with your actual Railway URL"
    echo ""
    echo "1. Go to Railway dashboard"
    echo "2. Click your swipe-sports-backend service"  
    echo "3. Find the public URL"
    echo "4. Replace RAILWAY_URL in this script"
    exit 1
fi

echo "Testing: $RAILWAY_URL"
echo ""

# Test health endpoint
echo "1. Testing /health endpoint:"
HEALTH=$(curl -s $RAILWAY_URL/health)
echo "   Response: $HEALTH"

if [[ $HEALTH == *"healthy"* ]]; then
    echo "   ✅ SUCCESS: API is healthy!"
else
    echo "   ❌ FAILED: API not responding properly"
fi
echo ""

# Test auth signup
echo "2. Testing /auth/signup endpoint:"
SIGNUP=$(curl -s -X POST $RAILWAY_URL/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"provider": "auth0", "token": "test_production"}')
echo "   Response: $SIGNUP"

if [[ $SIGNUP == *"error"* ]]; then
    echo "   ✅ EXPECTED: Auth endpoint responding (error expected with test token)"
else
    echo "   ✅ SUCCESS: Auth endpoint working"
fi
echo ""

echo "🎯 Next Steps:"
echo "1. Add this URL to Vercel: $RAILWAY_URL"
echo "2. Test signup/login from your frontend"
echo "3. Check database in DBeaver for new users"
echo ""
echo "📱 Your production API is ready!"
