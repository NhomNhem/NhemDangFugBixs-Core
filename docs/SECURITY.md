# Security Guide - Hollow Wilds Backend

## Overview

This document covers security features and best practices for the GameFeel backend.

---

## Security Features

### 1. PlayFab Token Validation ✅

**Status**: Implemented  
**Location**: `internal/services/auth_service.go`

All login requests must provide a valid PlayFab session token. The backend validates this token by calling PlayFab's Client API.

**How it works:**
1. Unity client logs in with PlayFab (gets session token)
2. Client sends token in `X-PlayFab-SessionToken` header
3. Backend calls PlayFab API to verify token
4. If valid, backend generates JWT token for subsequent requests

**Configuration:**
```bash
# Required for production
PLAYFAB_TITLE_ID=your-title-id
```

**Development mode:**
- If `PLAYFAB_TITLE_ID` is not set, validation is skipped
- Useful for local testing without PlayFab connection

**Testing:**
```bash
# Valid token (should succeed)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-PlayFab-SessionToken: <REAL_TOKEN>" \
  -d '{"playfabId": "ABC123"}'

# Invalid token (should return 401)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-PlayFab-SessionToken: FAKE_TOKEN" \
  -d '{"playfabId": "ABC123"}'
```

---

### 2. Rate Limiting ✅

**Status**: Implemented  
**Location**: `cmd/server/main.go`

Global rate limiter protects against abuse and DDoS attacks.

**Default limits:**
- 100 requests per minute per IP address
- Applies to all API endpoints

**Configuration:**
```bash
# Optional (defaults shown)
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_DURATION=60s
```

**Response when rate limit exceeded:**
```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests. Please try again later."
  }
}
```

**HTTP Status**: `429 Too Many Requests`

**Testing:**
```bash
# Send 101 requests in 60 seconds
for i in {1..101}; do
  curl http://localhost:8080/health
done
# 101st request should return 429
```

---

### 3. CORS Protection ✅

**Status**: Implemented  
**Location**: `cmd/server/main.go`

Cross-Origin Resource Sharing (CORS) restricts which domains can access the API.

**Configuration:**
```bash
# Production (restrict to your game domain)
ALLOWED_ORIGINS=https://yourgame.com,https://www.yourgame.com

# Development (allow localhost)
ALLOWED_ORIGINS=http://localhost:*,http://127.0.0.1:*

# Default (allow all - NOT RECOMMENDED for production)
ALLOWED_ORIGINS=*
```

**Allowed headers:**
- `Origin`
- `Content-Type`
- `Accept`
- `Authorization` (JWT tokens)
- `X-PlayFab-SessionToken` (PlayFab tokens)

**Testing:**
```bash
# Allowed origin (should succeed)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Origin: https://yourgame.com" \
  -H "Content-Type: application/json" \
  -d '{"playfabId": "ABC123"}'

# Blocked origin (should fail with CORS error)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Origin: https://malicious-site.com" \
  -H "Content-Type: application/json" \
  -d '{"playfabId": "ABC123"}'
```

---

### 4. JWT Authentication ✅

**Status**: Previously implemented  
**Location**: `internal/middleware/auth.go`

All protected endpoints require a valid JWT token in the `Authorization` header.

**Token format:**
```
Authorization: Bearer <JWT_TOKEN>
```

**Token lifetime**: 24 hours

**Claims:**
- `user_id`: Database user ID (UUID)
- `playfab_id`: PlayFab user ID
- `exp`: Expiration timestamp
- `iat`: Issued at timestamp

**Configuration:**
```bash
# Required - change in production!
JWT_SECRET=your-secret-key-here-change-in-production
```

**Protected endpoints:**
- `POST /api/v1/levels/complete`
- `GET /api/v1/talents`
- `POST /api/v1/talents/upgrade`

---

## Security Checklist

### Production Deployment

- [ ] Set `PLAYFAB_TITLE_ID` (enable token validation)
- [ ] Set `ALLOWED_ORIGINS` to your game domain(s)
- [ ] Change `JWT_SECRET` to strong random value
- [ ] Verify rate limiting is working
- [ ] Test with valid/invalid PlayFab tokens
- [ ] Review CORS allowed origins
- [ ] Enable HTTPS (handled by Fly.io)

### Environment Variables

**Required:**
```bash
DATABASE_URL=<supabase-connection-string>
JWT_SECRET=<strong-random-secret>
PLAYFAB_TITLE_ID=<playfab-title-id>
```

**Recommended:**
```bash
ALLOWED_ORIGINS=https://yourgame.com
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_DURATION=60s
```

---

## Attack Vectors & Mitigations

### 1. Fake Login Attempts

**Attack**: Attacker tries to login with fake PlayFab tokens  
**Mitigation**: ✅ PlayFab token validation rejects invalid tokens  
**Status**: Protected

### 2. Brute Force Attacks

**Attack**: Attacker sends thousands of requests  
**Mitigation**: ✅ Rate limiting (100 req/min per IP)  
**Status**: Protected

### 3. Unauthorized API Access

**Attack**: Attacker accesses protected endpoints without JWT  
**Mitigation**: ✅ JWT middleware returns 401 Unauthorized  
**Status**: Protected

### 4. Cross-Site Request Forgery (CSRF)

**Attack**: Malicious site makes API requests on behalf of user  
**Mitigation**: ✅ CORS restricts origins  
**Status**: Protected (when ALLOWED_ORIGINS configured)

### 5. Token Theft

**Attack**: Attacker steals JWT token from client  
**Mitigation**: 
- ✅ 24-hour token expiry (limits exposure)
- ✅ HTTPS encryption (prevents MITM)
- ⚠️ Consider: Token refresh mechanism  
**Status**: Partially protected

### 6. SQL Injection

**Attack**: Attacker injects SQL in API parameters  
**Mitigation**: ✅ Parameterized queries (pgx)  
**Status**: Protected

---

## Security Roadmap

### Implemented ✅
- [x] PlayFab token validation
- [x] Rate limiting (global)
- [x] CORS protection
- [x] JWT authentication
- [x] Parameterized SQL queries
- [x] HTTPS (Fly.io)

### Future Enhancements 🔮
- [ ] Per-endpoint rate limiting (stricter for login)
- [ ] Token refresh mechanism
- [ ] IP blacklist/whitelist
- [ ] Request signing (HMAC)
- [ ] Audit logging
- [ ] Anomaly detection
- [ ] DDoS protection (Cloudflare)

---

## Incident Response

### If PlayFab Token is Compromised

1. **Immediate**: No action needed - tokens expire automatically
2. **Verify**: Check if unauthorized access occurred
3. **Rotate**: Change `JWT_SECRET` if JWT tokens compromised

### If JWT Secret is Leaked

1. **Immediate**: Generate new `JWT_SECRET`
2. **Deploy**: Update Fly.io secrets
3. **Force**: All users must re-login

```bash
# Generate new secret
openssl rand -base64 32

# Update Fly.io
fly secrets set JWT_SECRET=<new-secret>
```

### If Database Credentials Leaked

1. **Immediate**: Rotate Supabase password
2. **Update**: All services with new `DATABASE_URL`
3. **Review**: Check for unauthorized database access

---

## Monitoring

### Key Metrics to Track

1. **401 Unauthorized** - Failed authentication attempts
2. **429 Too Many Requests** - Rate limit hits
3. **Failed PlayFab validations** - Potential attack attempts
4. **Unusual traffic patterns** - DDoS indicators

### Logging

Currently logs to stdout (visible in `fly logs`):
```
[2026-03-04] 401 - POST /api/v1/auth/login (15ms)
PlayFab token validation failed: invalid token
```

**Future**: Structured logging with log levels

---

## References

- [PlayFab Authentication](https://docs.microsoft.com/en-us/gaming/playfab/features/authentication/)
- [OWASP API Security Top 10](https://owasp.org/www-project-api-security/)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [Rate Limiting Patterns](https://cloud.google.com/architecture/rate-limiting-strategies-techniques)

---

## Support

For security concerns, contact the development team.

**Do not** disclose security vulnerabilities publicly.
