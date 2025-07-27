# Unit Testing Guide

## Overview
This project uses Go's built-in testing framework and [testify](https://github.com/stretchr/testify) for writing unit and integration tests. MongoDB is used for repository tests, and CI is set up with GitHub Actions.

## How to Run Tests Locally

1. Make sure MongoDB is running locally on `mongodb://localhost:27017`.
2. Run all tests:
   ```sh
   go test ./...
   ```
3. To run a specific test file:
   ```sh
   go test ./repository/mongo/mongo_task_repository_test.go
   ```

## Test Structure
- **Unit tests** cover business logic, controllers, and services.
- **Integration tests** (e.g., repository tests) interact with a real MongoDB instance.
- Test files are named with `_test.go` and use the `suite` package for setup/teardown.

## CI Integration
- Tests are automatically run on every push and pull request via GitHub Actions.
- The workflow spins up a MongoDB service and runs all tests.
- See `.github/workflows/go.yml` for details.

## Writing Tests
- Use `require` for critical assertions and `assert` for non-critical checks.
- Use `SetupSuite`, `SetupTest`, and `TearDownSuite` for test lifecycle management.
- Mock external dependencies for pure unit tests.

## Example Test
```go
func (s *TaskRepositorySuite) TestCreateTask_Success() {
    dueDate, err := time.Parse(time.RFC3339, "2025-07-30T00:00:00Z")
    s.Require().NoError(err)
    task := &domain.Task{ID: 1, Title: "Test", Description: "Desc", DueDate: dueDate, Status: "pending"}
    err = s.taskRepo.CreateTask(task)
    s.Require().NoError(err)
}
```

## Best Practices
- Clean up test data between tests.
- Use unique test databases to avoid conflicts.
- Keep tests fast and isolated.

## Troubleshooting
- If tests fail due to MongoDB, ensure the service is running and accessible.
- Check CI logs for details on failures.

## More Info
- See `docs/api_documentation.md` for API details.
- See each test file for more examples and patterns.
