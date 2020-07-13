/*
Copyright © 2020 Haris Chaniotakis

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/chanioxaris/json-server/handler"
	"github.com/chanioxaris/json-server/handler/common"
	"github.com/chanioxaris/json-server/logger"
	"github.com/chanioxaris/json-server/middleware"
	"github.com/chanioxaris/json-server/storage"
)

// rootCmd represents the base command when called without any sub commands.
var rootCmd = &cobra.Command{
	Use:   "json-server",
	Short: "Create a dummy REST API from a json file with zero coding within seconds",
	Long: `json-server is a cross-platform CLI tool to create within seconds a dummy REST API from a provided json 
			file. Depending on the provided file some default http endpoints are created. For array data (plural) a GET, 
			GET by ID, POST, PUT by ID, PATCH by ID and DELETE by ID endpoints are available. For object data (singular) 
			a GET endpoint is available`,
	RunE: run,
}

var (
	errFailedParseFlag     = errors.New("failed to parse flag")
	errFailedParseFile     = errors.New("failed to parse file")
	errFailedInitResources = errors.New("failed to initialize resources")
	errFileNotFound        = errors.New("unable to find requested file")
)

func init() {
	// Optional flag to set the server port.
	rootCmd.Flags().StringP("port", "p", "3000", "Port the server will listen to")
	// Optional flag to set the watch file.
	rootCmd.Flags().StringP("file", "f", "db.json", "File to watch")
	// Optional flag to enable logs.
	rootCmd.Flags().BoolP("logs", "l", false, "Enable logs")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, _ []string) error {
	rand.Seed(time.Now().UnixNano())

	// Parse command's flags.
	port, err := cmd.Flags().GetString("port")
	if err != nil {
		return fmt.Errorf("%w: port", errFailedParseFlag)
	}

	file, err := cmd.Flags().GetString("file")
	if err != nil {
		return fmt.Errorf("%w: file", errFailedParseFlag)
	}

	logs, err := cmd.Flags().GetBool("logs")
	if err != nil {
		return fmt.Errorf("%w: logs", errFailedParseFlag)
	}

	// Setup logger.
	setupLogger(logs)

	// Get storage resources.
	storageResources, err := getStorageResources(file)
	if err != nil {
		return err
	}

	// Setup router.
	router, err := SetupRouter(storageResources, file)
	if err != nil {
		return err
	}

	// Preview info about available resources and home page.
	displayInfo(storageResources, port)

	fmt.Println(http.ListenAndServe(":"+port, router))

	return nil
}

// SetupRouter based on provided resources.
func SetupRouter(storageResources map[string]bool, file string) (http.Handler, error) {
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
			router.HandleFunc(fmt.Sprintf("/%s", resource), handler.List(storageSvc)).Methods(http.MethodGet)
			router.HandleFunc(fmt.Sprintf("/%s/{id}", resource), handler.Read(storageSvc)).Methods(http.MethodGet)
			router.HandleFunc(fmt.Sprintf("/%s", resource), handler.Create(storageSvc)).Methods(http.MethodPost)
			router.HandleFunc(fmt.Sprintf("/%s/{id}", resource), handler.Replace(storageSvc)).Methods(http.MethodPut)
			router.HandleFunc(fmt.Sprintf("/%s/{id}", resource), handler.Update(storageSvc)).Methods(http.MethodPatch)
			router.HandleFunc(fmt.Sprintf("/%s/{id}", resource), handler.Delete(storageSvc)).Methods(http.MethodDelete)
			// Register default endpoint handler for singular resource.
		default:
			router.HandleFunc(fmt.Sprintf("/%s", resource), handler.Read(storageSvc)).Methods(http.MethodGet)
		}
	}

	// Default endpoint to list all resources.
	storageSvc, err := storage.NewStorage(file, "", false)
	if err != nil {
		return nil, errFailedInitResources
	}

	router.HandleFunc("/db", common.DB(storageSvc)).Methods(http.MethodGet)

	// Render a home page with useful info.
	router.HandleFunc("/", common.HomePage(storageResources)).Methods(http.MethodGet)

	return router, nil
}

func setupLogger(show bool) {
	logrus.SetFormatter(&logger.CustomFormatter{})

	if !show {
		logrus.SetOutput(ioutil.Discard)
	}
}

func getStorageResources(filename string) (map[string]bool, error) {
	// Read file contents used as storage.
	contentBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errFileNotFound, filename)
	}

	content := map[string]interface{}{}
	if err = json.Unmarshal(contentBytes, &content); err != nil {
		return nil, fmt.Errorf("%w: %s", errFailedParseFile, filename)
	}

	storageKeys := make(map[string]bool)

	// Range on content to retrieve resource keys and type (plural, singular).
	for resource, data := range content {
		switch reflect.TypeOf(data).Kind() {
		case reflect.Slice:
			storageKeys[resource] = false
		default:
			storageKeys[resource] = true
		}
	}

	return storageKeys, nil
}

func displayInfo(storageResources map[string]bool, port string) {
	fmt.Println("JSON Server successfully running")
	fmt.Println()

	fmt.Println("Resources")
	for resource := range storageResources {
		fmt.Printf("http://localhost:%s/%s\n", port, resource)
	}

	fmt.Printf("http://localhost:%s/db\n", port)
	fmt.Println()

	fmt.Println("Home")
	fmt.Printf("http://localhost:%s\n", port)
	fmt.Println()
}
