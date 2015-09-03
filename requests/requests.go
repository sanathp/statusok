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
	"os"
	"strconv"
	"time"
)

//TODO: increase menory efficiency by using pointers ?
//TODO: use https://github.com/parnurzeal/gorequest ?
var (
	RequestsList   []RequestConfig
	requestChannel chan RequestConfig
	throttle       chan int
)

const (
	ContentType     = "Content-Type"
	ContentLength   = "Content-Length"
	FormContentType = "application/x-www-form-urlencoded"
	JsonContentType = "application/json"
)

type RequestConfig struct {
	Id           int
	Url          string            `json:"url"`
	RequestType  string            `json:"requestType"`
	Headers      map[string]string `json:"headers"`
	FormParams   map[string]string `json:"formParams"`
	UrlParams    map[string]string `json:"urlParams"`
	ResponseCode int               `json:"responseCode"`
	ResponseTime int64             `json:"responseTime"`
	CheckEvery   time.Duration     `json:"checkEvery"`
}

func RequestsInit(data []RequestConfig, concurrency int) {
	RequestsList = data

	if concurrency == 0 {
		//TODO: decide default value
		throttle = make(chan int, 1)
	} else {
		throttle = make(chan int, concurrency)
	}

	//TODO: decide lenght of buffer ?
	requestChannel = make(chan RequestConfig, len(data))

	if len(data) == 0 {
		println("\nNo requests to monitor.Please add requests to you config file")
		os.Exit(3)
	}

	println("\nSending requests to apis.....making sure everything is right before we start monitoring")
	println("Api Count: ", len(data))
	for i, requestConfig := range data {
		println("Request #", i, " : ", requestConfig.RequestType, " ", requestConfig.Url)

		_, urlErr := url.Parse(requestConfig.Url)
		if urlErr != nil {
			println("Invalid Url ", requestConfig.Url, " given for Request #", i, " Please verify your config file")
		}

		reqErr := PerformRequest(requestConfig, nil)

		if reqErr != nil {
			//Request Failed
			println("\nFailed !!!! Not able to perfome below request")
			println("\n----Request Deatails---")
			println("Url :", requestConfig.Url)
			println("Type :", requestConfig.RequestType)
			println("Error Reason :", reqErr.Error())
			println("\nPlease check the config file and try again")
			os.Exit(3)
		}
	}

	println("All requests Successfull")
}

func StartMonitoring() {
	fmt.Println("\nStarted Monitoring all ", len(RequestsList), " apis .....")

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
			requestChannel <- requestConfig
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
func listenToRequestChannel() {

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
		if throttle != nil {
			<-throttle
		}
	}()

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
					Id:           requestConfig.Id,
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

		} else {

			formParams := GetUrlValues(requestConfig.FormParams)
			request, reqErr = http.NewRequest(requestConfig.RequestType,
				requestConfig.Url,
				bytes.NewBufferString(formParams.Encode()))
			request.Header.Add(ContentLength, strconv.Itoa(len(formParams.Encode())))

			if requestConfig.Headers[ContentType] != "" {
				//Add content type to header if user doesnt mention it config file
				//Default content type application/x-www-form-urlencoded
				request.Header.Add(ContentType, FormContentType)
			}
		}
	}

	if reqErr != nil {
		//Not able to create Request object.Add Error to Database
		go database.AddErrorInfo(database.ErrorInfo{
			Id:           requestConfig.Id,
			Url:          requestConfig.Url,
			RequestType:  requestConfig.RequestType,
			ResponseCode: 0,
			ResponseBody: "",
			Reason:       database.ErrCreateRequest,
			OtherInfo:    reqErr.Error(),
		})

		return reqErr
	}

	if len(requestConfig.UrlParams) != 0 {
		urlParams := GetUrlValues(requestConfig.UrlParams)
		request.URL.RawQuery = urlParams.Encode()
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
			Id:           requestConfig.Id,
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
	fmt.Println(convertResponseToString(getResponse))

	if getResponse.StatusCode != requestConfig.ResponseCode {
		//Response code is not the expected one .Add Error to database
		go database.AddErrorInfo(database.ErrorInfo{
			Id:           requestConfig.Id,
			Url:          requestConfig.Url,
			RequestType:  requestConfig.RequestType,
			ResponseCode: getResponse.StatusCode,
			ResponseBody: convertResponseToString(getResponse),
			Reason:       errResposeCode(getResponse.StatusCode, requestConfig.ResponseCode),
			OtherInfo:    "",
		})
		return errResposeCode(getResponse.StatusCode, requestConfig.ResponseCode)
	}

	elapsed := time.Since(start)
	//Request succesfull . Add infomartion to Database
	go database.AddRequestInfo(database.RequestInfo{
		Id:                   requestConfig.Id,
		Url:                  requestConfig.Url,
		RequestType:          requestConfig.RequestType,
		ResponseCode:         getResponse.StatusCode,
		ResponseTime:         elapsed.Nanoseconds() / 1000000,
		ExpectedResponseTime: requestConfig.ResponseTime,
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

func GetUrlValues(params map[string]string) url.Values {
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

func errResposeCode(status int, expectedStatus int) error {
	return errors.New(fmt.Sprintf("Got Response code %v .Expeceted Response Code %v ", status, expectedStatus))
}
