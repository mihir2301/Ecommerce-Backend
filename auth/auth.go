package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JwtWrapper struct {
	SecretKey      string
	Issuer         string
	ExpirationTime int64
}

type JwtClaim struct {
	UserId   primitive.ObjectID `json:"user_id"`
	Email    string             `json:"email"`
	UserType string             `json:"user_type"`
	jwt.StandardClaims
}

// this function will generate a token
func (j *JwtWrapper) GenerateToken(id primitive.ObjectID, email, userType string) (token string, err error) {
	claims := &JwtClaim{
		UserId:   id,
		UserType: userType,
		Email:    email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: &jwt.Time{time.Now().Add(time.Hour * time.Duration(j.ExpirationTime))},
			Issuer:    j.Issuer,
		},
	}
	//fmt.Println("secret key is", j.SecretKey)
	token1 := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = token1.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

// this function will validate a token
func (j *JwtWrapper) ValidateToken(signedToken string) (claims *JwtClaim, err error) {
	//fmt.Println(signedToken)
	//fmt.Println(j.SecretKey)
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = errors.New("could not parse cliams")
		return
	}

	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		err = errors.New("token is expired")
		return
	}
	return
}

/*
 *	This function will check the header contains
 *	Authorization or not and valdiating the token
 */
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, "you are not authorized")
			c.Abort()
			return
		}
		extractedToken := strings.Split(token, "Bearer ")
		if len(extractedToken) == 2 {
			token = strings.TrimSpace(extractedToken[1])
		} else {
			c.JSON(http.StatusUnauthorized, "you are not authorized")
			c.Abort()
			return
		}

		jwtWrapper := JwtWrapper{
			SecretKey: os.Getenv("JwtSecrets"),
			Issuer:    os.Getenv("JwtIssuer"),
		}
		//fmt.Println(jwtWrapper.SecretKey)
		//fmt.Println(token)
		claims, err := jwtWrapper.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserId)
		c.Set("email", claims.Email)
		c.Set("user_type", claims.UserType)
		c.Next()
	}
}
