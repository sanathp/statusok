package notify

import (
	"bytes"
	"fmt"
	"github.com/sanathp/StatusOk/requests"
	"net/http"
	"strconv"
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

func (httpNotify HttpNotify) SendNotification(message Notification) error {
	var request *http.Request
	var reqErr error
	if len(httpNotify.FormParams) == 0 {
		request, reqErr = http.NewRequest(httpNotify.RequestType,
			httpNotify.Url,
			nil)
	} else {
		if httpNotify.Headers[requests.ContentType] == requests.JsonContentType {
			request, reqErr = http.NewRequest(httpNotify.RequestType,
				httpNotify.Url,
				requests.GetJsonParamsBody(httpNotify.FormParams))

		} else if httpNotify.Headers[requests.ContentType] == requests.FormContentType {
			urlParams := requests.GetUrlParams(httpNotify.FormParams)
			request, reqErr = http.NewRequest(httpNotify.RequestType,
				httpNotify.Url,
				bytes.NewBufferString(urlParams.Encode()))
			request.Header.Add(requests.ContentLength, strconv.Itoa(len(urlParams.Encode())))
		} else {
			urlParams := requests.GetUrlParams(httpNotify.FormParams)
			request, reqErr = http.NewRequest(httpNotify.RequestType,
				httpNotify.Url,
				bytes.NewBufferString(urlParams.Encode()))

			request.Header.Add(requests.ContentType, requests.FormContentType)
			request.Header.Add(requests.ContentLength, strconv.Itoa(len(urlParams.Encode())))
		}
	}

	if reqErr != nil {
		fmt.Println("Request Error : " + reqErr.Error())
	}

	requests.AddHeaders(request, httpNotify.Headers)

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
