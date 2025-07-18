# Task Manager API

A simple RESTful API for managing tasks, built with Go and Gin.

## Features
- Create, read, update, and delete tasks
- MongoDB storage
- API documentation in `docs/api_documentation.md`

## Getting Started

### Prerequisites
- Go 1.24+

### Installation
1. Clone the repository:
   ```sh
   git clone <your-repo-url>
   cd task_manager_api_2
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```

### Running the API
1. Start the server:
   ```sh
   go run main.go
   ```
2. The API will be available at `http://localhost:8080/tasks`

### API Documentation
See [`docs/api_documentation.md`](docs/api_documentation.md) for details on endpoints, request/response formats, and examples.

## Project Structure
```
├── controllers/      # API endpoint handlers
├── data/             # In-memory data service
├── docs/             # Documentation
├── models/           # Data models
├── router/           # Route setup
├── main.go           # Entry point
├── go.mod            # Go module file
```


## Example Task Object (all fields required)
```
{
  "id": 1,
  "title": "Task Title",
  "description": "Task Description",
  "duedate": "2025-07-16T00:00:00Z",
  "status": "pending"
}
```

## MongoDB Notes
- All data is stored in MongoDB, not in-memory.
- The `id` field is an integer and must be unique for each task.
- All fields (`id`, `title`, `description`, `duedate`, `status`) are required when creating a task. If any are missing, the API will return an error.

## License
MIT
