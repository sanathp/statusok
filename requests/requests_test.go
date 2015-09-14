package requests

import (
	"testing"
)

func TestRequestsInit(t *testing.T) {

	data := make([]RequestConfig, 0)
	google := RequestConfig{1, "http://google.com", "GET", nil, nil, nil, 200, 100, 1}
	data = append(data, google)

	RequestsInit(data, 0)

	if len(RequestsList) != 1 {
		t.Error("Request initalize failed")
	}
}

func TestGetRequest(t *testing.T) {
	google := RequestConfig{1, "http://google.com", "GET", nil, nil, nil, 200, 100, 1}

	err := PerformRequest(google, nil)

	if err != nil {
		t.Error("Get Request Failed")
	}
}

func TestInvalidGetRequest(t *testing.T) {
	invalid := RequestConfig{1, "http://localhost:64521", "GET", nil, nil, nil, 200, 100, 1}

	err := PerformRequest(invalid, nil)

	if err == nil {
		t.Error("Invalid Get Request Succeded")
	}
}

func TestInvalidPostRequest(t *testing.T) {
	google := RequestConfig{1, "http://google.com", "POST", nil, nil, nil, 200, 100, 1}

	err := PerformRequest(google, nil)

	if err == nil {
		t.Error("Invalid POST Request Succeded")
	}
}
