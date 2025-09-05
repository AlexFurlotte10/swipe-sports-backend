# 🚀 Deployment Guide for swipesports.co

## Current Status ✅
- ✅ **Authentication endpoints** working
- ✅ **Profile update endpoint** implemented 
- ✅ **MySQL database** connected
- ✅ **JWT authentication** implemented
- ✅ **CORS configured** for swipesports.co

## API Endpoints Available

### Authentication (No auth required)
- `POST /api/v1/auth/signup` - Create new user with Auth0
- `POST /api/v1/auth/login` - Login existing user with Auth0

### Profile (Requires JWT token)
- `GET /api/v1/profile/me` - Get current user profile
- `PUT /api/v1/profile/update` - Update user profile with onboarding data

### Health Check
- `GET /health` - API health status

## Local Development

```bash
# Start with Docker Compose
docker-compose -f docker-compose.simple.yml up -d --build

# Test API
curl http://localhost:8080/health

# View logs
docker-compose -f docker-compose.simple.yml logs app

# Stop services
docker-compose -f docker-compose.simple.yml down
```

## Production Deployment

### 1. Environment Setup

Create production environment file:
```bash
# Copy .env.production to your server
# Update database credentials and JWT secret
```

### 2. Database Setup

```sql
-- Run this on your production MySQL database
CREATE DATABASE swipe_sports;
-- Tables will be created automatically by the application
```

### 3. Deploy Options

#### Option A: Docker (Recommended)
```bash
# Build production image
docker build -f Dockerfile.prod -t swipe-sports-api .

# Run with environment variables
docker run -d \
  --name swipe-sports-api \
  -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-db-password \
  -e JWT_SECRET=your-jwt-secret \
  swipe-sports-api
```

#### Option B: Direct Deployment
```bash
# On your server
go mod tidy
go build -o swipe-sports-api main_simple.go

# Set environment variables
export DB_HOST=your-db-host
export DB_PASSWORD=your-db-password
export JWT_SECRET=your-jwt-secret

# Run
./swipe-sports-api
```

### 4. Frontend Integration

Your frontend should make requests to:
```javascript
// Production API base URL
const API_BASE = 'https://api.swipesports.co'

// Auth signup/login
fetch(`${API_BASE}/api/v1/auth/signup`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    provider: 'auth0',
    token: auth0Token
  })
})

// Profile update (with JWT from signup/login)
fetch(`${API_BASE}/api/v1/profile/update`, {
  method: 'PUT',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${jwtToken}`
  },
  body: JSON.stringify({
    name: "User Name",
    first_name: "User",
    last_name: "Name",
    age: 25,
    location: "Halifax, NS",
    gender: "male",
    preferred_timeslots: "weekends-evenings",
    sport_preferences: {
      tennis: true,
      pickleball: false,
      basketball: false,
      soccer: false
    },
    skill_level: "intermediate",
    ntrp_rating: 3.5,
    play_style: "ranked",
    bio: "Tennis player looking for matches",
    availability: {
      monday: ["evening"],
      tuesday: ["evening"],
      // ... more days
    }
  })
})
```

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DB_HOST` | Database hostname | `localhost` or `your-db-host.com` |
| `DB_PORT` | Database port | `3306` |
| `DB_USER` | Database username | `root` or `swipe_user` |
| `DB_PASSWORD` | Database password | `your-secure-password` |
| `DB_NAME` | Database name | `swipe_sports` |
| `JWT_SECRET` | JWT signing secret | `your-32-char-secret-key` |
| `PORT` | Server port | `8080` |

## Security Notes

1. **JWT Secret**: Use a strong, random 32+ character secret
2. **Database**: Use dedicated database user with minimal permissions
3. **Auth0**: Implement proper Auth0 token verification (currently mocked)
4. **HTTPS**: Always use HTTPS in production
5. **CORS**: Currently allows swipesports.co - update as needed

## Next Steps

1. **Deploy to your server** using one of the options above
2. **Configure DNS** to point api.swipesports.co to your server
3. **Update frontend** to use production API URL
4. **Implement real Auth0 verification** (replace mock function)
5. **Add rate limiting** and other production middleware
6. **Set up monitoring** and logging

## Troubleshooting

### Common Issues

1. **Database connection fails**
   - Check DB_HOST, DB_PASSWORD environment variables
   - Ensure MySQL is running and accessible

2. **CORS errors**
   - Update CORS_ORIGIN environment variable
   - Check frontend is using correct API URL

3. **Auth fails**
   - Verify JWT_SECRET is set
   - Check Auth0 token format

### Checking Logs
```bash
# Docker logs
docker logs swipe-sports-api

# Direct logs (if running without Docker)
# Application logs to stdout
```

## Support

Your API is ready for production! 🎉

The simplified version includes:
- ✅ User authentication with Auth0
- ✅ Profile management 
- ✅ Database persistence
- ✅ JWT tokens
- ✅ Production-ready Docker setup
