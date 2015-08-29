package notify

import (
	"fmt"
	"reflect"
	"regexp"
)

//TODO: Decide Pattern for Notification message , what info is required
//TODO: type of notifications ? and wen to send ? how to take user input values
type Notification struct {
	Message string
}

var (
	errorCount        = 0
	notificationsList []Notify
)

type Notify interface {
	Initialize() error
	SendNotification(message Notification) error
}

type NotificationTypes struct {
	Mailgun MailgunNotify `json:"mailGun"`
	Slack   SlackNotify   `json:"slack"`
	Http    HttpNotify    `json:"http"`
}

func AddNew(notificationTypes NotificationTypes) {

	v := reflect.ValueOf(notificationTypes)

	for i := 0; i < v.NumField(); i++ {
		bytesCount, err := fmt.Print(v.Field(i).Interface().(Notify))
		//Check whether notify object is empty . if its not empty add to the list
		if bytesCount > 3 && err == nil {
			notificationsList = append(notificationsList, v.Field(i).Interface().(Notify))
		}
	}

	for _, value := range notificationsList {

		initErr := value.Initialize()

		if initErr != nil {
			panic(initErr)
		}

	}
}

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}
