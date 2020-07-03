package storage

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"strconv"
)

// Storage implements the service interface.
type Storage struct {
	File     string
	Key      string
	Singular bool
}

// NewStorage returns a new storage instance.
func NewStorage(file, key string, single bool) (*Storage, error) {
	return &Storage{File: file, Key: key, Singular: single}, nil
}

// Find all resources for the specific key.
func (s *Storage) Find() ([]Resource, error) {
	data, err := readFile(s.File)
	if err != nil {
		return nil, err
	}

	return data[s.Key].([]Resource), nil
}

// FindById a resource for the specific key.
func (s *Storage) FindById(id string) (Resource, error) {
	data, err := readFile(s.File)
	if err != nil {
		return nil, err
	}

	// Check if singular endpoint and construct the resource to be returned.
	if s.Singular {
		return Resource{s.Key: data[s.Key]}, nil
	}

	for _, resource := range data[s.Key].([]Resource) {
		if resource["id"] == id {
			return resource, nil
		}
	}

	return nil, ErrResourceNotFound
}

// Create a new resource for the specific key.
func (s *Storage) Create(newResource Resource) (Resource, error) {
	data, err := readFile(s.File)
	if err != nil {
		return nil, err
	}

	_, ok := newResource["id"]
	if !ok {
		newResource["id"] = generateNewId(data[s.Key].([]Resource))
	} else {
		for _, resource := range data[s.Key].([]Resource) {
			if resource["id"] == newResource["id"] {
				return nil, ErrResourceAlreadyExists
			}
		}
	}

	newData := append(data[s.Key].([]Resource), newResource)
	data[s.Key] = newData

	if err := updateFile(s.File, data); err != nil {
		return nil, err
	}

	return newResource, nil
}

// Update an existing resource for the specific key.
func (s *Storage) Update(id string, updatedResource Resource) (Resource, error) {
	data, err := readFile(s.File)
	if err != nil {
		return nil, err
	}

	updatedResource["id"] = id

	newResources := make([]Resource, 0)
	for _, d := range data[s.Key].([]Resource) {
		if d["id"] == id {
			newResources = append(newResources, updatedResource)
		} else {
			newResources = append(newResources, d)
		}
	}

	data[s.Key] = newResources

	if err := updateFile(s.File, data); err != nil {
		return nil, err
	}

	return updatedResource, nil
}

// Delete an existing resource for the specific key.
func (s *Storage) Delete(id string) error {
	data, err := readFile(s.File)
	if err != nil {
		return err
	}

	newResources := make([]Resource, 0)
	for _, d := range data[s.Key].([]Resource) {
		if d["id"] == id {
			continue
		}

		newResources = append(newResources, d)
	}

	data[s.Key] = newResources

	return updateFile(s.File, data)
}

// readFile returns all the data from the watch file.
func readFile(file string) (map[string]interface{}, error) {
	contentBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	content := map[string]interface{}{}
	if err = json.Unmarshal(contentBytes, &content); err != nil {
		return nil, err
	}

	contentResource := make(map[string]interface{}, 0)
	for key, val := range content {
		data := make([]Resource, 0)

		switch v := val.(type) {
		case []interface{}:
			for _, resource := range v {
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

			contentResource[key] = data
		default:
			contentResource[key] = v
		}
	}

	return contentResource, nil
}

// updateFile formats and writes the new data to the watch file.
func updateFile(file string, content map[string]interface{}) error {
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
	existingIds := make(map[string]bool, 0)
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
