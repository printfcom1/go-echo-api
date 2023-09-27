package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ToDoList struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

type ToDoListInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CostomerRepository interface {
	GetToDoAll() ([]ToDoList, error)
	GetToDoById(string) (*ToDoList, error)
	CreateToDo(ToDoListInput) (interface{}, error)
	UpdateToDo(string, ToDoListInput) (*ToDoList, error)
	DeleteToDo(string) (*ToDoList, error)
}
