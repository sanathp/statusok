package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/influxdb/influxdb/client"
	"log"
	"net/url"
	"time"
)

type InfluxDb struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	DatabaseName string `json:"databaseName"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

var (
	influxDBcon *client.Client
)

func (influxDb InfluxDb) Initialize() error {
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", influxDb.Host, influxDb.Port))
	if err != nil {
		return err
	}

	conf := client.Config{
		URL:      *u,
		Username: influxDb.Username,
		Password: influxDb.Password,
	}

	influxDBcon, err = client.NewClient(conf)
	if err != nil {
		return err
	}

	_, ver, err := influxDBcon.Ping()
	if err != nil {
		return err
	}
	createDbErr := createDatabase(influxDb.DatabaseName)
	if createDbErr != nil {
		return createDbErr
	}

	log.Printf("Successfuly connected to Influx Db! , %s", ver)

	return nil
}

//TODO: use limit insetead of time ?
//https://influxdb.com/docs/v0.8/api/query_language.html
func (influxDb InfluxDb) GetMeanResponseTime(Url string, span int) (float64, error) {

	q := fmt.Sprintf(`select mean(responseTime) from "%s" WHERE time > now() - %dm GROUP BY time(%dm)`, Url, span, span)

	res, err := queryDB(q, influxDb.DatabaseName)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	//Retrive the last record
	noOfRows := len(res[0].Series[0].Values)
	fmt.Println(q)
	if noOfRows != 0 {
		row := res[0].Series[0].Values[noOfRows-1]
		t, err := time.Parse(time.RFC3339, row[0].(string))
		if err != nil || row[1] == nil {

			fmt.Println("error ", err, " ", row[1])
			return 0, err
		}
		val, err2 := row[1].(json.Number).Float64()
		if err2 != nil {

			fmt.Println(err)
			return 0, err2
		}

		fmt.Println("[%2d] %s: %03d\n", 1, t.Format(time.Stamp), val, err2)
		return val, nil
	}
	return 0, errors.New("error")
}

func (influxDb InfluxDb) GetErrorsCount(Url string, span int) (int64, error) {
	//TODO:fix the value insde count for errors
	q := fmt.Sprintf(`select count(responseTime) from "%s" WHERE time > now() - %dm GROUP BY time(%dm)`, Url, span, span)

	res, err := queryDB(q, influxDb.DatabaseName)

	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	count := res[0].Series[0].Values[len(res[0].Series[0].Values)-1][1]

	if count == nil {
		return 0, nil
	}
	value, convErr := count.(json.Number).Int64()

	if convErr != nil {
		return 0, convErr
	}

	return value, nil
}

func createDatabase(databaseName string) error {

	//TODO: test this
	_, err := queryDB(fmt.Sprintf("create database %s", databaseName), "")

	return err
}

func (influxDb InfluxDb) AddRequestInfo(requestInfo RequestInfo) error {

	var pts = make([]client.Point, 0)
	point := client.Point{
		Measurement: requestInfo.Url,
		Tags: map[string]string{
			"requestType": requestInfo.RequestType,
		},
		Fields: map[string]interface{}{
			"responseTime": requestInfo.ResponseTime,
			"responseCode": requestInfo.ResponseCode,
		},
		Time:      time.Now(),
		Precision: "ms",
	}

	pts = append(pts, point)

	bps := client.BatchPoints{
		Points:          pts,
		Database:        influxDb.DatabaseName,
		RetentionPolicy: "default",
	}

	hi, err := influxDBcon.Write(bps)

	if err != nil {
		fmt.Println("Influx db ", err)
		return err
	}
	return nil
}

func (influxDb InfluxDb) AddErrorInfo(errorInfo ErrorInfo) error {

	var pts = make([]client.Point, 0)
	point := client.Point{
		Measurement: errorInfo.Url,
		Tags: map[string]string{
			"requestType": errorInfo.RequestType,
			"reason":      errorInfo.Reason.Error(),
		},
		Fields: map[string]interface{}{
			"responseBody": errorInfo.ResponseBody,
			"responseCode": errorInfo.ResponseCode,
			"otherInfo":    errorInfo.OtherInfo,
		},
		Time:      time.Now(),
		Precision: "ms",
	}

	pts = append(pts, point)

	bps := client.BatchPoints{
		Points:   pts,
		Database: influxDb.DatabaseName,
		//TODO: make this variable?
		RetentionPolicy: "default",
	}

	_, err := influxDBcon.Write(bps)

	if err != nil {
		fmt.Println("Influx db ", err)
	}
	return nil
}

func queryDB(cmd string, databaseName string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: databaseName,
	}
	if response, err := influxDBcon.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return
}
