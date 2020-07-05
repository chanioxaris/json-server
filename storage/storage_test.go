package storage_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/chanioxaris/json-server/storage"
)

var (
	keys     = []string{"key1", "key2"}
	testData = make(map[string][]storage.Resource, 0)
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())

	os.Exit(m.Run())
}

func TestFind(t *testing.T) {
	f, err := testGenerateStorageFile()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	randomKeyIndex := rand.Intn(len(keys))
	randomKey := keys[randomKeyIndex]

	type args struct {
		key      string
		singular bool
		filename string
	}
	testCases := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "List all resources of specific key",
			args: args{
				key:      randomKey,
				singular: false,
				filename: f.Name(),
			},
		},
		{
			name: "List all resources of invalid key",
			args: args{
				key:      "randomKey",
				singular: false,
				filename: f.Name(),
			},
			wantErr: true,
			err:     storage.ErrResourceNotFound,
		},
		{
			name: "List all resources of invalid file name",
			args: args{
				key:      randomKey,
				singular: false,
				filename: "randomFileName",
			},
			wantErr: true,
			err:     os.ErrNotExist,
		},
	}

	for _, tt := range testCases {
		storageSvc, err := storage.NewStorage(tt.args.filename, tt.args.key, tt.args.singular)
		if err != nil {
			t.Fatal(err)
		}

		got, err := storageSvc.Find()
		if err != nil && !tt.wantErr {
			t.Fatal(err)
		}

		if !tt.wantErr {
			if len(got) != len(testData[tt.args.key]) {
				t.Fatalf("expected data length %v, but got %v", len(testData[tt.args.key]), len(got))
			}

			if !reflect.DeepEqual(got, testData[tt.args.key]) {
				t.Fatalf("expected data %v, but got %v", testData[tt.args.key], got)
			}
		} else {
			if err == nil || !errors.Is(err, tt.err) {
				t.Fatalf("expected error %v, but got %v", tt.err, err)
			}
		}
	}
}

func TestFindById(t *testing.T) {
	f, err := testGenerateStorageFile()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	randomKeyIndex := rand.Intn(len(keys))
	randomKey := keys[randomKeyIndex]

	randomResourceIndex := rand.Intn(len(testData[randomKey]))
	randomResource := testData[randomKey][randomResourceIndex]

	type args struct {
		key      string
		singular bool
		filename string
	}
	testCases := []struct {
		name         string
		args         args
		id           string
		expectedData interface{}
		wantErr      bool
		err          error
	}{
		{
			name: "Read resource with id",
			args: args{
				key:      randomKey,
				singular: false,
				filename: f.Name(),
			},
			id:           randomResource["id"].(string),
			expectedData: randomResource,
		},
		{
			name: "Read resource of invalid id",
			args: args{
				key:      randomKey,
				singular: false,
				filename: f.Name(),
			},
			id:      "randomId",
			wantErr: true,
			err:     storage.ErrResourceNotFound,
		},
		{
			name: "Read resource of invalid resource key",
			args: args{
				key:      "randomKey",
				singular: false,
				filename: f.Name(),
			},
			wantErr: true,
			err:     storage.ErrResourceNotFound,
		},
		{
			name: "Read resource of invalid file name",
			args: args{
				key:      randomKey,
				singular: false,
				filename: "randomFileName",
			},
			wantErr: true,
			err:     os.ErrNotExist,
		},
	}

	for _, tt := range testCases {
		storageSvc, err := storage.NewStorage(tt.args.filename, tt.args.key, tt.args.singular)
		if err != nil {
			t.Fatal(err)
		}

		got, err := storageSvc.FindById(tt.id)
		if err != nil && !tt.wantErr {
			t.Fatal(err)
		}

		if !tt.wantErr {
			if !reflect.DeepEqual(got, tt.expectedData) {
				t.Fatalf("expected data %v, but got %v", tt.expectedData, got)
			}
		} else {
			if err == nil || !errors.Is(err, tt.err) {
				t.Fatalf("expected error %v, but got %v", tt.err, err)
			}
		}
	}
}

func TestCreate(t *testing.T) {
	f, err := testGenerateStorageFile()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	randomKeyIndex := rand.Intn(len(keys))
	randomKey := keys[randomKeyIndex]

	randomResourceIndex := rand.Intn(len(testData[randomKey]))
	randomResource := testData[randomKey][randomResourceIndex]

	type args struct {
		key      string
		singular bool
		filename string
	}
	testCases := []struct {
		name     string
		args     args
		resource storage.Resource
		wantErr  bool
		err      error
	}{
		{
			name: "Create resource with id provided",
			args: args{
				key:      randomKey,
				singular: false,
				filename: f.Name(),
			},
			resource: storage.Resource{
				"id":          "2020",
				"description": "test-description",
			},
		},
		{
			name: "Create resource without id provided",
			args: args{
				key:      randomKey,
				singular: false,
				filename: f.Name(),
			},
			resource: storage.Resource{
				"description": "new-description",
			},
		},
		{
			name: "Create invalid resource with existing id",
			args: args{
				key:      randomKey,
				singular: false,
				filename: f.Name(),
			},
			resource: storage.Resource{
				"id":          randomResource["id"].(string),
				"description": "new-description",
			},
			wantErr: true,
			err:     storage.ErrResourceAlreadyExists,
		},
		{
			name: "Create resource of invalid resource key",
			args: args{
				key:      "randomKey",
				singular: false,
				filename: f.Name(),
			},
			wantErr: true,
			err:     storage.ErrResourceNotFound,
		},
		{
			name: "Create resource of invalid file name",
			args: args{
				key:      randomKey,
				singular: false,
				filename: "randomFileName",
			},
			wantErr: true,
			err:     os.ErrNotExist,
		},
	}

	for _, tt := range testCases {
		if err := testResetData(f.Name()); err != nil {
			t.Fatal(err)
		}

		storageSvc, err := storage.NewStorage(tt.args.filename, tt.args.key, tt.args.singular)
		if err != nil {
			t.Fatal(err)
		}

		got, err := storageSvc.Create(tt.resource)
		if err != nil && !tt.wantErr {
			t.Fatal(err)
		}

		if !tt.wantErr {
			if !reflect.DeepEqual(got, tt.resource) {
				t.Fatalf("expected created %v, but got %v", tt.resource, got)
			}

			currData, err := storageSvc.Find()
			if err != nil {
				t.Fatal(err)
			}

			expectedTestData := append(testData[randomKey], tt.resource)

			if len(currData) != len(expectedTestData) {
				t.Fatalf("expected data length %v, but got %v", len(expectedTestData), len(currData))
			}

			if !reflect.DeepEqual(currData, expectedTestData) {
				t.Fatalf("expected data %v, but got %v", expectedTestData, currData)
			}
		} else {
			if err == nil || !errors.Is(err, tt.err) {
				t.Fatalf("expected error %v, but got %v", tt.err, err)
			}
		}
	}
}

func TestUpdate(t *testing.T) {
	f, err := testGenerateStorageFile()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	randomKeyIndex := rand.Intn(len(keys))
	randomKey := keys[randomKeyIndex]

	randomResourceIndex := rand.Intn(len(testData[randomKey]))
	randomResource := testData[randomKey][randomResourceIndex]

	type args struct {
		key      string
		singular bool
		filename string
	}
	testCases := []struct {
		name     string
		args     args
		id       string
		resource storage.Resource
		wantErr  bool
		err      error
	}{
		{
			name: "Update resource without id provided",
			args: args{
				key:      randomKey,
				singular: false,
				filename: f.Name(),
			},
			resource: storage.Resource{
				"description": "updated-description",
			},
			id: randomResource["id"].(string),
		},
		{
			name: "Update resource with id provided",
			args: args{
				key:      randomKey,
				singular: false,
				filename: f.Name(),
			},
			resource: storage.Resource{
				"id":          "2020",
				"description": "updated-description",
			},
			id: randomResource["id"].(string),
		},
		{
			name: "Update invalid resource with non existing id",
			args: args{
				key:      randomKey,
				singular: false,
				filename: f.Name(),
			},
			resource: storage.Resource{
				"description": "updated-description",
			},
			id:      "2020",
			wantErr: true,
			err:     storage.ErrResourceNotFound,
		},
		{
			name: "Update resource of invalid resource key",
			args: args{
				key:      "randomKey",
				singular: false,
				filename: f.Name(),
			},
			wantErr: true,
			err:     storage.ErrResourceNotFound,
		},
		{
			name: "Update resource of invalid file name",
			args: args{
				key:      randomKey,
				singular: false,
				filename: "randomFileName",
			},
			wantErr: true,
			err:     os.ErrNotExist,
		},
	}

	for _, tt := range testCases {
		if err := testResetData(f.Name()); err != nil {
			t.Fatal(err)
		}

		storageSvc, err := storage.NewStorage(tt.args.filename, tt.args.key, tt.args.singular)
		if err != nil {
			t.Fatal(err)
		}

		got, err := storageSvc.Update(tt.id, tt.resource)
		if err != nil && !tt.wantErr {
			t.Fatal(err)
		}

		if !tt.wantErr {
			if !reflect.DeepEqual(got, tt.resource) {
				t.Fatalf("expected updated %v, but got %v", tt.resource, got)
			}

			currData, err := storageSvc.Find()
			if err != nil {
				t.Fatal(err)
			}

			expectedTestData := testData[randomKey]
			expectedTestData[randomResourceIndex] = tt.resource

			if len(currData) != len(expectedTestData) {
				t.Fatalf("expected data length %v, but got %v", len(expectedTestData), len(currData))
			}

			if !reflect.DeepEqual(currData, expectedTestData) {
				t.Fatalf("expected data %v, but got %v", expectedTestData, currData)
			}
		} else {
			if err == nil || !errors.Is(err, tt.err) {
				t.Fatalf("expected error %v, but got %v", tt.err, err)
			}
		}
	}
}

func TestDelete(t *testing.T) {
	f, err := testGenerateStorageFile()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	randomKeyIndex := rand.Intn(len(keys))
	randomKey := keys[randomKeyIndex]

	randomResourceIndex := rand.Intn(len(testData[randomKey]))
	randomResource := testData[randomKey][randomResourceIndex]

	type args struct {
		key      string
		singular bool
		filename string
	}
	testCases := []struct {
		name    string
		args    args
		id      string
		wantErr bool
		err     error
	}{
		{
			name: "Delete resource with id",
			args: args{
				key:      randomKey,
				singular: false,
				filename: f.Name(),
			},
			id: randomResource["id"].(string),
		},
		{
			name: "Delete resource of invalid id",
			args: args{
				key:      randomKey,
				singular: false,
				filename: f.Name(),
			},
			id:      "randomId",
			wantErr: true,
			err:     storage.ErrResourceNotFound,
		},
		{
			name: "Delete resource of invalid resource key",
			args: args{
				key:      "randomKey",
				singular: false,
				filename: f.Name(),
			},
			wantErr: true,
			err:     storage.ErrResourceNotFound,
		},
		{
			name: "Delete resource of invalid file name",
			args: args{
				key:      randomKey,
				singular: false,
				filename: "randomFileName",
			},
			wantErr: true,
			err:     os.ErrNotExist,
		},
	}

	for _, tt := range testCases {
		storageSvc, err := storage.NewStorage(tt.args.filename, tt.args.key, tt.args.singular)
		if err != nil {
			t.Fatal(err)
		}

		err = storageSvc.Delete(tt.id)
		if err != nil && !tt.wantErr {
			t.Fatal(err)
		}

		if !tt.wantErr {
			currData, err := storageSvc.Find()
			if err != nil {
				t.Fatal(err)
			}

			testKeyData := testData[randomKey]
			expectedTestData := append(testKeyData[:randomResourceIndex], testKeyData[randomResourceIndex+1:]...)

			if len(currData) != len(expectedTestData) {
				t.Fatalf("expected data length %v, but got %v", len(expectedTestData), len(currData))
			}

			if !reflect.DeepEqual(currData, expectedTestData) {
				t.Fatalf("expected data %v, but got %v", expectedTestData, currData)
			}
		} else {
			if err == nil || !errors.Is(err, tt.err) {
				t.Fatalf("expected error %v, but got %v", tt.err, err)
			}
		}
	}
}

func testGenerateStorageFile() (*os.File, error) {
	f, err := ioutil.TempFile(".", "")
	if err != nil {
		return nil, err
	}

	contentBytes, err := testGenerateData()
	if err != nil {
		return nil, err
	}

	if err = ioutil.WriteFile(f.Name(), contentBytes, 0644); err != nil {
		return nil, err
	}

	return f, nil
}

func testGenerateData() ([]byte, error) {
	for _, key := range keys {
		resources := make([]storage.Resource, 0)
		for idx := 0; idx < rand.Intn(10)+1; idx++ {
			newResource := storage.Resource{
				"id":          strconv.Itoa(idx),
				"description": fmt.Sprintf("description-%s-%d", key, idx),
			}

			resources = append(resources, newResource)
		}

		testData[key] = resources
	}

	return json.Marshal(testData)
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
