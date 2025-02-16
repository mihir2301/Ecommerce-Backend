package helper

import (
	"fmt"
	"golang_project/constants"
	model "golang_project/models"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmailSendGrid(req model.Verification) (model.Verification, error) {
	apikey := os.Getenv("SENDGRID_API_KEY")
	if apikey == "" {
		return req, fmt.Errorf("error sendgrid")
	}

	//create client of sendGrid

	client := sendgrid.NewSendClient(apikey)

	//setup email message

	from := mail.NewEmail("Sender Name", constants.Sender)
	to := mail.NewEmail("Recipient Name", req.Email)
	subject := "OTP verification mail"
	otp := Randomnum()
	req.Otp = int64(otp)
	htmlContent := "<p>Your Login OTP is <strong>" + strconv.Itoa(otp) + "</strong></p>"
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	//send the email message
	_, err := client.Send(message)
	if err != nil {
		return req, err
	}

	return req, nil
}

func Randomnum() int {
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(1000) + 1000
	return randomInt
}
