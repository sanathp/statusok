package database

import (
	"fmt"
	"reflect"
	"time"
)

var (
	errorCount = 0
	dbMain     Database
	dbList     []Database
)

type Message struct {
	MessageText  string
	Url          string
	RequestType  string
	ResponseTime int64
}

type Database interface {
	Initialize() error
	AddToDatabase(message Message) error
	GetMeanResponseTime(url string, timeSpan int) (float64, error)
}

type DatabaseTypes struct {
	InfluxDb InfluxDb `json:"influxDb"`
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
func CheckStatus() {

	/*
		for _, req := range requests.RequestsList {
			//TODO: which one to lect for notifications
			times, err := dbMain.GetMeanResponseTime(req.Url, 5)
			fmt.Println("error validateAndSendNotification", err)
			if times > 50 {
				//sendNotification()
			}
		}
	*/
}

func AddToDatabase(message Message) {

}
