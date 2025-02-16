package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cart struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id,omitempty"`
	Checkout  bool               `json:"checkout" bson:"checkout"`
}
