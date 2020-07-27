package handler_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"testing"

	"github.com/chanioxaris/json-server/internal/storage"
)

func TestRead_Plural(t *testing.T) {
	randomPluralKeyIndex := rand.Intn(len(testResourceKeys))
	randomPluralKey := testResourceKeys[randomPluralKeyIndex]

	randomResourceIndex := rand.Intn(len(testData[randomPluralKey]))
	randomResource := testData[randomPluralKey][randomResourceIndex]

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
		testResetData(tt.key)

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
