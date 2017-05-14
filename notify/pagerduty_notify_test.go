package notify

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestCreatePagerdutyNotificationClient(t *testing.T) {
	pagerdutyNotify := PagerdutyNotify{"https://events.pagerduty.com/v2/enqueue", "abcdefghijklmnopqrstuvwxyz123456", "info"}

	err := pagerdutyNotify.Initialize()
	clientName := pagerdutyNotify.GetClientName()

	if err != nil {
		t.Error(err)
	}

	if clientName != "Pagerduty API v2 Endpoint" {
		t.Error("Client name error")
	}
}

func TestSendResponseTimeNotification(t *testing.T) {
	pagerdutyNotify := PagerdutyNotify{"https://events.pagerduty.com/v2/enqueue", "abcdefghijklmnopqrstuvwxyz123456", "info"}

	err := pagerdutyNotify.SendResponseTimeNotification(ResponseTimeNotification{"http://test.com", "GET", 700, 800})

	if err != nil {
		t.Error(err)
	}
}

func TestInvalidSendResponseTimeNotification(t *testing.T) {
	pagerdutyNotify := PagerdutyNotify{"https://events.pagerduty.com/v2/enqueue", "abcdefghijklmnopqrstuvwxyz123456", "info"}

	err := pagerdutyNotify.SendResponseTimeNotification(ResponseTimeNotification{})

	if !(strings.Contains(err.Error(), fmt.Sprintf("Pagerduty http response Status code expected: %v Got: %v ", http.StatusAccepted, http.StatusBadRequest))) {
		t.Error("Unexpected error message: " + err.Error())
	}
}

func TestSendErrorNotification(t *testing.T) {
	pagerdutyNotify := PagerdutyNotify{"https://events.pagerduty.com/v2/enqueue", "abcdefghijklmnopqrstuvwxyz123456", "info"}

	err := pagerdutyNotify.SendErrorNotification(ErrorNotification{"http://test.com", "GET", "This is test notification", "Test notiification", "test"})

	if err != nil {
		t.Error(err)
	}
}

func TestInvalidSendErrorNotification(t *testing.T) {
	pagerdutyNotify := PagerdutyNotify{"https://events.pagerduty.com/v2/enqueue", "abcdefghijklmnopqrstuvwxyz123456", "info"}

	err := pagerdutyNotify.SendErrorNotification(ErrorNotification{})

	if !(strings.Contains(err.Error(), fmt.Sprintf("Pagerduty http response Status code expected: %v Got: %v ", http.StatusAccepted, http.StatusBadRequest))) {
		t.Error("Unexpected error message: " + err.Error())
	}
}

func TestCreatePagerdutyRequest(t *testing.T) {
	pagerdutyUrl := "https://events.pagerduty.com/v2/enqueue"
	pagerdutyKey := "abcdefghijklmnopqrstuvwxyz123456"
	pagerdutySeverity := "info"
	notificationUrl := "http://test.com"
	notificationSummary := "Test summary to send to Pagerduty"

	pagerdutyNotify := PagerdutyNotify{pagerdutyUrl, pagerdutyKey, pagerdutySeverity}

	requestBody := CreatePagerdutyRequest(notificationUrl, notificationSummary, pagerdutyNotify)

	if requestBody.Payload.Summary != notificationSummary {
		t.Error("Request summary invalid")
	}

	if requestBody.Payload.Timestamp == "" {
		t.Error("Request time invalid")
	}

	if requestBody.Payload.Source != notificationUrl {
		t.Error("Request source url invalid")
	}

	if requestBody.Payload.Severity != pagerdutySeverity {
		t.Error("Request severity invalid")
	}

	if requestBody.RoutingKey != pagerdutyKey {
		t.Error("Request routing key invalid")
	}

	if requestBody.DedupKey != notificationUrl {
		t.Error("Request dedup key invalid")
	}

	if requestBody.EventAction != "trigger" {
		t.Error("Request event action invalid")
	}
}
