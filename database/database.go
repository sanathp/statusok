package database

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

const ()

var (
	errorCount = 0
	urlsMap    map[string]int
	dbMain     Database
	dbList     []Database

	ErrResposeCode   = errors.New("Response code do not Match")
	ErrTimeout       = errors.New("Request Time out Error")
	ErrCreateRequest = errors.New("Invalid Request Config.Not able to create request")
	ErrDoRequest     = errors.New("Request failed")
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
	AddRequestInfo(requestInfo RequestInfo) error
	AddErrorInfo(errorInfo ErrorInfo) error
	GetMeanResponseTime(url string, timeSpan int) (float64, error)
}

type DatabaseTypes struct {
	InfluxDb InfluxDb `json:"influxDb"`
}

func Initialize(urls map[string]int) {
	urlsMap = urls
}

func AddNew(databaseTypes DatabaseTypes) {

	v := reflect.ValueOf(databaseTypes)

	for i := 0; i < v.NumField(); i++ {
		bytesCount, err := fmt.Print(v.Field(i).Interface().(Database))
		//Check whether notify object is empty . if its not empty add to the list
		if bytesCount > 3 && err == nil {
			dbList = append(dbList, v.Field(i).Interface().(Database))
		}
	}

	for _, value := range dbList {

		initErr := value.Initialize()

		if initErr != nil {
			panic(initErr)
		}

	}

	//Set first database as primary database for monitoring
	//TODO: mention this in guide
	if len(dbList) != 0 {
		dbMain = dbList[0]
	} else {
		//TODO: how to monitor here
		fmt.Println("No Databse is selected")
	}
}

func CreateNotificationTickers() {

	var ticker *time.Ticker = time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	fmt.Println("CreateNotificationTickers")
	for {
		select {
		case <-ticker.C:
			CheckStatus()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func AddRequestInfo(requestInfo RequestInfo) {
	//Insert to all databses
	for _, db := range dbList {
		go db.AddRequestInfo(requestInfo)
	}
}

func AddErrorInfo(errorInfo ErrorInfo) {
	//Insert to all databses
	for _, db := range dbList {
		go db.AddErrorInfo(errorInfo)
	}
}

func CheckStatus() {

	for url, _ := range urlsMap {
		//TODO: which one to lect for notifications
		times, err := dbMain.GetMeanResponseTime(url, 5)
		fmt.Println("error validateAndSendNotification", err)
		if times > 50 {
			//sendNotification()
		}
	}

}
