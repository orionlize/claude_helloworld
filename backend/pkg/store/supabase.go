package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/supabase-community/supabase-go"
	"apihub/internal/model"
)

// SupabaseStore implements storage operations using Supabase
type SupabaseStore struct {
	client *supabase.Client
}

// NewSupabaseStore creates a new Supabase store
func NewSupabaseStore(client *supabase.Client) *SupabaseStore {
	return &SupabaseStore{
		client: client,
	}
}

// User operations
func (s *SupabaseStore) CreateUser(ctx context.Context, user *model.User) error {
	_, _, err := s.client.From("users").Insert(user, false, "", "", "").Execute()
	return err
}

func (s *SupabaseStore) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	data, _, err := s.client.From("users").Select("*", "", false).Eq("id", id).Execute()
	if err != nil {
		return nil, err
	}

	var users []model.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &users[0], nil
}

func (s *SupabaseStore) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	data, _, err := s.client.From("users").Select("*", "", false).Eq("username", username).Execute()
	if err != nil {
		return nil, err
	}

	var users []model.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &users[0], nil
}

// Project operations
func (s *SupabaseStore) CreateProject(ctx context.Context, project *model.Project) error {
	_, _, err := s.client.From("projects").Insert(project, false, "", "", "").Execute()
	return err
}

func (s *SupabaseStore) GetProjectByID(ctx context.Context, id string) (*model.Project, error) {
	data, _, err := s.client.From("projects").Select("*", "", false).Eq("id", id).Execute()
	if err != nil {
		return nil, err
	}

	var projects []model.Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		return nil, fmt.Errorf("project not found")
	}

	return &projects[0], nil
}

func (s *SupabaseStore) ListProjectsByUserID(ctx context.Context, userID string) ([]model.Project, error) {
	data, _, err := s.client.From("projects").Select("*", "", false).Eq("user_id", userID).Execute()
	if err != nil {
		return nil, err
	}

	var projects []model.Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *SupabaseStore) UpdateProject(ctx context.Context, project *model.Project) error {
	_, _, err := s.client.From("projects").Update(project, "", "").Eq("id", project.ID).Execute()
	return err
}

func (s *SupabaseStore) DeleteProject(ctx context.Context, id string) error {
	_, _, err := s.client.From("projects").Delete("", "").Eq("id", id).Execute()
	return err
}

// Collection operations
func (s *SupabaseStore) CreateCollection(ctx context.Context, collection *model.Collection) error {
	_, _, err := s.client.From("collections").Insert(collection, false, "", "", "").Execute()
	return err
}

func (s *SupabaseStore) GetCollectionByID(ctx context.Context, id string) (*model.Collection, error) {
	data, _, err := s.client.From("collections").Select("*", "", false).Eq("id", id).Execute()
	if err != nil {
		return nil, err
	}

	var collections []model.Collection
	if err := json.Unmarshal(data, &collections); err != nil {
		return nil, err
	}

	if len(collections) == 0 {
		return nil, fmt.Errorf("collection not found")
	}

	return &collections[0], nil
}

func (s *SupabaseStore) ListCollectionsByProjectID(ctx context.Context, projectID string) ([]model.Collection, error) {
	data, _, err := s.client.From("collections").Select("*", "", false).Eq("project_id", projectID).Execute()
	if err != nil {
		return nil, err
	}

	var collections []model.Collection
	if err := json.Unmarshal(data, &collections); err != nil {
		return nil, err
	}

	return collections, nil
}

func (s *SupabaseStore) UpdateCollection(ctx context.Context, collection *model.Collection) error {
	_, _, err := s.client.From("collections").Update(collection, "", "").Eq("id", collection.ID).Execute()
	return err
}

func (s *SupabaseStore) DeleteCollection(ctx context.Context, id string) error {
	_, _, err := s.client.From("collections").Delete("", "").Eq("id", id).Execute()
	return err
}

// Endpoint operations
func (s *SupabaseStore) CreateEndpoint(ctx context.Context, endpoint *model.Endpoint) error {
	_, _, err := s.client.From("endpoints").Insert(endpoint, false, "", "", "").Execute()
	return err
}

func (s *SupabaseStore) GetEndpointByID(ctx context.Context, id string) (*model.Endpoint, error) {
	data, _, err := s.client.From("endpoints").Select("*", "", false).Eq("id", id).Execute()
	if err != nil {
		return nil, err
	}

	var endpoints []model.Endpoint
	if err := json.Unmarshal(data, &endpoints); err != nil {
		return nil, err
	}

	if len(endpoints) == 0 {
		return nil, fmt.Errorf("endpoint not found")
	}

	return &endpoints[0], nil
}

func (s *SupabaseStore) ListEndpointsByCollectionID(ctx context.Context, collectionID string) ([]model.Endpoint, error) {
	data, _, err := s.client.From("endpoints").Select("*", "", false).Eq("collection_id", collectionID).Execute()
	if err != nil {
		return nil, err
	}

	var endpoints []model.Endpoint
	if err := json.Unmarshal(data, &endpoints); err != nil {
		return nil, err
	}

	return endpoints, nil
}

func (s *SupabaseStore) UpdateEndpoint(ctx context.Context, endpoint *model.Endpoint) error {
	_, _, err := s.client.From("endpoints").Update(endpoint, "", "").Eq("id", endpoint.ID).Execute()
	return err
}

func (s *SupabaseStore) DeleteEndpoint(ctx context.Context, id string) error {
	_, _, err := s.client.From("endpoints").Delete("", "").Eq("id", id).Execute()
	return err
}

// Environment operations
func (s *SupabaseStore) CreateEnvironment(ctx context.Context, env *model.Environment) error {
	_, _, err := s.client.From("environments").Insert(env, false, "", "", "").Execute()
	return err
}

func (s *SupabaseStore) GetEnvironmentByID(ctx context.Context, id string) (*model.Environment, error) {
	data, _, err := s.client.From("environments").Select("*", "", false).Eq("id", id).Execute()
	if err != nil {
		return nil, err
	}

	var environments []model.Environment
	if err := json.Unmarshal(data, &environments); err != nil {
		return nil, err
	}

	if len(environments) == 0 {
		return nil, fmt.Errorf("environment not found")
	}

	return &environments[0], nil
}

func (s *SupabaseStore) ListEnvironmentsByProjectID(ctx context.Context, projectID string) ([]model.Environment, error) {
	data, _, err := s.client.From("environments").Select("*", "", false).Eq("project_id", projectID).Execute()
	if err != nil {
		return nil, err
	}

	var environments []model.Environment
	if err := json.Unmarshal(data, &environments); err != nil {
		return nil, err
	}

	return environments, nil
}

func (s *SupabaseStore) UpdateEnvironment(ctx context.Context, env *model.Environment) error {
	_, _, err := s.client.From("environments").Update(env, "", "").Eq("id", env.ID).Execute()
	return err
}

func (s *SupabaseStore) DeleteEnvironment(ctx context.Context, id string) error {
	_, _, err := s.client.From("environments").Delete("", "").Eq("id", id).Execute()
	return err
}

// Health check
func (s *SupabaseStore) Ping(ctx context.Context) error {
	// Simple query to check connection
	_, _, err := s.client.From("users").Select("count", "", false).Limit(1, "").Execute()
	return err
}
