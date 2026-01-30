package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"apihub/internal/model"
	"apihub/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (h *Handler) ListProjects(c *gin.Context) {
	userID := c.GetString("user_id")

	ctx := context.Background()
	query := `
		SELECT id, name, description, user_id, created_at, updated_at
		FROM projects
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := h.db.Query(ctx, query, userID)
	if err != nil {
		response.Error(c, 500, "Failed to fetch projects")
		return
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var p model.Project
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.UserID, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			response.Error(c, 500, "Failed to scan project")
			return
		}
		projects = append(projects, p)
	}

	response.Success(c, 200, projects)
}

func (h *Handler) CreateProject(c *gin.Context) {
	userID := c.GetString("user_id")

	var req model.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	ctx := context.Background()
	projectID := generateUUID()

	query := `
		INSERT INTO projects (id, name, description, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, user_id, created_at, updated_at
	`

	var project model.Project
	err := h.db.QueryRow(ctx, query, projectID, req.Name, req.Description, userID).Scan(
		&project.ID, &project.Name, &project.Description, &project.UserID, &project.CreatedAt, &project.UpdatedAt,
	)
	if err != nil {
		response.Error(c, 500, "Failed to create project")
		return
	}

	response.Success(c, 201, project)
}

func (h *Handler) GetProject(c *gin.Context) {
	userID := c.GetString("user_id")
	projectID := c.Param("id")

	ctx := context.Background()
	query := `
		SELECT id, name, description, user_id, created_at, updated_at
		FROM projects
		WHERE id = $1 AND user_id = $2
	`

	var project model.Project
	err := h.db.QueryRow(ctx, query, projectID, userID).Scan(
		&project.ID, &project.Name, &project.Description, &project.UserID, &project.CreatedAt, &project.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(c, 404, "Project not found")
		} else {
			response.Error(c, 500, "Database error")
		}
		return
	}

	response.Success(c, 200, project)
}

func (h *Handler) UpdateProject(c *gin.Context) {
	userID := c.GetString("user_id")
	projectID := c.Param("id")

	var req model.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	ctx := context.Background()
	query := `
		UPDATE projects
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3 AND user_id = $4
		RETURNING id, name, description, user_id, created_at, updated_at
	`

	var project model.Project
	err := h.db.QueryRow(ctx, query, req.Name, req.Description, projectID, userID).Scan(
		&project.ID, &project.Name, &project.Description, &project.UserID, &project.CreatedAt, &project.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(c, 404, "Project not found")
		} else {
			response.Error(c, 500, "Database error")
		}
		return
	}

	response.Success(c, 200, project)
}

func (h *Handler) DeleteProject(c *gin.Context) {
	userID := c.GetString("user_id")
	projectID := c.Param("id")

	ctx := context.Background()
	query := `DELETE FROM projects WHERE id = $1 AND user_id = $2`

	result, err := h.db.Exec(ctx, query, projectID, userID)
	if err != nil {
		response.Error(c, 500, "Failed to delete project")
		return
	}

	if result.RowsAffected() == 0 {
		response.Error(c, 404, "Project not found")
		return
	}

	response.Success(c, 200, gin.H{"message": "Project deleted successfully"})
}

// Collections
func (h *Handler) ListCollections(c *gin.Context) {
	projectID := c.Param("project_id")
	userID := c.GetString("user_id")

	// First verify project ownership
	ctx := context.Background()
	var ownerID string
	err := h.db.QueryRow(ctx, "SELECT user_id FROM projects WHERE id = $1", projectID).Scan(&ownerID)
	if err != nil || ownerID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	query := `
		SELECT id, project_id, name, description, parent_id, sort_order, created_at, updated_at
		FROM collections
		WHERE project_id = $1
		ORDER BY sort_order ASC
	`

	rows, err := h.db.Query(ctx, query, projectID)
	if err != nil {
		response.Error(c, 500, "Failed to fetch collections")
		return
	}
	defer rows.Close()

	var collections []model.Collection
	for rows.Next() {
		var col model.Collection
		err := rows.Scan(&col.ID, &col.ProjectID, &col.Name, &col.Description, &col.ParentID, &col.SortOrder, &col.CreatedAt, &col.UpdatedAt)
		if err != nil {
			response.Error(c, 500, "Failed to scan collection")
			return
		}
		collections = append(collections, col)
	}

	response.Success(c, 200, collections)
}

func (h *Handler) CreateCollection(c *gin.Context) {
	projectID := c.Param("project_id")
	_userID := c.GetString("user_id")

	var req model.CreateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	ctx := context.Background()
	collectionID := generateUUID()

	query := `
		INSERT INTO collections (id, project_id, name, description, parent_id, sort_order)
		VALUES ($1, $2, $3, $4, $5, 0)
		RETURNING id, project_id, name, description, parent_id, sort_order, created_at, updated_at
	`

	var collection model.Collection
	err := h.db.QueryRow(ctx, query, collectionID, projectID, req.Name, req.Description, req.ParentID).Scan(
		&collection.ID, &collection.ProjectID, &collection.Name, &collection.Description, &collection.ParentID, &collection.SortOrder, &collection.CreatedAt, &collection.UpdatedAt,
	)
	if err != nil {
		response.Error(c, 500, "Failed to create collection")
		return
	}

	response.Success(c, 201, collection)
}

func (h *Handler) GetCollection(c *gin.Context) {
	collectionID := c.Param("id")
	userID := c.GetString("user_id")

	ctx := context.Background()
	query := `
		SELECT c.id, c.project_id, c.name, c.description, c.parent_id, c.sort_order, c.created_at, c.updated_at
		FROM collections c
		JOIN projects p ON c.project_id = p.id
		WHERE c.id = $1 AND p.user_id = $2
	`

	var collection model.Collection
	err := h.db.QueryRow(ctx, query, collectionID, userID).Scan(
		&collection.ID, &collection.ProjectID, &collection.Name, &collection.Description, &collection.ParentID, &collection.SortOrder, &collection.CreatedAt, &collection.UpdatedAt,
	)
	if err != nil {
		response.Error(c, 404, "Collection not found")
		return
	}

	response.Success(c, 200, collection)
}

func (h *Handler) UpdateCollection(c *gin.Context) {
	collectionID := c.Param("id")
	userID := c.GetString("user_id")

	var req model.UpdateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	ctx := context.Background()
	query := `
		UPDATE collections
		SET name = $1, description = $2, parent_id = $3, updated_at = NOW()
		FROM projects
		WHERE collections.id = $4 AND collections.project_id = projects.id AND projects.user_id = $5
		RETURNING collections.id, collections.project_id, collections.name, collections.description, collections.parent_id, collections.sort_order, collections.created_at, collections.updated_at
	`

	var collection model.Collection
	err := h.db.QueryRow(ctx, query, req.Name, req.Description, req.ParentID, collectionID, userID).Scan(
		&collection.ID, &collection.ProjectID, &collection.Name, &collection.Description, &collection.ParentID, &collection.SortOrder, &collection.CreatedAt, &collection.UpdatedAt,
	)
	if err != nil {
		response.Error(c, 404, "Collection not found")
		return
	}

	response.Success(c, 200, collection)
}

func (h *Handler) DeleteCollection(c *gin.Context) {
	collectionID := c.Param("id")
	userID := c.GetString("user_id")

	ctx := context.Background()
	query := `
		DELETE FROM collections
		USING projects
		WHERE collections.id = $1 AND collections.project_id = projects.id AND projects.user_id = $2
	`

	result, err := h.db.Exec(ctx, query, collectionID, userID)
	if err != nil {
		response.Error(c, 500, "Failed to delete collection")
		return
	}

	if result.RowsAffected() == 0 {
		response.Error(c, 404, "Collection not found")
		return
	}

	response.Success(c, 200, gin.H{"message": "Collection deleted successfully"})
}

// Endpoints
func (h *Handler) ListEndpoints(c *gin.Context) {
	collectionID := c.Param("collection_id")
	userID := c.GetString("user_id")

	ctx := context.Background()
	query := `
		SELECT e.id, e.collection_id, e.name, e.method, e.url, e.headers, e.body, e.description, e.sort_order, e.created_at, e.updated_at
		FROM endpoints e
		JOIN collections c ON e.collection_id = c.id
		JOIN projects p ON c.project_id = p.id
		WHERE e.collection_id = $1 AND p.user_id = $2
		ORDER BY e.sort_order ASC
	`

	rows, err := h.db.Query(ctx, query, collectionID, userID)
	if err != nil {
		response.Error(c, 500, "Failed to fetch endpoints")
		return
	}
	defer rows.Close()

	var endpoints []model.Endpoint
	for rows.Next() {
		var e model.Endpoint
		var headersJSON, bodyJSON []byte
		err := rows.Scan(&e.ID, &e.CollectionID, &e.Name, &e.Method, &e.URL, &headersJSON, &bodyJSON, &e.Description, &e.SortOrder, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			response.Error(c, 500, "Failed to scan endpoint")
			return
		}

		if len(headersJSON) > 0 {
			json.Unmarshal(headersJSON, &e.Headers)
		}
		if len(bodyJSON) > 0 {
			bodyStr := string(bodyJSON)
			e.Body = &bodyStr
		}

		endpoints = append(endpoints, e)
	}

	response.Success(c, 200, endpoints)
}

func (h *Handler) CreateEndpoint(c *gin.Context) {
	collectionID := c.Param("collection_id")
	// _userID := c.GetString("user_id")

	var req model.CreateEndpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	ctx := context.Background()
	endpointID := generateUUID()

	headersJSON, _ := json.Marshal(req.Headers)

	query := `
		INSERT INTO endpoints (id, collection_id, name, method, url, headers, body, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, collection_id, name, method, url, headers, body, description, sort_order, created_at, updated_at
	`

	var endpoint model.Endpoint
	err := h.db.QueryRow(ctx, query, endpointID, collectionID, req.Name, req.Method, req.URL, headersJSON, req.Body, req.Description).Scan(
		&endpoint.ID, &endpoint.CollectionID, &endpoint.Name, &endpoint.Method, &endpoint.URL, &headersJSON, &endpoint.Body, &endpoint.Description, &endpoint.SortOrder, &endpoint.CreatedAt, &endpoint.UpdatedAt,
	)
	if err != nil {
		response.Error(c, 500, "Failed to create endpoint")
		return
	}

	json.Unmarshal(headersJSON, &endpoint.Headers)

	response.Success(c, 201, endpoint)
}

func (h *Handler) GetEndpoint(c *gin.Context) {
	endpointID := c.Param("id")
	userID := c.GetString("user_id")

	ctx := context.Background()
	query := `
		SELECT e.id, e.collection_id, e.name, e.method, e.url, e.headers, e.body, e.description, e.sort_order, e.created_at, e.updated_at
		FROM endpoints e
		JOIN collections c ON e.collection_id = c.id
		JOIN projects p ON c.project_id = p.id
		WHERE e.id = $1 AND p.user_id = $2
	`

	var endpoint model.Endpoint
	var headersJSON, bodyJSON []byte
	err := h.db.QueryRow(ctx, query, endpointID, userID).Scan(
		&endpoint.ID, &endpoint.CollectionID, &endpoint.Name, &endpoint.Method, &endpoint.URL, &headersJSON, &bodyJSON, &endpoint.Description, &endpoint.SortOrder, &endpoint.CreatedAt, &endpoint.UpdatedAt,
	)
	if err != nil {
		response.Error(c, 404, "Endpoint not found")
		return
	}

	if len(headersJSON) > 0 {
		json.Unmarshal(headersJSON, &endpoint.Headers)
	}
	if len(bodyJSON) > 0 {
		bodyStr := string(bodyJSON)
		endpoint.Body = &bodyStr
	}

	response.Success(c, 200, endpoint)
}

func (h *Handler) UpdateEndpoint(c *gin.Context) {
	endpointID := c.Param("id")
	userID := c.GetString("user_id")

	var req model.UpdateEndpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	ctx := context.Background()
	headersJSON, _ := json.Marshal(req.Headers)

	query := `
		UPDATE endpoints
		SET name = $1, method = $2, url = $3, headers = $4, body = $5, description = $6, updated_at = NOW()
		FROM collections, projects
		WHERE endpoints.id = $7 AND endpoints.collection_id = collections.id AND collections.project_id = projects.id AND projects.user_id = $8
		RETURNING endpoints.id, endpoints.collection_id, endpoints.name, endpoints.method, endpoints.url, endpoints.headers, endpoints.body, endpoints.description, endpoints.sort_order, endpoints.created_at, endpoints.updated_at
	`

	var endpoint model.Endpoint
	var bodyJSON []byte
	err := h.db.QueryRow(ctx, query, req.Name, req.Method, req.URL, headersJSON, req.Body, endpointID, userID).Scan(
		&endpoint.ID, &endpoint.CollectionID, &endpoint.Name, &endpoint.Method, &endpoint.URL, &headersJSON, &bodyJSON, &endpoint.Description, &endpoint.SortOrder, &endpoint.CreatedAt, &endpoint.UpdatedAt,
	)
	if err != nil {
		response.Error(c, 404, "Endpoint not found")
		return
	}

	json.Unmarshal(headersJSON, &endpoint.Headers)
	if len(bodyJSON) > 0 {
		bodyStr := string(bodyJSON)
		endpoint.Body = &bodyStr
	}

	response.Success(c, 200, endpoint)
}

func (h *Handler) DeleteEndpoint(c *gin.Context) {
	endpointID := c.Param("id")
	userID := c.GetString("user_id")

	ctx := context.Background()
	query := `
		DELETE FROM endpoints
		USING collections, projects
		WHERE endpoints.id = $1 AND endpoints.collection_id = collections.id AND collections.project_id = projects.id AND projects.user_id = $2
	`

	result, err := h.db.Exec(ctx, query, endpointID, userID)
	if err != nil {
		response.Error(c, 500, "Failed to delete endpoint")
		return
	}

	if result.RowsAffected() == 0 {
		response.Error(c, 404, "Endpoint not found")
		return
	}

	response.Success(c, 200, gin.H{"message": "Endpoint deleted successfully"})
}

// Environments
func (h *Handler) ListEnvironments(c *gin.Context) {
	projectID := c.Param("project_id")
	userID := c.GetString("user_id")

	ctx := context.Background()
	query := `
		SELECT id, project_id, name, variables, is_default, created_at, updated_at
		FROM environments
		WHERE project_id = $1 AND EXISTS (SELECT 1 FROM projects WHERE id = $1 AND user_id = $2)
		ORDER BY is_default DESC, name ASC
	`

	rows, err := h.db.Query(ctx, query, projectID, userID)
	if err != nil {
		response.Error(c, 500, "Failed to fetch environments")
		return
	}
	defer rows.Close()

	var environments []model.Environment
	for rows.Next() {
		var env model.Environment
		var varsJSON []byte
		err := rows.Scan(&env.ID, &env.ProjectID, &env.Name, &varsJSON, &env.IsDefault, &env.CreatedAt, &env.UpdatedAt)
		if err != nil {
			response.Error(c, 500, "Failed to scan environment")
			return
		}

		if len(varsJSON) > 0 {
			json.Unmarshal(varsJSON, &env.Variables)
		}

		environments = append(environments, env)
	}

	response.Success(c, 200, environments)
}

func (h *Handler) CreateEnvironment(c *gin.Context) {
	projectID := c.Param("project_id")
	// _userID := c.GetString("user_id")

	var req model.CreateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	ctx := context.Background()
	envID := generateUUID()

	varsJSON, _ := json.Marshal(req.Variables)

	query := `
		INSERT INTO environments (id, project_id, name, variables, is_default)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, project_id, name, variables, is_default, created_at, updated_at
	`

	var environment model.Environment
	err := h.db.QueryRow(ctx, query, envID, projectID, req.Name, varsJSON, req.IsDefault).Scan(
		&environment.ID, &environment.ProjectID, &environment.Name, &varsJSON, &environment.IsDefault, &environment.CreatedAt, &environment.UpdatedAt,
	)
	if err != nil {
		response.Error(c, 500, "Failed to create environment")
		return
	}

	json.Unmarshal(varsJSON, &environment.Variables)

	response.Success(c, 201, environment)
}

func (h *Handler) GetEnvironment(c *gin.Context) {
	envID := c.Param("id")
	userID := c.GetString("user_id")

	ctx := context.Background()
	query := `
		SELECT e.id, e.project_id, e.name, e.variables, e.is_default, e.created_at, e.updated_at
		FROM environments e
		JOIN projects p ON e.project_id = p.id
		WHERE e.id = $1 AND p.user_id = $2
	`

	var environment model.Environment
	var varsJSON []byte
	err := h.db.QueryRow(ctx, query, envID, userID).Scan(
		&environment.ID, &environment.ProjectID, &environment.Name, &varsJSON, &environment.IsDefault, &environment.CreatedAt, &environment.UpdatedAt,
	)
	if err != nil {
		response.Error(c, 404, "Environment not found")
		return
	}

	if len(varsJSON) > 0 {
		json.Unmarshal(varsJSON, &environment.Variables)
	}

	response.Success(c, 200, environment)
}

func (h *Handler) UpdateEnvironment(c *gin.Context) {
	envID := c.Param("id")
	userID := c.GetString("user_id")

	var req model.UpdateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	ctx := context.Background()
	varsJSON, _ := json.Marshal(req.Variables)

	query := `
		UPDATE environments
		SET name = $1, variables = $2, is_default = $3, updated_at = NOW()
		FROM projects
		WHERE environments.id = $4 AND environments.project_id = projects.id AND projects.user_id = $5
		RETURNING environments.id, environments.project_id, environments.name, environments.variables, environments.is_default, environments.created_at, environments.updated_at
	`

	var environment model.Environment
	err := h.db.QueryRow(ctx, query, req.Name, varsJSON, req.IsDefault, envID, userID).Scan(
		&environment.ID, &environment.ProjectID, &environment.Name, &varsJSON, &environment.IsDefault, &environment.CreatedAt, &environment.UpdatedAt,
	)
	if err != nil {
		response.Error(c, 404, "Environment not found")
		return
	}

	json.Unmarshal(varsJSON, &environment.Variables)

	response.Success(c, 200, environment)
}

func (h *Handler) DeleteEnvironment(c *gin.Context) {
	envID := c.Param("id")
	userID := c.GetString("user_id")

	ctx := context.Background()
	query := `
		DELETE FROM environments
		USING projects
		WHERE environments.id = $1 AND environments.project_id = projects.id AND projects.user_id = $2
	`

	result, err := h.db.Exec(ctx, query, envID, userID)
	if err != nil {
		response.Error(c, 500, "Failed to delete environment")
		return
	}

	if result.RowsAffected() == 0 {
		response.Error(c, 404, "Environment not found")
		return
	}

	response.Success(c, 200, gin.H{"message": "Environment deleted successfully"})
}

// SendRequest - Send HTTP request for testing
func (h *Handler) SendRequest(c *gin.Context) {
	var req struct {
		Method  string            `json:"method" binding:"required"`
		URL     string            `json:"url" binding:"required"`
		Headers map[string]string `json:"headers"`
		Body    *string           `json:"body"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	// Create HTTP request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	var bodyReader *strings.Reader
	if req.Body != nil {
		bodyReader = strings.NewReader(*req.Body)
	} else {
		bodyReader = strings.NewReader("")
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, bodyReader)
	if err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	// Set headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Send request
	startTime := time.Now()
	resp, err := client.Do(httpReq)
	if err != nil {
		response.Error(c, 500, "Failed to send request: "+err.Error())
		return
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		response.Error(c, 500, "Failed to read response")
		return
	}

	duration := time.Since(startTime)

	// Build response headers map
	respHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			respHeaders[key] = values[0]
		}
	}

	result := gin.H{
		"status_code": resp.StatusCode,
		"status":      resp.Status,
		"headers":     respHeaders,
		"body":        string(respBody),
		"duration":    duration.Milliseconds(),
	}

	response.Success(c, 200, result)
}
