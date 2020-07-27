package storage

type Mock struct {
	data Database
	key  string
}

// NewMock returns a new mock instance.
func NewMock(data Database, key string) (*Mock, error) {
	return &Mock{data: data, key: key}, nil
}

// Find all mock resources for the specific key.
func (m *Mock) Find() ([]Resource, error) {
	return m.data[m.key], nil
}

// FindById a mock resource for the specific key.
func (m *Mock) FindById(id string) (Resource, error) {
	if err := checkResourceKeyExists(m.data, m.key); err != nil {
		return nil, ErrResourceNotFound
	}

	for _, resource := range m.data[m.key] {
		if resource["id"] == id {
			return resource, nil
		}
	}

	return nil, ErrResourceNotFound
}

// Create a new mock resource for the specific key.
func (m *Mock) Create(newResource Resource) (Resource, error) {
	if err := checkResourceKeyExists(m.data, m.key); err != nil {
		return nil, ErrResourceNotFound
	}

	_, ok := newResource["id"]
	if !ok {
		newResource["id"] = generateNewId(m.data[m.key])
	} else {
		for _, resource := range m.data[m.key] {
			if resource["id"] == newResource["id"] {
				return nil, ErrResourceAlreadyExists
			}
		}
	}

	newData := append(m.data[m.key], newResource)
	m.data[m.key] = newData

	return newResource, nil
}

// Replace an existing mock resource for the specific key.
func (m *Mock) Replace(id string, replaced Resource) (Resource, error) {
	if err := checkResourceKeyExists(m.data, m.key); err != nil {
		return nil, ErrResourceNotFound
	}

	// Check if resource with the requested id exists.
	if _, err := m.FindById(id); err != nil {
		return nil, err
	}

	replaced["id"] = id

	newResources := make([]Resource, 0)
	for _, d := range m.data[m.key] {
		if d["id"] == id {
			newResources = append(newResources, replaced)
		} else {
			newResources = append(newResources, d)
		}
	}

	m.data[m.key] = newResources

	return replaced, nil
}

// Update an existing mock resource for the specific key.
func (m *Mock) Update(id string, updatedReq Resource) (Resource, error) {
	if err := checkResourceKeyExists(m.data, m.key); err != nil {
		return nil, ErrResourceNotFound
	}

	// Check if resource with the requested id exists and retrieve it.
	updated, err := m.FindById(id)
	if err != nil {
		return nil, err
	}

	// Apply any changes to current resource.
	for key, val := range updatedReq {
		updated[key] = val
	}

	updated["id"] = id

	newResources := make([]Resource, 0)
	for _, d := range m.data[m.key] {
		if d["id"] == id {
			newResources = append(newResources, updated)
		} else {
			newResources = append(newResources, d)
		}
	}

	m.data[m.key] = newResources

	return updated, nil
}

// Delete an existing mock resource for the specific key.
func (m *Mock) Delete(id string) error {
	if err := checkResourceKeyExists(m.data, m.key); err != nil {
		return ErrResourceNotFound
	}

	// Check if resource with the requested id exists.
	if _, err := m.FindById(id); err != nil {
		return err
	}

	newResources := make([]Resource, 0)
	for _, d := range m.data[m.key] {
		if d["id"] == id {
			continue
		}

		newResources = append(newResources, d)
	}

	m.data[m.key] = newResources

	return nil
}

// DB returns all the mock resources.
func (m *Mock) DB() (Database, error) {
	return m.data, nil
}

func (m *Mock) SetData(data Database) {
	m.data = data
}
