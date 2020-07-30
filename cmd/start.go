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

var (
	errFailedParseFlag     = errors.New("failed to parse flag")
	errFailedParseFile     = errors.New("failed to parse file")
	errFileNotFound        = errors.New("unable to find requested file")
	errUnsupportedResource = errors.New("only array type resources are supported")
	errFailedStartServer   = errors.New("failed to start JSON server. Maybe port already in use")
	errFailedInitResources = errors.New("failed to initialize resources")
)

func newStartCmd() *cobra.Command {
	// startCmd represents the start command.
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the REST API server based on flags",
		Long: `
For every provided resource, specific http endpoints are created. 
Those include a GET, GET by ID, POST, PUT by ID, PATCH by ID and DELETE by ID. 
Also a '/db' endpoint is available that returns all the data. 
Please note that only array data type resources are supported`,
		RunE: runStart,
	}

	// Optional flag to set the server port.
	startCmd.Flags().StringP("port", "p", "3000", "Port the server will listen to")
	// Optional flag to set the watch file.
	startCmd.Flags().StringP("file", "f", "db.json", "File to watch")
	// Optional flag to enable logs.
	startCmd.Flags().BoolP("logs", "l", false, "Enable logs")

	return startCmd
}

func runStart(cmd *cobra.Command, _ []string) error {
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

func createResourceStorage(resourceKeys []string, filename string) (map[string]storage.Storage, error) {
	resourceStorage := make(map[string]storage.Storage)

	for _, resourceKey := range resourceKeys {
		storageSvc, err := storage.NewFile(filename, resourceKey)
		if err != nil {
			return nil, errFailedInitResources
		}

		resourceStorage[resourceKey] = storageSvc
	}

	// Create storage service for common db endpoint.
	storageSvcDB, err := storage.NewFile(filename, "")
	if err != nil {
		return nil, errFailedInitResources
	}

	resourceStorage["db"] = storageSvcDB

	return resourceStorage, nil
}

func displayInfo(resourceKeys []string, port string) {
	fmt.Printf("JSON Server successfully running\n\n")

	fmt.Println("Resources")
	for _, resource := range resourceKeys {
		fmt.Printf("http://localhost:%s/%s\n", port, resource)
	}

	fmt.Printf("http://localhost:%s/db\n\n", port)

	fmt.Println("Home")
	fmt.Printf("http://localhost:%s\n\n", port)
}
