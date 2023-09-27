package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/to-do-list/service"
)

type usersHandler struct {
	userHandler service.UserService
}

func NewUserHandler(userHandler service.UserService) usersHandler {
	return usersHandler{userHandler: userHandler}
}

func (h usersHandler) Login(c echo.Context) error {

	auth := &service.AuthInput{}

	if err := c.Bind(auth); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	token, err := h.userHandler.Login(*auth)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": *token,
	})
}

func (h usersHandler) RegisterUser(c echo.Context) error {

	register := &service.RegisterUser{}
	if err := c.Bind(register); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	message, err := h.userHandler.RegisterUser(*register)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message ": *message})
}
