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

func TestDelete(t *testing.T) {
	randomKeyIndex := rand.Intn(len(testResourceKeys))
	randomKey := testResourceKeys[randomKeyIndex]

	randomResourceIndex := rand.Intn(len(testData[randomKey]))
	randomResource := testData[randomKey][randomResourceIndex]

	type bodyError struct {
		Error string `json:"error"`
	}

	testCases := []struct {
		name       string
		statusCode int
		key        string
		id         string
		wantErr    bool
		err        error
	}{
		{
			name:       "Delete resource with id",
			statusCode: http.StatusOK,
			key:        randomKey,
			id:         randomResource["id"].(string),
		},
		{
			name:       "Delete resource not existing id",
			statusCode: http.StatusNotFound,
			key:        randomKey,
			id:         "randomId",
			wantErr:    true,
			err:        storage.ErrResourceNotFound,
		},
	}

	for _, tt := range testCases {
		testResetData(tt.key)

		url := fmt.Sprintf("%s/%s/%s", mockServer.URL, tt.key, tt.id)

		req, err := http.NewRequest(http.MethodDelete, url, nil)
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
			resources, err := testListResourcesByKey(tt.key)
			if err != nil {
				t.Fatal(err)
			}

			if expectedData := testData[randomKey]; !reflect.DeepEqual(resources, expectedData) {
				t.Fatalf("expected data %v, but got %v", expectedData, resources)
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
