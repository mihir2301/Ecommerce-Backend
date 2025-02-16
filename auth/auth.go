package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JwtWrapper struct {
	SecretKey      string
	Issuer         string
	Expirationtime int64
}

type JwtClaim struct {
	UserId   primitive.ObjectID
	Email    string
	UserType string
	jwt.StandardClaims
}

func (j *JwtWrapper) GenerateToken(id primitive.ObjectID, email, usertype string) (string, error) {
	claims := &JwtClaim{
		UserId:   id,
		Email:    email,
		UserType: usertype,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(j.Expirationtime)).Unix(),
			Issuer:    j.Issuer,
		},
	}
	token1 := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := token1.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (j *JwtWrapper) ValidateToken(signedToken string) (*JwtClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err := errors.New("could not parse claim")
		return nil, err
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err := errors.New("token is Expired")
		return nil, err
	}
	return claims, nil
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, "You are not authorized")
			c.Abort()
			return
		}
		extractedtoken := strings.Split(token, "Bearer ")
		if len(extractedtoken) == 2 {
			token = strings.TrimSpace(extractedtoken[1])
		} else {
			c.JSON(http.StatusUnauthorized, "you are not authorized")
			c.Abort()
			return
		}

		jwtWrapper := JwtWrapper{
			SecretKey: os.Getenv("Jwtsecrets"),
			Issuer:    os.Getenv("JwtIssuer"),
		}
		claims, err := jwtWrapper.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "you are not authorized")
			c.Abort()
			return
		}
		c.Set("user_id", claims.Id)
		c.Set("email", claims.Email)
		c.Set("user_type", claims.UserType)
		c.Next()
	}
}
