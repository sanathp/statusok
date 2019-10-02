package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"fmt"
)

type TelegramNotify struct {
	BotToken          string `json:"botToken"`
	ChatID            string `json:"chatID"`
}

type telegramPostMessage struct {
	ChatID   string `json:"chat_id"`
	Text     string `json:"text,omitempty"`
}

func (telegramNotify TelegramNotify) GetClientName() string {
	return "Telegram"
}

func (telegramNotify TelegramNotify) Initialize() error {

	if len(strings.TrimSpace(telegramNotify.BotToken)) == 0 {
		return errors.New("Telegram: botToken is a required field")
	}

	if len(strings.TrimSpace(telegramNotify.ChatID)) == 0 {
		return errors.New("Telegram: chatID is a required field")
	}

	return nil
}

func (telegramNotify TelegramNotify) SendResponseTimeNotification(responseTimeNotification ResponseTimeNotification) error {

	message := getMessageFromResponseTimeNotification(responseTimeNotification)

	payload, jsonErr := telegramNotify.getJsonParamBody(message)

	if jsonErr != nil {
		return jsonErr
	}

	requestURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegramNotify.BotToken)

	getResponse, respErr := http.Post(requestURL, "application/json", payload)

	if respErr != nil {
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != http.StatusOK {
		return errors.New("Telegram : Send notifaction failed. Response code " + strconv.Itoa(getResponse.StatusCode))
	}

	return nil
}

func (telegramNotify TelegramNotify) SendErrorNotification(errorNotification ErrorNotification) error {

	message := getMessageFromErrorNotification(errorNotification)

	payload, jsonErr := telegramNotify.getJsonParamBody(message)

	if jsonErr != nil {
		return jsonErr
	}

	requestURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegramNotify.BotToken)

	getResponse, respErr := http.Post(requestURL, "application/json", payload)

	if respErr != nil {
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != http.StatusOK {
		return errors.New("Telegram : Send notifaction failed. Response code " + strconv.Itoa(getResponse.StatusCode))
	}

	return nil
}

func (telegramNotify TelegramNotify) getJsonParamBody(message string) (io.Reader, error) {

	data, jsonErr := json.Marshal(telegramPostMessage{telegramNotify.ChatID,
		message,
	})

	if jsonErr != nil {

		jsonErr = errors.New("Invalid Parameters for Content-Type application/json : " + jsonErr.Error())

		return nil, jsonErr
	}

	return bytes.NewBuffer(data), nil
}
