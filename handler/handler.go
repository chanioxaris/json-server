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
func Setup(storageResources map[string]bool, file string) (http.Handler, error) {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(middleware.Recovery)
	router.Use(middleware.Logger)

	// For each resource create the appropriate endpoint handlers.
	for resource, singular := range storageResources {
		// Create storage service to access the 'database' for specific resource.
		storageSvc, err := storage.NewStorage(file, resource, singular)
		if err != nil {
			return nil, errFailedInitResources
		}

		switch singular {
		// Register all default endpoint handlers for plural resource.
		case false:
			router.HandleFunc(fmt.Sprintf("/%s", resource), List(storageSvc)).Methods(http.MethodGet)
			router.HandleFunc(fmt.Sprintf("/%s/{id}", resource), Read(storageSvc)).Methods(http.MethodGet)
			router.HandleFunc(fmt.Sprintf("/%s", resource), Create(storageSvc)).Methods(http.MethodPost)
			router.HandleFunc(fmt.Sprintf("/%s/{id}", resource), Replace(storageSvc)).Methods(http.MethodPut)
			router.HandleFunc(fmt.Sprintf("/%s/{id}", resource), Update(storageSvc)).Methods(http.MethodPatch)
			router.HandleFunc(fmt.Sprintf("/%s/{id}", resource), Delete(storageSvc)).Methods(http.MethodDelete)
			// Register default endpoint handler for singular resource.
		default:
			router.HandleFunc(fmt.Sprintf("/%s", resource), Read(storageSvc)).Methods(http.MethodGet)
		}
	}

	// Default endpoint to retrieve db contents.
	storageSvc, err := storage.NewStorage(file, "", false)
	if err != nil {
		return nil, errFailedInitResources
	}

	router.HandleFunc("/db", common.DB(storageSvc)).Methods(http.MethodGet)

	// Render a home page with useful info.
	router.HandleFunc("/", common.HomePage(storageResources)).Methods(http.MethodGet)

	return router, nil
}
