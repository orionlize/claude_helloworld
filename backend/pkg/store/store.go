package store

import "apihub/internal/model"

// Store defines the interface for data persistence
type Store interface {
	// User operations
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(id int64) (*model.User, error)
	UpdateUser(user *model.User) error

	// Project operations
	CreateProject(project *model.Project) error
	GetProjectsByUserID(userID int64) ([]model.Project, error)
	GetProjectByID(id string) (*model.Project, error)
	UpdateProject(project *model.Project) error
	DeleteProject(id string) error

	// Collection operations
	CreateCollection(collection *model.Collection) error
	GetCollectionsByProjectID(projectID string) ([]model.Collection, error)
	GetCollectionByID(id string) (*model.Collection, error)
	UpdateCollection(collection *model.Collection) error
	DeleteCollection(id string) error

	// Endpoint operations
	CreateEndpoint(endpoint *model.Endpoint) error
	GetEndpointsByCollectionID(collectionID string) ([]model.Endpoint, error)
	GetEndpointByID(id string) (*model.Endpoint, error)
	UpdateEndpoint(endpoint *model.Endpoint) error
	DeleteEndpoint(id string) error

	// Environment operations
	CreateEnvironment(env *model.Environment) error
	GetEnvironmentsByProjectID(projectID string) ([]model.Environment, error)
	GetEnvironmentByID(id string) (*model.Environment, error)
	UpdateEnvironment(env *model.Environment) error
	DeleteEnvironment(id string) error
}
