// Package common contains the set of handler functions and routes which by
// default supported from the web api.
package common

import (
	"net/http"

	"github.com/chanioxaris/json-server/storage"
	"github.com/chanioxaris/json-server/web"
)

// DB operates as a http handler, to list db content.
func DB(storageSvc storage.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := storageSvc.DB()
		if err != nil {
			web.Error(w, http.StatusInternalServerError, storage.ErrInternalServerError.Error())
			return
		}

		web.Success(w, http.StatusOK, data)
	}
}
