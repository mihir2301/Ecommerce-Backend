package controllers

import (
	"golang_project/auth"
	"golang_project/constants"
	"golang_project/database"
	"golang_project/helper"
	model "golang_project/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Website is working fine",
	})

}
func VerifyEmail(c *gin.Context) {

	var req model.Verification
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in decoding",
		})
	}
	resp := database.Mgr.GetSingleRecordByMail(req.Email, constants.Verification)
	if resp.Otp != 0 {
		sec := resp.CreatedAt + constants.Otpvalidation
		if sec < time.Now().Unix() {
			req, checkmail := helper.SendEmailSendGrid(req)
			if checkmail != nil {
				log.Println(checkmail)
				c.JSON(http.StatusBadGateway, gin.H{
					"message": "error in sending email",
				})
			}
			req.CreatedAt = time.Now().Unix()
			err := database.Mgr.UpdateVerification(req, constants.Verification)
			if err != nil {
				log.Println("error in updating", err)
			}
			c.JSON(200, gin.H{"message": "Mail sent successfully"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"error": false, "message": "Time to verify email is not over"})
		return
	}
	req, checkmail := helper.SendEmailSendGrid(req)
	if checkmail != nil {
		log.Println(checkmail)
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "error in sending email",
		})
	}
	req.CreatedAt = time.Now().Unix()
	_, errs := database.Mgr.Insert(req, constants.Verification)
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failure at insertion",
		})
	}
	c.JSON(200, gin.H{
		"message": "Email verification successful",
	})
}
func VerifyOtp(c *gin.Context) {

	var req model.Verification
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot bind json data",
		})
		return
	}
	if req.Otp <= 0 {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "No otp present",
		})
		return
	}
	resp := database.Mgr.GetSingleRecordByMail(req.Email, constants.Verification)
	if resp == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "NO data found"})
		return
	}
	if resp.Status {
		c.JSON(http.StatusOK, gin.H{"message": "Already verified"})
		return
	}
	if resp.Otp != req.Otp {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Otp"})
		return
	}
	sec := resp.CreatedAt + constants.Otpvalidation
	if sec < time.Now().Unix() {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Otp Expired"})
		return
	}
	req.Status = true
	req.CreatedAt = time.Now().Unix()
	errs := database.Mgr.UpdateEmailVerifiedStatus(req, constants.Verification)
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error at updating"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Otp verified Successfully",
	})
}
func ResendOtp(c *gin.Context) {
	var req model.Verification
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "error at binding"})
	}
	resp := database.Mgr.GetSingleRecordByMail(req.Email, constants.Verification)
	if resp == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "First Verify your Email"})
		return
	}
	req, checkmail := helper.SendEmailSendGrid(req)
	if checkmail != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error at sending Email"})
		return
	}
	req.CreatedAt = time.Now().Unix()
	errs := database.Mgr.UpdateVerification(req, constants.Verification)
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in updating"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfuly resend Otp"})
}
func RegisterUser(c *gin.Context) {
	var userclient model.UsersClient
	var dbclient model.Users
	err := c.BindJSON(&userclient)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in binding"})
		return
	}
	err = helper.CheckUservalidation(userclient)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	resp := database.Mgr.GetSingleRecordByMail(userclient.Email, constants.Verification)
	if !resp.Status {
		c.JSON(http.StatusBadRequest, gin.H{"message": "First verify your Email"})
		return
	}
	res := database.Mgr.GetSingleRecordForUser(userclient.Email, constants.UsersCollection)
	if res.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{"messsage": "user already registered"})
		return
	}
	dbclient.Email = userclient.Email
	dbclient.Name = userclient.Name
	dbclient.Phone = userclient.Phone
	encryptedpass := helper.GenPassHash(userclient.Password)
	dbclient.Password = encryptedpass
	dbclient.UserType = constants.NormalUser
	dbclient.CreatedAt = time.Now().Unix()
	dbclient.UpdatedAt = time.Now().Unix()
	id, errs := database.Mgr.Insert(dbclient, constants.UsersCollection)
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error at inserting"})
		return
	}

	//jwt struct prepare
	jwtWrapper := auth.JwtWrapper{
		SecretKey:      os.Getenv("JwtSecrets"),
		Issuer:         os.Getenv("JwtIssuer"),
		ExpirationTime: 48,
	}
	userId := id.(primitive.ObjectID)
	//gen token
	token, err := jwtWrapper.GenerateToken(userId, userclient.Email, dbclient.UserType)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	dbclient.Password = ""
	c.JSON(http.StatusOK, gin.H{"message": "user registered successfully", "data": dbclient, "token": token})

}
func UserLogin(c *gin.Context) {
	var userlogin model.Login
	err := c.BindJSON(&userlogin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error at binding at userlogin"})
		return
	}
	user := database.Mgr.GetSingleRecordForUser(userlogin.Email, constants.UsersCollection)
	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no user found at user login"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userlogin.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Password wrong"})
		return
	}
	jwtwrapper := auth.JwtWrapper{
		SecretKey:      os.Getenv("JwtSecrets"),
		Issuer:         os.Getenv("JwtIssuer"),
		ExpirationTime: 48,
	}
	token, err := jwtwrapper.GenerateToken(user.ID, user.Email, user.UserType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error at generating token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user, "token": token})
}

func AddToCart(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in getching email"})
		return
	}
	user := database.Mgr.GetSingleRecordByMail(email.(string), constants.UsersCollection)
	address, err := database.Mgr.GetSingleAddress(user.ID, constants.AddressCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error at fetching address"})
		return
	}
	if address.Address1 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "address required"})
		return
	}
	var cart model.UserCart
	var dbcart model.Cart
	err = c.BindJSON(&cart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "binding problem"})
		return
	}
	dbcart.ProductID, _ = primitive.ObjectIDFromHex(cart.ProductID)
	dbcart.UserId, _ = primitive.ObjectIDFromHex(cart.UserId)
	dbcart.Checkout = false
	_, err = database.Mgr.Insert(dbcart, constants.Cartcollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error at inserting"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "added to cart"})
}
func AddAddressOfUser(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Not authorized"})
		return
	}
	user := database.Mgr.GetSingleRecordByMail(email.(string), constants.UsersCollection)
	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not authorized"})
		return
	}
	var address model.AddressClient
	var addressdb model.Address
	err := c.BindJSON(&address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in binding"})
		return
	}
	id, err := primitive.ObjectIDFromHex(address.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in converting id"})
		return
	}
	addressdb.UserID = id
	addressdb.Address1 = address.Address1
	addressdb.City = address.City
	addressdb.Country = address.Country
	_, err = database.Mgr.Insert(addressdb, constants.AddressCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error at inserting"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mesage": "address successfully added"})
}
func GetSingleUser(c *gin.Context) {

}
func UpdateUser(c *gin.Context) {

}
func CheckOutOrder(c *gin.Context) {

}
