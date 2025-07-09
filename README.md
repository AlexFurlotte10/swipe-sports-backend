# Swipe Sports Backend

A comprehensive Go backend for the Swipe Sports dating-style sports matchmaking app. Users can swipe through player profiles to find teammates and chat with matched users.

## Features

- **OAuth Authentication**: Google, Apple, and Facebook signup/login
- **Profile Management**: User profiles with sports preferences, rankings, and availability
- **Matchmaking**: Swipe left/right on profiles with filtering by location, gender, and rank
- **Real-time Messaging**: WebSocket-based chat between matched users
- **Caching**: Redis for performance optimization
- **RESTful API**: Complete REST API with proper authentication and validation

## Tech Stack

- **Backend**: Go 1.21 with Gin framework
- **Database**: MySQL 8.0
- **Cache**: Redis 7
- **Authentication**: JWT tokens with OAuth 2.0
- **Real-time**: WebSockets (Gorilla)
- **Containerization**: Docker & Docker Compose

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- MySQL 8.0 (or use Docker)
- Redis 7 (or use Docker)

### Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd swipe-sports-backend
   ```

2. **Set up environment variables**
   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

3. **Start the development environment**
   ```bash
   make setup
   ```

4. **Verify the setup**
   - API: http://localhost:8080
   - Health check: http://localhost:8080/health
   - phpMyAdmin: http://localhost:8081

### Manual Setup

1. **Install dependencies**
   ```bash
   make deps
   ```

2. **Start services with Docker Compose**
   ```bash
   docker-compose up -d
   ```

3. **Run the application**
   ```bash
   make run
   ```

## API Endpoints

### Authentication
- `POST /api/v1/auth/signup` - OAuth signup
- `POST /api/v1/auth/login` - OAuth login
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/auth/logout` - Logout

### Profiles
- `GET /api/v1/profiles` - Get profiles to swipe (with filters)
- `GET /api/v1/profile/me` - Get current user's profile
- `PUT /api/v1/profile/me` - Update current user's profile
- `POST /api/v1/profile/picture` - Upload profile picture

### Matchmaking
- `POST /api/v1/swipe` - Record a swipe action
- `GET /api/v1/matches` - Get all user's matches
- `GET /api/v1/matches/:id` - Get specific match details

### Messaging
- `GET /api/v1/messages?match_id=123&page=1&limit=50` - Get messages for a match
- `POST /api/v1/messages` - Send a message
- `DELETE /api/v1/messages/:id` - Delete a message
- `GET /api/v1/messages/:match_id/latest` - Get latest message
- `GET /api/v1/messages/:match_id/unread-count` - Get unread count
- `POST /api/v1/messages/typing` - Send typing indicator

### WebSocket
- `GET /api/v1/ws/chat?user_id=123&match_id=456` - WebSocket connection for real-time chat

## Database Schema

### Users Table
```sql
CREATE TABLE users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  oauth_id VARCHAR(255) UNIQUE,
  oauth_provider VARCHAR(50),
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE,
  gender ENUM('male', 'female', 'other'),
  location VARCHAR(255),
  latitude DECIMAL(10, 8),
  longitude DECIMAL(11, 8),
  rank INT DEFAULT 1000,
  profile_pic_url VARCHAR(500),
  bio TEXT,
  sport_preferences JSON,
  skill_level VARCHAR(50),
  play_style VARCHAR(100),
  availability JSON,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Swipes Table
```sql
CREATE TABLE swipes (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  swiper_id BIGINT NOT NULL,
  swipee_id BIGINT NOT NULL,
  direction ENUM('left', 'right') NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (swiper_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (swipee_id) REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE KEY unique_swipe (swiper_id, swipee_id)
);
```

### Matches Table
```sql
CREATE TABLE matches (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user1_id BIGINT NOT NULL,
  user2_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user1_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (user2_id) REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE KEY unique_match (user1_id, user2_id)
);
```

### Messages Table
```sql
CREATE TABLE messages (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  match_id BIGINT NOT NULL,
  sender_id BIGINT NOT NULL,
  content TEXT NOT NULL,
  message_type ENUM('text', 'image', 'audio') DEFAULT 'text',
  media_url VARCHAR(500),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (match_id) REFERENCES matches(id) ON DELETE CASCADE,
  FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE
);
```

## Configuration

Environment variables are loaded from `.env` file:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=swipe_sports

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY=24h

# OAuth Configuration
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
APPLE_CLIENT_ID=your-apple-client-id
APPLE_CLIENT_SECRET=your-apple-client-secret
FACEBOOK_CLIENT_ID=your-facebook-client-id
FACEBOOK_CLIENT_SECRET=your-facebook-client-secret

# AWS Configuration
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-aws-access-key
AWS_SECRET_ACCESS_KEY=your-aws-secret-key
AWS_S3_BUCKET=swipe-sports-media
AWS_CLOUDFRONT_DOMAIN=your-cloudfront-domain

# Server Configuration
PORT=8080
ENV=development
CORS_ORIGIN=http://localhost:3000

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m
```

## Development

### Available Commands

```bash
make help          # Show all available commands
make build         # Build the application
make run           # Run the application locally
make test          # Run tests
make clean         # Clean build artifacts
make deps          # Download dependencies
make fmt           # Format code
make lint          # Run linter
make docker-build  # Build Docker image
make docker-run    # Run with Docker Compose
make docker-stop   # Stop Docker Compose services
make logs          # View application logs
make db            # Access MySQL database
make redis         # Access Redis CLI
```

### Project Structure

```
swipe-sports-backend/
├── internal/
│   ├── auth/          # Authentication and JWT
│   ├── config/        # Configuration management
│   ├── database/      # Database initialization
│   ├── handler/       # HTTP handlers
│   ├── models/        # Data models
│   ├── redis/         # Redis client and cache
│   ├── repository/    # Database operations
│   ├── server/        # Server setup and routes
│   └── service/       # Business logic
├── scripts/           # Database scripts
├── main.go           # Application entry point
├── go.mod            # Go module file
├── Dockerfile        # Docker configuration
├── docker-compose.yml # Development environment
├── Makefile          # Development commands
└── README.md         # This file
```

## Testing

Run tests with:
```bash
make test
```

## Deployment

### Docker Deployment

1. **Build the image**
   ```bash
   make docker-build
   ```

2. **Run with Docker Compose**
   ```bash
   docker-compose up -d
   ```

### Production Deployment

For production deployment, consider:

1. **Environment Variables**: Set all production environment variables
2. **Database**: Use managed MySQL service (AWS RDS, Google Cloud SQL)
3. **Redis**: Use managed Redis service (AWS ElastiCache, Google Memorystore)
4. **Load Balancer**: Use a load balancer for multiple instances
5. **SSL/TLS**: Configure HTTPS with proper certificates
6. **Monitoring**: Set up application and database monitoring
7. **Logging**: Configure structured logging and log aggregation

## Security Considerations

- JWT tokens are used for authentication
- OAuth tokens are validated with providers
- Input validation and sanitization
- Rate limiting to prevent abuse
- CORS configuration for frontend security
- Database queries use parameterized statements
- User data is validated before processing

## Performance Optimizations

- Redis caching for frequently accessed data
- Database indexing on frequently queried columns
- Connection pooling for database and Redis
- Pagination for large result sets
- Efficient queries with proper joins

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and linting
6. Submit a pull request

## License

This project is licensed under the MIT License.