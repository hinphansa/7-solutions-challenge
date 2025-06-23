# 7-Solutions Backend Challenge – Go implementation

A simple User-management service with **REST + gRPC** transports, **MongoDB** persistence, **JWT** auth and a clean (hexagonal) architecture.

## Project Structure
```
.
├── api/                # proto + generated gRPC code
├── cmd/                # entry-points
│ ├── http/             # :8080 REST server
│ └── grpc/             # :50051 gRPC server
├── internal/           # private application code
│ ├── domain/           # entities (no frameworks)
│ ├── ports/            # interfaces (repo, service)
│ ├── services/         # business logic
│ └── adapters/         # http / grpc / mongo / auth
├── pkg/                # reusable utility (logger)
├── config/             # YAML defaults (Viper)
├── docker/             # Dockerfile.migrate
└── docker-compose.yml  # Mongo + migration + API + HTTP + gRPC + Mongo Express
```

## 🏃‍♂️ Quick start

```bash
git clone https://github.com/hinphansa/7-solutions-challenge.git
cd 7-solutions-challenge
docker compose up --build -d
```

### Dev mode (no Docker)

```bash
go run cmd/migrate/main.go          # one-off schema + index
go run cmd/http/*.go                # REST  :8080
go run cmd/grpc/*.go                # gRPC :9090
```

## Code Generation

### Protobuf Generation

Install protoc and Go plugins:

```bash
# Install protoc
brew install protobuf

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Generate protobuf code:

```bash
protoc \
  --proto_path=./api/proto \
  --go_out=./api/gen/user \
  --go-grpc_out=./api/gen/user \
  api/proto/user.proto
```

### Mockgen: Generate mocks

```bash
go install github.com/golang/mock/mockgen@v1.6.0
```

```bash
mockgen -source=internal/ports/user_port.go -destination=internal/mocks/user_repo_mock.go -package=mocks UserRepository

mockgen -source=internal/services/user_service.go -destination=internal/mocks/user_service_mock.go -package=mocks UserService,PasswordHasher,TokenGenerator

mockgen -source=internal/services/auth_service.go -destination=internal/mocks/auth_service_mock.go -package=mocks AuthService
```

## Testing

```bash
go test ./...
```


# API Documentation

## HTTP API

### User Endpoints (Public)

Since, this is a simple API problem, I'll expose endpoints for listing all users to simplify 7-solutions team's work.
In the real-world, we should not expose endpoints for listing all users or any leakage of personal data.

#### POST `/api/v1/users` - Register a new user

```bash
curl -X POST http://localhost:8080/api/v1/users \
-H "Content-Type: application/json" \
-d '{"email": "test@example.com", "password": "password", "name": "John Doe"}'

# Response: 
# {
#   "id":"6857e9d3699a3ec29bfac36e"
# }
```

#### GET `/api/v1/users/{id}` - Get user by ID


```bash
curl -X GET http://localhost:8080/api/v1/users/$USER_ID

# Response:
# {
#   "id":"6857e9d3699a3ec29bfac36e",
#   "name":"John Doe",
#   "email":"test@example.com",
#   "created_at":"2025-06-22T10:49:12.93Z"
# }
```

#### GET `/api/v1/users?limit=<limit>&offset=<offset>` - List users

If none of query params are provided, it will return all users.

```bash
# Get all users
curl -X GET http://localhost:8080/api/v1/users

# Response:
# {
#   "users":[
#     {
#       "id":"6857e9d3699a3ec29bfac36e",
#       "name":"John Doe",
#       "email":"test@example.com",
#       "created_at":"2025-06-22T10:49:12.93Z"
#     }
#   ]
# }
```

```bash
# Get users with offset 0 and limit 10
curl -X GET http://localhost:8080/api/v1/users?limit=10&offset=0

# Response:
# {
#   "users":[
#     {
#       "id":"6857e9d3699a3ec29bfac36e",
#       "name":"John Doe",
#       "email":"test@example.com",
#       "created_at":"2025-06-22T10:49:12.93Z"
#     }
#   ]
# }
```

### Auth Endpoints (Public)

#### POST `/api/v1/auth/login` - Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password"}'

# Response: 
# {
#   "token":"<JWT_TOKEN>"
# }
```

### User Endpoints (Protected with JWT)

#### GET `/api/v1/users/{id}` - Get user by ID
```bash
curl -X GET http://localhost:8080/api/v1/users/<USER_ID> \
-H "Authorization: Bearer <JWT_TOKEN>"

# Response:
# {
#   "id":"6857e9d3699a3ec29bfac36e",
#   "name":"John Doe",
#   "email":"test@example.com",
#   "created_at":"2025-06-22T10:49:12.93Z"
# }
```

#### PUT `/api/v1/users/{id}` - Update user
```bash
curl -X PUT http://localhost:8080/api/v1/users/<USER_ID> \
-H "Authorization: Bearer <JWT_TOKEN>" \
-H "Content-Type: application/json" \
-d '{"email": "test2@example.com"}'

# Response:
# {
#   "message":"User updated successfully"
# }
```

#### DELETE `/api/v1/users/{id}` - Delete user
```bash
curl -X DELETE http://localhost:8080/api/v1/users/<USER_ID> \
-H "Authorization: Bearer <JWT_TOKEN>"

# Response:
# {
#   "message":"User deleted successfully"
# }
```

## gRPC API

The gRPC server runs on port 50051 and provides the same functionality as the HTTP API. You can use tools like [grpcurl](https://github.com/fullstorydev/grpcurl) or [BloomRPC](https://github.com/bloomrpc/bloomrpc) to interact with the gRPC server.

### Install grpcurl

```bash
brew install grpcurl
```

### List available services

```bash
grpcurl -plaintext localhost:50051 list
```

```bash
grpcurl -plaintext localhost:50051 list user.UserService
```

### User Endpoints (Public)

#### POST `/api/v1/users` - Register a new user

```bash
grpcurl -plaintext -d '{"name": "John Doe", "email": "test@example.com", "password": "password"}' \
  localhost:50051 user.UserService/CreateUser

# Response:
# {
#   "id": "6858e3f87fc09779d13be12e"
# }
```

#### GET `/api/v1/users/{id}` - Get user by ID

```bash
grpcurl -plaintext -d '{"id": "<USER_ID>"}' \
  localhost:50051 user.UserService/GetUserById

# Response:
# {
#   "id": "6858e3f87fc09779d13be12e",
#   "name": "John Doe",
#   "email": "test@example.com",
#   "created_at": "2025-06-22T10:49:12.93Z"
# }
```

#### GET `/api/v1/users?limit=<limit>&offset=<offset>` - List users

```bash
grpcurl -plaintext -d '{"limit": 10, "offset": 0}' \
  localhost:50051 user.UserService/ListUsers

# Response:
# {
#   "users": [
#     {
#       "id": "6858e3f87fc09779d13be12e",
#       "name": "John Doe",
#       "email": "test@example.com",
#       "created_at": "2025-06-22T10:49:12.93Z"
#     }
#   ]
# }
```

```bash
grpcurl -plaintext localhost:50051 user.UserService/ListUsers

# Response:
# {
#   "users": [
#     {
#       "id": "6858e3f87fc09779d13be12e",
#       "name": "John Doe",
#       "email": "test@example.com",
#       "created_at": "2025-06-22T10:49:12.93Z"
#     }
#   ]
# }
```


### Auth Endpoints (Public)

#### POST `/api/v1/auth/login` - Login

```bash
grpcurl -plaintext -d '{"email": "test@example.com", "password": "password"}' \
  localhost:50051 user.UserService/Login

# Response:
# {
#   "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxNzUwNzQyNDAzLCJzdWIiOiI2ODU4ZTNmODdmYzA5Nzc5ZDEzYmUxMmUifQ.e785In7z1Pf3zN5e65h-jGKWfhcId4sXcoLynhm6VVY"
# }
```

### User Endpoints (Protected with JWT)

#### PUT `/api/v1/users/{id}` - Update user

```bash
grpcurl -plaintext -d '{"id": "<USER_ID>", "email": "test2@example.com"}' \
-H "Authorization: Bearer <JWT_TOKEN>" \
localhost:50051 user.UserService/UpdateUser


# Response:
# {
#   "message": "User updated successfully"
# }
```

#### DELETE `/api/v1/users/{id}` - Delete user

```bash
grpcurl -plaintext -d '{"id": "<USER_ID>"}' \
-H "Authorization: Bearer <JWT_TOKEN>" \
localhost:50051 user.UserService/DeleteUser

# Response:
# {
#   "message": "User deleted successfully"
# }
```
