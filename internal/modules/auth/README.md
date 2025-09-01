# Authentication Module

A simple authentication system for M3 Storage with basic features like registration, login, and user profile management.

## Structure

The authentication module follows the Clean Architecture with the following layers:

```
internal/modules/auth/
├── domain/          # Domain entities and DTOs
│   ├── user.go      # User entity
│   └── dto.go       # Data transfer objects
├── port/            # Interfaces/contracts
│   └── interfaces.go # Repository and Service interfaces
├── service/         # Business logic
│   ├── auth_service.go      # Authentication service
│   └── user_repository.go   # User repository
├── handler/         # HTTP handlers
│   └── auth_handler.go      # HTTP request handlers
└── dependencies.go  # Dependency injection
```

## Features

### 1. User Management
- User entity with basic fields (email, password, name)
- Extended user profile with additional information
- User status management (active, inactive, suspended, pending)
- Account locking after multiple failed login attempts

### 2. Authentication
- New user registration
- Login with email/password
- JWT access and refresh tokens
- Token refresh endpoint
- Logout

### 3. Profile Management
- View profile information
- Update profile information
- Change password

### 4. Security Features
- Password hashing with bcrypt
- Failed login attempts tracking
- Account locking
- JWT token validation

## API Endpoints

### Public Endpoints (no authentication required)

#### POST /api/v1/auth/register
Register a new user

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "status": "active",
    "email_verified": false,
    "created_at": "2023-12-01T10:00:00Z"
  },
  "message": "User registered successfully"
}
```

#### POST /api/v1/auth/login
Login

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "access_token": "jwt_access_token",
    "refresh_token": "jwt_refresh_token",
    "token_type": "Bearer",
    "expires_in": 900,
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe"
    }
  },
  "message": "Login successful"
}
```

#### POST /api/v1/auth/refresh
Refresh token

**Request Body:**
```json
{
  "refresh_token": "jwt_refresh_token"
}
```

#### POST /api/v1/auth/forgot-password
Forgot password (not fully implemented)

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

### Protected Endpoints (Bearer token required)

#### GET /api/v1/auth/profile
View profile information

**Headers:**
```
Authorization: Bearer {access_token}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe"
    },
    "profile": {
      "user_id": "uuid",
      "avatar": "",
      "phone_number": "+1234567890",
      "timezone": "UTC",
      "language": "en"
    }
  },
  "message": "Profile retrieved successfully"
}
```

#### PUT /api/v1/auth/profile
Update profile

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body:**
```json
{
  "first_name": "Jane",
  "last_name": "Smith",
  "phone_number": "+1234567890",
  "timezone": "Asia/Ho_Chi_Minh",
  "language": "vi"
}
```

#### POST /api/v1/auth/change-password
Change password

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body:**
```json
{
  "current_password": "old_password",
  "new_password": "new_password123"
}
```

#### POST /api/v1/auth/logout
Logout

**Headers:**
```
Authorization: Bearer {access_token}
```

## Database Models

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    email_verified BOOLEAN NOT NULL DEFAULT false,
    last_login_at TIMESTAMP,
    failed_attempts INTEGER NOT NULL DEFAULT 0,
    locked_until TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);
```

### User Profiles Table
```sql
CREATE TABLE user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id),
    avatar TEXT,
    phone_number VARCHAR(20),
    date_of_birth TIMESTAMP,
    timezone VARCHAR(50) DEFAULT 'UTC',
    language VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);
```

## Usage

### 1. Initialize dependencies

```go
// In main.go or an initialization file
authDeps := auth.NewDependencies(db, jwtService, validator)

// Add to router config
routerConfig := &router.RouterConfig{
    AuthHandler: authDeps.AuthHandler,
    // ... other handlers
}
```

### 2. Using the authentication middleware

The authentication middleware is already integrated. To protect endpoints, use:

```go
protectedRoutes.Get("/protected", authMw.RequireAuth(), handler.ProtectedEndpoint)
```

### 3. Get user information from context

In protected handlers:

```go
func (h *Handler) ProtectedEndpoint(c *fiber.Ctx) error {
    userID, err := middleware.GetUserID(c)
    if err != nil {
        return err
    }
    
    claims, err := middleware.GetUserClaims(c)
    if err != nil {
        return err
    }
    
    // Use userID and claims
    return c.JSON(fiber.Map{"user_id": userID})
}
```

## Configuration

The system uses the following configurations:

- **JWT Secret**: Passed in when creating the JWT service
- **Database**: PostgreSQL with GORM
- **Password Hashing**: bcrypt with default cost
- **Token Expiry**: 15 minutes for access token, 7 days for refresh token
- **Account Locking**: 5 failed attempts will lock the account for 30 minutes

## Testing

To test the API endpoints:

1. **Register a new user:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }'
```

2. **Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

3. **View profile (token required):**
```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Future Improvements

1. **Email Verification**: Verify email upon registration
2. **Password Reset**: Complete the forgot password functionality
3. **2FA**: Two-Factor Authentication
4. **Role-Based Access Control**: Role-based permissions
5. **Token Blacklisting**: Blacklist tokens on logout
6. **Rate Limiting**: Limit the number of requests
7. **Audit Logging**: Log authentication activities

## Dependencies

- `gorm.io/gorm`: ORM for the database
- `golang.org/x/crypto/bcrypt`: Password hashing
- `github.com/golang-jwt/jwt/v5`: JWT tokens
- `github.com/gofiber/fiber/v2`: HTTP framework
- `github.com/google/uuid`: UUID generation
- `github.com/go-playground/validator/v10`: Input validation