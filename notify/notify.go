package notify

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

//TODO: Decide Pattern for Notification message , what info is required
//TODO: type of notifications ? and wen to send ? how to take user input values
//TODO: noftication using smtp request https://github.com/zbindenren/logrus_mail/blob/master/mail.go . testp with your gmail
type Notification struct {
	Message string
}

type ResponseTypeNotification struct {
	Url         string
	RequestType string
	Mean        int64
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
		notifyString := fmt.Sprint(v.Field(i).Interface().(Notify))
		fmt.Println(v.Field(i).Interface().(Notify), " ", notifyString, " ", len(notifyString))
		//Check whether notify object is empty . if its not empty add to the list
		notifyString = strings.Replace(notifyString, " ", "", -1)
		if len(notifyString) > 2 {
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

func SendResponseTimeNotification(responseTypeNotification ResponseTypeNotification) {
	//TODO: implement this with full data
	for _, value := range notificationsList {
		//TODO: exponential retry if fails ?
		err := value.SendNotification(Notification{"Hi this is notification from StatusOk .Your response time is low"})

		fmt.Println("Test Notivication error ", value, " ", err)
	}
}

func SendTestNotification() {

	for _, value := range notificationsList {
		err := value.SendNotification(Notification{"Hi this is a test notfocation from StatusOk .Noitifactions from statusOk are working cheers"})

		fmt.Println("Test Notivication error ", value, " ", err)
	}
}

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}
