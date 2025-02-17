package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Products struct {
	ID          primitive.ObjectID     `json:"_id" bson:"_id,omitempty"`
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	Price       float64                `json:"price" bson:"price"`
	ImageUrl    string                 `json:"image_url" bson:"image_url"`
	MetaInfo    map[string]interface{} `json:"meta_information" bson:"meta_information"`
	CreatedAt   int64                  `json:"created_at" bson:"created_at"`
	UpdatedAt   int64                  `json:"updated_at" bson:"updated_at"`
}

type ClientProducts struct {
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	Price       float64                `json:"price" bson:"price"`
	ImageUrl    string                 `json:"image_url" bson:"image_url"`
	MetaInfo    map[string]interface{} `json:"meta_information" bson:"meta_information"`
}

type UpdateProduct struct {
	Id          string  `json:"_id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Name        string  `json:"name"`
}
