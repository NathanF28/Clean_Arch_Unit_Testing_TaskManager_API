package mongo

import (
	"context"
	"fmt"
	"log"
	"task7/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// mongo db implementation of Task interface

type MongoTaskRepository struct { // one type of implementation
	TaskCollection *mongo.Collection
}

// constructor
func NewMongoTaskRepository(taskCol *mongo.Collection) *MongoTaskRepository { // create object for that
	return &MongoTaskRepository{
		TaskCollection: taskCol,
	}
}

func (m *MongoTaskRepository) GetAllTasks() ([]domain.Task, error) {
	var tasks []domain.Task

	filter := bson.D{}

	cursor, err := m.TaskCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var task domain.Task
		err := cursor.Decode(&task)
		if err != nil {
			log.Println("Error decoding task:", err)
			continue
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (m *MongoTaskRepository) GetTaskById(id int) (domain.Task, error) {
	filter := bson.M{"id": id}

	var task domain.Task
	err := m.TaskCollection.FindOne(context.TODO(), filter).Decode(&task)
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

func (m *MongoTaskRepository) CreateTask(newTask *domain.Task) error {
	if newTask.ID == 0 || newTask.Title == "" || newTask.Description == "" || newTask.Status == "" || newTask.DueDate.IsZero() {
		return fmt.Errorf("missing required field(s) in newTask")
	}
	filter := bson.M{"id": newTask.ID}
	var tempTask domain.Task
	err := m.TaskCollection.FindOne(context.TODO(), filter).Decode(&tempTask)
	if err == nil {
		return fmt.Errorf("id already exists")

	}
	_, err = m.TaskCollection.InsertOne(context.TODO(), newTask)
	return err
}

func (m *MongoTaskRepository) UpdateTask(id int, updatedTask *domain.Task) error {
	filter := bson.M{"id": id}

	updateFields := bson.M{}
	if updatedTask.Title != "" {
		updateFields["title"] = updatedTask.Title
	}
	if updatedTask.Description != "" {
		updateFields["description"] = updatedTask.Description
	}
	if !updatedTask.DueDate.IsZero() {
		updateFields["duedate"] = updatedTask.DueDate
	}
	if updatedTask.Status != "" {
		updateFields["status"] = updatedTask.Status
	}

	if len(updateFields) == 0 {
		return nil
	}

	update := bson.M{"$set": updateFields}
	res, err := m.TaskCollection.UpdateOne(context.TODO(), filter, update)
	if res.MatchedCount == 0 || err != nil {
		return fmt.Errorf("no task found with id %d", id)
	}
	return nil
}

func (m *MongoTaskRepository) DeleteTaskById(id int) error {
	filter := bson.M{"id": id}
	_, err := m.TaskCollection.DeleteOne(context.TODO(), filter)
	return err
}
