package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"github.com/sanathp/StatusOk/database"
	"github.com/sanathp/StatusOk/notifications"
	"github.com/sanathp/StatusOk/requests"
	"os"
)

type configParser struct {
	Requests      []requests.RequestConfig        `json:"requests"`
	Notifications notifications.NotificationTypes `json:"notifications"`
	Database      database.DatabaseTypes          `json:"database"`
}

func main() {

	//TODO: Use logrus for all logs
	//TODO: create an notify through http request use the same one used for requests
	//TODO: test all these , write test cases
	//TODO: format and comment code
	//TODO: run it for 1 day with all diffrenet type of requests
	//TODO: add cli tool
	//TODO: validations for Config.Json file

	r := gin.Default()
	configFile, err := os.Open("config.json")

	if err != nil {
		fmt.Println("Error opening config file", err.Error())
	}

	jsonParser := json.NewDecoder(configFile)

	var jsonData configParser
	if err = jsonParser.Decode(&jsonData); err != nil {

		fmt.Println("Error parsing config file", err.Error())
	}

	fmt.Println("notifications ", jsonData.Notifications.Slack)
	//TODO:
	//database.Initialize(jsonData.Requests)
	database.AddToDatabases(&jsonData.Database.InfluxDb)

	requests.RequestsInit(jsonData.Requests)
	requests.StartMonitoring()

	notifications.AddToNotification(&jsonData.Notifications.Slack)
	notifications.AddToNotification(&jsonData.Notifications.Mailgun)
	notifications.AddToNotification(&jsonData.Notifications.Http)

	r.Run(":3143")
}
