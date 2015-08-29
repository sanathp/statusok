package notify

import (
	"errors"
	"fmt"
	"github.com/mailgun/mailgun-go"
	"strings"
)

var mailGunClient mailgun.Mailgun

type MailgunNotify struct {
	Email        string `json:"email"`
	ApiKey       string `json:"apiKey"`
	Domain       string `json:"domain"`
	PublicApiKey string `json:"publicApiKey"`
}

func (mailgunNotify MailgunNotify) Initialize() error {
	if !validateEmail(mailgunNotify.Email) {
		return errors.New("Mailgun: Invalid Email Address")
	}

	if len(strings.TrimSpace(mailgunNotify.ApiKey)) == 0 {
		return errors.New("Mailgun: Invalid Api Key")
	}

	if len(strings.TrimSpace(mailgunNotify.Domain)) == 0 {
		return errors.New("Mailgun: Invalid Domain name")
	}

	if len(strings.TrimSpace(mailgunNotify.PublicApiKey)) == 0 {
		return errors.New("Mailgun: Invalid PublicApiKey")
	}

	mailGunClient = mailgun.NewMailgun(mailgunNotify.Domain, mailgunNotify.ApiKey, mailgunNotify.PublicApiKey)
	//Send Test Email ?
	//Check response

	return nil
}

func (mailgunNotify MailgunNotify) SendNotification(message Notification) error {

	if mailGunClient == nil {
		mailgunNotify.Initialize()
	}
	fmt.Println("Mailgun notify called", mailgunNotify)
	mail := mailGunClient.NewMessage("StatusOkNotifier <notify@StatusOk.com>",
		"Hi Proud Pichi Pulka",
		"To Proud pichi pulka , \n \n  You Received this ma", mailgunNotify.Email)

	response, id, _ := mailGunClient.Send(mail)
	fmt.Printf("Response ID: %s\n", id)
	fmt.Printf("Message from server: %s\n", response)

	return nil
}
