#!/bin/bash

echo "🚂 Railway Deployment Setup"
echo "=========================="

# Check if railway CLI is installed
if ! command -v railway &> /dev/null; then
    echo "Installing Railway CLI..."
    curl -fsSL https://railway.app/install.sh | sh
    echo "Railway CLI installed!"
else
    echo "✅ Railway CLI already installed"
fi

echo ""
echo "🚀 Quick Railway Setup Steps:"
echo "1. Go to https://railway.app"
echo "2. Sign up with GitHub"  
echo "3. Create new project from your GitHub repo"
echo "4. Add MySQL database service"
echo "5. Deploy automatically!"
echo ""
echo "Your API will be live at: https://your-app.up.railway.app"
echo "Add custom domain: api.swipesports.co"
echo ""
echo "💰 Cost: ~$4/month (covered by $5 monthly credit)"
