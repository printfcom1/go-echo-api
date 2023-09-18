package strc

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ToDo struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type JwtCustomClaims struct {
	Id       string `json:"_id"`
	UserName string `json:"username"`
	Admin    bool   `json:"admin"`
	jwt.RegisteredClaims
}

type AuthInput struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type ToDoInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ToDoDB struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName string             `json:"username"`
	Password string             `json:"password"`
}
