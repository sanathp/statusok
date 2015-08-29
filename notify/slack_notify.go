package notify

import (
	"errors"
	"fmt"
	"github.com/bluele/slack"
	"strings"
)

type SlackNotify struct {
	Token     string `json:"token"`
	ChannelId string `json:"channelId"`
}

func (slackNotify SlackNotify) Initialize() error {

	if len(strings.TrimSpace(slackNotify.Token)) == 0 {
		return errors.New("Slack: Invalid Token")
	}

	if len(strings.TrimSpace(slackNotify.ChannelId)) == 0 {
		return errors.New("Slack: Invalid ChannelId")
	}

	return nil
}

func (slackNotify SlackNotify) SendNotification(message Notification) error {
	fmt.Println("slack notify called")
	api := slack.New(slackNotify.Token)

	err := api.ChatPostMessage(slackNotify.ChannelId, "Hello, world!", nil)
	if err != nil {
		fmt.Println("Slack error", err)
	}

	return nil
}
