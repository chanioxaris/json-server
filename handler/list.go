package handler

import (
	"net/http"

	"github.com/chanioxaris/json-server/storage"
	"github.com/chanioxaris/json-server/web"
)

// List operates as a http handle, to return all available resources.
func List(storageSvc storage.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Find all resources.
		data, err := storageSvc.Find()
		if err != nil {
			web.Error(w, http.StatusInternalServerError, "internal Server Error")
			return
		}

		web.Success(w, http.StatusOK, data)
	}
}
