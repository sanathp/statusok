package database

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/sanathp/StatusOk/notify"
)

const ()

var (
	errorCount = 0
	dbMain     Database
	dbList     []Database

	ErrResposeCode   = errors.New("Response code do not Match")
	ErrTimeout       = errors.New("Request Time out Error")
	ErrCreateRequest = errors.New("Invalid Request Config.Not able to create request")
	ErrDoRequest     = errors.New("Request failed")

	meanUrlsCount    int
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
	AddRequestInfo(requestInfo RequestInfo) error
	AddErrorInfo(errorInfo ErrorInfo) error
	GetMeanResponseTime(url string, timeSpan int) (float64, error)
}

type DatabaseTypes struct {
	InfluxDb InfluxDb `json:"influxDb"`
}

func Initialize(urls map[string]int64, requestsMean int) {
	meanUrlsCount = requestsMean
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
		dbString = strings.Replace(dbString, " ", "", -1)
		if len(dbString) > 2 {
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
	fmt.Println("Got Reqest info ", requestInfo.Url, " ", requestInfo.ResponseTime)
	//Insert to all databses
	addResponseTimeToUrl(requestInfo.Url, requestInfo.ResponseTime)
	mean, meanErr := getMeanResponseTimeOfUrl(requestInfo.Url)
	if meanErr == nil {
		if mean > urlResponseTimes[requestInfo.Url] {
			notify.SendResponseTimeNotification(notify.ResponseTypeNotification{
				requestInfo.Url,
				requestInfo.RequestType,
				mean})
		}
	}
	for _, db := range dbList {
		go db.AddRequestInfo(requestInfo)
	}
}

func AddErrorInfo(errorInfo ErrorInfo) {
	//Insert to all databses
	fmt.Println("Got Error info ", errorInfo.Url, " ", errorInfo.Reason, " ", errorInfo.OtherInfo)
	for _, db := range dbList {
		go db.AddErrorInfo(errorInfo)
	}
}

func CheckStatus() {

	for url, _ := range responseMean {
		//TODO: which one to lect for notifications
		times, err := dbMain.GetMeanResponseTime(url, 5)
		fmt.Println("error validateAndSendNotification", err)
		if times > 50 {
			//sendNotification()
		}
	}

}

func addResponseTimeToUrl(url string, responseTime int64) {
	queue := responseMean[url]

	if len(queue) == meanUrlsCount {
		queue = queue[1:]
		queue = append(queue, responseTime)
	} else {
		queue = append(queue, responseTime)
	}

	responseMean[url] = queue
}

func getMeanResponseTimeOfUrl(url string) (int64, error) {

	queue := responseMean[url]

	if len(queue) < meanUrlsCount {
		return 0, errors.New("Stil the count has not been reached")

	}

	var sum int64

	for _, val := range queue {
		sum = sum + val
		fmt.Println("cuurent queue ", val)
	}

	return sum / int64(meanUrlsCount), nil
}
