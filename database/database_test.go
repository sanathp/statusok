package database

import (
	"errors"
	"testing"
)

func TestInitialize(t *testing.T) {

	ids := make(map[int]int64, 0)
	ids[1] = 10
	ids[2] = 2

	Initialize(ids, 10, 10)
	if len(responseMean) != len(ids) {
		t.Error("Ids not initialized")
	}

	if MeanResponseCount != 10 {
		t.Error("Mean Response Count Not Set")
	}

	if ErrorCount != 10 {
		t.Error("ErrorCount Not Set")
	}

}

func TestMeanResponseCalculation(t *testing.T) {
	ids := make(map[int]int64, 0)
	ids[1] = 10
	ids[2] = 2

	Initialize(ids, 1, 10)

	addResponseTimeToRequest(1, 10)
	mean, _ := getMeanResponseTimeOfUrl(1)

	if mean != 10 {
		t.Error("getMeanResponseTimeOfUrl Failed")
	}

}

func TestAddRequestAndErrorInfo(t *testing.T) {
	ids := make(map[int]int64, 0)
	ids[1] = 10
	ids[2] = 2

	Initialize(ids, 1, 10)

	requestInfo := RequestInfo{1, "http://test.com", "GET", 200, 10, 200}

	errorInfo := ErrorInfo{0, "http://test.com", "GET", 0, "test response", errors.New("test error"), "test other info"}

	AddErrorInfo(errorInfo)

	AddRequestInfo(requestInfo)

	mean, err := getMeanResponseTimeOfUrl(1)

	if mean != 10 {

		t.Error("Add Request Info Failed ", mean, err, responseMean[1], MeanResponseCount)
	}

}

func TestClearQueue(t *testing.T) {
	ids := make(map[int]int64, 0)
	ids[1] = 10
	ids[2] = 2

	Initialize(ids, 1, 10)

	addResponseTimeToRequest(1, 10)

	clearQueue(1)

	if len(responseMean[1]) != 0 {
		t.Error("ClearQueue Function is not working")
	}
}

func TestAddEmptyDatabase(t *testing.T) {
	ids := make(map[int]int64, 0)
	ids[1] = 10
	ids[2] = 2

	Initialize(ids, 1, 10)

	influxDb := InfluxDb{}

	databaseTypes := DatabaseTypes{influxDb}

	AddNew(databaseTypes)

	if len(dbList) != 0 {
		t.Error("Empty Databse should not be added to list")
	}
}

func TestAddValidDatabase(t *testing.T) {
	ids := make(map[int]int64, 0)
	ids[1] = 10
	ids[2] = 2

	Initialize(ids, 1, 10)

	influxDb := InfluxDb{"localhost", 8086, "statusok", "", ""}

	databaseTypes := DatabaseTypes{influxDb}

	AddNew(databaseTypes)

	if len(dbList) != 1 {
		t.Error("Not able to add databse to list")
	}
}
