package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type DingdingNotify struct {
	HttpNotify
}

type Txt struct {
	Content string `json:"content"`
}

type Message struct {
	MessageType string `json:"msgtype"`
	Text        Txt    `json:"text"`
}

func (dingdingNotify DingdingNotify) GetClientName() string {
	return "Dingding"
}

func (dingdingNotify DingdingNotify) Initialize() error {
	return nil
}

func (dingdingNotify DingdingNotify) SendResponseTimeNotification(responseTimeNotification ResponseTimeNotification) error {
	var request *http.Request
	var reqErr error

	msgParam := Message{
		MessageType: "text", // msgtype as text
		Text:        Txt{getMessageFromResponseTimeNotification(responseTimeNotification)},
	}
	fmt.Printf("%v", msgParam)

	if dingdingNotify.Headers[ContentType] == JsonContentType {

		jsonBody, jsonErr := getJsonParamsBodyDingding(msgParam)
		if jsonErr != nil {
			return jsonErr
		}
		request, reqErr = http.NewRequest(dingdingNotify.RequestType,
			dingdingNotify.Url,
			jsonBody)
	}

	if reqErr != nil {
		return reqErr
	}

	AddHeaders(request, dingdingNotify.Headers)

	client := &http.Client{}

	getResponse, respErr := client.Do(request)
	fmt.Printf("%v", respErr)

	if respErr != nil {
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Http response Status code expected: %v Got : %v ", http.StatusOK, getResponse.StatusCode))
	}

	return nil

}

func (dingdingNotify DingdingNotify) SendErrorNotification(errorNotification ErrorNotification) error {
	var request *http.Request
	var reqErr error

	msgParam := Message{
		"text", // msgtype as text
		Txt{getMessageFromErrorNotification(errorNotification)},
	}

	if dingdingNotify.Headers[ContentType] == JsonContentType {

		jsonBody, jsonErr := getJsonParamsBodyDingding(msgParam)
		if jsonErr != nil {
			return jsonErr
		}
		request, reqErr = http.NewRequest(dingdingNotify.RequestType,
			dingdingNotify.Url,
			jsonBody)

	}

	if reqErr != nil {
		fmt.Println(reqErr)
		return reqErr
	}

	AddHeaders(request, dingdingNotify.Headers)

	client := &http.Client{}

	getResponse, respErr := client.Do(request)

	if respErr != nil {
		fmt.Println(respErr, dingdingNotify)
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Http response Status code expected: %v Got : %v ", http.StatusOK, getResponse.StatusCode))
	}

	return nil
}

func getJsonParamsBodyDingding(msgParam Message) (io.Reader, error) {

	data, jsonErr := json.Marshal(msgParam)

	if jsonErr != nil {

		jsonErr = errors.New("Invalid Parameters for Content-Type application/json : " + jsonErr.Error())

		return nil, jsonErr
	}

	return bytes.NewBuffer(data), nil
}
