package requests

import (
	"testing"
)

func TestRequestsInit(t *testing.T) {

	data := make([]RequestConfig, 0)
	google := RequestConfig{Id: 1, Url: "http://google.com", RequestType: "GET", ResponseCode: 200, ResponseTime: 100, CheckEvery: 1}
	data = append(data, google)

	RequestsInit(data, 0)

	if len(RequestsList) != 1 {
		t.Error("Request initalize failed")
	}
}

func TestGetRequest(t *testing.T) {
	google := RequestConfig{Id: 1, Url: "http://google.com", RequestType: "GET", ResponseCode: 200, ResponseTime: 100, CheckEvery: 1}

	err := PerformRequest(google, nil)

	if err != nil {
		t.Error("Get Request Failed")
	}
}

func TestInvalidGetRequest(t *testing.T) {
	invalid := RequestConfig{Id: 1, Url: "http://localhost:64521", RequestType: "GET", ResponseCode: 200, ResponseTime: 100, CheckEvery: 1}

	err := PerformRequest(invalid, nil)

	if err == nil {
		t.Error("Invalid Get Request Succeded")
	}
}

func TestInvalidPostRequest(t *testing.T) {
	google := RequestConfig{Id: 1, Url: "http://google.com", RequestType: "POST", ResponseCode: 200, ResponseTime: 100, CheckEvery: 1}

	err := PerformRequest(google, nil)

	if err == nil {
		t.Error("Invalid POST Request Succeded")
	}
}

func TestRequestLimitOne(t *testing.T) {
	google := RequestConfig{Id: 1, Url: "http://google.com", RequestType: "GET", RequestLimit: 1}

	err := PerformRequest(google, nil)

	if err == nil {
		t.Error("Request limit 1 should not allow any redirects")
	}
}

func TestRequestLimitBig(t *testing.T) {
	google := RequestConfig{Id: 1, Url: "http://google.com", RequestType: "GET", RequestLimit: 10}

	err := PerformRequest(google, nil)

	if err != nil {
		t.Error("Invalid checking for request limit")
	}
}

func TestResponseCode(t *testing.T) {
	google := RequestConfig{Id: 1, Url: "http://google.com", RequestType: "GET", ResponseCode: 404}

	err := PerformRequest(google, nil)

	if err == nil {
		t.Error("Invalid checking for response code")
	}
}
