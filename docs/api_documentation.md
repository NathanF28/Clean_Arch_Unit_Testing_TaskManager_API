
---
# Task Management API Documentation

## Overview
This API provides endpoints for managing tasks and users, with authentication and authorization using JWT. Only authenticated users can access protected routes, and only admins can create, update, delete tasks or promote users.

---

## Authentication
All protected endpoints require a valid JWT token in the `Authorization` header:

```
Authorization: Bearer <your_jwt_token>
```

## Endpoints

### User Management

#### Register User
- **POST /register**
- **Request Body:**
  ```json
  {
    "username": "yourusername",
    "password": "yourpassword"
  }
  ```
- **Response:**
  - `201 Created` on success
  - `400 Bad Request` if invalid or password < 8 chars

#### Login User
- **POST /login**
- **Request Body:**
  ```json
  {
    "username": "yourusername",
    "password": "yourpassword"
  }
  ```
- **Response:**
  ```json
  {
    "token": "<jwt_token>"
  }
  ```
  - Use this token for all protected endpoints.

#### Promote User (Admin Only)
- **PUT /promote**
- **Headers:** `Authorization: Bearer <admin_jwt_token>`
- **Request Body:**
  ```json
  {
    "username": "targetuser"
  }
  ```
- **Response:**
  - `200 OK` on success
  - `403 Forbidden` if not admin

---

### Task Management

#### Get All Tasks
- **GET /tasks**
- **Headers:** `Authorization: Bearer <jwt_token>`
- **Response:** List of tasks

#### Get Task by ID
- **GET /tasks/:id**
- **Headers:** `Authorization: Bearer <jwt_token>`
- **Response:** Task object

#### Create Task (Admin Only)
- **POST /tasks**
- **Headers:** `Authorization: Bearer <admin_jwt_token>`
- **Request Body:**
  ```json
  {
    "title": "Task Title",
    "description": "Task Description",
    "due_date": "2024-08-01T17:00:00Z",
    "status": "pending"
  }
  ```
- **Response:** Created task

#### Update Task (Admin Only)
- **PUT /tasks/:id**
- **Headers:** `Authorization: Bearer <admin_jwt_token>`
- **Request Body:** (same as create)
- **Response:** Updated task

#### Delete Task (Admin Only)
- **DELETE /tasks/:id**
- **Headers:** `Authorization: Bearer <admin_jwt_token>`
- **Response:** Success message

---

## Roles
- **admin:** Can create, update, delete tasks, and promote users.
- **regular:** Can view tasks.

---

## Security
- Passwords are hashed before storage.
- JWT secret is stored in `.env` (not in version control).
- All protected endpoints require JWT authentication.

---

## Setup & Usage
1. **Clone the repo**
2. **Create a `.env` file:**
   ```
   JWT_SECRET=your_super_secret_key
   ```
3. **Run the server:**
   ```
   go run main.go
   ```
4. **Use Postman or similar tool to test endpoints.**

---

## Example Authorization Header
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## Notes
- First registered user becomes admin if no users exist.
- Only admins can promote other users.
- All errors are returned as JSON with an `error` field.

**POST** `/tasks`

**Request:**
- Body: JSON object (all fields required)
```
{
  "id": 1,
  "title": "Task Title",
  "description": "Task Description",
  "duedate": "2025-07-16T00:00:00Z",
  "status": "pending"
}
```

**Response:**
- Status: 201 Created
- Body: Created Task object

---
### 4. Update Task
**PUT** `/tasks/{id}`

**Request:**
- Path parameter: `id` (integer)
- Body: JSON object
```
{
  "title": "Updated Title",
  "description": "Updated Description",
  "duedate": "2025-07-17T00:00:00Z",
  "status": "completed"
}
```

**Response:**
- Status: 200 OK
- Body: Updated Task object
- Status: 400 Bad Request (if ID is invalid)
- Status: 404 Not Found (if task does not exist)

---
### 5. Delete Task
**DELETE** `/tasks/{id}`

**Request:**
- Path parameter: `id` (integer)

**Response:**
- Status: 204 No Content
- Body: `{ "message": "Deleted task" }`
- Status: 400 Bad Request (if ID is invalid)
- Status: 404 Not Found (if task does not exist)

---
## Task Object Format
```
{
  "id": 1,
  "title": "Task Title",
  "description": "Task Description",
  "duedate": "2025-07-16T00:00:00Z",
  "status": "pending"
}
```

- `id`: integer (must be provided and unique)
- `title`: string (required)
- `description`: string (required)
- `duedate`: string (ISO 8601 format, required)
- `status`: string (required, e.g., "pending", "completed")
