package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/chanioxaris/json-server/internal/handler"
	"github.com/chanioxaris/json-server/internal/storage"
)

var (
	mockServer *httptest.Server

	pluralKeys = []string{"plural_key_1", "plural_key_2"}

	testData = make(map[string]interface{})

	fileName string
)

func TestMain(m *testing.M) {
	// testMain wrapper is needed to support defers and panics.
	// os.Exit will ignore those and exit silently.
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	rand.Seed(time.Now().UnixNano())

	resourceKeys, err := testGenerateJSONFile()
	if err != nil {
		panic(err)
	}
	defer os.Remove(fileName)

	resourceStorage, err := testCreateResourceStorage(resourceKeys, fileName)
	if err != nil {
		panic(err)
	}

	router := handler.Setup(resourceStorage)

	mockServer = httptest.NewServer(router)
	defer mockServer.Close()

	return m.Run()
}

func testGenerateJSONFile() ([]string, error) {
	f, err := ioutil.TempFile(".", "")
	if err != nil {
		return nil, err
	}

	fileName = f.Name()

	contentBytes, resourceKeys, err := testGenerateData()
	if err != nil {
		return nil, err
	}

	if err = ioutil.WriteFile(f.Name(), contentBytes, 0644); err != nil {
		return nil, err
	}

	return resourceKeys, nil
}

func testGenerateData() ([]byte, []string, error) {
	resourceKeys := make([]string, 0)

	for _, key := range pluralKeys {
		resources := make([]storage.Resource, 0)
		for idx := 0; idx < rand.Intn(10)+1; idx++ {
			newResource := storage.Resource{
				"id":      strconv.Itoa(idx),
				"field_1": fmt.Sprintf("field_1-%s-%d", key, idx),
				"field_2": fmt.Sprintf("field_2-%s-%d", key, idx),
			}

			resources = append(resources, newResource)
		}

		testData[key] = resources
		resourceKeys = append(resourceKeys, key)
	}

	contentBytes, err := json.MarshalIndent(testData, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return contentBytes, resourceKeys, nil
}

func testCreateResourceStorage(resourceKeys []string, file string) (map[string]storage.Service, error) {
	resourceStorage := make(map[string]storage.Service)

	for _, resourceKey := range resourceKeys {
		storageSvc, err := storage.New(file, resourceKey)
		if err != nil {
			return nil, errors.New("failed to initialize resources")
		}

		resourceStorage[resourceKey] = storageSvc
	}

	storageSvcDB, err := storage.New(file, "")
	if err != nil {
		return nil, errors.New("failed to initialize resources")
	}

	resourceStorage["db"] = storageSvcDB

	return resourceStorage, nil
}

func testResetData(filename string) error {
	contentBytes, err := json.Marshal(testData)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filename, contentBytes, 0644); err != nil {
		return err
	}

	return nil
}

func testListResourcesByKey(key string) ([]storage.Resource, error) {
	url := fmt.Sprintf("%s/%s", mockServer.URL, key)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var body []storage.Resource
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}
