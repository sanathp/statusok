package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	ContentType     = "Content-Type"
	ContentLength   = "Content-Length"
	FormContentType = "application/x-www-form-urlencoded"
	JsonContentType = "application/json"
)

type HttpNotify struct {
	Url         string            `json:"url"`
	RequestType string            `json:"requestType"`
	Headers     map[string]string `json:"headers"`
}

type MessageParam struct {
	Message string `json:"message"`
}

func (httpNotify HttpNotify) GetClientName() string {
	return "Http End Point"
}

func (httpNotify HttpNotify) Initialize() error {
	return nil
}

func (httpNotify HttpNotify) SendResponseTimeNotification(responseTimeNotification ResponseTimeNotification) error {
	var request *http.Request
	var reqErr error

	msgParam := MessageParam{getMessageFromResponseTimeNotification(responseTimeNotification)}

	if httpNotify.Headers[ContentType] == JsonContentType {

		jsonBody, jsonErr := GetJsonParamsBody(msgParam)
		if jsonErr != nil {
			return jsonErr
		}
		request, reqErr = http.NewRequest(httpNotify.RequestType,
			httpNotify.Url,
			jsonBody)

	} else if httpNotify.Headers[ContentType] == FormContentType {
		urlParams := GetUrlValues(msgParam)
		request, reqErr = http.NewRequest(httpNotify.RequestType,
			httpNotify.Url,
			bytes.NewBufferString(urlParams.Encode()))
		request.Header.Add(ContentLength, strconv.Itoa(len(urlParams.Encode())))
	} else {
		urlParams := GetUrlValues(msgParam)
		request, reqErr = http.NewRequest(httpNotify.RequestType,
			httpNotify.Url,
			bytes.NewBufferString(urlParams.Encode()))

		request.Header.Add(ContentType, FormContentType)
		request.Header.Add(ContentLength, strconv.Itoa(len(urlParams.Encode())))
	}

	if reqErr != nil {
		return reqErr
	}

	AddHeaders(request, httpNotify.Headers)

	client := &http.Client{}

	getResponse, respErr := client.Do(request)

	if respErr != nil {
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Http response Status code expected: %v Got : %v ", http.StatusOK, getResponse.StatusCode))
	}

	return nil

}

func (httpNotify HttpNotify) SendErrorNotification(errorNotification ErrorNotification) error {
	var request *http.Request
	var reqErr error

	msgParam := MessageParam{getMessageFromErrorNotification(errorNotification)}

	if httpNotify.Headers[ContentType] == JsonContentType {

		jsonBody, jsonErr := GetJsonParamsBody(msgParam)
		if jsonErr != nil {
			return jsonErr
		}
		request, reqErr = http.NewRequest(httpNotify.RequestType,
			httpNotify.Url,
			jsonBody)

	} else if httpNotify.Headers[ContentType] == FormContentType {
		urlParams := GetUrlValues(msgParam)
		request, reqErr = http.NewRequest(httpNotify.RequestType,
			httpNotify.Url,
			bytes.NewBufferString(urlParams.Encode()))
		request.Header.Add(ContentLength, strconv.Itoa(len(urlParams.Encode())))
	} else {
		urlParams := GetUrlValues(msgParam)
		request, reqErr = http.NewRequest(httpNotify.RequestType,
			httpNotify.Url,
			bytes.NewBufferString(urlParams.Encode()))

		request.Header.Add(ContentType, FormContentType)
		request.Header.Add(ContentLength, strconv.Itoa(len(urlParams.Encode())))
	}

	if reqErr != nil {
		fmt.Println(reqErr)
		return reqErr
	}

	AddHeaders(request, httpNotify.Headers)

	client := &http.Client{}

	getResponse, respErr := client.Do(request)

	if respErr != nil {
		fmt.Println(respErr, httpNotify)
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Http response Status code expected: %v Got : %v ", http.StatusOK, getResponse.StatusCode))
	}

	return nil
}

func AddHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Add(key, value)
	}
}

func GetUrlValues(msgParam MessageParam) url.Values {
	urlParams := url.Values{}
	urlParams.Set("message", msgParam.Message)
	return urlParams
}

func GetJsonParamsBody(msgParam MessageParam) (io.Reader, error) {

	data, jsonErr := json.Marshal(msgParam)

	if jsonErr != nil {

		jsonErr = errors.New("Invalid Parameters for Content-Type application/json : " + jsonErr.Error())

		return nil, jsonErr
	}

	return bytes.NewBuffer(data), nil
}

func getStringFromResponseBody(body io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	return buf.String()
}
