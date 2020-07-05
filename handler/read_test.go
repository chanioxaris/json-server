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

func TestRead_Plural(t *testing.T) {
	randomPluralKeyIndex := rand.Intn(len(pluralKeys))
	randomPluralKey := pluralKeys[randomPluralKeyIndex]

	randomResourceIndex := rand.Intn(len(testData[randomPluralKey].([]storage.Resource)))
	randomResource := testData[randomPluralKey].([]storage.Resource)[randomResourceIndex]

	type bodyError struct {
		Error string `json:"error"`
	}

	testCases := []struct {
		name         string
		statusCode   int
		key          string
		id           string
		expectedData interface{}
		wantErr      bool
		err          error
	}{
		{
			name:         "Get plural resource with id",
			statusCode:   http.StatusOK,
			key:          randomPluralKey,
			id:           randomResource["id"].(string),
			expectedData: randomResource,
		},
		{
			name:       "Get plural resource invalid id",
			statusCode: http.StatusNotFound,
			key:        randomPluralKey,
			id:         "randomId",
			wantErr:    true,
			err:        storage.ErrResourceNotFound,
		},
	}

	for _, tt := range testCases {
		url := fmt.Sprintf("%s/%s/%s", mockServer.URL, tt.key, tt.id)

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

		if !tt.wantErr {
			var body storage.Resource
			if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(body, tt.expectedData) {
				t.Fatalf("expected body %v, but got %v", tt.expectedData, body)
			}
		} else {
			var body bodyError
			if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}

			if body.Error != tt.err.Error() {
				t.Fatalf("expected error message %v, but got %v", tt.err, body.Error)
			}
		}
	}
}

func TestRead_Singular(t *testing.T) {
	randomSingularKeyIndex := rand.Intn(len(singularKeys))
	randomSingularKey := singularKeys[randomSingularKeyIndex]

	randomResource := testData[randomSingularKey]

	type bodySuccess map[string]int

	testCases := []struct {
		name         string
		statusCode   int
		key          string
		expectedData interface{}
	}{
		{
			name:         "Get singular resource",
			statusCode:   http.StatusOK,
			key:          randomSingularKey,
			expectedData: randomResource,
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

		var body bodySuccess
		if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}

		if body[tt.key] != tt.expectedData {
			t.Fatalf("expected body %v, but got %v", tt.expectedData, body[tt.key])
		}
	}
}
