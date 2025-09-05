#!/bin/bash

echo "🚀 Testing Swipe Sports Backend API"
echo "=================================="

# Test health endpoint
echo "1. Testing health endpoint..."
curl -s http://localhost:8080/health | jq '.'
echo ""

# Test auth signup
echo "2. Testing auth signup..."
SIGNUP_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"provider": "auth0", "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjEyMzQ1Njc4OTAifQ.eyJzdWIiOiJhdXRoMHx1c2VyX2FiYzEyMyIsIm5hbWUiOiJUZXN0IFVzZXIiLCJlbWFpbCI6InRlc3RAdGVzdC5jb20iLCJpYXQiOjE2MzM2NjA4MDAsImV4cCI6MTYzMzc0NzIwMH0.example"}')

echo "Signup Response:"
echo $SIGNUP_RESPONSE | jq '.'
echo ""

# Extract token if signup worked
TOKEN=$(echo $SIGNUP_RESPONSE | jq -r '.token // empty')

if [ ! -z "$TOKEN" ] && [ "$TOKEN" != "null" ]; then
    echo "3. Testing profile endpoint with token..."
    curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/profile/me | jq '.'
    echo ""
    
    echo "4. Testing profile update..."
    curl -s -X PUT http://localhost:8080/api/v1/profile/update \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "name": "Alex F",
        "first_name": "Alex", 
        "last_name": "Furlotte",
        "age": 28,
        "location": "Halifax, NS",
        "gender": "male",
        "preferred_timeslots": "weekends-evenings",
        "sport_preferences": {
          "tennis": true,
          "pickleball": false,
          "basketball": true,
          "soccer": false
        },
        "skill_level": "intermediate",
        "ntrp_rating": 3.5,
        "play_style": "ranked",
        "bio": "Love playing tennis and basketball. Looking for competitive matches!",
        "availability": {
          "monday": ["evening"],
          "tuesday": ["evening"], 
          "wednesday": ["evening"],
          "thursday": ["evening"],
          "friday": ["evening"],
          "saturday": ["morning", "afternoon", "evening"],
          "sunday": ["morning", "afternoon", "evening"]
        }
      }' | jq '.'
else
    echo "❌ Signup failed - no token received"
    echo "Trying login instead..."
    
    LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
      -H "Content-Type: application/json" \
      -d '{"provider": "auth0", "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjEyMzQ1Njc4OTAifQ.eyJzdWIiOiJhdXRoMHx1c2VyX2FiYzEyMyIsIm5hbWUiOiJUZXN0IFVzZXIiLCJlbWFpbCI6InRlc3RAdGVzdC5jb20iLCJpYXQiOjE2MzM2NjA4MDAsImV4cCI6MTYzMzc0NzIwMH0.example"}')
    
    echo "Login Response:"
    echo $LOGIN_RESPONSE | jq '.'
fi

echo ""
echo "✅ API Testing Complete!"
echo "💡 Your endpoints are:"
echo "   - Health: http://localhost:8080/health"
echo "   - Signup: POST http://localhost:8080/api/v1/auth/signup"
echo "   - Login:  POST http://localhost:8080/api/v1/auth/login" 
echo "   - Profile: GET http://localhost:8080/api/v1/profile/me"
echo "   - Update:  PUT http://localhost:8080/api/v1/profile/update"
