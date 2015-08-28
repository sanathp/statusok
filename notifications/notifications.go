package notifications

import ()

type Message struct {
	message string
}

var (
	errorCount        = 0
	notificationsList []Notifications
)

type Notifications interface {
	SendNotification(message Message) error
}

type NotificationTypes struct {
	Mailgun MailgunNotify `json:"mailGun"`
	Slack   SlackNotify   `json:"slack"`
	Http    HttpNotify    `json:"http"`
}

func AddToNotification(notification Notifications) {
	if notificationsList == nil {
		//otificationsList = make(Notifications, 0)
	}

	notificationsList = append(notificationsList, notification)

	for _, value := range notificationsList {
		value.SendNotification(Message{"hi"})
	}
}
