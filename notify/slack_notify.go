package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type SlackNotify struct {
	Username          string `json:"username"`
	ChannelName       string `json:"channelName"` //Not mandatory field
	ChannelWebhookURL string `json:"channelWebhookURL"`
	IconUrl           string `json:"iconUrl"`
}

type postMessage struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Text     string `json:"text,omitempty"`
	//TODO: test if icon url is working
	Icon_url string `json:"icon_url"`
}

func (slackNotify SlackNotify) GetClientName() string {
	return "Slack"
}

func (slackNotify SlackNotify) Initialize() error {

	if len(strings.TrimSpace(slackNotify.Username)) == 0 {
		return errors.New("Slack: Username is a required field")
	}

	if len(strings.TrimSpace(slackNotify.ChannelWebhookURL)) == 0 {
		return errors.New("Slack: channelWebhookURL is a required field")
	}

	return nil
}

func (slackNotify SlackNotify) SendResponseTimeNotification(responseTimeNotification ResponseTimeNotification) error {

	message := fmt.Sprintf("Notifiaction From StatusOk\nOne of your apis response time is below than expected."+
		"\nPlease find the Details below"+
		"\nUrl: %v \nRequestType: %v \nCurrent Average Response Time: %v \n Expected Response Time: %v\n"+
		"\nThanks", responseTimeNotification.Url, responseTimeNotification.RequestType, responseTimeNotification.MeanResponseTime, responseTimeNotification.ExpectedResponsetime)

	payload, jsonErr := slackNotify.getJsonParamBody(message)

	if jsonErr != nil {
		return jsonErr
	}

	getResponse, respErr := http.Post(slackNotify.ChannelWebhookURL, "application/json", payload)

	if respErr != nil {
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != http.StatusOK {
		return errors.New("Slack : Send notifaction failed. Response code " + strconv.Itoa(getResponse.StatusCode))
	}

	return nil
}

func (slackNotify SlackNotify) SendErrorNotification(errorNotification ErrorNotification) error {
	//TODO: move this to util class or make local functions in all notifers
	message := fmt.Sprintf("Notifiaction From StatusOk\nWe are getting error when we try to send request to one of your apis"+
		"\nPlease find the Details below"+
		"\nUrl: %v \nRequestType: %v \nError Message: %v \n Response Body: %v\n Other Info:%v\n"+
		"\nThanks", errorNotification.Url, errorNotification.RequestType, errorNotification.Error, errorNotification.ResponseBody, errorNotification.OtherInfo)

	payload, jsonErr := slackNotify.getJsonParamBody(message)

	if jsonErr != nil {
		return jsonErr
	}

	getResponse, respErr := http.Post(slackNotify.ChannelWebhookURL, "application/json", payload)

	if respErr != nil {
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != http.StatusOK {
		return errors.New("Slack : Send notifaction failed. Response code " + strconv.Itoa(getResponse.StatusCode))
	}

	return nil
}

func (slackNotify SlackNotify) getJsonParamBody(message string) (io.Reader, error) {

	data, jsonErr := json.Marshal(postMessage{slackNotify.ChannelName,
		slackNotify.Username,
		message,
		slackNotify.IconUrl,
	})

	if jsonErr != nil {

		jsonErr = errors.New("Invalid Parameters for Content-Type application/json : " + jsonErr.Error())

		return nil, jsonErr
	}

	return bytes.NewBuffer(data), nil
}
