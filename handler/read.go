package handler

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/chanioxaris/json-server/storage"
	"github.com/chanioxaris/json-server/web"
)

// Read operates as a http handler, to return the requested resource by id.
func Read(storageSvc storage.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read request path parameter id.
		id := mux.Vars(r)["id"]

		// Find the resource with the requested id.
		data, err := storageSvc.FindById(id)
		if err != nil {
			// Resource not found.
			if errors.Is(err, storage.ErrResourceNotFound) {
				web.Error(w, http.StatusNotFound, err.Error())
				return
			}

			web.Error(w, http.StatusInternalServerError, storage.ErrInternalServerError.Error())
			return
		}

		web.Success(w, http.StatusOK, data)
	}
}
