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
	//Send Test Email ?
	//Check response

	return nil
}

func (mailgunNotify MailgunNotify) SendResponseTimeNotification(responseTimeNotification ResponseTimeNotification) error {

	subject := "Response Time Notification from StatusOK"
	message := fmt.Sprintf("Notifiaction From StatusOk\nOne of your apis response time is below than expected."+
		"\nPlease find the Details below"+
		"\nUrl: %v \nRequestType: %v \nCurrent Average Response Time: %v \n Expected Response Time: %v\n"+
		"\nThanks", responseTimeNotification.Url, responseTimeNotification.RequestType, responseTimeNotification.MeanResponseTime, responseTimeNotification.ExpectedResponsetime)

	mail := mailGunClient.NewMessage("StatusOkNotifier <notify@StatusOk.com>", subject, message, fmt.Sprintf("<%s>", mailgunNotify.Email))
	_, _, mailgunErr := mailGunClient.Send(mail)

	if mailgunErr != nil {
		return mailgunErr
	}

	return nil
}

func (mailgunNotify MailgunNotify) SendErrorNotification(errorNotification ErrorNotification) error {
	subject := "Error Time Notification from StatusOK"

	message := fmt.Sprintf("Notifiaction From StatusOk\nWe are getting error when we try to send request to one of your apis"+
		"\nPlease find the Details below"+
		"\nUrl: %v \nRequestType: %v \nError Message: %v \n Response Body: %v\n Other Info:%v\n"+
		"\nThanks", errorNotification.Url, errorNotification.RequestType, errorNotification.Error, errorNotification.ResponseBody, errorNotification.OtherInfo)

	mail := mailGunClient.NewMessage("StatusOkNotifier <notify@StatusOk.com>", subject, message, fmt.Sprintf("<%s>", mailgunNotify.Email))
	_, _, mailgunErr := mailGunClient.Send(mail)

	if mailgunErr != nil {
		return mailgunErr
	}

	return nil
}
