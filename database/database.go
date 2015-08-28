package database

import (
	"fmt"
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

func AddToDatabases(db Database) {

	dbErr := db.Initialize()

	if dbErr != nil {
		panic(dbErr)
	}

	dbList = append(dbList, db)
	dbMain = db

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
