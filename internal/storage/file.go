package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"strconv"
)

var (
	errResourceInvalidType = errors.New("invalid resource type")
)

// File implements the storage interface, and uses a file as 'database'.
type File struct {
	filename string
	key      string
}

// NewFile returns a new file instance.
func NewFile(filename, key string) (*File, error) {
	return &File{filename: filename, key: key}, nil
}

// Find all resources for the specific key.
func (f *File) Find() ([]Resource, error) {
	data, err := readFile(f.filename)
	if err != nil {
		return nil, err
	}

	if err = checkResourceKeyExists(data, f.key); err != nil {
		return nil, ErrResourceNotFound
	}

	return data[f.key], nil
}

// FindById a resource for the specific key.
func (f *File) FindById(id string) (Resource, error) {
	data, err := readFile(f.filename)
	if err != nil {
		return nil, err
	}

	if err = checkResourceKeyExists(data, f.key); err != nil {
		return nil, ErrResourceNotFound
	}

	for _, resource := range data[f.key] {
		if resource["id"] == id {
			return resource, nil
		}
	}

	return nil, ErrResourceNotFound
}

// Create a new resource for the specific key.
func (f *File) Create(newResource Resource) (Resource, error) {
	data, err := readFile(f.filename)
	if err != nil {
		return nil, err
	}

	if err = checkResourceKeyExists(data, f.key); err != nil {
		return nil, ErrResourceNotFound
	}

	_, ok := newResource["id"]
	if !ok {
		newResource["id"] = generateNewId(data[f.key])
	} else {
		for _, resource := range data[f.key] {
			if resource["id"] == newResource["id"] {
				return nil, ErrResourceAlreadyExists
			}
		}
	}

	newData := append(data[f.key], newResource)
	data[f.key] = newData

	if err := updateFile(f.filename, data); err != nil {
		return nil, err
	}

	return newResource, nil
}

// Replace an existing resource for the specific key.
func (f *File) Replace(id string, replaced Resource) (Resource, error) {
	data, err := readFile(f.filename)
	if err != nil {
		return nil, err
	}

	if err = checkResourceKeyExists(data, f.key); err != nil {
		return nil, ErrResourceNotFound
	}

	// Check if resource with the requested id exists.
	if _, err = f.FindById(id); err != nil {
		return nil, err
	}

	replaced["id"] = id

	newResources := make([]Resource, 0)
	for _, d := range data[f.key] {
		if d["id"] == id {
			newResources = append(newResources, replaced)
		} else {
			newResources = append(newResources, d)
		}
	}

	data[f.key] = newResources

	if err := updateFile(f.filename, data); err != nil {
		return nil, err
	}

	return replaced, nil
}

// Update an existing resource for the specific key.
func (f *File) Update(id string, updatedReq Resource) (Resource, error) {
	data, err := readFile(f.filename)
	if err != nil {
		return nil, err
	}

	if err = checkResourceKeyExists(data, f.key); err != nil {
		return nil, ErrResourceNotFound
	}

	// Check if resource with the requested id exists and retrieve it.
	updated, err := f.FindById(id)
	if err != nil {
		return nil, err
	}

	// Apply any changes to current resource.
	for key, val := range updatedReq {
		updated[key] = val
	}

	updated["id"] = id

	newResources := make([]Resource, 0)
	for _, d := range data[f.key] {
		if d["id"] == id {
			newResources = append(newResources, updated)
		} else {
			newResources = append(newResources, d)
		}
	}

	data[f.key] = newResources

	if err := updateFile(f.filename, data); err != nil {
		return nil, err
	}

	return updated, nil
}

// Delete an existing resource for the specific key.
func (f *File) Delete(id string) error {
	data, err := readFile(f.filename)
	if err != nil {
		return err
	}

	if err = checkResourceKeyExists(data, f.key); err != nil {
		return ErrResourceNotFound
	}

	// Check if resource with the requested id exists.
	if _, err = f.FindById(id); err != nil {
		return err
	}

	newResources := make([]Resource, 0)
	for _, d := range data[f.key] {
		if d["id"] == id {
			continue
		}

		newResources = append(newResources, d)
	}

	data[f.key] = newResources

	return updateFile(f.filename, data)
}

// DB returns all resources.
func (f *File) DB() (Database, error) {
	return readFile(f.filename)
}

// readFile returns all the data from the watch file.
func readFile(file string) (Database, error) {
	contentBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	content := make(map[string]interface{})
	if err = json.Unmarshal(contentBytes, &content); err != nil {
		return nil, err
	}

	database := make(Database)
	for key, val := range content {
		valResources, ok := val.([]interface{})
		if !ok {
			return nil, errResourceInvalidType
		}

		data := make([]Resource, 0)
		for _, resource := range valResources {
			resourceBytes, err := json.Marshal(resource)
			if err != nil {
				return nil, err
			}

			var newResource Resource
			if err := json.Unmarshal(resourceBytes, &newResource); err != nil {
				return nil, err
			}

			data = append(data, newResource)
		}

		database[key] = data
	}

	return database, nil
}

// updateFile formats and writes the new data to the watch file.
func updateFile(file string, content Database) error {
	contentBytes, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(file, contentBytes, 0644); err != nil {
		return err
	}

	return nil
}

// generateNewId and validate that is unique across provided data.
func generateNewId(data []Resource) string {
	existingIds := make(map[string]bool)
	for _, d := range data {
		existingIds[d["id"].(string)] = true
	}

	for {
		newId := strconv.Itoa(rand.Intn(1000))

		if !existingIds[newId] {
			return newId
		}
	}
}

// checkResourceKeyExists in the file data.
func checkResourceKeyExists(database Database, key string) error {
	if _, ok := database[key]; !ok {
		return ErrResourceNotFound
	}

	return nil
}
