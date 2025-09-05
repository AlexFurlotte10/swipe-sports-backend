# 🔐 Environment Variables Examples

## Local Development (.env.local)
```bash
# Local Docker setup (current)
DB_HOST=mysql
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=swipe_sports
JWT_SECRET=local-dev-secret-key-not-secure
PORT=8080
CORS_ORIGIN=http://localhost:3000
```

## Production Examples

### Self-Hosted MySQL (.env.production)
```bash
# Your own server with MySQL installed
DB_HOST=localhost
DB_PORT=3306
DB_USER=swipe_user
DB_PASSWORD=YourSecure123Password!
DB_NAME=swipe_sports
JWT_SECRET=super-secure-32-char-secret-key-here
PORT=8080
CORS_ORIGIN=https://swipesports.co
```

### DigitalOcean Managed Database
```bash
# DigitalOcean Database connection
DB_HOST=db-mysql-nyc1-12345-do-user-123456-0.b.db.ondigitalocean.com
DB_PORT=25060
DB_USER=doadmin
DB_PASSWORD=password_from_digitalocean_panel
DB_NAME=swipe_sports
JWT_SECRET=generate-secure-random-32-char-key
PORT=8080
CORS_ORIGIN=https://swipesports.co
```

### AWS RDS MySQL
```bash
# AWS RDS connection
DB_HOST=swipe-sports.abcd1234.us-east-1.rds.amazonaws.com
DB_PORT=3306
DB_USER=admin
DB_PASSWORD=your_rds_password
DB_NAME=swipe_sports
JWT_SECRET=aws-secure-jwt-secret-key-32-chars
PORT=8080
CORS_ORIGIN=https://swipesports.co
```

### PlanetScale (Serverless MySQL)
```bash
# PlanetScale connection
DB_HOST=aws.connect.psdb.cloud
DB_PORT=3306
DB_USER=your_planetscale_username
DB_PASSWORD=pscale_pw_generated_password
DB_NAME=swipe_sports
JWT_SECRET=planetscale-jwt-secret-32-characters
PORT=8080
CORS_ORIGIN=https://swipesports.co
```

## 🔒 Security Notes

### JWT Secret Generation
```bash
# Generate secure JWT secret (32+ characters)
openssl rand -base64 32

# Or use online generator:
# https://www.allkeysgenerator.com/Random/Security-Encryption-Key-Generator.aspx
```

### Database Password Requirements
- Minimum 12 characters
- Include uppercase, lowercase, numbers, symbols
- Avoid common words
- Example: `Swipe2024!Sports#DB`

## 🚀 Quick Setup Commands

### For DigitalOcean Droplet
```bash
# 1. Install MySQL
sudo apt update
sudo apt install mysql-server

# 2. Secure MySQL
sudo mysql_secure_installation

# 3. Create database and user
sudo mysql < database-setup.sql

# 4. Set environment variables
export DB_HOST=localhost
export DB_USER=swipe_user
export DB_PASSWORD=your_secure_password
export DB_NAME=swipe_sports
export JWT_SECRET=$(openssl rand -base64 32)
export CORS_ORIGIN=https://swipesports.co
```

### For Managed Database
```bash
# 1. Create database instance in provider dashboard
# 2. Note connection details
# 3. Set environment variables from provider info
export DB_HOST=provided_host_url
export DB_USER=provided_username
export DB_PASSWORD=provided_password
# etc.
```
