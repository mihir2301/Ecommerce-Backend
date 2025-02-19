package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UsersClient struct {
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email" binding:"required"`
	Phone    string `json:"phone" bson:"phone"`
	Password string `json:"password" bsom:"password"`
}
type UserUpdate struct {
	Id       string `json:"_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
type Users struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email" binding:"required"`
	Phone     string             `json:"phone" bson:"phone"`
	Password  string             `json:"password" bsom:"password"`
	UserType  string             `json:"user_type" bson:"user_type"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
	UpdatedAt int64              `json:"updated_at" bson:"updated_at"`
}

type Address struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Address1 string             `json:"address1" bson:"address1"`
	UserID   primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	City     string             `json:"city" bson:"city"`
	Country  string             `json:"country" bson:"country"`
}

type AddressClient struct {
	Address1 string `json:"address1" bson:"address1"`
	UserID   string `json:"user_id" bson:"user_id,omitempty"`
	City     string `json:"city" bson:"city"`
	Country  string `json:"country" bson:"country"`
}

type Verification struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email" binding:"required"`
	Otp       int64              `json:"otp" bson:"otp"`
	Status    bool               `json:"status" bson:"status"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
	UpdatedAt int64              `json:"updated_at bson:updated_at"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
