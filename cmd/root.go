// Package cmd contains the functionality for the set of commands
// currently supported by the CLI.
package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"time"

	"github.com/spf13/cobra"

	"github.com/chanioxaris/json-server/internal/handler"
	"github.com/chanioxaris/json-server/internal/logger"
	"github.com/chanioxaris/json-server/internal/storage"
)

// rootCmd represents the base command when called without any sub commands.
var rootCmd = &cobra.Command{
	Use:   "json-server",
	Short: "Create a dummy REST API from a json file with zero coding within seconds",
	Long: `json-server is a cross-platform CLI tool to create within seconds a dummy REST API from a provided json 
			file. Depending on the provided data, http endpoints are created which include GET, GET by ID, POST, 
			PUT by ID, PATCH by ID and DELETE by ID. Only array type data is supported`,
	RunE: run,
}

var (
	errFailedParseFlag     = errors.New("failed to parse flag")
	errFailedParseFile     = errors.New("failed to parse file")
	errFileNotFound        = errors.New("unable to find requested file")
	errUnsupportedResource = errors.New("only array type resources are supported")
	errFailedStartServer   = errors.New("failed to start JSON server. Maybe port already in use")
	errFailedInitResources = errors.New("failed to initialize resources")
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
	logger.Setup(logs)

	// Get resource keys.
	resourceKeys, err := getResourceKeys(file)
	if err != nil {
		return err
	}

	// Create storage service for each resource.
	resourceStorage, err := createResourceStorage(resourceKeys, file)
	if err != nil {
		return err
	}

	// Setup API server.
	api := &http.Server{
		Addr:    ":" + port,
		Handler: handler.Setup(resourceStorage),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	// Start REST API server.
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return errFailedStartServer
	}

	// nolint
	go api.Serve(listener)

	// Display info about available resources and home page.
	displayInfo(resourceKeys, port)

	gracefulShutdown(api)

	return nil
}

// gracefulShutdown handles any signal that interrupts the running server
func gracefulShutdown(server *http.Server) {
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("failed to gracefully shutdown server")
		return
	}

	fmt.Println("gracefully shutting down server")
}

func getResourceKeys(filename string) ([]string, error) {
	// Read file contents used as storage.
	contentBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errFileNotFound, filename)
	}

	content := map[string]interface{}{}
	if err = json.Unmarshal(contentBytes, &content); err != nil {
		return nil, fmt.Errorf("%w: %s", errFailedParseFile, filename)
	}

	resourceKeys := make([]string, 0)

	// Range on content to retrieve resource keys.
	for resource, data := range content {
		switch reflect.TypeOf(data).Kind() {
		case reflect.Slice:
			resourceKeys = append(resourceKeys, resource)
		default:
			return nil, errUnsupportedResource
		}
	}

	return resourceKeys, nil
}

func createResourceStorage(resourceKeys []string, file string) (map[string]storage.Service, error) {
	resourceStorage := make(map[string]storage.Service)

	for _, resourceKey := range resourceKeys {
		storageSvc, err := storage.New(file, resourceKey)
		if err != nil {
			return nil, errFailedInitResources
		}

		resourceStorage[resourceKey] = storageSvc
	}

	// Create storage service for common db endpoint.
	storageSvcDB, err := storage.New(file, "")
	if err != nil {
		return nil, errFailedInitResources
	}

	resourceStorage["db"] = storageSvcDB

	return resourceStorage, nil
}

func displayInfo(resourceKeys []string, port string) {
	fmt.Println("JSON Server successfully running")
	fmt.Println()

	fmt.Println("Resources")
	for _, resource := range resourceKeys {
		fmt.Printf("http://localhost:%s/%s\n", port, resource)
	}

	fmt.Printf("http://localhost:%s/db\n", port)
	fmt.Println()

	fmt.Println("Home")
	fmt.Printf("http://localhost:%s\n", port)
	fmt.Println()
}
