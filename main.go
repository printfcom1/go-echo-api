package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/to-do-list/handler"
	"github.com/to-do-list/repository"
	"github.com/to-do-list/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	initTimeZone()

	dbName := "golang"
	db := InitMongoDB(dbName)

	toDoListRepositoryDB := repository.NewToDoListRepositoryDB(db)
	toDoListService := service.NewToDolistService(toDoListRepositoryDB)
	toDoListHandler := handler.NewToDolistHandler(toDoListService)

	userRepository := repository.NewUserRepositoryDB(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	user := e.Group("/api/user")
	user.Use(middleware.BasicAuth(basicAuth))
	user.POST("/login", userHandler.Login)
	user.POST("/register", userHandler.RegisterUser)

	todoDB := e.Group("/api/todo")

	config := createConfigAuthJWT()
	todoDB.Use(echojwt.WithConfig(config))
	todoDB.GET("/ToDoListAll", toDoListHandler.GetToDoList)
	todoDB.GET("/ToDoListById/:id", toDoListHandler.GetToDoListById)
	todoDB.POST("/ToDoList", toDoListHandler.CreateToDoList)
	todoDB.PUT("/ToDoList/:id", toDoListHandler.UpdateToDoList)
	todoDB.DELETE("/ToDoList/:id", toDoListHandler.DeleteToDoToDoList)

	e.Logger.Fatal(e.Start(":3000"))
}

func initTimeZone() {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	time.Local = location
}

func InitMongoDB(dbName string) *mongo.Database {
	url, err := goDotEnvVariable("MONGODB_URL")
	if err != nil {
		panic(err)
	}
	clientOptions := options.Client().ApplyURI(*url)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	mongoDB := client.Database(dbName)
	fmt.Println("Connected to MongoDB!")
	return mongoDB
}

func createConfigAuthJWT() echojwt.Config {
	key, err := goDotEnvVariable("SECRET_KEY")
	if err != nil {
		panic(err)
	}
	config := echojwt.Config{
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
				return []byte(*key), nil
			})

			if err != nil {
				return nil, err

			} else {
				claims, _ := token.Claims.(jwt.MapClaims)
				c.Set("username", claims["username"].(string))
				return new(service.JwtCustomClaims), nil
			}
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(401, map[string]string{"error": err.Error()})
		},
		SigningKey: []byte(*key)}

	return config
}

func basicAuth(username, password string, c echo.Context) (bool, error) {
	usernameRegis, err := goDotEnvVariable("USERNAME_REGIS")
	if err != nil {
		return false, err
	}
	passwordRegis, err := goDotEnvVariable("PASSWORD_REGIS")
	if err != nil {
		return false, err
	}

	if username == *usernameRegis && password == *passwordRegis {
		return true, nil
	}
	return false, nil
}

func goDotEnvVariable(key string) (*string, error) {

	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}

	value := os.Getenv(key)

	return &value, nil
}
