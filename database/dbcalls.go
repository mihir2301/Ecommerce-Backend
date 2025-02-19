package database

import (
	"context"
	"fmt"
	"golang_project/constants"
	model "golang_project/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	connect := mgr.connection.Database(constants.Database).Collection(collection)
	err := connect.FindOne(context.TODO(), filter).Decode(&p)
	if err != nil {
		fmt.Println("error at getting product details")
		return p
	}
	return p
}
func (mgr *manager) GetListProducts(page, limit, offset int, collectionName string) (product []model.Products, count int64, err error) {
	//pagination
	skip := ((page - 1) * limit)
	if offset > 0 {
		skip = offset
	}
	connection := mgr.connection.Database(constants.Database).Collection(collectionName)
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	// break kar kar ke products detail return karega
	curr, err := connection.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		fmt.Println("Error in getting products")
		return nil, 0, nil
	}
	err = curr.All(context.TODO(), &product)
	if err != nil {
		fmt.Println("Error in getting products")
		return nil, 0, nil
	}
	count, err = connection.CountDocuments(context.TODO(), bson.M{})
	return product, count, err
}

func (mgr *manager) SearchProducts(page, limit, offset int, search, collectionName string) (products []model.Products, count int64, err error) {
	skip := ((page - 1) * limit)
	if offset > 0 {
		skip = offset
	}
	connection := mgr.connection.Database(constants.Database).Collection(collectionName)
	findoptions := options.Find()
	findoptions.SetSkip(int64(skip))
	findoptions.SetLimit(int64(limit))
	searchfilter := bson.M{}
	if len(search) >= 3 {
		searchfilter["$or"] = []bson.M{
			{"name": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
			{"description": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
		}
	}
	cur, err := connection.Find(context.TODO(), searchfilter, findoptions)
	if err != nil {
		fmt.Println("error in seraching products in database")
		return
	}
	err = cur.All(context.TODO(), &products)
	if err != nil {
		fmt.Println("error in binding products")
		return
	}
	count, err = connection.CountDocuments(context.TODO(), searchfilter)
	return products, count, err
}

func (mgr *manager) GetOneProduct(id string, collectionName string) (product model.Products, err error) {
	primitiveID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("canot change id to hex")
		return
	}
	filter := bson.M{"_id": primitiveID}
	connect := mgr.connection.Database(constants.Database).Collection(constants.ProductCollection)
	err = connect.FindOne(context.TODO(), filter).Decode(&product)
	return product, err
}
func (mgr *manager) UpdateProduct(product model.Products, collectionName string) error {
	filter := bson.M{"_id": product.ID}
	update := bson.M{"$set": product}
	connect := mgr.connection.Database(constants.Database).Collection(collectionName)
	_, err := connect.UpdateOne(context.TODO(), filter, update)
	return err
}

func (mgr *manager) DeleteOneProduct(id, collectionName string) error {
	primitiveId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": primitiveId}
	connect := mgr.connection.Database(constants.Database).Collection(collectionName)
	_, err = connect.DeleteOne(context.TODO(), filter)
	return err
}
func (mgr *manager) GetSingleAddress(id primitive.ObjectID, collectionName string) (address model.Address, err error) {
	connect := mgr.connection.Database(constants.Database).Collection(collectionName)
	filter := bson.M{"user_id": id}
	err = connect.FindOne(context.TODO(), filter).Decode(&address)
	return address, err
}
func (mgr *manager) GetOneUserByID(id primitive.ObjectID, collectionName string) (user model.Users) {
	connection := mgr.connection.Database(constants.Database).Collection(collectionName)
	filter := bson.M{"_id": id}
	err := connection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		fmt.Println("error in finding user data")
	}
	return user
}
func (mgr *manager) UpdateUser(user model.Users, collectionName string) error {
	connection := mgr.connection.Database(constants.Database).Collection(collectionName)
	filter := bson.M{"_id": user.ID}
	Update := bson.M{"$set": user}
	_, err := connection.UpdateOne(context.TODO(), filter, Update)
	if err != nil {
		return err
	}
	return nil
}

func (mgr *manager) GetCart(id primitive.ObjectID, collectionName string) (cart model.Cart, err error) {
	connect := mgr.connection.Database(constants.Database).Collection(collectionName)
	filter := bson.M{"_id": id}
	err = connect.FindOne(context.TODO(), filter).Decode(&cart)
	return cart, err
}

func (mgr *manager) UpdateCart(cart model.Cart, collectionName string) error {
	connect := mgr.connection.Database(constants.Database).Collection(collectionName)
	filter := bson.M{"_id": cart.ID}
	update := bson.M{"$set": cart}
	_, err := connect.UpdateOne(context.TODO(), filter, update)
	return err
}
