package domain

import "time"

type Task struct {
	ID          int       `bson:"id" json:"id"`
	Title       string    `bson:"title" json:"title"`
	Description string    `bson:"description" json:"description"`
	DueDate     time.Time `bson:"duedate" json:"duedate"`
	Status      string    `bson:"status" json:"status"`
}
