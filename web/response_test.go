package web_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/chanioxaris/json-server/web"
)

func TestSuccess(t *testing.T) {
	type body struct {
		ID    int    `json:"id"`
		Field string `json:"field"`
	}

	testCases := []struct {
		name       string
		statusCode int
		data       interface{}
	}{
		{
			name:       "Response success without data",
			statusCode: http.StatusOK,
			data:       nil,
		},
		{
			name:       "Response success with data",
			statusCode: http.StatusCreated,
			data: body{
				ID:    1,
				Field: "testing success response",
			},
		},
	}

	for _, tt := range testCases {
		handler := func(w http.ResponseWriter, r *http.Request) {
			web.Success(w, tt.statusCode, tt.data)
		}

		req := httptest.NewRequest(http.MethodGet, "/success", nil)
		w := httptest.NewRecorder()
		handler(w, req)

		resp := w.Result()

		if resp.StatusCode != tt.statusCode {
			t.Fatalf("expected status code %v, but got %v", tt.statusCode, resp.StatusCode)
		}

		if tt.data != nil {
			if header := w.Header().Get("Content-Type"); header != "application/json" {
				t.Fatalf("expected header Content-Type %v, but got %v", "application/json", header)
			}

			var respBody body
			if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(respBody, tt.data) {
				t.Fatalf("expected body %v, but got %v", tt.data, respBody)
			}
		}
	}
}

func TestError(t *testing.T) {
	type body struct {
		Error string `json:"error"`
	}

	testCases := []struct {
		name       string
		statusCode int
		error      string
	}{
		{
			name:       "Response error",
			statusCode: http.StatusBadRequest,
			error:      "expected error message",
		},
	}

	for _, tt := range testCases {
		handler := func(w http.ResponseWriter, r *http.Request) {
			web.Error(w, tt.statusCode, tt.error)
		}

		req := httptest.NewRequest(http.MethodGet, "/error", nil)
		w := httptest.NewRecorder()
		handler(w, req)

		resp := w.Result()

		if resp.StatusCode != tt.statusCode {
			t.Fatalf("expected status code %v, but got %v", tt.statusCode, resp.StatusCode)
		}

		var respBody body
		if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			t.Fatal(err)
		}

		if respBody.Error != tt.error {
			t.Fatalf("expected error message %v, but got %v", tt.error, respBody.Error)
		}
	}
}
