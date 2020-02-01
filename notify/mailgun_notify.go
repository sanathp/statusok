package notify

import (
	"errors"
	"fmt"
	"strings"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

var mailGunClient mailgun.Mailgun

type MailgunNotify struct {
	Email        string `json:"email"`
	ApiKey       string `json:"apiKey"`
	Domain       string `json:"domain"`
	PublicApiKey string `json:"publicApiKey"`
}

func (mailgunNotify MailgunNotify) GetClientName() string {
	return "Mailgun"
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

	return nil
}

func (mailgunNotify MailgunNotify) SendResponseTimeNotification(responseTimeNotification ResponseTimeNotification) error {

	subject := "Response Time Notification from StatusOK"
	message := getMessageFromResponseTimeNotification(responseTimeNotification)

	mail := mailGunClient.NewMessage("StatusOkNotifier <notify@StatusOk.com>", subject, message, fmt.Sprintf("<%s>", mailgunNotify.Email))
	_, _, mailgunErr := mailGunClient.Send(mail)

	if mailgunErr != nil {
		return mailgunErr
	}

	return nil
}

func (mailgunNotify MailgunNotify) SendErrorNotification(errorNotification ErrorNotification) error {
	subject := "Error Time Notification from StatusOK"

	message := getMessageFromErrorNotification(errorNotification)

	mail := mailGunClient.NewMessage("StatusOkNotifier <notify@StatusOk.com>", subject, message, fmt.Sprintf("<%s>", mailgunNotify.Email))
	_, _, mailgunErr := mailGunClient.Send(mail)

	if mailgunErr != nil {
		return mailgunErr
	}

	return nil
}
