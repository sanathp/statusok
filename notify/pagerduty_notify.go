package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PagerdutyNotify struct {
	Url        string `json:"url"`
	RoutingKey string `json:"routingKey"`
	Severity   string `json:"severity"`
}

type RequestBody struct {
	Payload     Payload `json:"payload"`
	RoutingKey  string  `json:"routing_key"`
	DedupKey    string  `json:"dedup_key"`
	EventAction string  `json:"event_action"`
}

type Payload struct {
	Summary   string `json:"summary"`
	Timestamp string `json:"timestamp"`
	Source    string `json:"source"`
	Severity  string `json:"severity"`
}

func (pagerdutyNotify PagerdutyNotify) GetClientName() string {
	return "Pagerduty API v2 Endpoint"
}

func (pagerdutyNotify PagerdutyNotify) Initialize() error {
	return nil
}

func (pagerdutyNotify PagerdutyNotify) SendResponseTimeNotification(responseTimeNotification ResponseTimeNotification) error {
	var request *http.Request
	var reqErr error

	msgParam := MessageParam{getMessageFromResponseTimeNotification(responseTimeNotification)}

	requestBody := CreatePagerdutyRequest(responseTimeNotification.Url, msgParam.Message, pagerdutyNotify)

	jsonBody, jsonErr := getJsonParamsBody(requestBody)
	if jsonErr != nil {
		return jsonErr
	}
	request, reqErr = http.NewRequest("POST", pagerdutyNotify.Url, jsonBody)

	if reqErr != nil {
		return reqErr
	}

	client := &http.Client{}

	getResponse, respErr := client.Do(request)

	if respErr != nil {
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != http.StatusAccepted {
		return errors.New(fmt.Sprintf("Pagerduty http response Status code expected: %v Got: %v ", http.StatusAccepted, getResponse.StatusCode))
	}

	return nil

}

func (pagerdutyNotify PagerdutyNotify) SendErrorNotification(errorNotification ErrorNotification) error {
	var request *http.Request
	var reqErr error

	msgParam := MessageParam{getMessageFromErrorNotification(errorNotification)}

	requestBody := CreatePagerdutyRequest(errorNotification.Url, msgParam.Message, pagerdutyNotify)

	jsonBody, jsonErr := getJsonParamsBody(requestBody)
	if jsonErr != nil {
		return jsonErr
	}
	request, reqErr = http.NewRequest("POST", pagerdutyNotify.Url, jsonBody)

	if reqErr != nil {
		fmt.Println(reqErr)
		return reqErr
	}

	client := &http.Client{}

	getResponse, respErr := client.Do(request)

	if respErr != nil {
		fmt.Println(respErr, pagerdutyNotify)
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != http.StatusAccepted {
		return errors.New(fmt.Sprintf("Pagerduty http response Status code expected: %v Got: %v ", http.StatusAccepted, getResponse.StatusCode))
	}

	return nil
}

func CreatePagerdutyRequest(url, summary string, config PagerdutyNotify) RequestBody {

	var requestBody RequestBody
	var payload Payload

	payload.Summary = summary
	payload.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000-0700")
	payload.Source = url
	payload.Severity = config.Severity

	requestBody.Payload = payload
	requestBody.RoutingKey = config.RoutingKey
	requestBody.DedupKey = url
	requestBody.EventAction = "trigger"
	return requestBody
}

func getJsonParamsBody(requestBody RequestBody) (io.Reader, error) {

	data, jsonErr := json.Marshal(requestBody)

	if jsonErr != nil {

		jsonErr = errors.New("Invalid Parameters for Content-Type application/json : " + jsonErr.Error())

		return nil, jsonErr
	}

	return bytes.NewBuffer(data), nil
}
