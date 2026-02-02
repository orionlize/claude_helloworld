package store

import (
	"context"
	"encoding/json"
	"errors"

	"apihub/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DatabaseStore implements Store interface using PostgreSQL
type DatabaseStore struct {
	db *pgxpool.Pool
}

// NewDatabaseStore creates a new database-backed store
func NewDatabaseStore(db *pgxpool.Pool) *DatabaseStore {
	return &DatabaseStore{db: db}
}

// User operations
func (s *DatabaseStore) CreateUser(user *model.User) error {
	ctx := context.Background()

	// Generate UUID for user
	userID := generateID()
	query := `INSERT INTO users (id, username, email, password_hash) VALUES ($1, $2, $3, $4) RETURNING id, username, email, created_at, updated_at`

	err := s.db.QueryRow(ctx, query, userID, user.Username, user.Email, user.Password).Scan(
		&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)

	return err
}

func (s *DatabaseStore) GetUserByEmail(email string) (*model.User, error) {
	ctx := context.Background()
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = $1`

	var user model.User
	err := s.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s *DatabaseStore) GetUserByID(id string) (*model.User, error) {
	ctx := context.Background()
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE id = $1`

	var user model.User
	err := s.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s *DatabaseStore) UpdateUser(user *model.User) error {
	ctx := context.Background()
	query := `UPDATE users SET username = $1, email = $2, password_hash = $3 WHERE id = $4`

	_, err := s.db.Exec(ctx, query, user.Username, user.Email, user.Password, user.ID)
	return err
}

// Project operations
func (s *DatabaseStore) CreateProject(project *model.Project) error {
	ctx := context.Background()

	if project.ID == "" {
		project.ID = generateID()
	}

	query := `INSERT INTO projects (id, name, description, user_id) VALUES ($1, $2, $3, $4) RETURNING id, name, description, user_id, created_at, updated_at`

	err := s.db.QueryRow(ctx, query, project.ID, project.Name, project.Description, project.UserID).Scan(
		&project.ID, &project.Name, &project.Description, &project.UserID, &project.CreatedAt, &project.UpdatedAt,
	)

	return err
}

func (s *DatabaseStore) GetProjectsByUserID(userID string) ([]model.Project, error) {
	ctx := context.Background()
	query := `SELECT id, name, description, user_id, created_at, updated_at FROM projects WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var p model.Project
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.UserID, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (s *DatabaseStore) GetProjectByID(id string) (*model.Project, error) {
	ctx := context.Background()
	query := `SELECT id, name, description, user_id, created_at, updated_at FROM projects WHERE id = $1`

	var project model.Project
	err := s.db.QueryRow(ctx, query, id).Scan(
		&project.ID, &project.Name, &project.Description, &project.UserID, &project.CreatedAt, &project.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &project, nil
}

func (s *DatabaseStore) UpdateProject(project *model.Project) error {
	ctx := context.Background()
	query := `UPDATE projects SET name = $1, description = $2 WHERE id = $3`

	_, err := s.db.Exec(ctx, query, project.Name, project.Description, project.ID)
	return err
}

func (s *DatabaseStore) DeleteProject(id string) error {
	ctx := context.Background()
	query := `DELETE FROM projects WHERE id = $1`

	_, err := s.db.Exec(ctx, query, id)
	return err
}

// Collection operations
func (s *DatabaseStore) CreateCollection(collection *model.Collection) error {
	ctx := context.Background()

	if collection.ID == "" {
		collection.ID = generateID()
	}

	query := `INSERT INTO collections (id, project_id, name, description, parent_id, sort_order) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, project_id, name, description, parent_id, sort_order, created_at, updated_at`

	err := s.db.QueryRow(ctx, query, collection.ID, collection.ProjectID, collection.Name, collection.Description, collection.ParentID, collection.SortOrder).Scan(
		&collection.ID, &collection.ProjectID, &collection.Name, &collection.Description, &collection.ParentID, &collection.SortOrder, &collection.CreatedAt, &collection.UpdatedAt,
	)

	return err
}

func (s *DatabaseStore) GetCollectionsByProjectID(projectID string) ([]model.Collection, error) {
	ctx := context.Background()
	query := `SELECT id, project_id, name, description, parent_id, sort_order, created_at, updated_at FROM collections WHERE project_id = $1 ORDER BY sort_order ASC`

	rows, err := s.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []model.Collection
	for rows.Next() {
		var col model.Collection
		err := rows.Scan(&col.ID, &col.ProjectID, &col.Name, &col.Description, &col.ParentID, &col.SortOrder, &col.CreatedAt, &col.UpdatedAt)
		if err != nil {
			return nil, err
		}
		collections = append(collections, col)
	}

	return collections, nil
}

func (s *DatabaseStore) GetCollectionByID(id string) (*model.Collection, error) {
	ctx := context.Background()
	query := `SELECT id, project_id, name, description, parent_id, sort_order, created_at, updated_at FROM collections WHERE id = $1`

	var collection model.Collection
	err := s.db.QueryRow(ctx, query, id).Scan(
		&collection.ID, &collection.ProjectID, &collection.Name, &collection.Description, &collection.ParentID, &collection.SortOrder, &collection.CreatedAt, &collection.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &collection, nil
}

func (s *DatabaseStore) UpdateCollection(collection *model.Collection) error {
	ctx := context.Background()
	query := `UPDATE collections SET name = $1, description = $2, parent_id = $3 WHERE id = $4`

	_, err := s.db.Exec(ctx, query, collection.Name, collection.Description, collection.ParentID, collection.ID)
	return err
}

func (s *DatabaseStore) DeleteCollection(id string) error {
	ctx := context.Background()
	query := `DELETE FROM collections WHERE id = $1`

	_, err := s.db.Exec(ctx, query, id)
	return err
}

// Endpoint operations
func (s *DatabaseStore) CreateEndpoint(endpoint *model.Endpoint) error {
	ctx := context.Background()

	if endpoint.ID == "" {
		endpoint.ID = generateID()
	}

	headersJSON, _ := json.Marshal(endpoint.Headers)
	requestParamsJSON, _ := json.Marshal(endpoint.RequestParams)
	requestBodyJSON, _ := json.Marshal(endpoint.RequestBody)
	responseParamsJSON, _ := json.Marshal(endpoint.ResponseParams)
	responseBodyJSON, _ := json.Marshal(endpoint.ResponseBody)

	query := `INSERT INTO endpoints (id, collection_id, name, method, url, headers, body, description, request_params, request_body, response_params, response_body) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, collection_id, name, method, url, headers, body, description, sort_order, request_params, request_body, response_params, response_body, created_at, updated_at`

	var headersJSONB, bodyStr, requestParamsJSONB, requestBodyJSONB, responseParamsJSONB, responseBodyJSONB []byte
	err := s.db.QueryRow(ctx, query,
		endpoint.ID, endpoint.CollectionID, endpoint.Name, endpoint.Method, endpoint.URL,
		headersJSON, endpoint.Body, endpoint.Description,
		requestParamsJSON, requestBodyJSON, responseParamsJSON, responseBodyJSON,
	).Scan(
		&endpoint.ID, &endpoint.CollectionID, &endpoint.Name, &endpoint.Method, &endpoint.URL,
		&headersJSONB, &bodyStr, &endpoint.Description, &endpoint.SortOrder,
		&requestParamsJSONB, &requestBodyJSONB, &responseParamsJSONB, &responseBodyJSONB,
		&endpoint.CreatedAt, &endpoint.UpdatedAt,
	)

	if err != nil {
		return err
	}

	if len(headersJSONB) > 0 {
		json.Unmarshal(headersJSONB, &endpoint.Headers)
	}
	if len(requestParamsJSONB) > 0 {
		json.Unmarshal(requestParamsJSONB, &endpoint.RequestParams)
	}
	if len(requestBodyJSONB) > 0 {
		json.Unmarshal(requestBodyJSONB, &endpoint.RequestBody)
	}
	if len(responseParamsJSONB) > 0 {
		json.Unmarshal(responseParamsJSONB, &endpoint.ResponseParams)
	}
	if len(responseBodyJSONB) > 0 {
		json.Unmarshal(responseBodyJSONB, &endpoint.ResponseBody)
	}

	return nil
}

func (s *DatabaseStore) GetEndpointsByCollectionID(collectionID string) ([]model.Endpoint, error) {
	ctx := context.Background()
	query := `SELECT id, collection_id, name, method, url, headers, body, description, sort_order, request_params, request_body, response_params, response_body, created_at, updated_at FROM endpoints WHERE collection_id = $1 ORDER BY sort_order ASC`

	rows, err := s.db.Query(ctx, query, collectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var endpoints []model.Endpoint
	for rows.Next() {
		var e model.Endpoint
		var headersJSON, bodyJSON, requestParamsJSON, requestBodyJSON, responseParamsJSON, responseBodyJSON []byte
		err := rows.Scan(
			&e.ID, &e.CollectionID, &e.Name, &e.Method, &e.URL,
			&headersJSON, &bodyJSON, &e.Description, &e.SortOrder,
			&requestParamsJSON, &requestBodyJSON, &responseParamsJSON, &responseBodyJSON,
			&e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if len(headersJSON) > 0 {
			json.Unmarshal(headersJSON, &e.Headers)
		}
		if len(bodyJSON) > 0 {
			bodyStr := string(bodyJSON)
			e.Body = &bodyStr
		}
		if len(requestParamsJSON) > 0 {
			json.Unmarshal(requestParamsJSON, &e.RequestParams)
		}
		if len(requestBodyJSON) > 0 {
			json.Unmarshal(requestBodyJSON, &e.RequestBody)
		}
		if len(responseParamsJSON) > 0 {
			json.Unmarshal(responseParamsJSON, &e.ResponseParams)
		}
		if len(responseBodyJSON) > 0 {
			json.Unmarshal(responseBodyJSON, &e.ResponseBody)
		}

		endpoints = append(endpoints, e)
	}

	return endpoints, nil
}

func (s *DatabaseStore) GetEndpointByID(id string) (*model.Endpoint, error) {
	ctx := context.Background()
	query := `SELECT id, collection_id, name, method, url, headers, body, description, sort_order, request_params, request_body, response_params, response_body, created_at, updated_at FROM endpoints WHERE id = $1`

	var endpoint model.Endpoint
	var headersJSON, bodyJSON, requestParamsJSON, requestBodyJSON, responseParamsJSON, responseBodyJSON []byte
	err := s.db.QueryRow(ctx, query, id).Scan(
		&endpoint.ID, &endpoint.CollectionID, &endpoint.Name, &endpoint.Method, &endpoint.URL,
		&headersJSON, &bodyJSON, &endpoint.Description, &endpoint.SortOrder,
		&requestParamsJSON, &requestBodyJSON, &responseParamsJSON, &responseBodyJSON,
		&endpoint.CreatedAt, &endpoint.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if len(headersJSON) > 0 {
		json.Unmarshal(headersJSON, &endpoint.Headers)
	}
	if len(bodyJSON) > 0 {
		bodyStr := string(bodyJSON)
		endpoint.Body = &bodyStr
	}
	if len(requestParamsJSON) > 0 {
		json.Unmarshal(requestParamsJSON, &endpoint.RequestParams)
	}
	if len(requestBodyJSON) > 0 {
		json.Unmarshal(requestBodyJSON, &endpoint.RequestBody)
	}
	if len(responseParamsJSON) > 0 {
		json.Unmarshal(responseParamsJSON, &endpoint.ResponseParams)
	}
	if len(responseBodyJSON) > 0 {
		json.Unmarshal(responseBodyJSON, &endpoint.ResponseBody)
	}

	return &endpoint, nil
}

func (s *DatabaseStore) UpdateEndpoint(endpoint *model.Endpoint) error {
	ctx := context.Background()

	headersJSON, _ := json.Marshal(endpoint.Headers)
	requestParamsJSON, _ := json.Marshal(endpoint.RequestParams)
	requestBodyJSON, _ := json.Marshal(endpoint.RequestBody)
	responseParamsJSON, _ := json.Marshal(endpoint.ResponseParams)
	responseBodyJSON, _ := json.Marshal(endpoint.ResponseBody)

	query := `UPDATE endpoints SET name = $1, method = $2, url = $3, headers = $4, body = $5, description = $6, request_params = $7, request_body = $8, response_params = $9, response_body = $10 WHERE id = $11`

	_, err := s.db.Exec(ctx, query,
		endpoint.Name, endpoint.Method, endpoint.URL, headersJSON, endpoint.Body, endpoint.Description,
		requestParamsJSON, requestBodyJSON, responseParamsJSON, responseBodyJSON,
		endpoint.ID)
	return err
}

func (s *DatabaseStore) DeleteEndpoint(id string) error {
	ctx := context.Background()
	query := `DELETE FROM endpoints WHERE id = $1`

	_, err := s.db.Exec(ctx, query, id)
	return err
}

// Environment operations
func (s *DatabaseStore) CreateEnvironment(env *model.Environment) error {
	ctx := context.Background()

	if env.ID == "" {
		env.ID = generateID()
	}

	varsJSON, _ := json.Marshal(env.Variables)

	query := `INSERT INTO environments (id, project_id, name, variables, is_default) VALUES ($1, $2, $3, $4, $5) RETURNING id, project_id, name, variables, is_default, created_at, updated_at`

	var varsJSONB []byte
	err := s.db.QueryRow(ctx, query, env.ID, env.ProjectID, env.Name, varsJSON, env.IsDefault).Scan(
		&env.ID, &env.ProjectID, &env.Name, &varsJSONB, &env.IsDefault, &env.CreatedAt, &env.UpdatedAt,
	)

	if err != nil {
		return err
	}

	if len(varsJSONB) > 0 {
		json.Unmarshal(varsJSONB, &env.Variables)
	}

	return nil
}

func (s *DatabaseStore) GetEnvironmentsByProjectID(projectID string) ([]model.Environment, error) {
	ctx := context.Background()
	query := `SELECT id, project_id, name, variables, is_default, created_at, updated_at FROM environments WHERE project_id = $1 ORDER BY is_default DESC, name ASC`

	rows, err := s.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var environments []model.Environment
	for rows.Next() {
		var env model.Environment
		var varsJSON []byte
		err := rows.Scan(&env.ID, &env.ProjectID, &env.Name, &varsJSON, &env.IsDefault, &env.CreatedAt, &env.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if len(varsJSON) > 0 {
			json.Unmarshal(varsJSON, &env.Variables)
		}

		environments = append(environments, env)
	}

	return environments, nil
}

func (s *DatabaseStore) GetEnvironmentByID(id string) (*model.Environment, error) {
	ctx := context.Background()
	query := `SELECT id, project_id, name, variables, is_default, created_at, updated_at FROM environments WHERE id = $1`

	var env model.Environment
	var varsJSON []byte
	err := s.db.QueryRow(ctx, query, id).Scan(
		&env.ID, &env.ProjectID, &env.Name, &varsJSON, &env.IsDefault, &env.CreatedAt, &env.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if len(varsJSON) > 0 {
		json.Unmarshal(varsJSON, &env.Variables)
	}

	return &env, nil
}

func (s *DatabaseStore) UpdateEnvironment(env *model.Environment) error {
	ctx := context.Background()

	varsJSON, _ := json.Marshal(env.Variables)
	query := `UPDATE environments SET name = $1, variables = $2, is_default = $3 WHERE id = $4`

	_, err := s.db.Exec(ctx, query, env.Name, varsJSON, env.IsDefault, env.ID)
	return err
}

func (s *DatabaseStore) DeleteEnvironment(id string) error {
	ctx := context.Background()
	query := `DELETE FROM environments WHERE id = $1`

	_, err := s.db.Exec(ctx, query, id)
	return err
}

