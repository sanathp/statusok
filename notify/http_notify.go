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
	FormParams  map[string]string `json:"params"` //Todo SHould be interface ?
}

func (httpNotify HttpNotify) Initialize() error {
	return nil
}

//TODO:post request to queue is not completly working test it
func (httpNotify HttpNotify) SendNotification(message Notification) error {
	var request *http.Request
	var reqErr error
	if len(httpNotify.FormParams) == 0 {
		request, reqErr = http.NewRequest(httpNotify.RequestType,
			httpNotify.Url,
			nil)
	} else {
		if httpNotify.Headers[ContentType] == JsonContentType {

			jsonBody, jsonErr := GetJsonParamsBody(httpNotify.FormParams)
			if jsonErr != nil {
				return jsonErr
			}
			request, reqErr = http.NewRequest(httpNotify.RequestType,
				httpNotify.Url,
				jsonBody)

		} else if httpNotify.Headers[ContentType] == FormContentType {
			urlParams := GetUrlParams(httpNotify.FormParams)
			request, reqErr = http.NewRequest(httpNotify.RequestType,
				httpNotify.Url,
				bytes.NewBufferString(urlParams.Encode()))
			request.Header.Add(ContentLength, strconv.Itoa(len(urlParams.Encode())))
		} else {
			urlParams := GetUrlParams(httpNotify.FormParams)
			request, reqErr = http.NewRequest(httpNotify.RequestType,
				httpNotify.Url,
				bytes.NewBufferString(urlParams.Encode()))

			request.Header.Add(ContentType, FormContentType)
			request.Header.Add(ContentLength, strconv.Itoa(len(urlParams.Encode())))
		}
	}

	if reqErr != nil {
		fmt.Println("Request Error : " + reqErr.Error())
	}

	AddHeaders(request, httpNotify.Headers)

	fmt.Println("PerformRequest")

	client := &http.Client{}

	getResponse, respErr := client.Do(request)

	if respErr != nil {
		fmt.Println("Response Error :" + respErr.Error())
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != http.StatusOK {
		fmt.Println("Request Status Error : Expected - ", httpNotify.Url, " Got %v", getResponse.Status)
	}

	return nil

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
