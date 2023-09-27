package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepositoryDB struct {
	db *mongo.Database
}

func NewUserRepositoryDB(db *mongo.Database) userRepositoryDB {
	return userRepositoryDB{db: db}
}

func (u userRepositoryDB) GetUser(username string) (*User, error) {
	filter := bson.M{"username": username}
	user := User{}
	err := u.db.Collection("User").FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u userRepositoryDB) CreateUser(user CreateUser) (interface{}, error) {
	userMap := bson.M{
		"username":  user.UserName,
		"password":  user.Password,
		"email":     user.Email,
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	}

	res, err := u.db.Collection("User").InsertOne(context.Background(), userMap)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID
	return id, nil
}
