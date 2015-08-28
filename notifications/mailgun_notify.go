package notifications

import (
	"fmt"
	"github.com/mailgun/mailgun-go"
)

var mailGunClient mailgun.Mailgun

type MailgunNotify struct {
	Email        string `json:"email"`
	ApiKey       string `json:"apiKey"`
	Domain       string `json:"domain"`
	PublicApiKey string `json:"publicApiKey"`
}

func (mailgunNotify *MailgunNotify) Initialize() error {
	mailGunClient = mailgun.NewMailgun(mailgunNotify.Domain, mailgunNotify.ApiKey, mailgunNotify.PublicApiKey)
	//Send Test Email ?
	//Check response

	return nil
}

func (mailgunNotify *MailgunNotify) SendNotification(message Message) error {

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
