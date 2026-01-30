package store

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
	"apihub/internal/model"
)

type MemoryStore struct {
	users        map[int64]*model.User
	projects     map[string]*model.Project
	collections  map[string]*model.Collection
	endpoints    map[string]*model.Endpoint
	environments map[string]*model.Environment
	mu           sync.RWMutex

	userIDCounter int64
}

func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		users:        make(map[int64]*model.User),
		projects:     make(map[string]*model.Project),
		collections:  make(map[string]*model.Collection),
		endpoints:    make(map[string]*model.Endpoint),
		environments: make(map[string]*model.Environment),
		userIDCounter: 1,
	}

	// Generate password hash for demo123
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("demo123"), bcrypt.DefaultCost)

	// Create demo user (password: demo123)
	demoUser := &model.User{
		ID:        1,
		Email:     "demo@example.com",
		Password:  string(hashedPassword),
		Name:      "Demo User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	store.users[1] = demoUser

	return store
}

// User operations
func (s *MemoryStore) CreateUser(user *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user.ID = s.userIDCounter
	s.userIDCounter++
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	s.users[user.ID] = user
	return nil
}

func (s *MemoryStore) GetUserByEmail(email string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (s *MemoryStore) GetUserByID(id int64) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (s *MemoryStore) UpdateUser(user *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user.UpdatedAt = time.Now()
	s.users[user.ID] = user
	return nil
}

// Project operations
func (s *MemoryStore) CreateProject(project *model.Project) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if project.ID == "" {
		project.ID = generateID()
	}
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()
	s.projects[project.ID] = project
	return nil
}

func (s *MemoryStore) GetProjectsByUserID(userID int64) ([]model.Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var projects []model.Project
	userIDStr := fmt.Sprintf("%d", userID)
	for _, project := range s.projects {
		// Compare userID as string
		if project.UserID == userIDStr {
			projects = append(projects, *project)
		}
	}
	return projects, nil
}

func (s *MemoryStore) GetProjectByID(id string) (*model.Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	project, ok := s.projects[id]
	if !ok {
		return nil, nil
	}
	return project, nil
}

func (s *MemoryStore) UpdateProject(project *model.Project) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	project.UpdatedAt = time.Now()
	s.projects[project.ID] = project
	return nil
}

func (s *MemoryStore) DeleteProject(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.projects, id)
	return nil
}

// Collection operations
func (s *MemoryStore) CreateCollection(collection *model.Collection) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if collection.ID == "" {
		collection.ID = generateID()
	}
	collection.CreatedAt = time.Now()
	collection.UpdatedAt = time.Now()
	s.collections[collection.ID] = collection
	return nil
}

func (s *MemoryStore) GetCollectionsByProjectID(projectID string) ([]model.Collection, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var collections []model.Collection
	for _, collection := range s.collections {
		if collection.ProjectID == projectID {
			collections = append(collections, *collection)
		}
	}
	return collections, nil
}

func (s *MemoryStore) GetCollectionByID(id string) (*model.Collection, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	collection, ok := s.collections[id]
	if !ok {
		return nil, nil
	}
	return collection, nil
}

func (s *MemoryStore) UpdateCollection(collection *model.Collection) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	collection.UpdatedAt = time.Now()
	s.collections[collection.ID] = collection
	return nil
}

func (s *MemoryStore) DeleteCollection(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.collections, id)
	return nil
}

// Endpoint operations
func (s *MemoryStore) CreateEndpoint(endpoint *model.Endpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if endpoint.ID == "" {
		endpoint.ID = generateID()
	}
	endpoint.CreatedAt = time.Now()
	endpoint.UpdatedAt = time.Now()
	s.endpoints[endpoint.ID] = endpoint
	return nil
}

func (s *MemoryStore) GetEndpointsByCollectionID(collectionID string) ([]model.Endpoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var endpoints []model.Endpoint
	for _, endpoint := range s.endpoints {
		if endpoint.CollectionID == collectionID {
			endpoints = append(endpoints, *endpoint)
		}
	}
	return endpoints, nil
}

func (s *MemoryStore) GetEndpointByID(id string) (*model.Endpoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	endpoint, ok := s.endpoints[id]
	if !ok {
		return nil, nil
	}
	return endpoint, nil
}

func (s *MemoryStore) UpdateEndpoint(endpoint *model.Endpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	endpoint.UpdatedAt = time.Now()
	s.endpoints[endpoint.ID] = endpoint
	return nil
}

func (s *MemoryStore) DeleteEndpoint(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.endpoints, id)
	return nil
}

// Environment operations
func (s *MemoryStore) CreateEnvironment(env *model.Environment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if env.ID == "" {
		env.ID = generateID()
	}
	env.CreatedAt = time.Now()
	env.UpdatedAt = time.Now()
	s.environments[env.ID] = env
	return nil
}

func (s *MemoryStore) GetEnvironmentsByProjectID(projectID string) ([]model.Environment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var environments []model.Environment
	for _, env := range s.environments {
		if env.ProjectID == projectID {
			environments = append(environments, *env)
		}
	}
	return environments, nil
}

func (s *MemoryStore) GetEnvironmentByID(id string) (*model.Environment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	env, ok := s.environments[id]
	if !ok {
		return nil, nil
	}
	return env, nil
}

func (s *MemoryStore) UpdateEnvironment(env *model.Environment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	env.UpdatedAt = time.Now()
	s.environments[env.ID] = env
	return nil
}

func (s *MemoryStore) DeleteEnvironment(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.environments, id)
	return nil
}

func generateID() string {
	return time.Now().Format("20060102150405") + randomString(6)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
