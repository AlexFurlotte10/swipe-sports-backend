# 🚂 Find Your Railway API URL

## Your Railway Deployment Info
- **Project ID**: `d5d9f688-9cb5-4776-a900-8522f42dc5dd`
- **Service**: `swipe-sports-backend`
- **Environment**: `production`

## How to Find Your API URL

### Method 1: Railway Dashboard
1. Go to your Railway project dashboard
2. Click on your `swipe-sports-backend` service
3. Look for the **"Domains"** section
4. You'll see a URL like: `https://swipe-sports-backend-production-xxxx.up.railway.app`

### Method 2: Railway CLI (if installed)
```bash
railway status
```

### Method 3: Check Railway Logs
1. In Railway dashboard → your service
2. Click "Logs" tab
3. Look for startup messages showing the server URL

## Expected URL Format
Your API URL will look like:
```
https://swipe-sports-backend-production-d5d9f688.up.railway.app
```

## Test Your Deployed API
Once you have the URL, test it:
```bash
curl https://your-railway-url.up.railway.app/health
# Should return: {"status":"healthy"}
```
