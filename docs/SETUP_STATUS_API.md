# 📋 Setup Status API Documentation

## Overview
API endpoint สำหรับเช็คว่า user ทำ initial profile setup เสร็จหรือยัง

---

## 🔗 Endpoint

### **GET** `/api/v1/users/setup-status`

เช็คสถานะการทำ profile setup ของ user ที่ login อยู่

---

## 🔐 Authentication
**Required**: Bearer Token (JWT)

```
Authorization: Bearer <your_jwt_token>
```

---

## 📥 Request

### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### No Request Body
GET request ไม่ต้องส่ง body

---

## 📤 Response

### ✅ Success Response (200 OK)

```json
{
  "data": {
    "setup_completed": true
  },
  "message": "Setup status retrieved successfully"
}
```

**หรือ** ถ้ายังไม่ได้ทำ setup:

```json
{
  "data": {
    "setup_completed": false
  },
  "message": "Setup status retrieved successfully"
}
```

### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `data.setup_completed` | `boolean` | `true` = setup เสร็จแล้ว, `false` = ยังไม่เสร็จ |
| `message` | `string` | Status message |

---

## ❌ Error Responses

### 401 Unauthorized
ไม่มี token หรือ token ไม่ valid

```json
{
  "error": "Unauthorized",
  "message": "User not authenticated"
}
```

### 500 Internal Server Error
Database error

```json
{
  "error": "Failed to check setup status",
  "message": "database error: ..."
}
```

---

## 🎯 Setup Completion Logic

User ถือว่า **setup เสร็จ** เมื่อ:

มี **อย่างน้อยหนึ่งอย่าง** ใน profile:
- ✅ **Bio** ที่ไม่ว่างเปล่า
- ✅ **Gender** (male, female, other, prefer_not_to_say)
- ✅ **Languages** ที่ไม่ว่างเปล่า

ถ้าไม่มีอะไรเลย → `setup_completed = false`

---

## 📝 Example Usage

### JavaScript (Axios)

```javascript
const checkSetupStatus = async () => {
  try {
    const response = await axios.get(
      'https://api.tindertrip.phitik.com/api/v1/users/setup-status',
      {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      }
    );
    
    const { setup_completed } = response.data.data;
    
    if (!setup_completed) {
      // Redirect to setup page
      window.location.href = '/setup';
    } else {
      // Continue to main app
      window.location.href = '/home';
    }
  } catch (error) {
    console.error('Failed to check setup status:', error);
  }
};
```

### JavaScript (Fetch API)

```javascript
fetch('https://api.tindertrip.phitik.com/api/v1/users/setup-status', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
})
  .then(response => response.json())
  .then(data => {
    if (data.data.setup_completed) {
      console.log('Setup completed! ✅');
    } else {
      console.log('Please complete setup ⚠️');
    }
  })
  .catch(error => console.error('Error:', error));
```

### TypeScript (React with Axios)

```typescript
import axios from 'axios';

interface SetupStatusResponse {
  data: {
    setup_completed: boolean;
  };
  message: string;
}

const useSetupStatus = () => {
  const [setupCompleted, setSetupCompleted] = useState<boolean | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const checkStatus = async () => {
      try {
        const { data } = await axios.get<SetupStatusResponse>(
          '/api/v1/users/setup-status'
        );
        setSetupCompleted(data.data.setup_completed);
      } catch (error) {
        console.error('Failed to check setup status:', error);
      } finally {
        setLoading(false);
      }
    };

    checkStatus();
  }, []);

  return { setupCompleted, loading };
};

// Usage in component
const App = () => {
  const { setupCompleted, loading } = useSetupStatus();

  if (loading) return <Spinner />;
  
  if (setupCompleted === false) {
    return <Navigate to="/setup" />;
  }

  return <MainApp />;
};
```

---

## 🔄 Related Endpoints

### Update Profile (to complete setup)
```
PUT /api/v1/users/profile
```

เมื่อ user update profile ด้วย bio, gender, หรือ languages แล้ว
การเรียก `/setup-status` ครั้งต่อไปจะได้ `setup_completed: true`

**Example Request:**
```json
{
  "bio": "I love traveling and trying new foods!",
  "gender": "male",
  "languages": "English, Thai"
}
```

---

## 🎨 Frontend Integration Flow

```mermaid
graph TD
    A[User Login] --> B[Get JWT Token]
    B --> C[Check Setup Status]
    C --> D{setup_completed?}
    D -->|true| E[Go to Home Page]
    D -->|false| F[Redirect to Setup Page]
    F --> G[User Fills Profile]
    G --> H[Update Profile API]
    H --> I[Check Status Again]
    I --> E
```

---

## ⚠️ Important Notes

1. **Authentication Required**: ต้อง login ก่อนถึงจะเรียกได้
2. **Token in Header**: ใส่ JWT token ใน `Authorization` header
3. **Check After Login**: เรียก API นี้ทันทีหลัง login เพื่อ redirect user
4. **Cache Strategy**: อาจ cache ผลลัพธ์ไว้ใน session (แต่อย่าเกิน 5 นาที)
5. **Error Handling**: Handle 401 → redirect to login

---

## 🧪 Testing

### Using cURL

```bash
# Get setup status
curl -X GET "https://api.tindertrip.phitik.com/api/v1/users/setup-status" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Expected Response
```json
{
  "data": {
    "setup_completed": false
  },
  "message": "Setup status retrieved successfully"
}
```

---

## 📞 Support

หากมีปัญหาการใช้งาน API นี้:
- ตรวจสอบว่า token ยังใช้งานได้อยู่
- ตรวจสอบ CORS settings
- ดู API logs สำหรับ error details

**Base URL**: `https://api.tindertrip.phitik.com`

---

## ✅ Checklist for Frontend

- [ ] เพิ่ม API call หลัง login สำเร็จ
- [ ] Handle `setup_completed = false` → redirect to setup page
- [ ] Handle `setup_completed = true` → redirect to home
- [ ] Handle 401 error → redirect to login
- [ ] Show loading state ระหว่างเช็ค
- [ ] Test with different user states

---

**Status**: ✅ Ready to use
**Last Updated**: 2025-10-27

