package database

import (
	"context"
	model "golang_project/models"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type manager struct {
	connection *mongo.Client
	//ctx        context.Context
	//cancel     context.CancelFunc
}

var Mgr Manager

type Manager interface {
	Insert(interface{}, string) (interface{}, error)
	GetSingleRecordByMail(string, string) *model.Verification
	UpdateVerification(model.Verification, string) error
	UpdateEmailVerifiedStatus(model.Verification, string) error
	GetSingleRecordForUser(string, string) *model.Users
	GetSingleRecordByProductName(string, string) *model.Products
}

func ConnectDb() {

	uri := os.Getenv("DB")

	/*ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		clientoptions := options.Client().ApplyURI(uri)
		client, err := mongo.Connect(ctx, clientoptions)
		if err != nil {
			log.Printf("Mongo Db Setup failed")
		}

		log.Printf("Successfully connected to database at %s\n", uri)

		Mgr = &manager{connection: client, ctx: ctx, cancel: cancel}
	}*/
	clientoptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientoptions)
	if err != nil {
		log.Fatal("Error connecting database", err)
	}
	Mgr = &manager{connection: client}
}
