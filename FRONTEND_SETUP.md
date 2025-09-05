# 🌐 Frontend Setup for Railway Backend

## 🚀 Add These to Your Vercel Dashboard

### Your Frontend Repository: `swipe-sports`

Go to your Vercel dashboard → `swipe-sports` project → Settings → Environment Variables

Add these **exact** environment variables:

### 1. Production Environment
```bash
# Variable Name: VITE_API_BASE_URL
# Value: https://swipe-sports-backend-production-d5d9f688.up.railway.app
# Environments: ✓ Production

# Variable Name: VITE_WS_BASE_URL  
# Value: wss://swipe-sports-backend-production-d5d9f688.up.railway.app
# Environments: ✓ Production
```

### 2. Preview Environment (Optional)
```bash
# Variable Name: VITE_API_BASE_URL
# Value: https://swipe-sports-backend-production-d5d9f688.up.railway.app
# Environments: ✓ Preview

# Variable Name: VITE_WS_BASE_URL
# Value: wss://swipe-sports-backend-production-d5d9f688.up.railway.app  
# Environments: ✓ Preview
```

### 3. Development Environment (Keep Local)
```bash
# Variable Name: VITE_API_BASE_URL
# Value: http://localhost:8080
# Environments: ✓ Development

# Variable Name: VITE_WS_BASE_URL
# Value: ws://localhost:8080
# Environments: ✓ Development
```

## 🔧 **Step-by-Step Vercel Setup**

### 1. Go to Vercel Dashboard
- Navigate to: https://vercel.com/dashboard
- Select your `swipe-sports` project

### 2. Add Environment Variables
```bash
1. Click "Settings" tab
2. Click "Environment Variables" in sidebar
3. Add new variable:
   - Name: VITE_API_BASE_URL
   - Value: https://swipe-sports-backend-production-d5d9f688.up.railway.app
   - Environments: Production ✓
4. Click "Save"
5. Repeat for VITE_WS_BASE_URL
```

### 3. Redeploy Your Frontend
```bash
1. In Vercel dashboard
2. Go to "Deployments" tab  
3. Click "Redeploy" on latest deployment
4. Or push a new commit to trigger deployment
```

## 🧪 **Test the Connection**

After redeployment, your frontend will:
1. **Development**: Connect to `http://localhost:8080` (your local API)
2. **Production**: Connect to Railway API automatically

### Test Production Connection
1. Open your live site: `https://swipesports.co`
2. Try to sign up/login
3. Check browser Network tab for API calls to Railway

## ⚠️ **Important Notes**

### Railway URL Format
Your Railway URL should be:
```
https://swipe-sports-backend-production-d5d9f688.up.railway.app
```

**Note**: The exact URL might be slightly different. Get the real URL from your Railway dashboard.

### CORS Configuration
Your Railway backend is already configured for:
- ✅ `https://swipesports.co` (production)
- ✅ `http://localhost:3000` (development)

### Custom Domain (Optional)
If you want `api.swipesports.co`:
1. **Railway**: Add custom domain
2. **DNS**: Point subdomain to Railway
3. **Vercel**: Update environment variable to custom domain

## 🎉 **Final Result**

After setup:
- **Your frontend** (`swipesports.co`) → **Your Railway API** 
- **Sign up/login** will create users in your Railway MySQL database
- **Profile updates** will work with your onboarding flow

## 🆘 **If Something Doesn't Work**

1. **Check Railway logs** for errors
2. **Verify Vercel environment variables** are set correctly  
3. **Test API directly**: `curl https://your-railway-url.up.railway.app/health`
4. **Check browser console** for CORS errors
