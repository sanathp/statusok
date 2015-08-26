package request

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	//Define errors here
	//ErrHeaderTooLong = &ProtocolError{"header too long"}
	requests []RequestConfig
)

func createRequestsConfigArray(data []RequestConfig) {
	requests = make([]RequestConfig, 0)
	for _, value := range data {
		requests = append(requests, value)
	}
}

func PerformRequest(requestConfig RequestConfig) {

	client := &http.Client{}
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

	getResponse, respErr := client.Do(request)

	if respErr != nil {
		fmt.Println("Response Error :" + respErr.Error())
	}

	if getResponse.StatusCode != requestConfig.ResponseCode {
		fmt.Println("Request Status Error : Expected - ", requestConfig.Url, requestConfig.ResponseCode, " Got %v", getResponse.Status)
	}

	fmt.Println("Success")
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
			urlParams.Add("url", "http://google.com")
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
