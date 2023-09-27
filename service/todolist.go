package service

import (
	"time"

	"github.com/to-do-list/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ToDoListRes struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	CreatedAt   time.Time          `json:"createdAt"`
}

type ToDoListService interface {
	GetToDoListAllService() ([]ToDoListRes, error)
	GetToDoListByIdService(string) (*ToDoListRes, error)
	CreateToDoListService(repository.ToDoListInput) (*string, error)
	UpdateToDoListService(string, repository.ToDoListInput) (*ToDoListRes, error)
	DeleteToDoListService(string) (*ToDoListRes, error)
}
