package common_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/chanioxaris/json-server/storage"
)

func TestDB(t *testing.T) {
	testCases := []struct {
		name         string
		statusCode   int
		expectedData storage.Database
	}{
		{
			name:         "Get db",
			statusCode:   http.StatusOK,
			expectedData: testData,
		},
	}

	for _, tt := range testCases {
		if err := testResetData(fileName); err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("%s/db", mockServer.URL)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != tt.statusCode {
			t.Fatalf("expected status code %v, but got %v", tt.statusCode, resp.StatusCode)
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		expectedDataBytes, err := json.Marshal(tt.expectedData)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(bodyBytes, expectedDataBytes) {
			t.Fatalf("expected body %v, but got %v", expectedDataBytes, bodyBytes)
		}
	}
}
