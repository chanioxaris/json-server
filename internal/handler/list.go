package handler

import (
	"net/http"

	"github.com/chanioxaris/json-server/internal/storage"
	"github.com/chanioxaris/json-server/internal/web"
)

// List operates as a http handler, to return all available resources.
func List(storageSvc storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Find all resources.
		data, err := storageSvc.Find()
		if err != nil {
			web.Error(w, http.StatusInternalServerError, storage.ErrInternalServerError.Error())
			return
		}

		web.Success(w, http.StatusOK, data)
	}
}
