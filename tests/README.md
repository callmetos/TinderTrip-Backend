# TinderTrip Backend Tests

โครงสร้างของ test suite สำหรับ TinderTrip Backend

## 📁 โครงสร้าง

```
tests/
├── unit/                    # Unit tests
│   ├── utils/              # Utils tests
│   │   ├── jwt_test.go
│   │   └── password_test.go
│   ├── middleware/         # Middleware tests
│   │   └── auth_test.go
│   ├── models/             # Model tests
│   │   ├── user_test.go
│   │   └── event_test.go
│   ├── service/            # Service tests
│   │   ├── auth_service_test.go
│   │   ├── event_service_test.go
│   │   └── user_service_test.go
│   └── handlers/           # Handler tests
│       ├── auth_handler_test.go
│       └── event_handler_test.go
├── integration/            # Integration tests
│   └── api_test.go
├── mocks/                  # Mock implementations
│   └── (mock files)
└── README.md              # This file
```

## 🧪 การรัน Tests

### รันทุก tests
```bash
go test ./tests/...
```

### รัน tests แบบ verbose
```bash
go test -v ./tests/...
```

### รัน tests ที่เฉพาะเจาะจง
```bash
# Unit tests only
go test ./tests/unit/...

# Specific package
go test ./tests/unit/utils/
go test ./tests/unit/middleware/
go test ./tests/unit/models/
go test ./tests/unit/service/
go test ./tests/unit/handlers/

# Integration tests only
go test ./tests/integration/...
```

### รัน tests พร้อม coverage
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./tests/...

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Open in browser
open coverage.html  # macOS
xdg-open coverage.html  # Linux
start coverage.html  # Windows
```

### รัน specific test function
```bash
go test -v -run TestGenerateToken ./tests/unit/utils/
go test -v -run TestAuthMiddleware ./tests/unit/middleware/
```

### รัน tests แบบ parallel
```bash
go test -parallel 4 ./tests/...
```

### รัน benchmarks
```bash
go test -bench=. ./tests/unit/utils/
go test -bench=BenchmarkHashPassword ./tests/unit/utils/
```

## 📊 Test Coverage Goals

- **Utils**: 90%+
- **Middleware**: 90%+
- **Models**: 90%+
- **Services**: 80%+
- **Handlers**: 80%+
- **Integration**: 70%+

## 🔍 Test Types

### Unit Tests
ทดสอบแต่ละ component แยกกันโดยไม่ต้องพึ่งพา external dependencies

**ตัวอย่าง:**
- JWT token generation/validation
- Password hashing/verification
- Model methods
- Middleware functions

### Integration Tests
ทดสอบการทำงานร่วมกันของหลาย components

**ตัวอย่าง:**
- API endpoints
- Database operations
- Authentication flow
- Complete user journeys

## 📝 Test Naming Conventions

### Function Names
```go
func TestFunctionName(t *testing.T)              // Basic test
func TestFunctionName_Scenario(t *testing.T)     // Specific scenario
func BenchmarkFunctionName(b *testing.B)         // Benchmark test
```

### Test Cases
```go
tests := []struct {
    name    string  // Test case description
    input   Type    // Input data
    want    Type    // Expected output
    wantErr bool    // Expect error?
}{
    {
        name: "Valid input",
        input: ...,
        want: ...,
        wantErr: false,
    },
}
```

## 🛠️ Testing Tools & Libraries

### Core
- **testing** - Go standard testing package
- **github.com/stretchr/testify** - Assertions and mocks

### Assertions
```go
assert.Equal(t, expected, actual)
assert.NoError(t, err)
assert.True(t, condition)
assert.Contains(t, str, substr)
require.NoError(t, err)  // Stops test on failure
```

### HTTP Testing
```go
// Create test context
w := httptest.NewRecorder()
c, _ := gin.CreateTestContext(w)

// Create test request
req := httptest.NewRequest("GET", "/test", nil)
c.Request = req
```

## 🎯 Test Scenarios ที่ต้องครอบคลุม

### Authentication Tests
- [x] Token generation
- [x] Token validation
- [x] Token refresh
- [x] Token expiration
- [x] Token extraction from header
- [x] Password hashing
- [x] Password verification
- [x] Password strength validation
- [x] Auth middleware (valid/invalid tokens)
- [x] Optional auth middleware

### Model Tests
- [x] User model methods
- [x] Event model methods
- [x] Model validation
- [x] BeforeCreate hooks
- [x] TableName methods

### Service Tests
- [ ] User registration
- [ ] User login
- [ ] Password reset flow
- [ ] Email verification
- [ ] Event creation
- [ ] Event updates
- [ ] Event deletion
- [ ] Join/leave events
- [ ] Swipe functionality

### Handler Tests
- [ ] Request validation
- [ ] Response formatting
- [ ] Error handling
- [ ] Success scenarios
- [ ] Edge cases

### Integration Tests
- [ ] Complete auth flow
- [ ] Event lifecycle
- [ ] User profile management
- [ ] Chat functionality

## 🚀 การเขียน Test ใหม่

### 1. สร้างไฟล์ test
```bash
# ตั้งชื่อไฟล์ตามไฟล์ที่ test
# ถ้า test file: internal/service/auth_service.go
# สร้าง: tests/unit/service/auth_service_test.go
```

### 2. Import packages
```go
package service_test

import (
    "testing"
    "TinderTrip-Backend/internal/service"
    "github.com/stretchr/testify/assert"
)
```

### 3. เขียน test cases
```go
func TestServiceMethod(t *testing.T) {
    tests := []struct {
        name    string
        setup   func() // Setup test data
        input   Type
        want    Type
        wantErr bool
    }{
        // Test cases here
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## 🐛 Debug Tests

### Print debug info
```go
t.Logf("Debug info: %v", value)
```

### Skip tests
```go
t.Skip("Skipping this test temporarily")
```

### Run only specific tests
```bash
go test -v -run TestSpecificFunction
```

## 📈 Continuous Integration

Tests จะรันอัตโนมัติใน CI pipeline:
- ✅ On every push
- ✅ On pull requests
- ✅ Before deployment

## 💡 Best Practices

1. **Test Isolation**: แต่ละ test ต้องเป็นอิสระ
2. **Clear Names**: ใช้ชื่อที่อธิบายว่า test อะไร
3. **Arrange-Act-Assert**: จัดเรียง test logic ชัดเจน
4. **Mock External Dependencies**: ใช้ mocks สำหรับ database, API calls
5. **Test Edge Cases**: ทดสอบทั้ง happy path และ error cases
6. **Keep Tests Fast**: Unit tests ต้องรันเร็ว
7. **One Assertion Per Test**: พยายามให้แต่ละ test มี assertion หลักเดียว
8. **Clean Up**: ทำความสะอาด test data หลังแต่ละ test

## 📚 Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Testing in Go](https://go.dev/doc/tutorial/add-a-test)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)

## ✅ Test Checklist

เมื่อเขียน test ใหม่ ควรตรวจสอบว่า:

- [ ] ทดสอบ happy path (case ที่ทำงานถูกต้อง)
- [ ] ทดสอบ error cases (case ที่เกิด error)
- [ ] ทดสอบ edge cases (กรณีพิเศษ)
- [ ] ทดสอบ boundary conditions (ขอบเขต)
- [ ] มี assertions ที่ชัดเจน
- [ ] ไม่มี test ที่ depend กัน
- [ ] รันผ่านทุก test
- [ ] Coverage อยู่ในเกณฑ์ที่กำหนด

---

**Note**: Tests จะถูก update และเพิ่มเติมอยู่เสมอ ตรวจสอบ coverage report เป็นประจำ

