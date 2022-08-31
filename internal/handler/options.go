package handler

import (
	"github.com/chanioxaris/json-server/internal/storage"
	"net/http"
)

func Options(storageSvc storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
	}
}
