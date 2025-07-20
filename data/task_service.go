package data

import (
	"context"
	"fmt"
	"log"
	"task6/models"
	"go.mongodb.org/mongo-driver/bson"
)


 


func GetAllTasks() ([]models.Task, error) {
	var tasks []models.Task

	filter := bson.D{}

	cursor, err := TaskCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var task models.Task
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

func GetTaskById(id int) (models.Task, error) {
	filter := bson.M{"id": id}

	var task models.Task
	err := TaskCollection.FindOne(context.TODO(), filter).Decode(&task)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func CreateTask(newTask *models.Task) error {
	// Check for required fields
	if newTask.ID == 0 || newTask.Title == "" || newTask.Description == "" || newTask.Status == "" || newTask.DueDate.IsZero() {
		fmt.Println("ooooooo",newTask)
		return fmt.Errorf("missing required field(s) in newTask")
	}
	_, err := TaskCollection.InsertOne(context.TODO(), newTask)
	return err
}

func UpdateTask(id int, updatedTask *models.Task) error {
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
		return nil // No fields to update
	}

	update := bson.M{"$set": updateFields}
	res, err := TaskCollection.UpdateOne(context.TODO(), filter, update)
	if res.MatchedCount == 0 || err != nil {
		return fmt.Errorf("no task found with id %d", id)
	}
	return nil
}

func RemoveTasks(id int) error {
	filter := bson.M{"id": id}
	_, err := TaskCollection.DeleteOne(context.TODO(), filter)
	return err
}
