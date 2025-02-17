package main

import (
	"golang_project/constants"
	"golang_project/database"
	"golang_project/helper"
	model "golang_project/models"
	"golang_project/routes"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error while loading .env file")
	}
	database.ConnectDb()
	//crating system admin
	hashpassword := helper.GenPassHash("1234")
	user := model.Users{
		Name:     "Admin",
		Email:    "admin@gmail.com",
		Password: hashpassword,
		UserType: constants.AdminUser,
	}
	u := database.Mgr.GetSingleRecordForUser(user.Email, constants.UsersCollection)
	if u.Email == "" {
		var data interface{} = user
		_, err := database.Mgr.Insert(data, constants.UsersCollection)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {

	routes.ClientRoutes()
}
