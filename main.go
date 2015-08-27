package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"github.com/sanathp/StatusOk/database"
	"github.com/sanathp/StatusOk/request"
	"os"
)

func main() {

	// You library shoulbe be like this
	//The only problem is influx db doesnot have alerts
	//https://github.com/AcalephStorage/consul-alerts/blob/48c6d89c980f498117d394965af5cf9e7abecb7a/consul-alerts.go

	//TODO:
	//Write a timeChecker using this query SELECT mean(responseTime) FROM "https://facebook.com" WHERE time > now() - 5m GROUP BY time(5m);
	//if response time greater than given config send notification

	r := gin.Default()
	configFile, err := os.Open("config.json")

	if err != nil {
		fmt.Println("Error opening config file", err.Error())
	}

	jsonParser := json.NewDecoder(configFile)

	var jsonData request.RequestJson
	if err = jsonParser.Decode(&jsonData); err != nil {

		fmt.Println("Error parsing config file", err.Error())
	}
	database.DatabaseInit()

	request.StartMonitoring(jsonData)

	r.Run(":3143")
}
