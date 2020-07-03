package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/chanioxaris/json-server/storage"
	"github.com/chanioxaris/json-server/web"
)

// Update operates as a http handle, to update an existing resource.
func Update(storageSvc storage.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read request path parameter id.
		id := mux.Vars(r)["id"]

		// Check if resource with the requested id exists.
		if _, err := storageSvc.FindById(id); err != nil {
			// Resource not found.
			if errors.Is(err, storage.ErrResourceNotFound) {
				web.Error(w, http.StatusNotFound, err.Error())
				return
			}

			web.Error(w, http.StatusInternalServerError, "internal Server Error")
			return
		}

		// Read and decode request body.
		var newResource storage.Resource
		if err := json.NewDecoder(r.Body).Decode(&newResource); err != nil {
			web.Error(w, http.StatusBadRequest, "bad request")
			return
		}

		// Check if request body is empty, or contains only id.
		if _, ok := newResource["id"]; len(newResource) == 0 || (len(newResource) == 1 && ok) {
			web.Error(w, http.StatusBadRequest, "bad request")
			return
		}

		// Update the resource.
		data, err := storageSvc.Update(id, newResource)
		if err != nil {
			web.Error(w, http.StatusInternalServerError, "internal Server Error")
			return
		}

		web.Success(w, http.StatusOK, data)
	}
}
