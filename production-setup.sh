#!/bin/bash

echo "🚀 Swipe Sports Production Setup Script"
echo "======================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}This script helps you deploy to production${NC}"
echo ""

# Check if running on production server
echo "1. Environment Check"
echo "==================="
if [ "$PWD" = "/Users/alex.furlotte/Documents/Repos/swipe-sports-backend" ]; then
    echo -e "${YELLOW}⚠️  You're running this on your local machine${NC}"
    echo "This script is meant for your production server"
    echo ""
    echo "Local URLs (current):"
    echo "  Frontend: http://localhost:3000"
    echo "  Backend:  http://localhost:8080"
    echo "  Database: localhost:3306"
    echo ""
    echo "Production URLs (target):"
    echo "  Frontend: https://swipesports.co"
    echo "  Backend:  https://api.swipesports.co"
    echo "  Database: your-production-db-host"
    echo ""
    echo "To continue with production setup:"
    echo "1. Copy this project to your production server"
    echo "2. Run this script on the production server"
    echo "3. Set up DNS records for api.swipesports.co"
    exit 1
fi

# Production setup steps
echo "2. Setting up production environment..."
echo "======================================"

# Create production environment file
cat > .env.production << EOF
# Production Environment for swipesports.co
DB_HOST=localhost
DB_PORT=3306
DB_USER=swipe_user
DB_PASSWORD=CHANGE_THIS_PASSWORD
DB_NAME=swipe_sports
JWT_SECRET=CHANGE_THIS_32_CHAR_SECRET_KEY
PORT=8080
CORS_ORIGIN=https://swipesports.co
EOF

echo -e "${GREEN}✅ Created .env.production${NC}"

# Check if Go is installed
if command -v go &> /dev/null; then
    echo -e "${GREEN}✅ Go is installed${NC}"
    go version
else
    echo -e "${RED}❌ Go is not installed${NC}"
    echo "Install Go: https://golang.org/doc/install"
    exit 1
fi

# Check if MySQL is available
if command -v mysql &> /dev/null; then
    echo -e "${GREEN}✅ MySQL client available${NC}"
else
    echo -e "${YELLOW}⚠️  MySQL client not found${NC}"
    echo "Install: sudo apt install mysql-client"
fi

echo ""
echo "3. Building application..."
echo "========================="

# Build the application
if go mod tidy && go build -o swipe-api main_simple.go; then
    echo -e "${GREEN}✅ Application built successfully${NC}"
else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

echo ""
echo "4. Setup instructions:"
echo "====================="

echo -e "${YELLOW}Manual steps required:${NC}"
echo ""
echo "1. Database Setup:"
echo "   - Install MySQL server"
echo "   - Create database: swipe_sports"
echo "   - Create user: swipe_user"
echo "   - Update DB_PASSWORD in .env.production"
echo ""
echo "2. Security:"
echo "   - Generate secure JWT_SECRET (32+ characters)"
echo "   - Update .env.production with real values"
echo ""
echo "3. DNS Configuration:"
echo "   - Point api.swipesports.co to this server's IP"
echo "   - Set up SSL certificate (Let's Encrypt)"
echo ""
echo "4. Run the API:"
echo "   source .env.production && ./swipe-api"
echo ""
echo "5. Test endpoints:"
echo "   curl https://api.swipesports.co/health"
echo ""

echo -e "${GREEN}🎉 Production setup ready!${NC}"
echo ""
echo "Your API will run on:"
echo "  Local:      http://localhost:8080"
echo "  Production: https://api.swipesports.co"
