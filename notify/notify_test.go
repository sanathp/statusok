package notify

import (
	"testing"
)

func TestAddEmptyNotifyObject(t *testing.T) {

	notificationTypes := NotificationTypes{MailNotify{},
		MailgunNotify{},
		SlackNotify{},
		HttpNotify{},
		PagerdutyNotify{}}

	AddNew(notificationTypes)

	if len(notificationsList) != 0 {
		t.Error("Empty Notification Object should not be added to list")
	}
}

func TestAddValidNotifyObject(t *testing.T) {

	notificationTypes := NotificationTypes{MailNotify{},
		MailgunNotify{},
		SlackNotify{},
		HttpNotify{"http://statusOk.com", "GET", nil},
		PagerdutyNotify{}}

	AddNew(notificationTypes)

	if len(notificationsList) != 1 {
		t.Error("Failed to Add Notification Object to list")
	}
}
