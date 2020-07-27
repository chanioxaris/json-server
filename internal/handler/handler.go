// Package handler contains the full set of handler functions and routes
// supported by the web api.
package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/chanioxaris/json-server/internal/handler/common"
	"github.com/chanioxaris/json-server/internal/storage"
	"github.com/chanioxaris/json-server/internal/web/middleware"
)

// Setup API handler based on provided resources.
func Setup(resourceStorage map[string]storage.Service) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(middleware.Recovery)
	router.Use(middleware.Logger)

	// For each resource create the appropriate endpoint handlers.
	for resourceKey, storageSvc := range resourceStorage {
		// Common endpoint to retrieve db contents.
		if resourceKey == "db" {
			router.HandleFunc("/db", common.DB(storageSvc)).Methods(http.MethodGet)
			continue
		}

		// Register all default endpoint handlers for resource.
		router.HandleFunc(fmt.Sprintf("/%s", resourceKey), List(storageSvc)).Methods(http.MethodGet)
		router.HandleFunc(fmt.Sprintf("/%s/{id}", resourceKey), Read(storageSvc)).Methods(http.MethodGet)
		router.HandleFunc(fmt.Sprintf("/%s", resourceKey), Create(storageSvc)).Methods(http.MethodPost)
		router.HandleFunc(fmt.Sprintf("/%s/{id}", resourceKey), Replace(storageSvc)).Methods(http.MethodPut)
		router.HandleFunc(fmt.Sprintf("/%s/{id}", resourceKey), Update(storageSvc)).Methods(http.MethodPatch)
		router.HandleFunc(fmt.Sprintf("/%s/{id}", resourceKey), Delete(storageSvc)).Methods(http.MethodDelete)
	}

	// Render a home page with useful info.
	router.HandleFunc("/", common.HomePage(resourceStorage)).Methods(http.MethodGet)

	return router
}
