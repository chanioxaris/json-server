package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/chanioxaris/json-server/storage"
	"github.com/chanioxaris/json-server/web"
)

// Create operates as a http handle, to add a new resource.
func Create(storageSvc storage.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read and decode request body.
		var newResource storage.Resource
		if err := json.NewDecoder(r.Body).Decode(&newResource); err != nil {
			web.Error(w, http.StatusBadRequest, storage.ErrBadRequest.Error())
			return
		}

		// Check if request body is empty, or contains only id.
		if _, ok := newResource["id"]; len(newResource) == 0 || (len(newResource) == 1 && ok) {
			web.Error(w, http.StatusBadRequest, storage.ErrBadRequest.Error())
			return
		}

		// Create the new resource.
		data, err := storageSvc.Create(newResource)
		if err != nil {
			// Already exists with the requested id.
			if errors.Is(err, storage.ErrResourceAlreadyExists) {
				web.Error(w, http.StatusConflict, err.Error())
				return
			}

			web.Error(w, http.StatusInternalServerError, storage.ErrInternalServerError.Error())
			return
		}

		web.Success(w, http.StatusCreated, data)
	}
}
