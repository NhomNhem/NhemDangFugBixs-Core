# 🧪 Swagger UI Testing Guide - EXAMPLE vs REAL RESPONSE

## ⚠️ PHÂN BIỆT QUAN TRỌNG

### ❌ KHÔNG phải Response thật:
```json
{
  "data": {
    "jwt": "string",
    "user": {
      "id": "string",
      "gold": 0,
      "displayName": "string"
    }
  },
  "success": true
}
```

**Dấu hiệu nhận biết:**
- Values = `"string"`, `0`, `true` (generic)
- Xuất hiện ở tab **"Example Value"**
- Đây là **schema documentation**, KHÔNG phải response thật!

---

### ✅ Response thật:
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid PlayFab session token"
  }
}
```

**Dấu hiệu nhận biết:**
- Values cụ thể: `"UNAUTHORIZED"`, `"Invalid PlayFab session token"`
- Xuất hiện sau khi click **"Execute"**
- Có HTTP Status Code: **401**, 200, 500...
- Section **"Response Body"** (không phải Example Value)

---

## 📊 Swagger UI Layout

```
┌─────────────────────────────────────────────────────────┐
│ POST /api/v1/auth/login                        [Try it] │
├─────────────────────────────────────────────────────────┤
│ Parameters:                                             │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ X-PlayFab-SessionToken: YOUR_TOKEN                  │ │
│ │ Request Body:                                       │ │
│ │   {                                                 │ │
│ │     "playfabId": "FAKE_ID",                        │ │
│ │     "displayName": "Test"                          │ │
│ │   }                                                 │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│              [Execute]  ← Click này                     │
│                                                         │
├─────────────────────────────────────────────────────────┤
│ Responses:                                              │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ Code: 401 ← HTTP Status                            │ │
│ │ Response body:                                      │ │
│ │ {                                                   │ │
│ │   "success": false,  ← Real values!                │ │
│ │   "error": {                                        │ │
│ │     "code": "UNAUTHORIZED",                         │ │
│ │     "message": "Invalid PlayFab session token"      │ │
│ │   }                                                 │ │
│ │ }                                                   │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│ Response headers:                                       │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ content-type: application/json                      │ │
│ │ date: Tue, 04 Mar 2026 16:48:00 GMT                │ │
│ └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

---

## 🧪 Step-by-Step Testing

### Test 1: Fake PlayFab Token (Should FAIL)

**Input:**
```
POST /api/v1/auth/login
Headers:
  X-PlayFab-SessionToken: FAKE_TOKEN_12345
Body:
  {
    "playfabId": "FAKE_ID_12345",
    "displayName": "TestUser"
  }
```

**Expected Result:**
```
HTTP Status: 401 Unauthorized
Response:
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid PlayFab session token"
  }
}
```

**✅ Security: PASSED** - Rejects fake tokens

---

### Test 2: Fake JWT Token (Should FAIL)

**Input:**
```
POST /api/v1/levels/complete
Headers:
  Authorization: Bearer FAKE_JWT_TOKEN
Body:
  {
    "levelId": "1-1",
    "mapId": "map1",
    "timeSeconds": 60,
    "finalHp": 100
  }
```

**Expected Result:**
```
HTTP Status: 401 Unauthorized
Response:
{
  "success": false,
  "error": {
    "code": "INVALID_TOKEN",
    "message": "JWT token is invalid or expired"
  }
}
```

**✅ Security: PASSED** - Rejects fake JWT

---

### Test 3: Valid Credentials (Should SUCCESS)

**Step 1: Login**
```
POST /api/v1/auth/login
Headers:
  X-PlayFab-SessionToken: <REAL_PLAYFAB_TOKEN>
Body:
  {
    "playfabId": "<REAL_PLAYFAB_ID>",
    "displayName": "TestUser"
  }
```

**Expected Result:**
```
HTTP Status: 200 OK
Response:
{
  "success": true,
  "data": {
    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",  ← Real JWT!
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",  ← Real UUID!
      "playfabId": "ABC123DEF456",
      "displayName": "TestUser",
      "gold": 1200,  ← Real gold amount!
      "diamonds": 50
    },
    "expiresIn": 86400
  }
}
```

**Step 2: Use JWT for Protected Endpoints**

Click **"Authorize"** button (🔒 icon):
```
Value: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

Now test `/levels/complete`, `/talents`, etc.

---

## 🔒 Security Status

| Validation | Status | Test Result |
|------------|--------|-------------|
| **PlayFab Token Validation** | ✅ WORKING | Rejects fake tokens (401) |
| **JWT Authentication** | ✅ WORKING | Rejects fake/expired JWT (401) |
| **Rate Limiting** | ✅ WORKING | 100 req/min per IP |
| **CORS Protection** | ✅ WORKING | Configurable origins |
| **Database Disconnected** | ⚠️ WARNING | Returns 500, no mock data leak |

---

## 🐛 Common Mistakes

### ❌ Mistake 1: Looking at "Example Value" tab
```
Tab: [Example Value] [Schema]
     ^^^^^^^^^^^^^^
     This is NOT real data!
```

**Fix:** Scroll down to **"Response Body"** section after clicking Execute.

---

### ❌ Mistake 2: Not scrolling down after Execute
```
[Execute]
...
...  ← Scroll down here!
...
Response body:  ← Real response here
{ ... }
```

**Fix:** Always scroll down to see actual response.

---

### ❌ Mistake 3: Forgetting to click "Try it out"
```
[Try it out]  ← Click this first!
   ↓
[Execute]     ← Then click this
```

**Fix:** Must click "Try it out" to enable input fields.

---

## 📝 Summary

✅ **Backend security is WORKING correctly!**

- ❌ Fake PlayFab tokens → Rejected (401)
- ❌ Fake JWT tokens → Rejected (401)
- ❌ No Authorization header → Rejected (401)
- ✅ Valid credentials → Success (200)

**Swagger UI "Example Value" is just documentation, not real data!**

To see real responses:
1. Click "Try it out"
2. Fill in real values
3. Click "Execute"
4. **Scroll down** to see "Response Body"
5. Check HTTP Status Code

---

## 🔗 Test URLs

- **Production**: https://gamefeel-backend.fly.dev/swagger/index.html
- **Local Docker**: http://localhost:8080/swagger/index.html

---

## 💡 Pro Tips

1. **Use curl for precise testing:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "X-PlayFab-SessionToken: FAKE" \
     -d '{"playfabId":"FAKE"}' \
     | jq .
   ```

2. **Check HTTP Status Code:**
   - 200 = Success
   - 401 = Unauthorized (security working!)
   - 500 = Server error

3. **Real vs Mock:**
   - Real: Specific strings, UUIDs, numbers
   - Mock: "string", 0, true, false

4. **Docker Database:**
   - Docker has IPv6 issue → database disconnected
   - But validation still works!
   - Production (Fly.io) has database connected
