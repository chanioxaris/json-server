package handler_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"testing"

	"github.com/chanioxaris/json-server/storage"
)

func TestList(t *testing.T) {
	randomKeyIndex := rand.Intn(len(pluralKeys))
	randomKey := pluralKeys[randomKeyIndex]

	testCases := []struct {
		name         string
		statusCode   int
		key          string
		expectedData interface{}
	}{
		{
			name:         "List resources",
			statusCode:   http.StatusOK,
			key:          randomKey,
			expectedData: testData[randomKey],
		},
	}

	for _, tt := range testCases {
		if err := testResetData(fileName); err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("%s/%s", mockServer.URL, tt.key)

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

		var body []storage.Resource
		if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(body, tt.expectedData) {
			t.Fatalf("expected body %v, but got %v", tt.expectedData, body)
		}
	}
}
