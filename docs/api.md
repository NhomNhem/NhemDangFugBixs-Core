# Hollow Wilds Backend API Documentation

**Base URL**: `http://localhost:8080/api/v1` (development)  
**Production URL**: `https://api.yourgame.com/api/v1`

## Authentication

Most endpoints require JWT authentication.

**Header:**
```
Authorization: Bearer <jwt_token>
```

Get JWT token via `/auth/login` endpoint.

---

## Endpoints

### Health Check

**GET /health**

Check if server is running.

**Response:**
```json
{
  "status": "ok",
  "message": "Hollow Wilds Backend is running",
  "version": "1.0.0"
}
```

---

### Authentication

**POST /api/v1/auth/login**

Exchange PlayFab session token for JWT.

**Headers:**
```
X-PlayFab-SessionToken: <playfab_session_token>
```

**Request Body:**
```json
{
  "playfabId": "ABC123DEF456",
  "displayName": "Player123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "userId": "uuid-1234-5678",
      "playfabId": "ABC123DEF456",
      "displayName": "Player123",
      "gold": 1000,
      "diamonds": 50
    },
    "expiresIn": 86400
  }
}
```

---

### Level Completion

**POST /api/v1/levels/complete**

Submit level completion for validation.

**Authentication:** Required

**Request Body:**
```json
{
  "levelId": "map1-5",
  "mapId": "map1",
  "timeSeconds": 58.3,
  "finalHP": 75.0,
  "maxHP": 100.0,
  "dashCount": 12,
  "counterCount": 8,
  "vulnerableKills": 6
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "validation": {
      "isValid": true,
      "starsEarned": 2
    },
    "rewards": {
      "goldEarned": 100,
      "newTotalGold": 1250
    },
    "progression": {
      "levelUnlocked": "map1-6"
    }
  }
}
```

---

### Talent Upgrade

**POST /api/v1/talents/upgrade**

Upgrade a talent.

**Authentication:** Required

**Request Body:**
```json
{
  "talentType": "health",
  "targetLevel": 6
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "newLevel": 6,
    "goldSpent": 249,
    "newGoldBalance": 901
  }
}
```

---

### Payments

**POST /api/v1/payments/create-session**

Create Stripe checkout session.

**Authentication:** Required

**Request Body:**
```json
{
  "sku": "diamonds_650",
  "successUrl": "unitydl://payment-success",
  "cancelUrl": "unitydl://payment-cancel"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "sessionId": "cs_test_1234567890",
    "checkoutUrl": "https://checkout.stripe.com/c/pay/...",
    "expiresAt": "2026-03-04T12:28:00Z"
  }
}
```

---

### Analytics

**POST /api/v1/analytics/events**

Submit analytics events (batch).

**Authentication:** Required

**Request Body:**
```json
{
  "events": [
    {
      "eventType": "level_start",
      "timestamp": "2026-03-04T11:30:00Z",
      "properties": {
        "levelId": "map1-5",
        "playerGold": 1150
      }
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "eventsProcessed": 1
  }
}
```

---

## Error Responses

All errors follow this format:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message"
  }
}
```

**Common Error Codes:**
- `INVALID_TOKEN` (401) - JWT invalid or expired
- `VALIDATION_FAILED` (400) - Request validation failed
- `INSUFFICIENT_FUNDS` (400) - Not enough gold/diamonds
- `NOT_FOUND` (404) - Resource not found
- `RATE_LIMIT_EXCEEDED` (429) - Too many requests
- `INTERNAL_ERROR` (500) - Server error

---

For detailed workflow diagrams, see [Phase 0 Documentation](../../../.copilot/session-state/*/files/backend-workflows.md)
