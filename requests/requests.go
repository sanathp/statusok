package requests

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

var (
	RequestsList []RequestConfig
)

const (
	ContentType     = "Content-Type"
	ContentLength   = "Content-Length"
	FormContentType = "application/x-www-form-urlencoded"
	JsonContentType = "application/json"
)

type RequestConfig struct {
	Url          string            `json:"url"`
	RequestType  string            `json:"requestType"`
	Headers      map[string]string `json:"headers"`
	FormParams   map[string]string `json:"params"` //Todo SHould be interface ?
	ResponseCode int               `json:"responseCode"`
	Time         time.Duration     `json:"time"`
}

func RequestsInit(data []RequestConfig) {
	RequestsList = data
}

func StartMonitoring() {
	for _, requestConfig := range RequestsList {
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
			go requestConfig.PerformRequest()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func (requestConfig *RequestConfig) PerformRequest() {

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
				GetJsonParamsBody(requestConfig.FormParams))

		} else if requestConfig.Headers[ContentType] == FormContentType {
			urlParams := GetUrlParams(requestConfig.FormParams)
			request, reqErr = http.NewRequest(requestConfig.RequestType,
				requestConfig.Url,
				bytes.NewBufferString(urlParams.Encode()))
			request.Header.Add(ContentLength, strconv.Itoa(len(urlParams.Encode())))
		} else {
			urlParams := GetUrlParams(requestConfig.FormParams)
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

	AddHeaders(request, requestConfig.Headers)

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

	fmt.Println("Time Taken took %s", elapsed)
	database.AddToDatabase(database.Message{"hi", "gsdg", "gdgf", 12})

}

func AddHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Add(key, value)
	}
}

func GetUrlParams(params map[string]string) url.Values {
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

func GetJsonParamsBody(params map[string]string) io.Reader {
	data, jsonErr := json.Marshal(params)
	fmt.Println("json data ", string(data), " ", jsonErr)
	return bytes.NewBuffer(data)
}
