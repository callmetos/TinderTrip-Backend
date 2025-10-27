# 🔐 Google OAuth Error Handling - Frontend Integration

## Overview
เมื่อ Google OAuth login เกิด error ตอนนี้ Backend จะ **redirect กลับไป Frontend** พร้อม error parameters แทนที่จะ return JSON

---

## 🔄 Flow Diagram

```
User clicks "Login with Google"
         ↓
Frontend → Backend (/auth/google)
         ↓
Google OAuth consent screen
         ↓
[User approves/denies]
         ↓
Google → Backend (/auth/google/callback)
         ↓
    [Success?]
         ↓
    ╔════╩════╗
    ↓         ↓
  ERROR    SUCCESS
    ↓         ↓
FE/callback  FE/callback
?error=...   ?token=...
```

---

## ✅ Success Callback

### URL Format
```
{FRONTEND_URL}/callback?token={jwt}&user_id={id}&email={email}&display_name={name}&provider=google&is_verified=true
```

### Example
```
http://localhost:8081/callback?token=eyJhbGci...&user_id=123e4567...&email=user@gmail.com&display_name=John+Doe&provider=google&is_verified=true
```

### Parameters
| Parameter | Type | Description |
|-----------|------|-------------|
| `token` | string | JWT authentication token |
| `user_id` | string | User UUID |
| `email` | string | User email |
| `display_name` | string | User display name |
| `provider` | string | Always "google" |
| `is_verified` | boolean | Email verification status |

---

## ❌ Error Callback

### URL Format
```
{FRONTEND_URL}/callback?error={error_type}&message={error_message}
```

### Error Types

#### 1. **missing_parameters**
```
/callback?error=missing_parameters&message=Authorization+code+and+state+are+required
```
**Cause**: ไม่มี `code` หรือ `state` จาก Google  
**Action**: แสดงข้อความให้ user ลองใหม่

---

#### 2. **invalid_state**
```
/callback?error=invalid_state&message=Invalid+or+expired+state+parameter
```
**Cause**: State parameter ไม่ valid หรือหมดอายุ (CSRF protection)  
**Action**: แสดงข้อความให้ user ลองใหม่

---

#### 3. **token_exchange_failed**
```
/callback?error=token_exchange_failed&message=Failed+to+exchange+authorization+code
```
**Cause**: ไม่สามารถแลก authorization code เป็น access token ได้  
**Action**: แสดงข้อความ "Authentication failed" และให้ลองใหม่

---

#### 4. **user_info_failed**
```
/callback?error=user_info_failed&message=Failed+to+get+user+information+from+Google
```
**Cause**: ไม่สามารถดึงข้อมูล user จาก Google ได้  
**Action**: แสดงข้อความ "Cannot retrieve user information"

---

#### 5. **user_creation_failed**
```
/callback?error=user_creation_failed&message=Failed+to+create+or+update+user+account
```
**Cause**: Database error ขณะสร้าง/อัพเดท user  
**Action**: แสดงข้อความ "Account creation failed" และติดต่อ support

---

#### 6. **token_generation_failed**
```
/callback?error=token_generation_failed&message=Failed+to+generate+authentication+token
```
**Cause**: ไม่สามารถสร้าง JWT token ได้  
**Action**: แสดงข้อความ "Authentication failed" และลองใหม่

---

## 💻 Frontend Implementation

### React/Next.js Example

```typescript
// pages/callback.tsx
import { useEffect } from 'react';
import { useRouter } from 'next/router';
import { toast } from 'react-hot-toast';

const OAuthCallback = () => {
  const router = useRouter();
  const { token, error, message, user_id, email, display_name } = router.query;

  useEffect(() => {
    if (error) {
      // Handle error
      handleOAuthError(error as string, message as string);
    } else if (token) {
      // Handle success
      handleOAuthSuccess(
        token as string,
        user_id as string,
        email as string,
        display_name as string
      );
    }
  }, [error, token, message, user_id, email, display_name]);

  const handleOAuthError = (errorType: string, errorMessage: string) => {
    // Map error types to user-friendly messages
    const errorMessages: Record<string, string> = {
      missing_parameters: 'Login incomplete. Please try again.',
      invalid_state: 'Session expired. Please try again.',
      token_exchange_failed: 'Authentication failed. Please try again.',
      user_info_failed: 'Cannot retrieve your information. Please try again.',
      user_creation_failed: 'Account creation failed. Please contact support.',
      token_generation_failed: 'Authentication failed. Please try again.',
    };

    const userMessage = errorMessages[errorType] || 'Login failed. Please try again.';
    
    // Show error toast
    toast.error(userMessage);
    
    // Log for debugging
    console.error('OAuth Error:', { errorType, errorMessage });
    
    // Redirect to login page
    setTimeout(() => {
      router.push('/login');
    }, 2000);
  };

  const handleOAuthSuccess = (
    token: string,
    userId: string,
    email: string,
    displayName: string
  ) => {
    // Save token to localStorage
    localStorage.setItem('token', token);
    localStorage.setItem('user_id', userId);
    localStorage.setItem('email', email);
    localStorage.setItem('display_name', displayName);
    
    // Show success message
    toast.success(`Welcome, ${displayName}!`);
    
    // Check setup status
    checkSetupStatus(token);
  };

  const checkSetupStatus = async (token: string) => {
    try {
      const response = await fetch('/api/v1/users/setup-status', {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      const data = await response.json();
      
      if (data.data.setup_completed) {
        router.push('/home');
      } else {
        router.push('/setup');
      }
    } catch (error) {
      console.error('Failed to check setup status:', error);
      router.push('/setup'); // Default to setup page
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
        <p className="mt-4 text-gray-600">Processing login...</p>
      </div>
    </div>
  );
};

export default OAuthCallback;
```

---

### Vue.js Example

```vue
<!-- pages/callback.vue -->
<template>
  <div class="callback-container">
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <p>Processing login...</p>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useToast } from 'vue-toastification';

const router = useRouter();
const route = useRoute();
const toast = useToast();
const loading = ref(true);

onMounted(() => {
  const { token, error, message, user_id, email, display_name } = route.query;
  
  if (error) {
    handleOAuthError(error, message);
  } else if (token) {
    handleOAuthSuccess(token, user_id, email, display_name);
  } else {
    toast.error('Invalid callback');
    router.push('/login');
  }
});

const handleOAuthError = (errorType, errorMessage) => {
  const errorMessages = {
    missing_parameters: 'Login incomplete. Please try again.',
    invalid_state: 'Session expired. Please try again.',
    token_exchange_failed: 'Authentication failed. Please try again.',
    user_info_failed: 'Cannot retrieve your information. Please try again.',
    user_creation_failed: 'Account creation failed. Please contact support.',
    token_generation_failed: 'Authentication failed. Please try again.',
  };
  
  const userMessage = errorMessages[errorType] || 'Login failed. Please try again.';
  toast.error(userMessage);
  
  console.error('OAuth Error:', { errorType, errorMessage });
  
  setTimeout(() => {
    router.push('/login');
  }, 2000);
};

const handleOAuthSuccess = async (token, userId, email, displayName) => {
  // Save to localStorage
  localStorage.setItem('token', token);
  localStorage.setItem('user_id', userId);
  localStorage.setItem('email', email);
  localStorage.setItem('display_name', displayName);
  
  toast.success(`Welcome, ${displayName}!`);
  
  // Check setup status
  try {
    const response = await fetch('/api/v1/users/setup-status', {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    
    const data = await response.json();
    
    if (data.data.setup_completed) {
      router.push('/home');
    } else {
      router.push('/setup');
    }
  } catch (error) {
    console.error('Failed to check setup status:', error);
    router.push('/setup');
  }
};
</script>

<style scoped>
.callback-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
}

.loading {
  text-align: center;
}

.spinner {
  width: 48px;
  height: 48px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #3498db;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>
```

---

## 🔧 Configuration

### Backend Environment Variables

```bash
# .env
FRONTEND_URL=http://localhost:8081

# Production
FRONTEND_URL=https://tindertrip.phitik.com
```

### Frontend Routes Required

1. **`/callback`** - OAuth callback handler (ต้องมีแน่นอน)
2. **`/login`** - Redirect เมื่อ error
3. **`/setup`** - Redirect เมื่อ setup ไม่เสร็จ
4. **`/home`** - Redirect เมื่อ setup เสร็จแล้ว

---

## 🎯 Error Handling Best Practices

### 1. **Show User-Friendly Messages**
```typescript
// Bad ❌
toast.error(error); // "token_exchange_failed"

// Good ✅
const message = errorMessages[error] || 'Login failed. Please try again.';
toast.error(message);
```

### 2. **Log for Debugging**
```typescript
console.error('OAuth Error:', {
  error: errorType,
  message: errorMessage,
  timestamp: new Date().toISOString(),
  userAgent: navigator.userAgent
});
```

### 3. **Provide Retry Option**
```typescript
toast.error(message, {
  action: {
    label: 'Try Again',
    onClick: () => router.push('/login')
  }
});
```

### 4. **Handle Network Errors**
```typescript
try {
  await checkSetupStatus(token);
} catch (error) {
  if (error.message === 'Network Error') {
    toast.error('Connection failed. Please check your internet.');
  } else {
    toast.error('Something went wrong. Please try again.');
  }
}
```

---

## 📊 Analytics & Monitoring

### Track OAuth Errors

```typescript
const handleOAuthError = (errorType: string, errorMessage: string) => {
  // Send to analytics
  analytics.track('OAuth Error', {
    error_type: errorType,
    error_message: errorMessage,
    timestamp: new Date().toISOString(),
    url: window.location.href
  });
  
  // Show user message
  // ...
};
```

---

## 🧪 Testing

### Test Success Callback
```
http://localhost:8081/callback?token=test_token&user_id=123&email=test@example.com&display_name=Test+User&provider=google&is_verified=true
```

### Test Error Callbacks
```
http://localhost:8081/callback?error=missing_parameters&message=Test+error
http://localhost:8081/callback?error=invalid_state&message=Session+expired
http://localhost:8081/callback?error=user_creation_failed&message=Database+error
```

---

## ⚠️ Important Notes

1. **URL Decode**: Parameters มาใน URL-encoded format ต้อง decode
2. **State Management**: Clear sensitive data หลัง redirect สำเร็จ
3. **Security**: ห้าม log token ใน production
4. **Timeout**: ให้ timeout callback page หาก processing นานเกิน 30 วินาที
5. **Mobile Deep Links**: อาจต้อง handle deep link scheme สำหรับ mobile app

---

**📞 Support**: หากมีปัญหา check logs ที่ Backend และ Frontend console

**Last Updated**: 2025-10-27

