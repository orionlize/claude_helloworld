package handler

import (
	"io"
	"net/http"
	"strings"
	"time"

	"apihub/internal/model"
	"apihub/pkg/response"

	"github.com/gin-gonic/gin"
)

// Projects

func (h *Handler) ListProjects(c *gin.Context) {
	userID := c.GetString("user_id")

	projects, err := h.store.GetProjectsByUserID(userID)
	if err != nil {
		response.Error(c, 500, "Failed to fetch projects")
		return
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

	project := &model.Project{
		Name:        req.Name,
		Description: req.Description,
		UserID:      userID,
	}

	if err := h.store.CreateProject(project); err != nil {
		response.Error(c, 500, "Failed to create project")
		return
	}

	response.Success(c, 201, project)
}

func (h *Handler) GetProject(c *gin.Context) {
	userID := c.GetString("user_id")
	projectID := c.Param("pid")

	project, err := h.store.GetProjectByID(projectID)
	if err != nil || project == nil {
		response.Error(c, 404, "Project not found")
		return
	}

	// Verify ownership
	if project.UserID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	response.Success(c, 200, project)
}

func (h *Handler) UpdateProject(c *gin.Context) {
	userID := c.GetString("user_id")
	projectID := c.Param("pid")

	var req model.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	project, err := h.store.GetProjectByID(projectID)
	if err != nil || project == nil {
		response.Error(c, 404, "Project not found")
		return
	}

	// Verify ownership
	if project.UserID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	project.Name = req.Name
	project.Description = req.Description

	if err := h.store.UpdateProject(project); err != nil {
		response.Error(c, 500, "Failed to update project")
		return
	}

	response.Success(c, 200, project)
}

func (h *Handler) DeleteProject(c *gin.Context) {
	userID := c.GetString("user_id")
	projectID := c.Param("pid")

	project, err := h.store.GetProjectByID(projectID)
	if err != nil || project == nil {
		response.Error(c, 404, "Project not found")
		return
	}

	// Verify ownership
	if project.UserID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	if err := h.store.DeleteProject(projectID); err != nil {
		response.Error(c, 500, "Failed to delete project")
		return
	}

	response.Success(c, 200, gin.H{"message": "Project deleted successfully"})
}

// Collections

func (h *Handler) ListCollections(c *gin.Context) {
	projectID := c.Param("pid")
	userID := c.GetString("user_id")

	// Verify project ownership
	project, err := h.store.GetProjectByID(projectID)
	if err != nil || project == nil {
		response.Error(c, 404, "Project not found")
		return
	}

	if project.UserID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	collections, err := h.store.GetCollectionsByProjectID(projectID)
	if err != nil {
		response.Error(c, 500, "Failed to fetch collections")
		return
	}

	response.Success(c, 200, collections)
}

func (h *Handler) CreateCollection(c *gin.Context) {
	projectID := c.Param("pid")
	userID := c.GetString("user_id")

	// Verify project ownership
	project, err := h.store.GetProjectByID(projectID)
	if err != nil || project == nil {
		response.Error(c, 404, "Project not found")
		return
	}

	if project.UserID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	var req model.CreateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	collection := &model.Collection{
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
	}

	if err := h.store.CreateCollection(collection); err != nil {
		response.Error(c, 500, "Failed to create collection")
		return
	}

	response.Success(c, 201, collection)
}

func (h *Handler) GetCollection(c *gin.Context) {
	collectionID := c.Param("cid")
	userID := c.GetString("user_id")

	collection, err := h.store.GetCollectionByID(collectionID)
	if err != nil || collection == nil {
		response.Error(c, 404, "Collection not found")
		return
	}

	// Verify project ownership
	project, err := h.store.GetProjectByID(collection.ProjectID)
	if err != nil || project == nil || project.UserID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	response.Success(c, 200, collection)
}

func (h *Handler) UpdateCollection(c *gin.Context) {
	collectionID := c.Param("cid")
	userID := c.GetString("user_id")

	collection, err := h.store.GetCollectionByID(collectionID)
	if err != nil || collection == nil {
		response.Error(c, 404, "Collection not found")
		return
	}

	// Verify project ownership
	project, err := h.store.GetProjectByID(collection.ProjectID)
	if err != nil || project == nil || project.UserID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	var req model.UpdateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	collection.Name = req.Name
	collection.Description = req.Description
	collection.ParentID = req.ParentID

	if err := h.store.UpdateCollection(collection); err != nil {
		response.Error(c, 500, "Failed to update collection")
		return
	}

	response.Success(c, 200, collection)
}

func (h *Handler) DeleteCollection(c *gin.Context) {
	collectionID := c.Param("cid")
	userID := c.GetString("user_id")

	collection, err := h.store.GetCollectionByID(collectionID)
	if err != nil || collection == nil {
		response.Error(c, 404, "Collection not found")
		return
	}

	// Verify project ownership
	project, err := h.store.GetProjectByID(collection.ProjectID)
	if err != nil || project == nil || project.UserID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	if err := h.store.DeleteCollection(collectionID); err != nil {
		response.Error(c, 500, "Failed to delete collection")
		return
	}

	response.Success(c, 200, gin.H{"message": "Collection deleted successfully"})
}

// Endpoints

func (h *Handler) ListEndpoints(c *gin.Context) {
	collectionID := c.Param("cid")
	userID := c.GetString("user_id")

	// Verify collection access through project ownership
	collection, err := h.store.GetCollectionByID(collectionID)
	if err != nil || collection == nil {
		response.Error(c, 404, "Collection not found")
		return
	}

	project, err := h.store.GetProjectByID(collection.ProjectID)
	if err != nil || project == nil || project.UserID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	endpoints, err := h.store.GetEndpointsByCollectionID(collectionID)
	if err != nil {
		response.Error(c, 500, "Failed to fetch endpoints")
		return
	}

	response.Success(c, 200, endpoints)
}

func (h *Handler) CreateEndpoint(c *gin.Context) {
	collectionID := c.Param("cid")
	userID := c.GetString("user_id")

	// Verify collection access through project ownership
	collection, err := h.store.GetCollectionByID(collectionID)
	if err != nil || collection == nil {
		response.Error(c, 404, "Collection not found")
		return
	}

	project, err := h.store.GetProjectByID(collection.ProjectID)
	if err != nil || project == nil || project.UserID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	var req model.CreateEndpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	endpoint := &model.Endpoint{
		CollectionID: collectionID,
		Name:         req.Name,
		Method:       req.Method,
		URL:          req.URL,
		Headers:      req.Headers,
		Body:         req.Body,
		Description:  req.Description,
	}

	if err := h.store.CreateEndpoint(endpoint); err != nil {
		response.Error(c, 500, "Failed to create endpoint")
		return
	}

	response.Success(c, 201, endpoint)
}

func (h *Handler) GetEndpoint(c *gin.Context) {
	endpointID := c.Param("epid")
	userID := c.GetString("user_id")

	endpoint, err := h.store.GetEndpointByID(endpointID)
	if err != nil || endpoint == nil {
		response.Error(c, 404, "Endpoint not found")
		return
	}

	// Verify access through project ownership
	collection, err := h.store.GetCollectionByID(endpoint.CollectionID)
	if err != nil || collection == nil {
		response.Error(c, 404, "Collection not found")
		return
	}

	project, err := h.store.GetProjectByID(collection.ProjectID)
	if err != nil || project == nil || project.UserID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	response.Success(c, 200, endpoint)
}

func (h *Handler) UpdateEndpoint(c *gin.Context) {
	endpointID := c.Param("epid")
	userID := c.GetString("user_id")

	endpoint, err := h.store.GetEndpointByID(endpointID)
	if err != nil || endpoint == nil {
		response.Error(c, 404, "Endpoint not found")
		return
	}

	// Verify access through project ownership
	collection, err := h.store.GetCollectionByID(endpoint.CollectionID)
	if err != nil || collection == nil {
		response.Error(c, 404, "Collection not found")
		return
	}

	project, err := h.store.GetProjectByID(collection.ProjectID)
	if err != nil || project == nil || project.UserID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	var req model.UpdateEndpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	endpoint.Name = req.Name
	endpoint.Method = req.Method
	endpoint.URL = req.URL
	endpoint.Headers = req.Headers
	endpoint.Body = req.Body
	endpoint.Description = req.Description

	if err := h.store.UpdateEndpoint(endpoint); err != nil {
		response.Error(c, 500, "Failed to update endpoint")
		return
	}

	response.Success(c, 200, endpoint)
}

func (h *Handler) DeleteEndpoint(c *gin.Context) {
	endpointID := c.Param("epid")
	userID := c.GetString("user_id")

	endpoint, err := h.store.GetEndpointByID(endpointID)
	if err != nil || endpoint == nil {
		response.Error(c, 404, "Endpoint not found")
		return
	}

	// Verify access through project ownership
	collection, err := h.store.GetCollectionByID(endpoint.CollectionID)
	if err != nil || collection == nil {
		response.Error(c, 404, "Collection not found")
		return
	}

	project, err := h.store.GetProjectByID(collection.ProjectID)
	if err != nil || project == nil || project.UserID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	if err := h.store.DeleteEndpoint(endpointID); err != nil {
		response.Error(c, 500, "Failed to delete endpoint")
		return
	}

	response.Success(c, 200, gin.H{"message": "Endpoint deleted successfully"})
}

// Environments

func (h *Handler) ListEnvironments(c *gin.Context) {
	projectID := c.Param("pid")
	userID := c.GetString("user_id")

	// Verify project ownership
	project, err := h.store.GetProjectByID(projectID)
	if err != nil || project == nil {
		response.Error(c, 404, "Project not found")
		return
	}

	if project.UserID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	environments, err := h.store.GetEnvironmentsByProjectID(projectID)
	if err != nil {
		response.Error(c, 500, "Failed to fetch environments")
		return
	}

	response.Success(c, 200, environments)
}

func (h *Handler) CreateEnvironment(c *gin.Context) {
	projectID := c.Param("pid")
	userID := c.GetString("user_id")

	// Verify project ownership
	project, err := h.store.GetProjectByID(projectID)
	if err != nil || project == nil {
		response.Error(c, 404, "Project not found")
		return
	}

	if project.UserID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	var req model.CreateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	env := &model.Environment{
		ProjectID: projectID,
		Name:      req.Name,
		Variables: req.Variables,
		IsDefault: req.IsDefault,
	}

	if err := h.store.CreateEnvironment(env); err != nil {
		response.Error(c, 500, "Failed to create environment")
		return
	}

	response.Success(c, 201, env)
}

func (h *Handler) GetEnvironment(c *gin.Context) {
	envID := c.Param("eid")
	userID := c.GetString("user_id")

	env, err := h.store.GetEnvironmentByID(envID)
	if err != nil || env == nil {
		response.Error(c, 404, "Environment not found")
		return
	}

	// Verify project ownership
	project, err := h.store.GetProjectByID(env.ProjectID)
	if err != nil || project == nil || project.UserID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	response.Success(c, 200, env)
}

func (h *Handler) UpdateEnvironment(c *gin.Context) {
	envID := c.Param("eid")
	userID := c.GetString("user_id")

	env, err := h.store.GetEnvironmentByID(envID)
	if err != nil || env == nil {
		response.Error(c, 404, "Environment not found")
		return
	}

	// Verify project ownership
	project, err := h.store.GetProjectByID(env.ProjectID)
	if err != nil || project == nil || project.UserID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	var req model.UpdateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	env.Name = req.Name
	env.Variables = req.Variables
	env.IsDefault = req.IsDefault

	if err := h.store.UpdateEnvironment(env); err != nil {
		response.Error(c, 500, "Failed to update environment")
		return
	}

	response.Success(c, 200, env)
}

func (h *Handler) DeleteEnvironment(c *gin.Context) {
	envID := c.Param("eid")
	userID := c.GetString("user_id")

	env, err := h.store.GetEnvironmentByID(envID)
	if err != nil || env == nil {
		response.Error(c, 404, "Environment not found")
		return
	}

	// Verify project ownership
	project, err := h.store.GetProjectByID(env.ProjectID)
	if err != nil || project == nil || project.UserID != userID {
		response.Error(c, 403, "Access denied")
		return
	}

	if err := h.store.DeleteEnvironment(envID); err != nil {
		response.Error(c, 500, "Failed to delete environment")
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
