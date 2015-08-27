package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sanathp/StatusOk/database"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	ContentType     = "Content-Type"
	ContentLength   = "Content-Length"
	FormContentType = "application/x-www-form-urlencoded"
	JsonContentType = "application/json"
)

type RequestJson struct {
	Requests []RequestConfig `json:"requests"`
}

type RequestConfig struct {
	Url          string            `json:"url"`
	RequestType  string            `json:"requestType"`
	Headers      map[string]string `json:"headers"`
	FormParams   map[string]string `json:"params"` //Todo SHould be interface ?
	ResponseCode int               `json:"responseCode"`
	Time         time.Duration     `json:"time"`
}

var (
	requests []RequestConfig
)

func StartMonitoring(data RequestJson) {
	fmt.Println("lenght ", len(data.Requests))
	for _, requestConfig := range data.Requests {
		go createTicker(requestConfig)
	}
}

func createTicker(requestConfig RequestConfig) {

	fmt.Println("createTicker")
	var ticker *time.Ticker = time.NewTicker(requestConfig.Time * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			go PerformRequest(requestConfig)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func PerformRequest(requestConfig RequestConfig) {

	fmt.Println("PerformRequest")
	var request *http.Request
	var reqErr error
	if len(requestConfig.FormParams) == 0 {
		request, reqErr = http.NewRequest(requestConfig.RequestType,
			requestConfig.Url,
			nil)
	} else {
		if requestConfig.Headers[ContentType] == JsonContentType {
			request, reqErr = http.NewRequest(requestConfig.RequestType,
				requestConfig.Url,
				getJsonParamsBody(requestConfig.FormParams))

		} else if requestConfig.Headers[ContentType] == FormContentType {
			urlParams := getUrlParams(requestConfig.FormParams)
			request, reqErr = http.NewRequest(requestConfig.RequestType,
				requestConfig.Url,
				bytes.NewBufferString(urlParams.Encode()))
			request.Header.Add(ContentLength, strconv.Itoa(len(urlParams.Encode())))
		} else {
			urlParams := getUrlParams(requestConfig.FormParams)
			request, reqErr = http.NewRequest(requestConfig.RequestType,
				requestConfig.Url,
				bytes.NewBufferString(urlParams.Encode()))

			request.Header.Add(ContentType, FormContentType)
			request.Header.Add(ContentLength, strconv.Itoa(len(urlParams.Encode())))
		}
	}

	if reqErr != nil {
		fmt.Println("Request Error : " + reqErr.Error())
	}

	addHeaders(request, requestConfig.Headers)

	fmt.Println("PerformRequest")
	timeout := time.Duration(10 * time.Second)

	client := &http.Client{
		Timeout: timeout,
	}

	start := time.Now()

	getResponse, respErr := client.Do(request)

	if respErr != nil {
		fmt.Println("Response Error :" + respErr.Error())
		return
	}

	elapsed := time.Since(start)

	defer getResponse.Body.Close()

	if getResponse.StatusCode != requestConfig.ResponseCode {
		fmt.Println("Request Status Error : Expected - ", requestConfig.Url, requestConfig.ResponseCode, " Got %v", getResponse.Status)
	}

	database.WritePoints(requestConfig.Url, elapsed.Nanoseconds()/1000000, requestConfig.RequestType)

	fmt.Println("Time Taken took %s", elapsed)
	database.GetMeanResponseTime(requestConfig.Url, 5)

	database.GetErrorsCount(requestConfig.Url, 5)

}

func addHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Add(key, value)
	}
}

func getUrlParams(params map[string]string) url.Values {
	urlParams := url.Values{}
	i := 0
	for key, value := range params {
		if i == 0 {
			urlParams.Set(key, value)
		} else {
			urlParams.Add(key, value)
		}
	}

	fmt.Println("url values", urlParams)

	return urlParams
}

func getJsonParamsBody(params map[string]string) io.Reader {
	data, jsonErr := json.Marshal(params)
	fmt.Println("json data ", string(data), " ", jsonErr)
	return bytes.NewBuffer(data)
}
