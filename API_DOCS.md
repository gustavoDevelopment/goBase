# API Documentation

## Base URL
`http://localhost:8080/api/v1`

## Authentication
Currently, the API does not implement authentication. In a production environment, you should secure your endpoints with proper authentication.

## Endpoints

### List Users
- **URL**: `/users`
- **Method**: `GET`
- **Query Parameters**:
  - `page` (optional): Page number (default: 1)
  - `limit` (optional): Number of items per page (default: 10, max: 100)

**Example Request**:
```
GET /api/v1/users?page=1&limit=10
```

**Success Response**:
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "email": "user@example.com",
      "date_created": "2023-01-01T00:00:00Z",
      "date_updated": "2023-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10
  }
}
```

### Get User by ID
- **URL**: `/users/:id`
- **Method**: `GET`

**Example Request**:
```
GET /api/v1/users/507f1f77bcf86cd799439011
```

**Success Response**:
```json
{
  "id": "507f1f77bcf86cd799439011",
  "email": "user@example.com",
  "date_created": "2023-01-01T00:00:00Z",
  "date_updated": "2023-01-01T00:00:00Z"
}
```

**Error Response (404)**:
```json
{
  "error": "user not found"
}
```

### Create User
- **URL**: `/users`
- **Method**: `POST`
- **Request Body**:
  - `email` (string, required): User's email address
  - `password` (string, required): User's password

**Example Request**:
```
POST /api/v1/users
Content-Type: application/json

{
  "email": "newuser@example.com",
  "password": "securepassword123"
}
```

**Success Response (201)**:
```json
{
  "id": "507f1f77bcf86cd799439011",
  "email": "newuser@example.com",
  "date_created": "2023-01-01T00:00:00Z",
  "date_updated": "2023-01-01T00:00:00Z"
}
```

**Error Response (400)**:
```json
{
  "error": "email is required"
}
```

## Running the Application

1. Make sure you have Go installed (v1.16+)
2. Set up your environment variables in `.env` file
3. Run the application:
   ```bash
   go run cmd/server/main.go
   ```
4. The API will be available at `http://localhost:8080`

## Testing

To run the tests:
```bash
go test -v ./...
```

## Environment Variables

- `MONGO_URI`: MongoDB connection string
- `MONGO_DATABASE`: MongoDB database name
- `MONGO_COLLECTION`: MongoDB collection name
- `TIMEOUT`: Request timeout duration (default: 10s)
