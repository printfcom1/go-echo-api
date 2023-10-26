package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type toDoListRepositoryDB struct {
	db *mongo.Database
}

func NewToDoListRepositoryDB(db *mongo.Database) toDoListRepositoryDB {
	return toDoListRepositoryDB{db: db}
}

func (t toDoListRepositoryDB) GetToDoAll() ([]ToDoList, error) {
	filter := bson.M{}

	cursor, err := t.db.Collection("ToDoList").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	todoList := []ToDoList{}
	for cursor.Next(context.Background()) {
		todo := ToDoList{}
		cursor.Decode(&todo)
		todoList = append(todoList, todo)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return todoList, nil
}

func (t toDoListRepositoryDB) GetToDoById(id string) (*ToDoList, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID}
	toDo := ToDoList{}
	err = t.db.Collection("ToDoList").FindOne(context.Background(), filter).Decode(&toDo)
	if err != nil {
		return nil, err
	}
	return &toDo, nil
}

func (t toDoListRepositoryDB) CreateToDo(toDo ToDoListInput) (*ToDoListResponseCreate, error) {

	todoMap := bson.M{
		"title":       toDo.Title,
		"description": toDo.Description,
		"createdAt":   time.Now(),
		"updatedAt":   time.Now(),
	}

	res, err := t.db.Collection("ToDoList").InsertOne(context.Background(), todoMap)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID
	objIDString := id.(primitive.ObjectID)
	todo := ToDoListResponseCreate{
		ID:          objIDString,
		Title:       toDo.Title,
		Description: toDo.Description,
	}

	return &todo, nil
}

func (t toDoListRepositoryDB) UpdateToDo(id string, toDo ToDoListInput) (*ToDoList, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"title":       toDo.Title,
		"description": toDo.Description,
		"updatedAt":   time.Now(),
	}}

	result := ToDoList{}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = t.db.Collection("ToDoList").FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (t toDoListRepositoryDB) DeleteToDo(id string) (*ToDoList, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID}
	toDo := ToDoList{}
	err = t.db.Collection("ToDoList").FindOneAndDelete(context.Background(), filter).Decode(&toDo)
	if err != nil {
		return nil, err
	}
	return &toDo, nil
}
