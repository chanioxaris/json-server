package common_test

import (
	"errors"
	"fmt"
	"math/rand"
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

	testResourceKeys    = []string{"resource_key_1", "resource_key_2"}
	testData            = make(storage.Database)
	testResourceStorage = make(map[string]*storage.Mock)
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())

	if err := testGenerateData(); err != nil {
		panic(err)
	}

	resourceStorage, err := testCreateResourceStorage()
	if err != nil {
		panic(err)
	}

	router := handler.Setup(resourceStorage)

	mockServer = httptest.NewServer(router)
	defer mockServer.Close()

	os.Exit(m.Run())
}

func testGenerateData() error {
	for _, key := range testResourceKeys {
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
	}

	return nil
}

func testCreateResourceStorage() (map[string]storage.Storage, error) {
	resourceStorage := make(map[string]storage.Storage)

	storageSvcDB, err := storage.NewMock(testData, "")
	if err != nil {
		return nil, errors.New("failed to initialize resources")
	}

	resourceStorage["db"] = storageSvcDB

	for key, storageSvc := range resourceStorage {
		testResourceStorage[key] = storageSvc.(*storage.Mock)
	}

	return resourceStorage, nil
}

func testResetData(key string) {
	testResourceStorage[key].SetData(testData)
}
