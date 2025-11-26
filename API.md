# 📖 Dejavu API Documentation

REST API documentation untuk Dejavu platform.

## Base URL

```
Development: http://localhost:8080/api
Production: https://api.dejavu.id/api
```

## Authentication

Semua endpoint (kecuali auth) membutuhkan JWT token di header:

```
Authorization: Bearer <token>
```

---

## Authentication

### Register

Create a new user account.

**Endpoint:** `POST /auth/register`

**Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:** `201 Created`
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Login

Login to existing account.

**Endpoint:** `POST /auth/login`

**Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:** `200 OK`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## Projects

### List Projects

Get all projects for authenticated user.

**Endpoint:** `GET /projects`

**Response:** `200 OK`
```json
[
  {
    "id": "uuid",
    "user_id": "uuid",
    "name": "My Project",
    "repo_url": "https://github.com/user/repo",
    "build_command": "npm run build",
    "output_dir": "dist",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

### Create Project

Create a new project.

**Endpoint:** `POST /projects`

**Body:**
```json
{
  "name": "My Project",
  "repo_url": "https://github.com/user/repo",
  "build_command": "npm run build",
  "output_dir": "dist"
}
```

**Response:** `201 Created`
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "name": "My Project",
  "repo_url": "https://github.com/user/repo",
  "build_command": "npm run build",
  "output_dir": "dist",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Get Project

Get project details by ID.

**Endpoint:** `GET /projects/:id`

**Response:** `200 OK`
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "name": "My Project",
  "repo_url": "https://github.com/user/repo",
  "build_command": "npm run build",
  "output_dir": "dist",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Update Project

Update project configuration.

**Endpoint:** `PUT /projects/:id`

**Body:**
```json
{
  "name": "Updated Name",
  "build_command": "yarn build"
}
```

**Response:** `200 OK`
```json
{
  "message": "Project updated successfully"
}
```

### Delete Project

Delete a project and all its deployments.

**Endpoint:** `DELETE /projects/:id`

**Response:** `200 OK`
```json
{
  "message": "Project deleted successfully"
}
```

---

## Deployments

### Trigger Deployment

Start a new deployment for a project.

**Endpoint:** `POST /deploy`

**Body:**
```json
{
  "project_id": "uuid",
  "commit_hash": "abc123" // optional
}
```

**Response:** `201 Created`
```json
{
  "id": "uuid",
  "project_id": "uuid",
  "status": "pending",
  "subdomain": "app-xyz123",
  "commit_hash": "abc123",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Get Deployment Status

Get deployment details and status.

**Endpoint:** `GET /deploy/:id`

**Response:** `200 OK`
```json
{
  "id": "uuid",
  "project_id": "uuid",
  "status": "ready",
  "subdomain": "app-xyz123",
  "image_url": "registry.dejavu.id/dejavu/project:tag",
  "commit_hash": "abc123",
  "build_logs": "Building...\nSuccess!",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:05:00Z"
}
```

**Status values:**
- `pending` - Waiting to start
- `building` - Building application
- `deploying` - Deploying to Kubernetes
- `ready` - Live and accessible
- `error` - Deployment failed

### Stream Deployment Logs

Get real-time build logs via WebSocket.

**Endpoint:** `GET /deploy/:id/logs` (WebSocket)

**Example (JavaScript):**
```javascript
const ws = new WebSocket('ws://localhost:8080/api/deploy/uuid/logs')

ws.onmessage = (event) => {
  console.log('Log:', event.data)
}
```

---

## Billing

### Get Balance

Get current credit balance.

**Endpoint:** `GET /billing/balance`

**Response:** `200 OK`
```json
{
  "balance": 100.00
}
```

### Get Usage History

Get billing usage history.

**Endpoint:** `GET /billing/usage`

**Response:** `200 OK`
```json
[
  {
    "id": "uuid",
    "user_id": "uuid",
    "deployment_id": "uuid",
    "type": "deployment",
    "amount": 0.50,
    "description": "Deployment charge",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

### Add Credits

Add credits to account.

**Endpoint:** `POST /billing/credits`

**Body:**
```json
{
  "amount": 50.00
}
```

**Response:** `200 OK`
```json
{
  "message": "Credits added successfully"
}
```

---

## Error Responses

All errors follow this format:

```json
{
  "error": "Error message here"
}
```

**HTTP Status Codes:**
- `400` - Bad Request (invalid input)
- `401` - Unauthorized (missing/invalid token)
- `404` - Not Found
- `500` - Internal Server Error

---

## Rate Limiting

API rate limits:
- **Free Plan:** 100 requests/hour
- **Pro Plan:** 1000 requests/hour
- **Enterprise:** Unlimited

Rate limit headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1609459200
```

---

## Webhooks (Coming Soon)

Configure webhooks to receive deployment notifications.

**Events:**
- `deployment.started`
- `deployment.building`
- `deployment.ready`
- `deployment.failed`

---

## SDK Examples

### JavaScript/TypeScript

```typescript
import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  headers: {
    Authorization: `Bearer ${token}`
  }
})

// Trigger deployment
const deploy = await api.post('/deploy', {
  project_id: 'uuid'
})

console.log(deploy.data)
```

### Go

```go
import (
  "bytes"
  "encoding/json"
  "net/http"
)

type DeployRequest struct {
  ProjectID string `json:"project_id"`
}

func triggerDeploy(token, projectID string) error {
  body, _ := json.Marshal(DeployRequest{ProjectID: projectID})
  
  req, _ := http.NewRequest("POST", "http://localhost:8080/api/deploy", bytes.NewBuffer(body))
  req.Header.Set("Authorization", "Bearer "+token)
  req.Header.Set("Content-Type", "application/json")
  
  client := &http.Client{}
  resp, err := client.Do(req)
  return err
}
```

### cURL

```bash
# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Get projects
curl http://localhost:8080/api/projects \
  -H "Authorization: Bearer <token>"

# Trigger deployment
curl -X POST http://localhost:8080/api/deploy \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"project_id":"uuid"}'
```

---

## Postman Collection

Import Postman collection:
```
https://api.dejavu.id/postman.json
```

---

## Support

- Documentation: https://docs.dejavu.id
- GitHub: https://github.com/yourusername/dejavu
- Email: support@dejavu.id

