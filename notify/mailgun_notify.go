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

	mailGunClient = mailgun.NewMailgun(mailgunNotify.Domain, mailgunNotify.ApiKey)

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
