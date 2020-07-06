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
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/chanioxaris/json-server/handler"
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

func init() {
	// Optional flag to set the server port.
	rootCmd.Flags().StringP("port", "p", "3000", "Port the server will listen to")
	// Optional flag to set the watch file.
	rootCmd.Flags().StringP("file", "f", "db.json", "File to watch")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, _ []string) error {
	// Parse command's flags.
	port, err := cmd.Flags().GetString("port")
	if err != nil {
		return err
	}

	file, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	rand.Seed(time.Now().UnixNano())

	// Read file contents used as 'database'.
	contentBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	content := map[string]interface{}{}
	if err = json.Unmarshal(contentBytes, &content); err != nil {
		return err
	}

	// Setup router.
	router, err := SetupRouter(content, file)
	if err != nil {
		return err
	}

	fmt.Printf("Server listening on http://localhost:%s\n\n", port)

	fmt.Println("Available endpoints:")
	for key := range content {
		fmt.Printf("\033[32m/%s \033[0m\n", key)
	}

	fmt.Println(http.ListenAndServe(":"+port, router))

	return nil
}

func SetupRouter(content map[string]interface{}, file string) (http.Handler, error) {
	router := mux.NewRouter().StrictSlash(true)

	// For each resource create the appropriate endpoint handlers.
	for key, val := range content {
		switch reflect.TypeOf(val).Kind() {
		// If there is an array, register all default endpoint handlers.
		case reflect.Slice:
			// Create storage service to access the 'database' for specific resource.
			storageSvc, err := storage.NewStorage(file, key, false)
			if err != nil {
				return nil, err
			}

			router.HandleFunc(fmt.Sprintf("/%s", key), handler.List(storageSvc)).Methods(http.MethodGet)
			router.HandleFunc(fmt.Sprintf("/%s/{id}", key), handler.Read(storageSvc)).Methods(http.MethodGet)
			router.HandleFunc(fmt.Sprintf("/%s", key), handler.Create(storageSvc)).Methods(http.MethodPost)
			router.HandleFunc(fmt.Sprintf("/%s/{id}", key), handler.Replace(storageSvc)).Methods(http.MethodPut)
			router.HandleFunc(fmt.Sprintf("/%s/{id}", key), handler.Update(storageSvc)).Methods(http.MethodPatch)
			router.HandleFunc(fmt.Sprintf("/%s/{id}", key), handler.Delete(storageSvc)).Methods(http.MethodDelete)
		// Otherwise register only a read endpoint handler.
		default:
			// Create storage service to access the 'database' for specific resource.
			storageSvc, err := storage.NewStorage(file, key, true)
			if err != nil {
				return nil, err
			}

			router.HandleFunc(fmt.Sprintf("/%s", key), handler.Read(storageSvc)).Methods(http.MethodGet)
		}
	}

	return router, nil
}
