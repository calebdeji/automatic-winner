package services

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"

	sendinblue "github.com/sendinblue/APIv3-go-library/lib"
	"gopkg.in/gomail.v2"
)

type EmailType string

const (
	Auth EmailType = "auth"
)

func getEmailCredentials() (address string, port int, host string, password string, api string) {
	EMAIL_ADDRESS := os.Getenv("EMAIL_ADDRESS")
	EMAIL_PORT := os.Getenv("EMAIL_PORT")
	EMAIL_HOST := os.Getenv("EMAIL_HOST")
	EMAIL_PASSWORD := os.Getenv("EMAIL_APP_PASSWORD")
	EMAIL_API := os.Getenv("EMAIL_API")

	PARSED_EMAIL_PORT, err := strconv.Atoi(EMAIL_PORT)

	if err != nil {
		log.Panic(err.Error())
	}

	return EMAIL_ADDRESS, PARSED_EMAIL_PORT, EMAIL_HOST, EMAIL_PASSWORD, EMAIL_API
}

func SendEmail(subject string, recipient string, body string) (err error) {

	address, port, host, password, _ := getEmailCredentials()

	msg := gomail.NewMessage()

	msg.SetHeader("From", address)
	msg.SetHeader("To", recipient)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	dialer := gomail.NewDialer(host, port, address, password)

	if err := dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}

func SendEmailV2(subject string, recipient string, body string) (err error) {

	var ctx context.Context = context.Background()

	address, _, _, _, api := getEmailCredentials()

	cfg := sendinblue.NewConfiguration()
	cfg.AddDefaultHeader("api-key", api)

	payload := sendinblue.SendSmtpEmail{
		To: []sendinblue.SendSmtpEmailTo{
			{
				Email: recipient,
			},
		},
		Sender: &sendinblue.SendSmtpEmailSender{
			Email: address,
			Name:  "Zate",
		},
		Subject:     subject,
		HtmlContent: body,
	}

	sib := sendinblue.NewAPIClient(cfg)

	_, response, err := sib.TransactionalEmailsApi.SendTransacEmail(ctx, payload)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 && response.StatusCode != 201 {
		err = errors.New("Invalid status code")
		return err
	}

	return nil
}
