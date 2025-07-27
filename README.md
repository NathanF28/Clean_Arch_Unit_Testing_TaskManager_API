

# Task Management API (Clean Architecture)

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


---

## Unit Testing & CI

This project uses Go's built-in testing framework and [testify](https://github.com/stretchr/testify) for unit and integration tests. Repository tests use a real MongoDB instance. All tests are run automatically on every push and pull request via GitHub Actions.

- See [docs/unit_testing.md](docs/unit_testing.md) for a full guide on running, writing, and troubleshooting tests.
- The CI workflow is defined in `.github/workflows/go.yml` and spins up a MongoDB service for integration tests.

To run all tests locally:
```sh
go test ./...
```

## Authentication
All protected endpoints require a valid JWT token in the `Authorization` header:

```
Authorization: Bearer <your_jwt_token>
```

---

## Endpoints

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

---

## Contact
For questions, contact the maintainer.
