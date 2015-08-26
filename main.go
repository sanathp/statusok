package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"github.com/sanathp/StatusOk/request"
	"os"
	"time"
)

func main() {

	r := gin.New()
	configFile, err := os.Open("config.json")

	if err != nil {
		fmt.Println("Error opening config file", err.Error())
	}

	jsonParser := json.NewDecoder(configFile)

	var jsonData request.RequestJson
	if err = jsonParser.Decode(&jsonData); err != nil {

		fmt.Println("Error parsing config file", err.Error())
	}

	for _, value := range jsonData.Requests {
		fmt.Println("%v", value.Url)
		go createTicker(value)
	}

	r.Run(":3143")
}

func createTicker(value request.RequestConfig) {
	var ticker *time.Ticker = time.NewTicker(value.Time * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			go request.PerformRequest(value)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
func getRequestConfig() request.RequestConfig {

	return request.RequestConfig{
		"http://google.com",
		"GET",
		nil,
		nil,
		200,
		30,
	}
}
