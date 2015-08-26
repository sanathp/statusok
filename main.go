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
	database.DatabaseInit()

	request.StartMonitoring(jsonData)

	r.Run(":3143")
}
