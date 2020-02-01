package notify

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
)

//Diffrent types of clients to deliver notifications
type NotificationTypes struct {
	MailNotify MailNotify      `json:"mail"`
	Mailgun    MailgunNotify   `json:"mailGun"`
	Slack      SlackNotify     `json:"slack"`
	Http       HttpNotify      `json:"httpEndPoint"`
	Dingding   DingdingNotify  `json:"dingding"`
	Pagerduty  PagerdutyNotify `json:"pagerduty"`
}

type ResponseTimeNotification struct {
	Url                  string
	RequestType          string
	ExpectedResponsetime int64
	MeanResponseTime     int64
}

type ErrorNotification struct {
	Url          string
	RequestType  string
	ResponseBody string
	Error        string
	OtherInfo    string
}

var (
	errorCount        = 0
	notificationsList []Notify
)

type Notify interface {
	GetClientName() string
	Initialize() error
	SendResponseTimeNotification(notification ResponseTimeNotification) error
	SendErrorNotification(notification ErrorNotification) error
}

//Add notification clients given by user in config file to notificationsList
func AddNew(notificationTypes NotificationTypes) {

	v := reflect.ValueOf(notificationTypes)

	for i := 0; i < v.NumField(); i++ {
		notifyString := fmt.Sprint(v.Field(i).Interface().(Notify))
		//Check whether notify object is empty . if its not empty add to the list
		if !isEmptyObject(notifyString) {
			notificationsList = append(notificationsList, v.Field(i).Interface().(Notify))
		}
	}

	if len(notificationsList) == 0 {
		println("No clients Registered for Notifications")
	} else {
		println("Initializing Notification Clients....")
	}

	for _, value := range notificationsList {
		initErr := value.Initialize()

		if initErr != nil {
			println("Notifications : Failed to Initialize ", value.GetClientName(), ".Please check the details in config file ")
			println("Error Details :", initErr.Error())
		} else {
			println("Notifications :", value.GetClientName(), " Intialized")
		}

	}
}

//Send response time notification to all clients registered
func SendResponseTimeNotification(responseTimeNotification ResponseTimeNotification) {

	for _, value := range notificationsList {
		err := value.SendResponseTimeNotification(responseTimeNotification)

		//TODO: exponential retry if fails ? what to do when error occurs ?
		if err != nil {

		}
	}
}

//Send Error notification to all clients registered
func SendErrorNotification(errorNotification ErrorNotification) {

	for _, value := range notificationsList {
		err := value.SendErrorNotification(errorNotification)

		//TODO: exponential retry if fails ? what to do when error occurs ?
		if err != nil {

		}
	}
}

//Send Test notification to all registered clients .To make sure everything is working
func SendTestNotification() {

	println("Sending Test notifications to the registered clients")

	for _, value := range notificationsList {
		err := value.SendResponseTimeNotification(ResponseTimeNotification{"http://test.com", "GET", 700, 800})

		if err != nil {
			println("Failed to Send Response Time notification to ", value.GetClientName(), " Please check the details entered in the config file")
			println("Error Details :", err.Error())
			os.Exit(3)
		} else {
			println("Sent Test Response Time notification to ", value.GetClientName(), ". Make sure you received it")
		}

		err1 := value.SendErrorNotification(ErrorNotification{"http://test.com", "GET", "This is test notification", "Test notification", "test"})

		if err1 != nil {
			println("Failed to Send Error notification to ", value.GetClientName(), " Please check the details entered in the config file")
			println("Error Details :", err1.Error())
			os.Exit(3)
		} else {
			println("Sent Test Error notification to ", value.GetClientName(), ". Make sure you received it")
		}
	}
}

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}

func isEmptyObject(objectString string) bool {
	objectString = strings.Replace(objectString, "0", "", -1)
	objectString = strings.Replace(objectString, "map", "", -1)
	objectString = strings.Replace(objectString, "[]", "", -1)
	objectString = strings.Replace(objectString, " ", "", -1)
	objectString = strings.Replace(objectString, "{", "", -1)
	objectString = strings.Replace(objectString, "}", "", -1)

	if len(objectString) > 0 {
		return false
	} else {
		return true
	}
}

//A readable message string from responseTimeNotification
func getMessageFromResponseTimeNotification(responseTimeNotification ResponseTimeNotification) string {

	message := fmt.Sprintf("Notification From StatusOk\n\nOne of your apis response time is below than expected."+
		"\n\nPlease find the Details below"+
		"\n\nUrl: %v \nRequestType: %v \nCurrent Average Response Time: %v ms\nExpected Response Time: %v ms\n"+
		"\n\nThanks", responseTimeNotification.Url, responseTimeNotification.RequestType, responseTimeNotification.MeanResponseTime, responseTimeNotification.ExpectedResponsetime)

	return message
}

//A readable message string from errorNotification
func getMessageFromErrorNotification(errorNotification ErrorNotification) string {

	message := fmt.Sprintf("Notification From StatusOk\n\nWe are getting error when we try to send request to one of your apis"+
		"\n\nPlease find the Details below"+
		"\n\nUrl: %v \nRequestType: %v \nError Message: %v \nResponse Body: %v\nOther Info:%v\n"+
		"\n\nThanks", errorNotification.Url, errorNotification.RequestType, errorNotification.Error, errorNotification.ResponseBody, errorNotification.OtherInfo)

	return message
}
