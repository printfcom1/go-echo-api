package service

import (
	"errors"

	"github.com/to-do-list/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type toDolistService struct {
	toDoListRepo repository.CostomerRepository
}

func NewToDolistService(toDoListRepo repository.CostomerRepository) toDolistService {
	return toDolistService{toDoListRepo: toDoListRepo}
}

func (s toDolistService) GetToDoListAllService() ([]ToDoListRes, error) {
	toDoList, err := s.toDoListRepo.GetToDoAll()
	if err != nil {
		return nil, err
	}

	toDoListRes := []ToDoListRes{}
	for _, todo := range toDoList {
		toDoRes := ToDoListRes{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			CreatedAt:   todo.CreatedAt,
		}
		toDoListRes = append(toDoListRes, toDoRes)
	}
	return toDoListRes, nil
}

func (s toDolistService) GetToDoListByIdService(id string) (*ToDoListRes, error) {
	toDo, err := s.toDoListRepo.GetToDoById(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ToDo item with ID " + id + " not found.")
		}
		return nil, err
	}
	toDoRes := ToDoListRes{
		ID:          toDo.ID,
		Title:       toDo.Title,
		Description: toDo.Description,
		CreatedAt:   toDo.CreatedAt,
	}
	return &toDoRes, nil
}

func (s toDolistService) CreateToDoListService(toDoInput repository.ToDoListInput) (*string, error) {
	id, err := s.toDoListRepo.CreateToDo(toDoInput)

	if err != nil {
		return nil, err
	}

	objIDString := id.(primitive.ObjectID).Hex()

	message := "ToDo id " + objIDString + " created successfully"

	return &message, nil
}

func (s toDolistService) UpdateToDoListService(id string, toDoInput repository.ToDoListInput) (*ToDoListRes, error) {
	toDo, err := s.toDoListRepo.UpdateToDo(id, toDoInput)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ToDo item with ID " + id + " not found.")
		}
		return nil, err
	}

	toDoRes := ToDoListRes{
		ID:          toDo.ID,
		Title:       toDo.Title,
		Description: toDo.Description,
	}

	return &toDoRes, nil
}

func (s toDolistService) DeleteToDoListService(id string) (*ToDoListRes, error) {
	toDo, err := s.toDoListRepo.DeleteToDo(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ToDo item with ID " + id + " not found.")
		}
		return nil, err
	}
	toDoRes := ToDoListRes{
		ID:          toDo.ID,
		Title:       toDo.Title,
		Description: toDo.Description,
		CreatedAt:   toDo.CreatedAt,
	}
	return &toDoRes, nil
}
