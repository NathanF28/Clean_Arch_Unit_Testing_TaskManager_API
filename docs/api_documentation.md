
---

# Task Management API Documentation (Clean Architecture)


## Overview
This API provides endpoints for managing tasks and users, with authentication and authorization using JWT. The project is structured according to Clean Architecture principles, ensuring clear separation of concerns, testability, and maintainability.

**Key Clean Architecture Principles:**
- **Separation of Concerns:** Each layer has a distinct responsibility (domain, use case, delivery, infrastructure).
- **Dependency Inversion:** Core business logic depends on abstractions, not concrete implementations.
- **Decoupling:** Domain models and use cases are independent of frameworks and external dependencies.
- **Testability:** Interfaces and dependency injection allow for easy mocking and unit testing.

**Codebase Structure:**
```
├── domain/         # Core business models (e.g., User, Task)
├── usecases/       # Application logic (services, interactors)
├── repository/     # Interfaces and implementations for data access
│   ├── interfaces/ # Repository interfaces (abstractions)
│   └── mongo/      # MongoDB implementations
├── delivery/       # HTTP controllers and router
├── data/           # Infrastructure (MongoDB client, etc.)
├── docs/           # Documentation
└── main.go         # Application entry point (wires dependencies)
```

**Design Decisions:**
- All business logic is in the usecase layer, not in controllers or repositories.
- Repositories are injected as interfaces, allowing for easy substitution and testing.
- Controllers handle HTTP, validate input, and delegate to usecases.
- Infrastructure (e.g., MongoDB) is abstracted behind interfaces.

---


## Authentication
All protected endpoints require a valid JWT token in the `Authorization` header:

```
Authorization: Bearer <your_jwt_token>
```


## Endpoints


---
### User Management


#### Register User (Public)
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


#### Login User (Public)
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


---
### Task Management


#### Get All Tasks (Protected)
- **GET /tasks**
- **Headers:** `Authorization: Bearer <jwt_token>`
- **Response:** List of tasks


#### Get Task by ID (Protected)
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
    "duedate": "2024-08-01T17:00:00Z",
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
- Passwords are hashed before storage (handled in usecase layer).
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
   go run delivery/main.go
   ```
4. **Use Postman or similar tool to test endpoints.**
5. **Run unit tests:**
   ```
   go test ./...
   ```

---


## Example Authorization Header
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---


## Notes
- First registered user becomes admin if no users exist (handled in usecase layer).
- Only admins can promote other users.
- All errors are returned as JSON with an `error` field.

---

## Clean Architecture Guidelines & Testing

- **Domain Layer:** Contains only business entities and logic, no dependencies on other layers.
- **Usecase Layer:** Implements application-specific business rules, orchestrates domain logic, and is decoupled from frameworks and infrastructure.
- **Repository Layer:** Defines interfaces for data access; concrete implementations (e.g., MongoDB) are injected via dependency inversion.
- **Delivery Layer:** Handles HTTP requests, validates input, and delegates to usecases. No business logic here.
- **Infrastructure Layer:** Contains external dependencies (e.g., database clients, middleware).

**Testing:**
- Unit tests are provided for usecases and repositories using mocks.
- To run all tests:
  ```
  go test ./...
  ```
- Use dependency injection to substitute real implementations with mocks for testing.

**Design Decisions:**
- All dependencies flow inward (inversion of control).
- No direct references to infrastructure in domain/usecase layers.
- Controllers and routers are thin and only coordinate requests.

**For future development:**
- Add new features by defining interfaces in the repository layer and implementing them in infrastructure.
- Keep business logic in usecases, not in controllers or repositories.
- Write tests for all new usecases and repository implementations.

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
