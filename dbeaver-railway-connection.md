# 🗄️ DBeaver Railway MySQL Connection

## Connection Details for Your Railway Database

### Database Connection Settings:

```bash
Connection Type: MySQL
Server Host: mysql.railway.internal  
Port: 3306
Database: railway
Username: root
Password: AtngiQJPqmVkftIGLDsOPcHWERMswqzE
```

## Step-by-Step DBeaver Setup:

### 1. Create New Connection
1. Open DBeaver
2. Click "New Database Connection" (+ icon)
3. Select "MySQL"
4. Click "Next"

### 2. Connection Details
Fill in these **exact** values:
```
Server Host: mysql.railway.internal
Port: 3306
Database: railway
Username: root
Password: AtngiQJPqmVkftIGLDsOPcHWERMswqzE
```

### 3. Test Connection
1. Click "Test Connection"
2. Should show "Connected" ✅
3. Click "Finish"

### 4. Browse Your Database
After connecting, you should see:
- Database: `railway`
- Tables: `users` (created by your app)
- User data from your app signups

## Alternative: Railway Public URL
If `mysql.railway.internal` doesn't work, try:
```
Server Host: (Check Railway for public MySQL URL)
Port: (Check Railway for public port)
```

## View Your Data
Once connected, you can:
- Browse `users` table
- See user signups from your frontend
- Run queries to check data
- Monitor database activity

## Security Note
This password gives full access to your production database.
Keep it secure and don't share it publicly.
