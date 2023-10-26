package repository_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/to-do-list/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestGetToDoAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("success", func(mt *mtest.T) {
		idFrist := primitive.NewObjectID()
		idSecond := primitive.NewObjectID()
		createdAt := time.Date(2023, time.October, 26, 6, 44, 38, 0, time.UTC)
		updatedAt := time.Date(2023, time.October, 26, 6, 44, 38, 0, time.UTC)

		first := mtest.CreateCursorResponse(1, "golang.ToDoList", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: idFrist},
			{Key: "title", Value: "Test Title 1"},
			{Key: "description", Value: "Test Description 1"},
			{Key: "createdAt", Value: createdAt},
			{Key: "updatedAt", Value: updatedAt},
		})

		second := mtest.CreateCursorResponse(1, "golang.ToDoList", mtest.NextBatch, bson.D{
			{Key: "_id", Value: idSecond},
			{Key: "title", Value: "Test Title 2"},
			{Key: "description", Value: "Test Description 2"},
			{Key: "createdAt", Value: createdAt},
			{Key: "updatedAt", Value: updatedAt},
		})

		killCursors := mtest.CreateCursorResponse(0, "golang.ToDoList", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		todoRepo := repository.NewToDoListRepositoryDB(mt.DB)
		todoList, _ := todoRepo.GetToDoAll()
		expectedToDoAll := []repository.ToDoList{
			{ID: idFrist,
				Title:       "Test Title 1",
				Description: "Test Description 1",
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt},
			{ID: idSecond,
				Title:       "Test Title 2",
				Description: "Test Description 2",
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt},
		}
		assert.Equal(t, expectedToDoAll, todoList)
	})

	mt.Run("error DB", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "error DB",
			Name:    "error",
			Labels:  []string{"test"},
		}))
		todoRepo := repository.NewToDoListRepositoryDB(mt.DB)
		_, err := todoRepo.GetToDoAll()
		assert.Error(t, err)
	})

}
func TestGetToDoById(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		id := "650407582e3a536ae6bd03c6"
		objectID, _ := primitive.ObjectIDFromHex(id)
		createdAt := time.Date(2023, time.October, 26, 6, 44, 38, 0, time.UTC)
		updatedAt := time.Date(2023, time.October, 26, 6, 44, 38, 0, time.UTC)
		expectedToDo := repository.ToDoList{
			ID:          objectID,
			Title:       "Test Title",
			Description: "Test Description",
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "golang.ToDoList", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expectedToDo.ID},
			{Key: "title", Value: expectedToDo.Title},
			{Key: "description", Value: expectedToDo.Description},
			{Key: "createdAt", Value: expectedToDo.CreatedAt},
			{Key: "updatedAt", Value: expectedToDo.UpdatedAt},
		}))
		todoRepo := repository.NewToDoListRepositoryDB(mt.DB)
		todo, _ := todoRepo.GetToDoById(id)
		assert.Equal(t, &expectedToDo, todo)
	})

	mt.Run("error ObjectIDFromHex", func(mt *mtest.T) {
		id := "invalidObjectId"
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "golang.ToDoList", mtest.FirstBatch))
		todoRepo := repository.NewToDoListRepositoryDB(mt.DB)
		_, err := todoRepo.GetToDoById(id)
		assert.Error(t, err)
	})

	mt.Run("error DB", func(mt *mtest.T) {
		id := "650407582e3a536ae6bd03c6"
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "error DB",
			Name:    "error",
			Labels:  []string{"test"},
		}))
		todoRepo := repository.NewToDoListRepositoryDB(mt.DB)
		_, err := todoRepo.GetToDoById(id)
		assert.Error(t, err)
	})
}

func TestCreateToDo(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		todoRepo := repository.NewToDoListRepositoryDB(mt.DB)
		todoCreate := repository.ToDoListInput{
			Title:       "Test Title",
			Description: "Test Description",
		}
		todo, _ := todoRepo.CreateToDo(todoCreate)
		expectedToDo := repository.ToDoListResponseCreate{
			ID:          todo.ID,
			Title:       todoCreate.Title,
			Description: todoCreate.Description,
		}
		assert.Equal(t, &expectedToDo, todo)
	})

}
