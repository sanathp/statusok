package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/sanathp/StatusOk/database"
	"github.com/sanathp/StatusOk/notify"
	"github.com/sanathp/StatusOk/requests"
	"io"
	"net/http"
	"os"
	"strconv"
)

type configParser struct {
	NotifyWhen    NotifyWhen               `json:"notifyWhen"`
	Requests      []requests.RequestConfig `json:"requests"`
	Notifications notify.NotificationTypes `json:"notifications"`
	Database      database.DatabaseTypes   `json:"database"`
	Concurrency   int                      `json:"concurrency"`
	Port          int                      `json:"port"`
}

type NotifyWhen struct {
	MeanResponseCount int `json:"meanResponseCount"`
	ErrorCount        int `json:"errorCount"`
}

func main() {

	//TODO: Use logrus for all logs
	//TODO: test all these , write test cases
	//TODO: format and comment code
	//TODO: run it for 1 day with all diffrenet type of requests
	//TODO: validations for Config.Json file
	//TODO: build a website using github pages
	//TODO: create Docker file with complete setup
	//TODO: run a deamon using upstaart . learn how to do it https://github.com/zaf/agitator/tree/master/init
	//TODO: gracefull shutdown when user stops the app

	app := cli.NewApp()

	app.Name = "StatusOk"
	app.Usage = "Monitor your api.Get notifications when its down"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "config.json",
			Usage: "location of config file",
		},
	}
	app.Action = func(c *cli.Context) {
		fileName := ""
		if len(c.Args()) > 0 {
			fileName = c.Args()[0]
		}
		if len(c.String("config")) == 0 {
			println("Config parameter is mandatory ")
		} else {
			if fileExists(c.String("config")) {
				println(" File location is", c.String("config"), fileName)
				startServer(c.String("config"))
			} else {
				println("Config file not present at the given location: ", c.String("config"), "\nPlease give correct file location using --config parameter")
			}
		}
	}

	app.Run(os.Args)

}

func startServer(fileName string) {

	configFile, err := os.Open(fileName)

	if err != nil {
		fmt.Println("Error opening config file", err.Error())
	}

	jsonParser := json.NewDecoder(configFile)

	var config configParser
	if err = jsonParser.Decode(&config); err != nil {
		fmt.Println("Error parsing config file .Please check format of the file \nParse Error:", err.Error())
		os.Exit(3)
	}

	notify.AddNew(config.Notifications)
	notify.SendTestNotification()

	//Initialze urls map for monitoring
	urls := make(map[string]int64)

	for _, value := range config.Requests {
		urls[value.Url] = value.ResponseTime
	}

	database.AddNew(config.Database)
	database.Initialize(urls, config.NotifyWhen.MeanResponseCount, config.NotifyWhen.ErrorCount)

	requests.RequestsInit(config.Requests, config.Concurrency)
	requests.StartMonitoring()

	//Tells whether Status Handler is running or not
	http.HandleFunc("/", statusHandler)

	if config.Port == 0 {
		//Default
		http.ListenAndServe(":7321", nil)
	} else {
		http.ListenAndServe(":"+strconv.Itoa(config.Port), nil)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "StatusOk is running")
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
