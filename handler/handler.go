package handler

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/to-do-list/struc"
	"golang.org/x/crypto/bcrypt"
)

func QueryToDoByID(toDoList []struc.ToDo, idToQuery string) *struc.ToDo {
	for _, item := range toDoList {
		if item.Id == idToQuery {
			return &item
		}
	}
	return nil
}

func UpdateToDoMap(toDoList *[]struc.ToDo, idToQuery string, newTitle string, newDescription string) *struc.ToDo {
	for i, item := range *toDoList {
		if item.Id == idToQuery {
			(*toDoList)[i].Title = newTitle
			(*toDoList)[i].Description = newDescription
			return &(*toDoList)[i]
		}
	}
	return nil
}

func DeleteToDoByID(toDoList *[]struc.ToDo, idToQuery string) *struc.ToDo {
	for i, item := range *toDoList {
		if item.Id == idToQuery {
			deletedItem := (*toDoList)[i]
			*toDoList = append((*toDoList)[:i], (*toDoList)[i+1:]...)
			return &deletedItem
		}
	}
	return nil
}

func GoDotEnvVariable(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println(map[string]string{"error": err.Error()})
	}

	return os.Getenv(key)
}

func ErrorResponse(c echo.Context, statusCode int, err error) error {
	return c.JSON(statusCode, map[string]string{"error": err.Error()})
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckHashPassword(password string, passwordDB string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordDB), []byte(password))
	if err == nil {
		return true, nil
	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	} else {
		return false, err
	}
}
