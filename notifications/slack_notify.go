package notifications

import (
	"fmt"
	"github.com/bluele/slack"
)

type SlackNotify struct {
	Token     string `json:"token"`
	ChannelId string `json:"channelId"`
}

func (slackNotify *SlackNotify) SendNotification(message Message) error {
	fmt.Println("slack notify called")
	api := slack.New(slackNotify.Token)

	err := api.ChatPostMessage(slackNotify.ChannelId, "Hello, world!", nil)
	if err != nil {
		fmt.Println("Slack error", err)
	}

	return nil
}
