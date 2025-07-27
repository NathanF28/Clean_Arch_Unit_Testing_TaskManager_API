package mongo_test

import (
	"context"
	"fmt"
	"task7/domain"
	"task7/repository/mongo"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepositorySuite struct {
	suite.Suite
	mongoClient    *mongodriver.Client
	taskCollection *mongodriver.Collection
	taskRepo       *mongo.MongoTaskRepository
	databaseName   string
}

func TestTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(TaskRepositorySuite))
}

func (s *TaskRepositorySuite) SetupSuite() {

	s.databaseName = "task7_test_tasks_db"
	mongoURI := "mongodb://localhost:27017"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongodriver.Connect(ctx, options.Client().ApplyURI(mongoURI))
	s.Require().NoError(err, "Failed to connect to local MongoDB at "+mongoURI)
	s.mongoClient = client

	err = client.Ping(ctx, nil)
	s.Require().NoError(err, "Failed to ping local MongoDB. Is it running?")

	s.taskCollection = client.Database(s.databaseName).Collection("tasks")
	s.taskRepo = mongo.NewMongoTaskRepository(s.taskCollection)

	_, err = s.taskCollection.Indexes().CreateOne(context.Background(), mongodriver.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	s.Require().NoError(err, "Failed to create unique index on task ID")

 
}

func (s *TaskRepositorySuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if s.mongoClient != nil {
		err := s.mongoClient.Database(s.databaseName).Drop(ctx)
		s.NoError(err, "Failed to drop test database")
	}

	if s.mongoClient != nil {
		err := s.mongoClient.Disconnect(ctx)
		s.NoError(err, "Failed to disconnect MongoDB client")
	}
 
}

func (s *TaskRepositorySuite) SetupTest() {
	_, err := s.taskCollection.DeleteMany(context.Background(), bson.D{})
	s.Require().NoError(err, "Failed to clear tasks collection")
}

func (s *TaskRepositorySuite) TestCreateTask_Success() {
	dueDate, err := time.Parse(time.RFC3339, "2025-07-30T00:00:00Z")
	s.Require().NoError(err, "Failed to parse due date string")

	task := &domain.Task{
		ID:          1,
		Title:       "Test Task",
		Description: "Test Description",
		DueDate:     dueDate.Truncate(time.Millisecond),
		Status:      "pending",
	}
	err = s.taskRepo.CreateTask(task)
	s.Require().NoError(err, "Failed to create task")

	var result domain.Task
	err = s.taskCollection.FindOne(context.Background(), bson.M{"id": task.ID}).Decode(&result)
	s.Require().NoError(err, "Task not found after creation")
	s.Equal(task.ID, result.ID)
	s.Equal(task.Title, result.Title)
	s.Equal(task.Description, result.Description)
	s.WithinDuration(task.DueDate, result.DueDate, time.Second)
	s.Equal(task.Status, result.Status)
}

func (s *TaskRepositorySuite) TestGetAllTasks() {
	dueDate1, err := time.Parse(time.RFC3339, "2025-07-30T00:00:00Z")
	s.Require().NoError(err, "Failed to parse due date for Task1")
	dueDate2, err := time.Parse(time.RFC3339, "2025-07-31T00:00:00Z")
	s.Require().NoError(err, "Failed to parse due date for Task2")

	task1 := &domain.Task{ID: 10, Title: "Task1", Description: "Desc1", DueDate: dueDate1.Truncate(time.Millisecond), Status: "pending"}
	task2 := &domain.Task{ID: 11, Title: "Task2", Description: "Desc2", DueDate: dueDate2.Truncate(time.Millisecond), Status: "completed"}

	s.Require().NoError(s.taskRepo.CreateTask(task1), "Failed to insert Task1")
	s.Require().NoError(s.taskRepo.CreateTask(task2), "Failed to insert Task2")

	allTasks, err := s.taskRepo.GetAllTasks()
	s.Require().NoError(err, "Failed to get all tasks")
	s.Len(allTasks, 2)

	s.Contains(allTasks, *task1)
	s.Contains(allTasks, *task2)
}

func (s *TaskRepositorySuite) TestGetTaskById_Success() {
	dueDate, err := time.Parse(time.RFC3339, "2025-07-30T00:00:00Z")
	s.Require().NoError(err, "Failed to parse due date for FindMe")

	task := &domain.Task{
		ID:          101,
		Title:       "FindMe",
		Description: "Find this task",
		DueDate:     dueDate.Truncate(time.Millisecond),
		Status:      "pending",
	}
	s.Require().NoError(s.taskRepo.CreateTask(task), "Failed to insert task for GetTaskById")

	found, err := s.taskRepo.GetTaskById(101)
	s.Require().NoError(err, "Failed to get task by ID")
	s.Equal(task.ID, found.ID)
	s.Equal(task.Title, found.Title)
	s.Equal(task.Description, found.Description)
	s.WithinDuration(task.DueDate, found.DueDate, time.Second)
	s.Equal(task.Status, found.Status)
}

func (s *TaskRepositorySuite) TestGetTaskById_NotFound() {
	_, err := s.taskRepo.GetTaskById(9999)
	s.Error(err, "Expected error for non-existent task ID")
	s.Contains(err.Error(), "mongo: no documents in result")
}

func (s *TaskRepositorySuite) TestCreateTask_MissingFields() {
	task := &domain.Task{ID: 0, Title: "", Description: "", Status: "", DueDate: time.Time{}}
	err := s.taskRepo.CreateTask(task)
	s.Error(err, "Expected error for missing required fields")
	s.Contains(err.Error(), "missing required field(s) in newTask")
}

func (s *TaskRepositorySuite) TestUpdateTask_NoFields() {
	dueDate, err := time.Parse(time.RFC3339, "2025-07-30T00:00:00Z")
	s.Require().NoError(err, "Failed to parse due date for NoUpdate")

	task := &domain.Task{ID: 202, Title: "NoUpdate", Description: "Nothing", DueDate: dueDate.Truncate(time.Millisecond), Status: "pending"}
	s.Require().NoError(s.taskRepo.CreateTask(task), "Failed to insert task for no update")

	emptyUpdate := &domain.Task{}
	err = s.taskRepo.UpdateTask(202, emptyUpdate)
	s.NoError(err, "No error expected when no fields to update")

	fetchedTask, err := s.taskRepo.GetTaskById(202)
	s.Require().NoError(err)
	s.Equal(task.Title, fetchedTask.Title)
	s.Equal(task.Description, fetchedTask.Description)
	s.WithinDuration(task.DueDate, fetchedTask.DueDate, time.Second)
	s.Equal(task.Status, fetchedTask.Status)
}

func (s *TaskRepositorySuite) TestUpdateTask_NotFound() {
	update := &domain.Task{Title: "ShouldNotUpdate"}
	err := s.taskRepo.UpdateTask(9999, update)
	s.Error(err, "Expected error for updating non-existent task")
	s.Contains(err.Error(), "no task found with id 9999")
}

func (s *TaskRepositorySuite) TestDeleteTaskById_NotFound() {
	err := s.taskRepo.DeleteTaskById(9999)
	s.NoError(err, "Delete on non-existent ID should not error")
}

func (s *TaskRepositorySuite) TestUpdateTask_Success() {
	originalDueDate, err := time.Parse(time.RFC3339, "2025-07-30T00:00:00Z")
	s.Require().NoError(err, "Failed to parse original due date")
	taskID := 300
	task := &domain.Task{
		ID:          taskID,
		Title:       "Old Title",
		Description: "Old Description",
		DueDate:     originalDueDate.Truncate(time.Millisecond),
		Status:      "pending",
	}
	s.Require().NoError(s.taskRepo.CreateTask(task), "Failed to insert task for update test")

	updatedDueDate, err := time.Parse(time.RFC3339, "2025-08-01T00:00:00Z")
	s.Require().NoError(err, "Failed to parse updated due date")

	update := &domain.Task{
		Title:       "Updated Title",
		Description: "New Description",
		DueDate:     updatedDueDate.Truncate(time.Millisecond),
		Status:      "completed",
	}
	err = s.taskRepo.UpdateTask(taskID, update)
	s.Require().NoError(err, "Failed to update task")

	var result domain.Task
	err = s.taskCollection.FindOne(context.Background(), bson.M{"id": taskID}).Decode(&result)
	s.Require().NoError(err, "Updated task not found")
	s.Equal(update.Title, result.Title)
	s.Equal(update.Description, result.Description)
	s.WithinDuration(update.DueDate, result.DueDate, time.Second)
	s.Equal(update.Status, result.Status)
}

func (s *TaskRepositorySuite) TestDeleteTaskById_Success() {
	dueDate, err := time.Parse(time.RFC3339, "2025-07-30T00:00:00Z")
	s.Require().NoError(err, "Failed to parse due date for DeleteMe")
	taskID := 400
	task := &domain.Task{
		ID:          taskID,
		Title:       "DeleteMe",
		Description: "ToDelete",
		DueDate:     dueDate.Truncate(time.Millisecond),
		Status:      "pending",
	}
	s.Require().NoError(s.taskRepo.CreateTask(task), "Failed to insert task for deletion")

	err = s.taskRepo.DeleteTaskById(taskID)
	s.Require().NoError(err, "Failed to delete task")

	err = s.taskCollection.FindOne(context.Background(), bson.M{"id": taskID}).Err()
	s.Equal(mongodriver.ErrNoDocuments, err, "Task was not deleted or wrong error returned")
}
