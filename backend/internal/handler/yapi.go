package handler

import (
	"strconv"

	"apihub/internal/model"
	"apihub/pkg/response"
	"apihub/pkg/yapi"

	"github.com/gin-gonic/gin"
)

// SyncFromYAPI syncs interfaces from YAPI to the project
func (h *Handler) SyncFromYAPI(c *gin.Context) {
	userID := c.GetString("user_id")
	projectID := c.Param("pid")

	var req model.YAPISyncConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	// Override project_id from URL
	req.ProjectID = projectID

	// Validate config
	if err := yapi.ValidateConfig(req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

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

	// Create YAPI client
	client := yapi.NewClient(req.YAPIUrl, req.YAPIToken)

	// Test connection
	if err := client.TestConnection(req.YAPIProjectID); err != nil {
		response.Error(c, 400, "Failed to connect to YAPI: "+err.Error())
		return
	}

	// Get categories (collections)
	categories, err := client.GetCategories(req.YAPIProjectID)
	if err != nil {
		response.Error(c, 500, "Failed to fetch categories: "+err.Error())
		return
	}

	// Get interfaces
	interfaces, err := client.GetInterfaces(req.YAPIProjectID)
	if err != nil {
		response.Error(c, 500, "Failed to fetch interfaces: "+err.Error())
		return
	}

	// Sync data
	stats, err := h.syncYAPIData(projectID, categories, interfaces, client)
	if err != nil {
		response.Error(c, 500, "Failed to sync data: "+err.Error())
		return
	}

	response.Success(c, 200, gin.H{
		"message": "YAPI sync completed successfully",
		"stats":   stats,
	})
}

// syncYAPIData syncs YAPI categories and interfaces to the local database
func (h *Handler) syncYAPIData(projectID string, categories []model.YAPICategory, interfaces []model.YAPIInterface, client *yapi.Client) (*SyncStats, error) {
	stats := &SyncStats{}

	// Create a map of catID to collection ID
	catToCollection := make(map[int]string)

	// First, sync categories (collections)
	for _, cat := range categories {
		collection := &model.Collection{
			ProjectID:   projectID,
			Name:        cat.Name,
			Description: "Synced from YAPI",
			SortOrder:   cat.Order,
		}

		// Check if collection exists by name (using YAPI cat ID as reference)
		existingCollections, _ := h.store.GetCollectionsByProjectID(projectID)
		var existingID *string
		for _, col := range existingCollections {
			if col.Name == cat.Name {
				id := col.ID
				existingID = &id
				break
			}
		}

		if existingID != nil {
			// Update existing collection
			collection.ID = *existingID
			h.store.UpdateCollection(collection)
			stats.UpdatedCollections++
		} else {
			// Create new collection
			h.store.CreateCollection(collection)
			stats.CreatedCollections++
		}

		catToCollection[cat.ID] = collection.ID
	}

	// Then, sync interfaces (endpoints)
	for _, yInt := range interfaces {
		// Skip if no category found
		collectionID, ok := catToCollection[yInt.CatID]
		if !ok {
			continue
		}

		// Convert YAPI interface to Endpoint
		endpoint, err := yapi.ConvertYAPIInterfaceToEndpoint(yInt)
		if err != nil {
			continue
		}
		endpoint.CollectionID = collectionID

		// Check if endpoint exists by name and method in this collection
		existingEndpoints, _ := h.store.GetEndpointsByCollectionID(collectionID)
		var existingID *string
		for _, ep := range existingEndpoints {
			if ep.Name == endpoint.Name && ep.Method == endpoint.Method {
				id := ep.ID
				existingID = &id
				break
			}
		}

		if existingID != nil {
			// Update existing endpoint
			endpoint.ID = *existingID
			h.store.UpdateEndpoint(endpoint)
			stats.UpdatedEndpoints++
		} else {
			// Create new endpoint
			h.store.CreateEndpoint(endpoint)
			stats.CreatedEndpoints++
		}
	}

	return stats, nil
}

// TestYAPIConnection tests the YAPI connection
func (h *Handler) TestYAPIConnection(c *gin.Context) {
	var req model.YAPISyncConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	// For testing connection, only validate YAPI-specific fields (not ProjectID which is for sync)
	if req.YAPIUrl == "" {
		response.Error(c, 400, "yapi_url is required")
		return
	}
	if req.YAPIToken == "" {
		response.Error(c, 400, "yapi_token is required")
		return
	}
	if req.YAPIProjectID == 0 {
		response.Error(c, 400, "yapi_project_id is required")
		return
	}

	client := yapi.NewClient(req.YAPIUrl, req.YAPIToken)

	yapiProject, err := client.GetProject(req.YAPIProjectID)
	if err != nil {
		response.Error(c, 400, "Failed to connect to YAPI: "+err.Error())
		return
	}

	// Get categories to verify token works
	categories, err := client.GetCategories(req.YAPIProjectID)
	if err != nil {
		response.Error(c, 400, "Failed to fetch categories. Please check your token permissions.")
		return
	}

	interfaces, err := client.GetInterfaces(req.YAPIProjectID)
	if err != nil {
		response.Error(c, 400, "Failed to fetch interfaces. Please check your token permissions.")
		return
	}

	response.Success(c, 200, gin.H{
		"message":           "Connection successful",
		"yapi_project":      yapiProject,
		"total_categories":  len(categories),
		"total_interfaces":  len(interfaces),
	})
}

// GetYAPIProjectInfo fetches YAPI project info
func (h *Handler) GetYAPIProjectInfo(c *gin.Context) {
	yapiURL := c.Query("yapi_url")
	token := c.Query("token")
	projectIDStr := c.Query("project_id")

	if yapiURL == "" || token == "" || projectIDStr == "" {
		response.Error(c, 400, "Missing required parameters: yapi_url, token, project_id")
		return
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		response.Error(c, 400, "Invalid project_id")
		return
	}

	client := yapi.NewClient(yapiURL, token)

	yapiProject, err := client.GetProject(projectID)
	if err != nil {
		response.Error(c, 400, "Failed to fetch YAPI project: "+err.Error())
		return
	}

	categories, err := client.GetCategories(projectID)
	if err != nil {
		response.Error(c, 400, "Failed to fetch categories: "+err.Error())
		return
	}

	interfaces, err := client.GetInterfaces(projectID)
	if err != nil {
		response.Error(c, 400, "Failed to fetch interfaces: "+err.Error())
		return
	}

	response.Success(c, 200, gin.H{
		"project":    yapiProject,
		"categories": categories,
		"interfaces": len(interfaces),
	})
}

// SyncStats represents sync statistics
type SyncStats struct {
	CreatedCollections  int `json:"created_collections"`
	UpdatedCollections  int `json:"updated_collections"`
	CreatedEndpoints    int `json:"created_endpoints"`
	UpdatedEndpoints    int `json:"updated_endpoints"`
}
