package database

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
	"log"
	"net/url"
	"os"
	_ "strconv"
	"time"
)

const (
	MyHost        = "localhost"
	MyPort        = 8086
	MyDB          = "statusok"
	MyMeasurement = "shapes"
)

var (
	influxDBcon *client.Client
)

func DatabaseInit() {
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", MyHost, MyPort))
	if err != nil {
		log.Fatal(err)
	}

	conf := client.Config{
		URL:      *u,
		Username: os.Getenv("INFLUX_USER"),
		Password: os.Getenv("INFLUX_PWD"),
	}

	influxDBcon, err = client.NewClient(conf)
	if err != nil {
		log.Fatal(err)
	}

	dur, ver, err := influxDBcon.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Happy as a Hippo! %v, %s", dur, ver)
}

func WritePoints(url string, responseTime int64, requestType string) {

	var pts = make([]client.Point, 0)
	point := client.Point{
		Measurement: url,
		Tags: map[string]string{
			"requestType":  requestType,
			"responseTime": "responseTime",
		},
		Fields: map[string]interface{}{
			"responseTime": responseTime,
		},
		Time:      time.Now(),
		Precision: "ms",
	}

	pts = append(pts, point)

	bps := client.BatchPoints{
		Points:          pts,
		Database:        MyDB,
		RetentionPolicy: "default",
	}

	_, err := influxDBcon.Write(bps)

	if err != nil {
		log.Fatal(err)
	}
}
