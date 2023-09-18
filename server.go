package main

import (
	"time"

	"github.com/golang-jwt/jwt"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	ctrl "github.com/to-do-list/src"
	strc "github.com/to-do-list/struct"
	handler "github.com/to-do-list/util"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	time.Local = location

	e.POST("/login", ctrl.Login)
	e.POST("/db/login", ctrl.LoginDB)

	todo := e.Group("/todo")

	key := handler.GoDotEnvVariable("SECRET_KEY")
	config := echojwt.Config{
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
				return []byte(key), nil
			})

			if err != nil {
				return nil, err

			} else {
				claims, _ := token.Claims.(jwt.MapClaims)
				c.Set("username", claims["username"].(string))
				return new(strc.JwtCustomClaims), nil
			}
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(401, map[string]string{"error": err.Error()})
		},
		SigningKey: []byte(key)}

	todo.Use(echojwt.WithConfig(config))
	todo.GET("/GetToDoListAll", ctrl.GetToDolist)
	todo.GET("/GetToDoListById/:id", ctrl.GetToDolistById)
	todo.POST("/AddToDoList", ctrl.AddToDoList)
	todo.PUT("/UpdateToDoList/:id", ctrl.UpdateToDoList)
	todo.DELETE("/DeleteToDoList/:id", ctrl.DeleteToDolistById)

	todoDB := e.Group("/todo/db")
	todoDB.Use(echojwt.WithConfig(config))
	todoDB.GET("/GetToDoListAll", ctrl.GetToDoListDB)
	todoDB.POST("/CreateToDoList", ctrl.CreateToDoListDB)
	todoDB.GET("/GetToDoListById/:id", ctrl.GetToDolistDBById)
	todoDB.PUT("/UpdateToDoList/:id", ctrl.UpdateToDoListDB)
	todo.DELETE("/DeleteToDoList/:id", ctrl.DeleteToDolistDBById)

	e.Logger.Fatal(e.Start(":1323"))
}
