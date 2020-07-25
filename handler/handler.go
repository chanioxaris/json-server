// Package handler contains the full set of handler functions and routes
// supported by the web api.
package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/chanioxaris/json-server/handler/common"
	"github.com/chanioxaris/json-server/middleware"
	"github.com/chanioxaris/json-server/storage"
)

var (
	errFailedInitResources = errors.New("failed to initialize resources")
)

// Setup API handler based on provided resources.
func Setup(resourceKeys []string, file string) (http.Handler, error) {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(middleware.Recovery)
	router.Use(middleware.Logger)

	// For each resource create the appropriate endpoint handlers.
	for _, resourceKey := range resourceKeys {
		// Create storage service to access the 'database' for specific resource.
		storageSvc, err := storage.NewStorage(file, resourceKey)
		if err != nil {
			return nil, errFailedInitResources
		}

		// Register all default endpoint handlers for resource.
		router.HandleFunc(fmt.Sprintf("/%s", resourceKey), List(storageSvc)).Methods(http.MethodGet)
		router.HandleFunc(fmt.Sprintf("/%s/{id}", resourceKey), Read(storageSvc)).Methods(http.MethodGet)
		router.HandleFunc(fmt.Sprintf("/%s", resourceKey), Create(storageSvc)).Methods(http.MethodPost)
		router.HandleFunc(fmt.Sprintf("/%s/{id}", resourceKey), Replace(storageSvc)).Methods(http.MethodPut)
		router.HandleFunc(fmt.Sprintf("/%s/{id}", resourceKey), Update(storageSvc)).Methods(http.MethodPatch)
		router.HandleFunc(fmt.Sprintf("/%s/{id}", resourceKey), Delete(storageSvc)).Methods(http.MethodDelete)
	}

	// Default endpoint to retrieve db contents.
	storageSvc, err := storage.NewStorage(file, "")
	if err != nil {
		return nil, errFailedInitResources
	}

	router.HandleFunc("/db", common.DB(storageSvc)).Methods(http.MethodGet)

	// Render a home page with useful info.
	router.HandleFunc("/", common.HomePage(resourceKeys)).Methods(http.MethodGet)

	return router, nil
}
