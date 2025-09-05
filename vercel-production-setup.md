# 🌐 Vercel Production Environment Variables

## Add to Your `swipe-sports` Project in Vercel

### Environment Variables to Add:

```bash
# Variable 1
Name: VITE_API_BASE_URL
Value: https://your-railway-url.up.railway.app
Environments: ✓ Production ✓ Preview

# Variable 2
Name: VITE_WS_BASE_URL  
Value: wss://your-railway-url.up.railway.app
Environments: ✓ Production ✓ Preview
```

### Steps:
1. Go to Vercel Dashboard
2. Select "swipe-sports" project
3. Settings → Environment Variables
4. Add both variables above
5. Click "Redeploy" or push new commit

### Replace URL:
Get your actual Railway URL from Railway dashboard and replace:
`https://your-railway-url.up.railway.app`

### Test After Deployment:
1. Visit https://swipesports.co
2. Try sign up/login
3. Check browser Network tab for API calls to Railway
4. Should see calls to your Railway URL instead of localhost

### Keep Development Variables:
```bash
# Keep these for local development
Name: VITE_API_BASE_URL  
Value: http://localhost:8080
Environments: ✓ Development
```
