package handler_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/chanioxaris/json-server/cmd"
	"github.com/chanioxaris/json-server/storage"
)

var (
	mockServer *httptest.Server

	pluralKeys   = []string{"plural_key_1", "plural_key_2"}
	singularKeys = []string{"singular_key"}

	testData = make(map[string]interface{}, 0)

	fileName string
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())

	var err error
	fileName, err = testGenerateJSONFile()
	if err != nil {
		panic(err)
	}
	defer os.Remove(fileName)

	router, err := cmd.SetupRouter(testData, fileName)
	if err != nil {
		panic(err)
	}

	mockServer = httptest.NewServer(router)
	defer mockServer.Close()

	m.Run()
}

func testGenerateJSONFile() (string, error) {
	f, err := ioutil.TempFile(".", "")
	if err != nil {
		return "", err
	}

	contentBytes, err := testGenerateData()
	if err != nil {
		return "", err
	}

	if err = ioutil.WriteFile(f.Name(), contentBytes, 0644); err != nil {
		return "", err
	}

	return f.Name(), nil
}

func testGenerateData() ([]byte, error) {
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
	}

	for _, key := range singularKeys {
		testData[key] = rand.Intn(1000)
	}

	return json.MarshalIndent(testData, "", "  ")
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
