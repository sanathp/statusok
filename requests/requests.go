package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sanathp/StatusOk/database"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//TODO: increase menory efficiency by using pointers ?
var (
	RequestsList   []RequestConfig
	requestChannel chan RequestConfig
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
	ResponseTime int64             `json:"responseTime"`
	CheckEvery   time.Duration     `json:"checkEvery"`
}

func RequestsInit(data []RequestConfig) {
	RequestsList = data
	//TODO: decide lenght of buffer ?
	requestChannel = make(chan RequestConfig, len(data))
}

func StartMonitoring() {
	fmt.Println("Number of requests = ", len(RequestsList))
	go listenToRequestChannel()
	for _, requestConfig := range RequestsList {
		go createTicker(requestConfig)
	}
}

func createTicker(requestConfig RequestConfig) {

	var ticker *time.Ticker = time.NewTicker(requestConfig.CheckEvery * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			//TODO: instead of directly permofrming a request to write to a channel.
			//decide how many requests to keep in channel. take value from config
			//go requestConfig.PerformRequest()
			fmt.Println("add to channel ", requestConfig.Url)
			requestChannel <- requestConfig
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
func listenToRequestChannel() {
	//var wg sync.WaitGroup
	var throttle = make(chan int, 1)

	for {
		select {
		case requect := <-requestChannel:
			throttle <- 1
			go PerformRequest(requect, throttle)
		}
	}

}

func PerformRequest(requestConfig RequestConfig, throttle chan int) error {
	defer func() {
		<-throttle
	}()
	fmt.Println("PerformRequest ", requestConfig.Url)
	var request *http.Request
	var reqErr error

	if len(requestConfig.FormParams) == 0 {
		request, reqErr = http.NewRequest(requestConfig.RequestType,
			requestConfig.Url,
			nil)
	} else {
		if requestConfig.Headers[ContentType] == JsonContentType {
			jsonBody, jsonErr := GetJsonParamsBody(requestConfig.FormParams)
			if jsonErr != nil {
				//Not able to create Request object.Add Error to Database
				go database.AddErrorInfo(database.ErrorInfo{
					Url:          requestConfig.Url,
					RequestType:  requestConfig.RequestType,
					ResponseCode: 0,
					ResponseBody: "",
					Reason:       database.ErrCreateRequest,
					OtherInfo:    jsonErr.Error(),
				})

				return jsonErr
			}
			request, reqErr = http.NewRequest(requestConfig.RequestType,
				requestConfig.Url,
				jsonBody)

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
		//Not able to create Request object.Add Error to Database
		go database.AddErrorInfo(database.ErrorInfo{
			Url:          requestConfig.Url,
			RequestType:  requestConfig.RequestType,
			ResponseCode: 0,
			ResponseBody: "",
			Reason:       database.ErrCreateRequest,
			OtherInfo:    reqErr.Error(),
		})

		return reqErr
	}

	AddHeaders(request, requestConfig.Headers)

	//TODO: put timeout ?
	/*
		timeout := 10 * requestConfig.ResponseTime

		client := &http.Client{
			Timeout: timeout,
		}
	*/

	client := &http.Client{}
	start := time.Now()

	getResponse, respErr := client.Do(request)

	if respErr != nil {
		//Request failed . Add error info to database
		var statusCode int
		if getResponse == nil {
			statusCode = 0
		} else {
			statusCode = getResponse.StatusCode
		}
		go database.AddErrorInfo(database.ErrorInfo{
			Url:          requestConfig.Url,
			RequestType:  requestConfig.RequestType,
			ResponseCode: statusCode,
			ResponseBody: convertResponseToString(getResponse),
			Reason:       database.ErrDoRequest,
			OtherInfo:    respErr.Error(),
		})
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != requestConfig.ResponseCode {
		//Response code is not the expected one .Add Error to database
		go database.AddErrorInfo(database.ErrorInfo{
			Url:          requestConfig.Url,
			RequestType:  requestConfig.RequestType,
			ResponseCode: getResponse.StatusCode,
			ResponseBody: convertResponseToString(getResponse),
			Reason:       database.ErrResposeCode,
			OtherInfo:    "",
		})
		return database.ErrResposeCode
	}

	elapsed := time.Since(start)
	//Request succesfull . Add infomartion to Database
	go database.AddRequestInfo(database.RequestInfo{
		Url:          requestConfig.Url,
		RequestType:  requestConfig.RequestType,
		ResponseCode: getResponse.StatusCode,
		ResponseTime: elapsed.Nanoseconds() / 1000000,
	})

	return nil
}

//convert response body to string
func convertResponseToString(resp *http.Response) string {
	if resp == nil {
		return " "
	}
	buf := new(bytes.Buffer)
	_, bufErr := buf.ReadFrom(resp.Body)

	if bufErr != nil {
		return " "
	}

	return buf.String()
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

	return urlParams
}

func GetJsonParamsBody(params map[string]string) (io.Reader, error) {

	data, jsonErr := json.Marshal(params)

	if jsonErr != nil {

		jsonErr = errors.New("Invalid Parameters for Content-Type application/json : " + jsonErr.Error())

		return nil, jsonErr
	}

	return bytes.NewBuffer(data), nil
}
