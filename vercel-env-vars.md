# 🔐 Vercel Environment Variables

## Environment Variables for Your Vercel Dashboard

Set these in your Vercel project settings:

### For Production Deployment

```bash
# Backend API URL (Railway will provide this)
VITE_API_BASE_URL=https://your-app-name.up.railway.app

# WebSocket URL (if you add WebSocket later)
VITE_WS_BASE_URL=wss://your-app-name.up.railway.app

# Custom domain (after you set it up)
# VITE_API_BASE_URL=https://api.swipesports.co
# VITE_WS_BASE_URL=wss://api.swipesports.co
```

### For Development/Preview

```bash
# Local development (keep these for preview deployments)
VITE_API_BASE_URL=http://localhost:8080
VITE_WS_BASE_URL=ws://localhost:8080
```

## How to Set in Vercel

1. Go to your Vercel dashboard
2. Select your project
3. Go to Settings → Environment Variables
4. Add each variable:
   - **Name**: `VITE_API_BASE_URL`
   - **Value**: `https://your-app-name.up.railway.app`
   - **Environments**: Production ✓

## After Railway Deployment

1. Railway will give you a URL like: `https://swipe-sports-backend-production-xxxx.up.railway.app`
2. Update `VITE_API_BASE_URL` in Vercel with this URL
3. Test your frontend → backend connection

## Custom Domain Setup (Optional)

If you want `api.swipesports.co`:

1. **In Railway**: Add custom domain `api.swipesports.co`
2. **In DNS**: Add CNAME record `api.swipesports.co` → Railway URL
3. **In Vercel**: Update `VITE_API_BASE_URL=https://api.swipesports.co`
