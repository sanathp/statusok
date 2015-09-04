package notify

import (
	"bytes"
	"encoding/json"
	"errors"
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

	message := getMessageFromResponseTimeNotification(responseTimeNotification)

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

	message := getMessageFromErrorNotification(errorNotification)

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
