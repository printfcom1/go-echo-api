package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/to-do-list/handler"
	"github.com/to-do-list/mongc"
	"github.com/to-do-list/struc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var TodoList []struc.ToDo

var db *mongo.Collection = mongc.InitMongoClient().Database("golang").Collection("ToDoList")
var user *mongo.Collection = mongc.InitMongoClient().Database("golang").Collection("User")

func Login(c echo.Context) error {
	auth := new(struc.AuthInput)

	if err := c.Bind(auth); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if auth.UserName != "jon" || auth.Password != "123456" {
		return echo.ErrUnauthorized
	}

	claims := &struc.JwtCustomClaims{
		UserName: auth.UserName,
		Admin:    true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key := handler.GoDotEnvVariable("SECRET_KEY")
	t, err := token.SignedString([]byte(key))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func LoginDB(c echo.Context) error {
	auth := new(struc.AuthInput)

	if err := c.Bind(auth); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	filter := bson.M{"username": auth.UserName}
	var result struc.User
	err := user.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return echo.ErrUnauthorized
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	checkPass, err := handler.CheckHashPassword(auth.Password, result.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if !checkPass {
		return echo.ErrUnauthorized
	}

	claims := &struc.JwtCustomClaims{
		Id:       result.ID.Hex(),
		UserName: result.UserName,
		Admin:    true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key := handler.GoDotEnvVariable("SECRET_KEY")
	t, err := token.SignedString([]byte(key))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func GetToDolist(c echo.Context) error {
	return c.JSON(http.StatusOK, TodoList)
}

func AddToDoList(c echo.Context) error {
	todo := new(struc.ToDo)
	fmt.Println(c.Get("username"))
	if err := c.Bind(todo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	TodoList = append(TodoList, *todo)

	return c.JSON(http.StatusCreated, todo)
}

func UpdateToDoList(c echo.Context) error {
	id := c.Param("id")
	todo := new(struc.ToDoInput)

	if err := c.Bind(todo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	result := handler.QueryToDoByID(TodoList, id)
	if result != nil {
		result = handler.UpdateToDoMap(&TodoList, id, todo.Title, todo.Description)
		return c.JSON(http.StatusOK, result)
	} else {
		response := map[string]interface{}{"message": "ToDo item with ID " + id + " not found."}
		return c.JSON(http.StatusBadRequest, response)
	}

}

func GetToDolistById(c echo.Context) error {
	id := c.Param("id")
	result := handler.QueryToDoByID(TodoList, id)
	if result != nil {
		return c.JSON(http.StatusOK, result)
	} else {
		response := map[string]interface{}{"message": "ToDo item with ID " + id + " not found."}
		return c.JSON(http.StatusBadRequest, response)
	}
}

func DeleteToDolistById(c echo.Context) error {
	id := c.Param("id")
	result := handler.QueryToDoByID(TodoList, id)
	if result != nil {
		result = handler.DeleteToDoByID(&TodoList, id)
		return c.JSON(http.StatusOK, result)
	} else {
		response := map[string]interface{}{"message": "ToDo item with ID " + id + " not found."}
		return c.JSON(http.StatusBadRequest, response)
	}

}

func GetToDoListDB(c echo.Context) error {
	filter := bson.M{}
	cursor, err := db.Find(context.Background(), filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	defer cursor.Close(context.Background())
	var todoList []struc.ToDoDB
	for cursor.Next(context.Background()) {
		var todo struc.ToDoDB
		if err := cursor.Decode(&todo); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		todoList = append(todoList, todo)
	}

	if err := cursor.Err(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, todoList)
}

func CreateToDoListDB(c echo.Context) error {

	todo := new(struc.ToDoInput)

	if err := c.Bind(todo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	todoMap := bson.M{
		"title":       todo.Title,
		"description": todo.Description,
		"createdAt":   time.Now(),
		"updatedAt":   time.Now(),
	}

	res, err := db.InsertOne(context.Background(), todoMap)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Todo created successfully",
		"id":      res.InsertedID,
	})
}

func GetToDolistDBById(c echo.Context) error {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	filter := bson.M{"_id": objectID}
	var result struc.ToDoDB
	err = db.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			response := map[string]interface{}{"message": "ToDo item with ID " + id + " not found."}
			return c.JSON(http.StatusBadRequest, response)
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func UpdateToDoListDB(c echo.Context) error {
	id := c.Param("id")
	todo := new(struc.ToDoInput)

	if err := c.Bind(todo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return handler.ErrorResponse(c, http.StatusBadRequest, err)
	}
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"title":       todo.Title,
		"description": todo.Description,
		"updatedAt":   time.Now(),
	}}

	var result struc.ToDoDB
	err = db.FindOneAndUpdate(context.Background(), filter, update).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			response := map[string]interface{}{"message": "ToDo item with ID " + id + " not found."}
			return c.JSON(http.StatusBadRequest, response)
		}
		return handler.ErrorResponse(c, http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, result)

}

func DeleteToDolistDBById(c echo.Context) error {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return handler.ErrorResponse(c, http.StatusBadRequest, err)
	}

	filter := bson.M{"_id": objectID}

	var result struc.ToDoDB
	err = db.FindOneAndDelete(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			response := map[string]interface{}{"message": "ToDo item with ID " + id + " not found."}
			return c.JSON(http.StatusBadRequest, response)
		}
		return handler.ErrorResponse(c, http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, result)
}
