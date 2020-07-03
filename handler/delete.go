package handler

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/chanioxaris/json-server/storage"
	"github.com/chanioxaris/json-server/web"
)

// Delete operates as a http handle, to delete an existing resource.
func Delete(storageSvc storage.Service) http.HandlerFunc {
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

		// Delete resource.
		if err := storageSvc.Delete(id); err != nil {
			web.Error(w, http.StatusInternalServerError, "internal Server Error")
			return
		}

		web.Success(w, http.StatusOK, nil)
	}
}
