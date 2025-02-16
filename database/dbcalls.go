package database

import (
	"context"
	"fmt"
	"golang_project/constants"
	model "golang_project/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func (mgr *manager) Insert(data interface{}, collectionName string) (interface{}, error) {
	log.Println(data)
	connect := mgr.connection.Database(constants.Database).Collection(collectionName)

	result, err := connect.InsertOne(context.TODO(), data)
	if err != nil {
		fmt.Printf("Error in inserting in database")
		return nil, err
	}
	log.Println(result.InsertedID)
	return result.InsertedID, nil
}

func (mgr *manager) GetSingleRecordByMail(email string, collection string) *model.Verification {
	resp := &model.Verification{}
	filter := bson.M{"email": email}
	connect := mgr.connection.Database(constants.Database).Collection(collection)

	err := connect.FindOne(context.TODO(), filter).Decode(&resp)
	if err != nil {
		fmt.Println("Error in grtting single record")
		return resp
	}
	return resp
}
func (mgr *manager) UpdateVerification(data model.Verification, collection string) error {
	connect := mgr.connection.Database(constants.Database).Collection(collection)
	filter := bson.M{"email": data.Email}
	update := bson.M{"$set": data}

	_, err := connect.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("Error in updating verification")
		return err
	}
	return nil
}
func (mgr *manager) UpdateEmailVerifiedStatus(data model.Verification, collection string) error {
	connect := mgr.connection.Database(constants.Database).Collection(collection)
	filter := bson.M{"email": data.Email}
	update := bson.M{"$set": data}

	_, err := connect.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("error while updating status")
		return err
	}
	return nil
}
func (mgr *manager) GetSingleRecordForUser(email string, collectionName string) *model.Users {
	resp := &model.Users{}
	filter := bson.M{"email": email}
	connect := mgr.connection.Database(constants.Database).Collection(collectionName)
	err := connect.FindOne(context.TODO(), filter).Decode(&resp)
	if err != nil {
		log.Println("error while getting user")
		return resp
	}
	return resp
}
func (mgr *manager) GetSingleRecordByProductName(name, collection string) *model.Products {
	p := &model.Products{}
	filter := bson.M{"name": name}
	connect := mgr.connection.Database(constants.Database).Collection(constants.ProductCollection)
	err := connect.FindOne(context.TODO(), filter).Decode(&p)
	if err != nil {
		fmt.Println("error at getting product details")
		return p
	}
	return p
}
