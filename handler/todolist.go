package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/to-do-list/repository"
	"github.com/to-do-list/service"
)

type toDoListHandler struct {
	toDoHandler service.ToDoListService
}

func NewToDolistHandler(toDoHandler service.ToDoListService) toDoListHandler {
	return toDoListHandler{toDoHandler: toDoHandler}
}

func (h toDoListHandler) GetToDoList(c echo.Context) error {

	toDoList, err := h.toDoHandler.GetToDoListAllService()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toDoList)
}

func (h toDoListHandler) GetToDoListById(c echo.Context) error {
	id := c.Param("id")
	toDoList, err := h.toDoHandler.GetToDoListByIdService(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toDoList)
}

func (h toDoListHandler) CreateToDoList(c echo.Context) error {
	toDo := &repository.ToDoListInput{}

	if err := c.Bind(toDo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	todolist, err := h.toDoHandler.CreateToDoListService(*toDo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, *todolist)
}

func (h toDoListHandler) UpdateToDoList(c echo.Context) error {
	id := c.Param("id")
	toDo := &repository.ToDoListInput{}

	if err := c.Bind(toDo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	toDoList, err := h.toDoHandler.UpdateToDoListService(id, *toDo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toDoList)
}

func (h toDoListHandler) DeleteToDoToDoList(c echo.Context) error {
	id := c.Param("id")
	toDoList, err := h.toDoHandler.DeleteToDoListService(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toDoList)
}
