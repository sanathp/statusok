package database

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/sanathp/StatusOk/notify"
)

const ()

var (
	ErrorCount        = 1 //Default Value
	MeanResponseCount = 5 //Default value
	dbMain            Database
	dbList            []Database

	ErrResposeCode   = errors.New("Response code do not Match")
	ErrTimeout       = errors.New("Request Time out Error")
	ErrCreateRequest = errors.New("Invalid Request Config.Not able to create request")
	ErrDoRequest     = errors.New("Request failed")

	urlResponseTimes map[string]int64
	responseMean     map[string][]int64
)

type RequestInfo struct {
	Url          string
	RequestType  string
	ResponseCode int
	ResponseTime int64
}

type ErrorInfo struct {
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

func Initialize(urls map[string]int64, mMeanResponseCount int, mErrorCount int) {

	if mMeanResponseCount != 0 {
		MeanResponseCount = mMeanResponseCount
	}

	if mErrorCount != 0 {
		ErrorCount = mErrorCount
	}
	//TODO: try to make all slices as pointers
	responseMean = make(map[string][]int64)

	for url, _ := range urls {
		queue := make([]int64, 0)
		responseMean[url] = queue
	}
}

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

	for _, value := range dbList {

		initErr := value.Initialize()

		if initErr != nil {
			println("Failed to Intialize Database ")
			os.Exit(3)
		}

	}

	//Set first database as primary database for monitoring
	//TODO: mention this in guide
	if len(dbList) != 0 {
		dbMain = dbList[0]
		addTestErrorAndRequestInfo()
	} else {
		//TODO: how to monitor here
		fmt.Println("No Database selected.")
	}
}

func addTestErrorAndRequestInfo() {

	println("Adding Test data to your database ....")

	requestInfo := RequestInfo{"http://test.com", "GET", 0, 0}

	errorInfo := ErrorInfo{"http://test.com", "GET", 0, "test response", errors.New("test error"), "test other info"}

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

func AddRequestInfo(requestInfo RequestInfo) {
	//fmt.Println("Got Reqest info ", requestInfo.Url, " ", requestInfo.ResponseTime)
	//Insert to all databses
	addResponseTimeToUrl(requestInfo.Url, requestInfo.ResponseTime)
	mean, meanErr := getMeanResponseTimeOfUrl(requestInfo.Url)
	if meanErr == nil {
		if mean > urlResponseTimes[requestInfo.Url] {
			clearQueue(requestInfo.Url)
			//TODO :error retry  exponential?
			notify.SendResponseTimeNotification(notify.ResponseTimeNotification{
				requestInfo.Url,
				requestInfo.RequestType,
				requestInfo.ResponseTime,
				mean})
		}
	}
	for _, db := range dbList {
		go db.AddRequestInfo(requestInfo)
	}
}

func AddErrorInfo(errorInfo ErrorInfo) {
	//TODO :error retry  exponential?
	notify.SendErrorNotification(notify.ErrorNotification{
		errorInfo.Url,
		errorInfo.RequestType,
		errorInfo.ResponseBody,
		errorInfo.Reason.Error(),
		errorInfo.OtherInfo})

	for _, db := range dbList {
		go db.AddErrorInfo(errorInfo)
	}
}

func addResponseTimeToUrl(url string, responseTime int64) {
	queue := responseMean[url]

	if len(queue) == MeanResponseCount {
		queue = queue[1:]
		queue = append(queue, responseTime)
	} else {
		queue = append(queue, responseTime)
	}

	responseMean[url] = queue
}

func getMeanResponseTimeOfUrl(url string) (int64, error) {

	queue := responseMean[url]

	if len(queue) < MeanResponseCount {
		return 0, errors.New("Stil the count has not been reached")

	}

	var sum int64

	for _, val := range queue {
		sum = sum + val
		fmt.Println("cuurent queue ", val)
	}

	return sum / int64(MeanResponseCount), nil
}

func clearQueue(url string) {
	responseMean[url] = make([]int64, 0)
}

//TODO: add to util class
func isEmptyObject(objectString string) bool {
	objectString = strings.Replace(objectString, "map", "", -1)
	objectString = strings.Replace(objectString, "[]", "", -1)
	objectString = strings.Replace(objectString, " ", "", -1)

	if len(objectString) > 2 {
		return false
	} else {
		return true
	}
}
