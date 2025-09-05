# 🚀 Frontend Integration Guide - Swipe Sports Backend

## 📊 **Complete API Endpoints for Frontend**

### **Base Configuration**
```javascript
const API_BASE_URL = 'https://your-backend-url.com'; // Production
// const API_BASE_URL = 'http://localhost:8080'; // Local development

const WS_BASE_URL = 'wss://your-backend-url.com'; // Production WebSocket
// const WS_BASE_URL = 'ws://localhost:8080'; // Local development
```

---

## 🔐 **Authentication Flow**

### **Step 1: OAuth Login/Signup**

**Both endpoints create users if they don't exist**

```javascript
// POST /auth/signup OR /auth/login (same behavior)
const authUser = async (oauthToken, provider = 'auth0') => {
  const response = await fetch(`${API_BASE_URL}/auth/signup`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      provider: provider,  // "auth0", "google", "apple", "facebook"
      token: oauthToken    // JWT token from OAuth provider
    })
  });

  const data = await response.json();
  
  if (response.ok) {
    // Store JWT token for future requests
    localStorage.setItem('authToken', data.token);
    return data.user; // Basic user profile
  } else {
    throw new Error(data.error);
  }
};
```

**Response:**
```javascript
{
  "token": "eyJhbGciOiJIUzI1NiIs...",  // Your backend JWT token
  "user": {
    "id": 123,
    "oauth_id": "auth0|abc123",
    "oauth_provider": "auth0", 
    "name": "John Doe",                // ✅ Real name from OAuth
    "email": "john.doe@gmail.com",     // ✅ Real email from OAuth
    "gender": "other",                 // ❌ Default value
    "skill_level": "beginner",         // ❌ Default value  
    "play_style": "fun",               // ❌ Default value
    "preferred_timeslots": "anytime-anywhere", // ❌ Default value
    "rank": 1000,
    "first_name": null,                // ❌ Not set yet
    "last_name": null,                 // ❌ Not set yet
    "age": null,                       // ❌ Not set yet
    "location": null,                  // ❌ Not set yet
    "bio": null,                       // ❌ Not set yet
    "sport_preferences": {},           // ❌ Empty
    "availability": {},                // ❌ Empty
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## 👤 **Profile Management**

### **Step 2: Get Current User Profile**

```javascript
// GET /profile/me
const getCurrentUser = async () => {
  const token = localStorage.getItem('authToken');
  
  const response = await fetch(`${API_BASE_URL}/profile/me`, {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });

  if (response.ok) {
    return await response.json();
  } else {
    throw new Error('Failed to get user profile');
  }
};
```

### **Step 3: Complete Profile Setup (Onboarding)**

**This is where all the detailed information gets stored!**

```javascript
// PUT /profile/me
const completeProfile = async (profileData) => {
  const token = localStorage.getItem('authToken');
  
  const response = await fetch(`${API_BASE_URL}/profile/me`, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(profileData)
  });

  const data = await response.json();
  
  if (response.ok) {
    // Update stored token with new user data
    localStorage.setItem('authToken', data.token);
    return data.user; // Complete user profile
  } else {
    throw new Error(data.error);
  }
};
```

---

## 📝 **Frontend Profile Form Data Structure**

### **Complete Profile Data to Send:**

```javascript
const profileFormData = {
  // Required Fields
  "name": "John Doe",                    // Full display name
  "first_name": "John",                  // ✅ Gets stored in database
  "last_name": "Doe",                    // ✅ Gets stored in database  
  "age": 28,                             // ✅ Gets stored in database
  "location": "New York, NY",            // ✅ Gets stored in database
  "gender": "male",                      // ✅ "male", "female", "other"
  
  // Sports Preferences (Required)
  "sport_preferences": {                 // ✅ JSON object
    "tennis": true,
    "basketball": false, 
    "soccer": true,
    "volleyball": false
  },
  
  // Skill & Play Style (Required)
  "skill_level": "intermediate",         // ✅ "beginner", "intermediate", "advanced"
  "ntrp_rating": 3.5,                   // ✅ Tennis rating (1.0 - 5.5)
  "play_style": "ranked",               // ✅ "ranked" or "fun"
  "preferred_timeslots": "weekends-evenings", // ✅ See valid options below
  
  // Availability Schedule (Required)
  "availability": {                      // ✅ JSON object
    "monday": ["18:00", "20:00"],
    "tuesday": [],
    "wednesday": ["18:00", "20:00"], 
    "thursday": [],
    "friday": ["19:00", "21:00"],
    "saturday": ["10:00", "14:00"],
    "sunday": ["10:00", "12:00"]
  },
  
  // Optional Fields
  "bio": "Tennis enthusiast looking for doubles partners. Love competitive games but also enjoy casual matches on weekends!"
};

// Call the API
const updatedUser = await completeProfile(profileFormData);
```

### **Valid Options for Form Fields:**

```javascript
// Gender Options
const GENDER_OPTIONS = ["male", "female", "other"];

// Skill Level Options  
const SKILL_LEVELS = ["beginner", "intermediate", "advanced"];

// Play Style Options
const PLAY_STYLES = ["ranked", "fun"];

// Preferred Timeslots Options
const TIMESLOT_OPTIONS = [
  "weekends-evenings",
  "anytime-anywhere", 
  "weekends-only",
  "weekdays-only"
];

// Sports Options
const SPORTS_OPTIONS = [
  "tennis",
  "basketball",
  "soccer", 
  "volleyball"
];
```

---

## 🎯 **Complete Frontend Flow Example**

```javascript
// 1. User Authentication
const handleOAuthSuccess = async (oauthToken) => {
  try {
    const user = await authUser(oauthToken, 'auth0');
    
    // Check if profile is complete
    if (!user.first_name || !user.age || !user.location) {
      // Redirect to onboarding/profile setup
      navigate('/onboarding');
    } else {
      // Profile complete, go to main app
      navigate('/dashboard');
    }
  } catch (error) {
    console.error('Auth failed:', error);
  }
};

// 2. Profile Completion Form Submit
const handleProfileSubmit = async (formData) => {
  try {
    const updatedUser = await completeProfile(formData);
    
    // Profile complete! 
    navigate('/dashboard');
  } catch (error) {
    console.error('Profile update failed:', error);
    // Show validation errors to user
  }
};

// 3. Check Profile Completion Status
const checkProfileComplete = (user) => {
  return !!(
    user.first_name && 
    user.last_name && 
    user.age && 
    user.location && 
    user.gender !== 'other' && // Default value
    user.skill_level !== 'beginner' && // Default value
    Object.keys(user.sport_preferences || {}).length > 0 &&
    Object.keys(user.availability || {}).length > 0
  );
};
```

---

## 🔄 **Profile Update Response**

**After successful profile completion:**

```javascript
{
  "token": "eyJhbGciOiJIUzI1NiIs...",  // Updated JWT token
  "user": {
    "id": 123,
    "name": "John Doe",
    "first_name": "John",              // ✅ Now populated
    "last_name": "Doe",                // ✅ Now populated
    "age": 28,                         // ✅ Now populated
    "email": "john.doe@gmail.com",
    "gender": "male",                  // ✅ Real selection
    "location": "New York, NY",        // ✅ Now populated
    "latitude": 44.6488,               // ✅ Auto-set (Halifax default)
    "longitude": -63.5752,             // ✅ Auto-set (Halifax default)
    "rank": 1000,
    "profile_pic_url": null,
    "bio": "Tennis enthusiast...",     // ✅ Now populated
    "sport_preferences": {             // ✅ Real preferences
      "tennis": true,
      "basketball": false,
      "soccer": true,
      "volleyball": false
    },
    "skill_level": "intermediate",     // ✅ Real skill level
    "ntrp_rating": 3.5,               // ✅ Real rating
    "play_style": "ranked",           // ✅ Real preference
    "preferred_timeslots": "weekends-evenings", // ✅ Real preference
    "availability": {                  // ✅ Real schedule
      "monday": ["18:00", "20:00"],
      "wednesday": ["18:00", "20:00"],
      "friday": ["19:00", "21:00"],
      "saturday": ["10:00", "14:00"],
      "sunday": ["10:00", "12:00"]
    },
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:30:00Z"  // ✅ Updated timestamp
  }
}
```

---

## 🛡️ **Authentication Headers**

**For all protected endpoints:**

```javascript
const headers = {
  'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
  'Content-Type': 'application/json'
};
```

---

## 🎨 **Suggested Frontend User Flow**

### **1. Landing Page**
- OAuth login buttons (Auth0, Google, etc.)

### **2. After OAuth Success**
```javascript
if (checkProfileComplete(user)) {
  navigate('/dashboard');  // Profile already complete
} else {
  navigate('/onboarding'); // Need to complete profile
}
```

### **3. Onboarding Flow (Multi-step form)**
- **Step 1:** Personal Info (first name, last name, age, location)
- **Step 2:** Sports Preferences (which sports, skill level, NTRP rating)
- **Step 3:** Play Style & Schedule (ranked vs fun, availability)
- **Step 4:** Bio (optional)
- **Submit:** Call `completeProfile()` → Navigate to dashboard

### **4. Dashboard**
- Show user profile
- Access to swipe/match features
- Profile edit options

---

## 🚨 **Error Handling**

```javascript
const handleApiError = (error, response) => {
  if (response.status === 401) {
    // Token expired, redirect to login
    localStorage.removeItem('authToken');
    navigate('/login');
  } else if (response.status === 400) {
    // Validation error, show to user
    console.error('Validation error:', error);
  } else {
    // General error
    console.error('API error:', error);
  }
};
```

---

## 📱 **Summary for Frontend Team**

### **Key Points:**
1. **OAuth creates basic user** with real email/name but default values for other fields
2. **Profile completion stores all detailed info** (first name, last name, age, etc.)
3. **Use profile completion status** to determine user flow
4. **Store JWT token** and include in all protected requests
5. **Real user data flows into database** after profile completion

### **Required Frontend Components:**
- [ ] OAuth login integration
- [ ] Profile completion check
- [ ] Multi-step onboarding form
- [ ] Profile completion API calls
- [ ] Error handling for validation
- [ ] JWT token management

Your backend is ready to receive and store all this detailed profile information! 🎉
