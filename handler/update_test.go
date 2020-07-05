package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"testing"

	"github.com/chanioxaris/json-server/storage"
)

func TestUpdate(t *testing.T) {
	randomKeyIndex := rand.Intn(len(pluralKeys))
	randomKey := pluralKeys[randomKeyIndex]

	randomResourceIndex := rand.Intn(len(testData[randomKey].([]storage.Resource)))
	randomResource := testData[randomKey].([]storage.Resource)[randomResourceIndex]

	type bodyError struct {
		Error string `json:"error"`
	}

	testCases := []struct {
		name       string
		statusCode int
		key        string
		id         string
		body       storage.Resource
		wantErr    bool
		err        error
	}{
		{
			name:       "Update resource with id provided in body",
			statusCode: http.StatusOK,
			key:        randomKey,
			id:         randomResource["id"].(string),
			body: storage.Resource{
				"id":          randomResource["id"].(string),
				"description": "updated-description",
			},
		},
		{
			name:       "Update resource with different id provided in body",
			statusCode: http.StatusOK,
			key:        randomKey,
			id:         randomResource["id"].(string),
			body: storage.Resource{
				"id":          "2020",
				"description": "updated-description",
			},
		},
		{
			name:       "Update resource without id provided in body",
			statusCode: http.StatusOK,
			key:        randomKey,
			id:         randomResource["id"].(string),
			body: storage.Resource{
				"description": "updated-description",
			},
		},
		{
			name:       "Update resource with empty body",
			statusCode: http.StatusBadRequest,
			key:        randomKey,
			body:       nil,
			id:         randomResource["id"].(string),
			wantErr:    true,
			err:        storage.ErrBadRequest,
		},
		{
			name:       "Update resource with body contains only id",
			statusCode: http.StatusBadRequest,
			key:        randomKey,
			body: storage.Resource{
				"id": randomResource["id"].(string),
			},
			id:      randomResource["id"].(string),
			wantErr: true,
			err:     storage.ErrBadRequest,
		},
		{
			name:       "Update resource with not existing id",
			statusCode: http.StatusNotFound,
			key:        randomKey,
			body: storage.Resource{
				"description": "updated-description",
			},
			id:      "randomId",
			wantErr: true,
			err:     storage.ErrResourceNotFound,
		},
	}

	for _, tt := range testCases {
		if err := testResetData(fileName); err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("%s/%s/%s", mockServer.URL, tt.key, tt.id)

		bodyBytes, err := json.Marshal(tt.body)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(bodyBytes))
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
			var got storage.Resource
			if err = json.NewDecoder(resp.Body).Decode(&got); err != nil {
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

			var body []storage.Resource
			if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}

			expectedData := testData[randomKey].([]storage.Resource)
			expectedData[randomResourceIndex] = got

			if !reflect.DeepEqual(body, expectedData) {
				t.Fatalf("expected data %v, but got %v", expectedData, body)
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
