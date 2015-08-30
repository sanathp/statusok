package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"github.com/sanathp/StatusOk/database"
	"github.com/sanathp/StatusOk/notify"
	"github.com/sanathp/StatusOk/requests"
	"os"
)

type configParser struct {
	Requests      []requests.RequestConfig `json:"requests"`
	Notifications notify.NotificationTypes `json:"notifications"`
	Database      database.DatabaseTypes   `json:"database"`
}

func main() {

	//TODO: Use logrus for all logs
	//TODO: test all these , write test cases
	//TODO: format and comment code
	//TODO: run it for 1 day with all diffrenet type of requests
	//TODO: add cli tool
	//TODO: validations for Config.Json file
	//TODO: build a website using github pages
	//TODO: create Docker file with complete setup

	r := gin.Default()
	configFile, err := os.Open("test_config.json")

	if err != nil {
		fmt.Println("Error opening config file", err.Error())
	}

	jsonParser := json.NewDecoder(configFile)

	var jsonData configParser
	if err = jsonParser.Decode(&jsonData); err != nil {
		fmt.Println("Error parsing config file", err.Error())
	}

	database.AddNew(jsonData.Database)
	notify.AddNew(jsonData.Notifications)

	notify.SendTestNotification()

	//TODO: send all requests for one time and stop of there are any errors
	requests.RequestsInit(jsonData.Requests)
	requests.StartMonitoring()

	//Initialze urls map for monitoring

	urls := make(map[string]int64)

	for _, value := range jsonData.Requests {
		urls[value.Url] = value.ResponseTime
	}

	database.Initialize(urls, 3)

	r.Run(":3143")
}
