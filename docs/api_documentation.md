
# Task Manager API Documentation (MongoDB)


Base URL: `http://localhost:8080/tasks`

**Note:** All data is stored in MongoDB. The `id` field is an integer and must be unique for each task. All fields are required when creating a task.

## Endpoints

---
### 1. Get All Tasks
**GET** `/tasks`

**Request:**
- No body required

**Response:**
- Status: 200 OK
- Body: Array of Task objects

```
[
  {
    "id": 1,
    "title": "Task Title",
    "description": "Task Description",
    "duedate": "2025-07-16T00:00:00Z",
    "status": "pending"
  },
  ...
]
```

---
### 2. Get Task By ID
**GET** `/tasks/{id}`

**Request:**
- Path parameter: `id` (integer)

**Response:**
- Status: 200 OK
- Body: Task object
- Status: 400 Bad Request (if ID is invalid)
- Status: 404 Not Found (if task does not exist)

```
{
  "id": 1,
  "title": "Task Title",
  "description": "Task Description",
  "duedate": "2025-07-16T00:00:00Z",
  "status": "pending"
}
```

---
### 3. Create Task

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
