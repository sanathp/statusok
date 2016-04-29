package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

type InfluxDb struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	DatabaseName string `json:"databaseName"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

var (
	influxDBcon client.Client
)

const (
	DatabaseName = "InfluxDB"
)

//Return database name
func (influxDb InfluxDb) GetDatabaseName() string {
	return DatabaseName
}

//Intiliaze influx db
func (influxDb InfluxDb) Initialize() error {
	println("InfluxDB : Trying to Connect to database ")

	u, err := url.Parse(fmt.Sprintf("http://%s:%d", influxDb.Host, influxDb.Port))
	if err != nil {
		println("InfluxDB : Invalid Url,Please check domain name given in config file \nError Details: ", err.Error())
		return err
	}

	conf := client.HTTPConfig{
		Addr:     u.String(),
		Username: influxDb.Username,
		Password: influxDb.Password,
	}

	influxDBcon, err = client.NewHTTPClient(conf)

	if err != nil {
		println("InfluxDB : Failed to connect to Database . Please check the details entered in the config file\nError Details: ", err.Error())
		return err
	}

	_, ver, err := influxDBcon.Ping(10 * time.Second)

	if err != nil {
		println("InfluxDB : Failed to connect to Database . Please check the details entered in the config file\nError Details: ", err.Error())
		return err
	}

	createDbErr := createDatabase(influxDb.DatabaseName)

	if createDbErr != nil {
		if createDbErr.Error() != "database already exists" {
			println("InfluxDB : Failed to create Database")
			return createDbErr
		}

	}

	println("InfluxDB: Successfuly connected . Version:", ver)

	return nil
}

//Add request information to database
func (influxDb InfluxDb) AddRequestInfo(requestInfo RequestInfo) error {

	tags := map[string]string{
		"requestId":   strconv.Itoa(requestInfo.Id),
		"requestType": requestInfo.RequestType,
	}
	fields := map[string]interface{}{
		"responseTime": requestInfo.ResponseTime,
		"responseCode": requestInfo.ResponseCode,
	}

	bps, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influxDb.DatabaseName,
		Precision: "ms",
	})

	if err != nil {
		return err
	}

	point, err := client.NewPoint(
		requestInfo.Url,
		tags,
		fields,
		time.Now(),
	)

	if err != nil {
		return err
	}

	bps.AddPoint(point)

	err = influxDBcon.Write(bps)

	if err != nil {
		return err
	}

	return nil
}

//Add Error information to database
func (influxDb InfluxDb) AddErrorInfo(errorInfo ErrorInfo) error {

	tags := map[string]string{
		"requestId":   strconv.Itoa(errorInfo.Id),
		"requestType": errorInfo.RequestType,
		"reason":      errorInfo.Reason.Error(),
	}
	fields := map[string]interface{}{
		"responseBody": errorInfo.ResponseBody,
		"responseCode": errorInfo.ResponseCode,
		"otherInfo":    errorInfo.OtherInfo,
	}

	point, err := client.NewPoint(
		errorInfo.Url,
		tags,
		fields,
		time.Now(),
	)

	if err != nil {
		return err
	}
	bps, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influxDb.DatabaseName,
		Precision: "ms",
	})

	if err != nil {
		return err
	}

	bps.AddPoint(point)

	err = influxDBcon.Write(bps)

	if err != nil {
		return err
	}

	return nil
}

//Returns mean response time of url in given time .Currentlt not used
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

func createDatabase(databaseName string) error {

	_, err := queryDB(fmt.Sprintf("create database %s", databaseName), "")

	return err
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
