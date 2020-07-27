package handler

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/chanioxaris/json-server/internal/storage"
	"github.com/chanioxaris/json-server/internal/web"
)

// Delete operates as a http handler, to delete an existing resource.
func Delete(storageSvc storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read request path parameter id.
		id := mux.Vars(r)["id"]

		// Delete resource.
		if err := storageSvc.Delete(id); err != nil {
			// Resource not found.
			if errors.Is(err, storage.ErrResourceNotFound) {
				web.Error(w, http.StatusNotFound, err.Error())
				return
			}

			web.Error(w, http.StatusInternalServerError, storage.ErrInternalServerError.Error())
			return
		}

		web.Success(w, http.StatusOK, nil)
	}
}
