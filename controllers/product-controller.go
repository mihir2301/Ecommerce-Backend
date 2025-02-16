package controllers

import (
	"golang_project/constants"
	"golang_project/database"
	"golang_project/helper"
	model "golang_project/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Registerproducts(c *gin.Context) {
	useremail, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "cant fetch email from context"})
		return
	}
	user := database.Mgr.GetSingleRecordForUser(useremail.(string), constants.UsersCollection)
	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
		return
	}
	if user.UserType != constants.AdminUser {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not a admin so cannot register product"})
		return
	}
	var product model.ClientProducts
	var p model.Products
	err := c.BindJSON(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error at binding in register products"})
		return
	}
	err = helper.CheckProductValidation(product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	dataproduct := database.Mgr.GetSingleRecordByProductName(p.Name, constants.ProductCollection)
	if dataproduct.Name != "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Product already registerd"})
		return
	}
	p.Name = product.Name
	p.ImageUrl = product.ImageUrl
	p.Description = product.Description
	p.Price = product.Price
	p.MetaInfo = product.MetaInfo
	p.CreatedAt = time.Now().Unix()
	p.UpdatedAt = time.Now().Unix()
	id, err := database.Mgr.Insert(p, constants.ProductCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error at inserting in database"})
		return
	}
	p.ID = id.(primitive.ObjectID)
	c.JSON(http.StatusOK, gin.H{"message": "product registerd succesfuly", "data": p})
}
