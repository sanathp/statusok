package database

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/sanathp/statusok/notify"
)

var (
	MeanResponseCount = 5 //Default number of response times to calcuate mean response time
	ErrorCount        = 1 //Default number of errors should occur to send notification

	dbList       []Database      //list of databases registered
	responseMean map[int][]int64 //A map of queues to calculate mean response time
	dbMain       Database

	ErrResposeCode   = errors.New("Response code do not Match")
	ErrTimeout       = errors.New("Request Time out Error")
	ErrCreateRequest = errors.New("Invalid Request Config.Not able to create request")
	ErrDoRequest     = errors.New("Request failed")

	isLoggingEnabled = false //default
)

type RequestInfo struct {
	Id                   int
	Url                  string
	RequestType          string
	ResponseCode         int
	ResponseTime         int64
	ExpectedResponseTime int64
}

type ErrorInfo struct {
	Id           int
	Url          string
	RequestType  string
	ResponseCode int
	ResponseBody string
	Reason       error
	OtherInfo    string
}

type Database interface {
	Initialize() error
	GetDatabaseName() string
	AddRequestInfo(requestInfo RequestInfo) error
	AddErrorInfo(errorInfo ErrorInfo) error
}

type DatabaseTypes struct {
	InfluxDb InfluxDb `json:"influxDb"`
}

//Intialize responseMean app and counts
func Initialize(ids map[int]int64, mMeanResponseCount int, mErrorCount int) {

	if mMeanResponseCount != 0 {
		MeanResponseCount = mMeanResponseCount
	}

	if mErrorCount != 0 {
		ErrorCount = mErrorCount
	}
	//TODO: try to make all slices as pointers
	responseMean = make(map[int][]int64)

	for id, _ := range ids {
		queue := make([]int64, 0)
		responseMean[id] = queue
	}

}

//Add database to the database List
func AddNew(databaseTypes DatabaseTypes) {

	v := reflect.ValueOf(databaseTypes)

	for i := 0; i < v.NumField(); i++ {
		dbString := fmt.Sprint(v.Field(i).Interface().(Database))

		//Check whether notify object is empty . if its not empty add to the list
		if !isEmptyObject(dbString) {
			dbList = append(dbList, v.Field(i).Interface().(Database))
		}
	}

	if len(dbList) != 0 {
		println("Intializing Database....")
	}

	//Intialize all databases given by user by calling the initialize method
	for _, value := range dbList {

		initErr := value.Initialize()

		if initErr != nil {
			println("Failed to Intialize Database ")
			os.Exit(3)
		}

	}

	//Set first database as primary database
	if len(dbList) != 0 {
		dbMain = dbList[0]
		addTestErrorAndRequestInfo()
	} else {
		fmt.Println("No Database selected.")
	}
}

//Insert test data to database
func addTestErrorAndRequestInfo() {

	println("Adding Test data to your database ....")

	requestInfo := RequestInfo{0, "http://test.com", "GET", 0, 0, 0}

	errorInfo := ErrorInfo{0, "http://test.com", "GET", 0, "test response", errors.New("test error"), "test other info"}

	for _, db := range dbList {
		reqErr := db.AddRequestInfo(requestInfo)
		if reqErr != nil {
			println(db.GetDatabaseName, ": Failed to insert Request Info to database.Please check whether database is installed properly")
		}

		errErr := db.AddErrorInfo(errorInfo)

		if errErr != nil {
			println(db.GetDatabaseName, ": Failed to insert Error Info to database.Please check whether database is installed properly")
		}

	}
}

//This function is called by requests package when request has been successfully performed
//Request data is inserted to all the registered databases
func AddRequestInfo(requestInfo RequestInfo) {
	logRequestInfo(requestInfo)

	//Insert to all databses
	for _, db := range dbList {
		go db.AddRequestInfo(requestInfo)
	}

	//Response time to queue
	addResponseTimeToRequest(requestInfo.Id, requestInfo.ResponseTime)

	//calculate current mean response time . if its less than expected send notitifcation
	mean, meanErr := getMeanResponseTimeOfUrl(requestInfo.Id)

	if meanErr == nil {
		if mean > requestInfo.ExpectedResponseTime {
			clearQueue(requestInfo.Id)
			//TODO :error retry  exponential?
			notify.SendResponseTimeNotification(notify.ResponseTimeNotification{
				requestInfo.Url,
				requestInfo.RequestType,
				requestInfo.ExpectedResponseTime,
				mean})
		}
	}

}

//This function is called by requests package when a reuquest fails
//Error Information is inserted to all the registered databases
func AddErrorInfo(errorInfo ErrorInfo) {
	logErrorInfo(errorInfo)

	//Request failed send notification
	//TODO :error retry  exponential?
	notify.SendErrorNotification(notify.ErrorNotification{
		errorInfo.Url,
		errorInfo.RequestType,
		errorInfo.ResponseBody,
		errorInfo.Reason.Error(),
		errorInfo.OtherInfo})

	//Add Error information to database
	for _, db := range dbList {
		go db.AddErrorInfo(errorInfo)
	}
}

func addResponseTimeToRequest(id int, responseTime int64) {
	if responseMean != nil {
		queue := responseMean[id]

		if len(queue) == MeanResponseCount {
			queue = queue[1:]
			queue = append(queue, responseTime)
		} else {
			queue = append(queue, responseTime)
		}

		responseMean[id] = queue
	}
}

//Calculate current  mean response time for the given request id
func getMeanResponseTimeOfUrl(id int) (int64, error) {

	queue := responseMean[id]

	if len(queue) < MeanResponseCount {
		return 0, errors.New("Stil the count has not been reached")
	}

	var sum int64

	for _, val := range queue {
		sum = sum + val
	}

	return sum / int64(MeanResponseCount), nil
}

func clearQueue(id int) {
	responseMean[id] = make([]int64, 0)
}

func isEmptyObject(objectString string) bool {

	objectString = strings.Replace(objectString, "0", "", -1)
	objectString = strings.Replace(objectString, "map", "", -1)
	objectString = strings.Replace(objectString, "[]", "", -1)
	objectString = strings.Replace(objectString, " ", "", -1)

	if len(objectString) > 2 {
		return false
	} else {
		return true
	}
}

func EnableLogging(fileName string) {

	isLoggingEnabled = true

	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if len(fileName) == 0 {
		// Output to stderr instead of stdout, could also be a file.
		logrus.SetOutput(os.Stderr)
	} else {
		f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

		if err != nil {
			println("Invalid File Path given for parameter --log")
			os.Exit(3)
		}

		logrus.SetOutput(f)
	}

}

func logErrorInfo(errorInfo ErrorInfo) {

	if isLoggingEnabled {
		logrus.WithFields(logrus.Fields{
			"id":           errorInfo.Id,
			"url":          errorInfo.Url,
			"requestType":  errorInfo.RequestType,
			"responseCode": errorInfo.ResponseCode,
			"responseBody": errorInfo.ResponseBody,
			"reason":       errorInfo.Reason.Error(),
			"otherInfo":    errorInfo.Reason,
		}).Error("Status Ok Error occurred for url " + errorInfo.Url)
	}

}

func logRequestInfo(requestInfo RequestInfo) {

	if isLoggingEnabled {
		logrus.WithFields(logrus.Fields{
			"id":                   requestInfo.Id,
			"url":                  requestInfo.Url,
			"requestType":          requestInfo.RequestType,
			"responseCode":         requestInfo.ResponseCode,
			"responseTime":         requestInfo.ResponseTime,
			"expectedResponseTime": requestInfo.ExpectedResponseTime,
		}).Info("")
	}
}
