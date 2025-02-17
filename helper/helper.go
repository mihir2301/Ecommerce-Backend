package helper

import (
	"errors"
	"fmt"
	model "golang_project/models"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func CheckUservalidation(u model.UsersClient) error {
	if u.Email == "" {
		return errors.New("email cant be empty")
	}
	if u.Name == "" {
		return errors.New("name cant be empty")
	}
	if u.Password == "" {
		return errors.New("password cant be empty")
	}
	if u.Phone == "" {
		return errors.New("phone cant be empty")
	}
	return nil
}

func GenPassHash(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func CheckProductValidation(p model.ClientProducts) error {
	if p.Description == "" {
		return errors.New("product desription cannot be empty")
	}
	if p.Name == "" {
		return errors.New("name cannot be empty")
	}
	if p.ImageUrl == "" {
		return errors.New("image cannot be empty")
	}
	if p.Price == 0 {
		return errors.New("price cannot be empty")
	}
	return nil
}

func ConvertStringToInteger(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Error at conversion")
		return 0
	}
	return val
}
