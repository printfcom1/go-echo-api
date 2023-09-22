package main

import (
	"time"

	"github.com/golang-jwt/jwt"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/to-do-list/controller"
	"github.com/to-do-list/handler"
	"github.com/to-do-list/struc"
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

	e.POST("/login", controller.Login)
	e.POST("/db/login", controller.LoginDB)

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
				return new(struc.JwtCustomClaims), nil
			}
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(401, map[string]string{"error": err.Error()})
		},
		SigningKey: []byte(key)}

	todo.Use(echojwt.WithConfig(config))
	todo.GET("/GetToDoListAll", controller.GetToDolist)
	todo.GET("/GetToDoListById/:id", controller.GetToDolistById)
	todo.POST("/AddToDoList", controller.AddToDoList)
	todo.PUT("/UpdateToDoList/:id", controller.UpdateToDoList)
	todo.DELETE("/DeleteToDoList/:id", controller.DeleteToDolistById)

	todoDB := e.Group("/todo/db")
	todoDB.Use(echojwt.WithConfig(config))
	todoDB.GET("/GetToDoListAll", controller.GetToDoListDB)
	todoDB.POST("/CreateToDoList", controller.CreateToDoListDB)
	todoDB.GET("/GetToDoListById/:id", controller.GetToDolistDBById)
	todoDB.PUT("/UpdateToDoList/:id", controller.UpdateToDoListDB)
	todo.DELETE("/DeleteToDoList/:id", controller.DeleteToDolistDBById)

	e.Logger.Fatal(e.Start(":1323"))
}
