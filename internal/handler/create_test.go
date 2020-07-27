package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"testing"

	"github.com/chanioxaris/json-server/internal/storage"
)

func TestCreate(t *testing.T) {
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
		body       storage.Resource
		wantErr    bool
		err        error
	}{
		{
			name:       "Create resource with id provided",
			statusCode: http.StatusCreated,
			key:        randomKey,
			body: storage.Resource{
				"id":      "2020",
				"field_1": "new-field_1",
				"field_2": "new-field_2",
			},
		},
		{
			name:       "Create resource without id provided",
			statusCode: http.StatusCreated,
			key:        randomKey,
			body: storage.Resource{
				"field_1": "new-field_1",
				"field_2": "new-field_2",
			},
		},
		{
			name:       "Create resource with empty body",
			statusCode: http.StatusBadRequest,
			key:        randomKey,
			body:       nil,
			wantErr:    true,
			err:        storage.ErrBadRequest,
		},
		{
			name:       "Create resource with body contains only id",
			statusCode: http.StatusBadRequest,
			key:        randomKey,
			body: storage.Resource{
				"id": "2020",
			},
			wantErr: true,
			err:     storage.ErrBadRequest,
		},
		{
			name:       "Create invalid resource with existing id",
			statusCode: http.StatusConflict,
			key:        randomKey,
			body: storage.Resource{
				"id":      randomResource["id"],
				"field_1": "new-field_1",
				"field_2": "new-field_2",
			},
			wantErr: true,
			err:     storage.ErrResourceAlreadyExists,
		},
	}

	for _, tt := range testCases {
		testResetData(tt.key)

		url := fmt.Sprintf("%s/%s", mockServer.URL, tt.key)

		bodyBytes, err := json.Marshal(tt.body)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
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
