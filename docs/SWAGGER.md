# 📊 Swagger API Documentation

## 🌐 Access Swagger UI

### Production:
```
https://gamefeel-backend.fly.dev/swagger/index.html
```

### Local Development:
```
http://localhost:8080/swagger/index.html
```

---

## 🔧 How to Update Swagger Docs

### 1️⃣ Add/Update API Annotations

Edit your handler functions with Swagger comments:

```go
// @Summary Your endpoint summary
// @Description Detailed description
// @Tags YourTag
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param request body models.YourRequest true "Request body"
// @Success 200 {object} models.APIResponse{data=models.YourResponse}
// @Failure 400 {object} models.APIResponse{error=models.APIError}
// @Router /your/endpoint [post]
// @Security BearerAuth
func (h *YourHandler) YourEndpoint(c *fiber.Ctx) error {
    // ...
}
```

### 2️⃣ Regenerate Docs

```bash
# Install swag CLI (first time only)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
cd I:\unityVers\GameFeel-Backend
swag init -g cmd/server/main.go --output ./docs
```

### 3️⃣ Deploy

```bash
git add .
git commit -m "docs: Update Swagger documentation"
git push origin main

# Auto-deploy via GitHub Actions or manual:
fly deploy
```

---

## 📚 Current APIs

### 🔐 Authentication
- **POST** `/api/v1/auth/login` - Login với PlayFab token

### 🎮 Levels
- **POST** `/api/v1/levels/complete` - Complete level với anti-cheat

### ⭐ Talents
- **GET** `/api/v1/talents` - Get user talents
- **POST** `/api/v1/talents/upgrade` - Upgrade talent

---

## 🧪 Testing với Swagger UI

### Step 1: Login
1. Expand **Authentication** → **POST /auth/login**
2. Click **"Try it out"**
3. Fill in:
   - Header: `X-PlayFab-SessionToken: YOUR_PLAYFAB_TOKEN`
   - Body:
     ```json
     {
       "playfabId": "YOUR_PLAYFAB_ID",
       "displayName": "TestUser"
     }
     ```
4. Click **"Execute"**
5. Copy JWT token từ response

### Step 2: Test Protected Endpoints
1. Click **"Authorize"** button (🔒 icon ở trên)
2. Paste: `Bearer YOUR_JWT_TOKEN`
3. Click **"Authorize"**
4. Bây giờ bạn có thể test Level Complete và Talent APIs

---

## 🔒 Security Definitions

### BearerAuth (JWT)
- **Type**: API Key
- **In**: Header
- **Name**: `Authorization`
- **Format**: `Bearer {jwt_token}`

### PlayFabToken
- **Type**: API Key
- **In**: Header
- **Name**: `X-PlayFab-SessionToken`
- **Required for**: `/auth/login` only

---

## 📖 Swagger Annotations Reference

### Common Annotations:

```go
// @Summary        - Short description (< 120 chars)
// @Description    - Detailed description
// @Tags           - Group endpoints (Auth, Levels, Talents)
// @Accept         - Input format (json, xml, form)
// @Produce        - Output format (json, xml)
// @Param          - Request parameters
// @Success        - Success response
// @Failure        - Error response
// @Router         - Endpoint path and method
// @Security       - Required authentication
```

### Parameter Types:

```go
// @Param name query string true "Description"      // ?name=value
// @Param id path int true "User ID"                // /users/{id}
// @Param Auth header string true "JWT token"       // Authorization: Bearer xxx
// @Param body body models.Request true "Body"      // JSON body
```

---

## 🎨 Custom Configuration

Edit `cmd/server/main.go`:

```go
// @title Hollow Wilds Backend API
// @version 1.0
// @description Your description
// @host gamefeel-backend.fly.dev
// @BasePath /api/v1
```

---

## 📝 Tips

✅ **DO**:
- Keep summaries concise
- Use proper HTTP status codes
- Document all parameters
- Group related endpoints with Tags
- Test changes locally before deploying

❌ **DON'T**:
- Forget to regenerate docs after changes
- Use sensitive data in examples
- Mix different API versions in same docs

---

## 🔗 Resources

- **Swagger/OpenAPI Spec**: https://swagger.io/specification/
- **swaggo/swag**: https://github.com/swaggo/swag
- **Fiber Swagger**: https://github.com/gofiber/swagger

---

## 🐛 Troubleshooting

### Swagger UI shows 404
- Check if docs folder is in Docker image (not in `.dockerignore`)
- Verify route is registered: `app.Get("/swagger/*", swagger.HandlerDefault)`

### Changes not showing
- Regenerate docs: `swag init -g cmd/server/main.go --output ./docs`
- Clear browser cache
- Redeploy

### Build fails "no required module provides package docs"
- Ensure docs folder is NOT in `.dockerignore`
- Run `go mod tidy`
- Commit all docs files (docs.go, swagger.json, swagger.yaml)
