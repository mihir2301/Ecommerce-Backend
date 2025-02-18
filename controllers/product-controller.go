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

func ListProductsController(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0") //kaha se start karo

	pageInt := helper.ConvertStringToInteger(page)
	limitInt := helper.ConvertStringToInteger(limit)
	offsetInt := helper.ConvertStringToInteger(offset)

	products, count, err := database.Mgr.GetListProducts(pageInt, limitInt, offsetInt, constants.ProductCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": products, "count": count})
}

func SearchProducts(c *gin.Context) {

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")
	search := c.Query("search")
	pageInt := helper.ConvertStringToInteger(page)
	limitInt := helper.ConvertStringToInteger(limit)
	offsetInt := helper.ConvertStringToInteger(offset)

	products, count, err := database.Mgr.SearchProducts(pageInt, limitInt, offsetInt, search, constants.ProductCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": products, "count": count})
}
func UpdateProducts(c *gin.Context) {
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
	var product model.UpdateProduct
	var req model.Products
	err := c.BindJSON(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in binding"})
		return
	}
	productdatabse, err := database.Mgr.GetOneProduct(product.Id, constants.ProductCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in fetching from database"})
		return
	}
	req.ID = productdatabse.ID
	req.Name = productdatabse.Name
	req.Price = productdatabse.Price
	req.ImageUrl = productdatabse.ImageUrl
	req.Description = productdatabse.Description
	req.MetaInfo = productdatabse.MetaInfo
	req.CreatedAt = productdatabse.CreatedAt
	req.UpdatedAt = time.Now().Unix()
	if product.Name != "" {
		req.Name = product.Name
	}
	if product.Description != "" {
		req.Description = product.Description
	}
	if product.Price > 0 {
		req.Price = product.Price
	}
	err = database.Mgr.UpdateProduct(req, constants.ProductCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in updating product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}
func DeleteProduct(c *gin.Context) {
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
	id := c.Query("id")
	product, err := database.Mgr.GetOneProduct(id, constants.ProductCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in geting data from databse"})
		return
	}
	if product.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no product with such name present"})
		return
	}
	err = database.Mgr.DeleteOneProduct(id, constants.ProductCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in deleting databse"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deletion successful"})
}
