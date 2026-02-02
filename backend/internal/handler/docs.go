package handler

import (
	"fmt"
	"strings"

	"apihub/internal/model"
	"apihub/pkg/response"

	"github.com/gin-gonic/gin"
)

// GenerateMarkdown generates Markdown documentation for a project
func (h *Handler) GenerateDocs(c *gin.Context) {
	userID := c.GetString("user_id")
	projectID := c.Param("pid")
	format := c.Query("format") // markdown, html, openapi

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

	// Get all collections
	collections, err := h.store.GetCollectionsByProjectID(projectID)
	if err != nil {
		response.Error(c, 500, "Failed to fetch collections")
		return
	}

	// Build documentation
	switch format {
	case "openapi", "json":
		doc := h.generateOpenAPIDoc(project, collections)
		c.Header("Content-Type", "application/json")
		c.String(200, doc)
	case "html":
		doc := h.generateHTMLDoc(project, collections)
		c.Header("Content-Type", "text/html")
		c.String(200, doc)
	default:
		doc := h.generateMarkdownDoc(project, collections)
		c.Header("Content-Type", "text/markdown")
		c.String(200, doc)
	}
}

// generateMarkdownDoc generates Markdown documentation
func (h *Handler) generateMarkdownDoc(project *model.Project, collections []model.Collection) string {
	var sb strings.Builder

	projectName := project.Name
	if projectName == "" {
		projectName = "API Documentation"
	}

	sb.WriteString(fmt.Sprintf("# %s\n\n", projectName))
	if project.Description != "" {
		sb.WriteString(fmt.Sprintf("%s\n\n", project.Description))
	}

	sb.WriteString("## Base URL\n\n")
	sb.WriteString("```\nhttps://api.example.com\n```\n\n")

	// Process each collection
	for _, col := range collections {
		sb.WriteString(fmt.Sprintf("## %s\n\n", col.Name))
		if col.Description != "" {
			sb.WriteString(fmt.Sprintf("%s\n\n", col.Description))
		}

		// Get endpoints for this collection
		endpoints, _ := h.store.GetEndpointsByCollectionID(col.ID)

		for _, ep := range endpoints {
			sb.WriteString(fmt.Sprintf("### %s %s\n\n", ep.Method, ep.Name))
			if ep.Description != "" {
				sb.WriteString(fmt.Sprintf("**Description:** %s\n\n", ep.Description))
			}

			sb.WriteString(fmt.Sprintf("**Endpoint:** `%s`\n\n", ep.URL))
			sb.WriteString(fmt.Sprintf("**Method:** `%s`\n\n", ep.Method))

			if len(ep.Headers) > 0 {
				sb.WriteString("**Headers:**\n\n")
				sb.WriteString("```\n")
				for key, value := range ep.Headers {
					sb.WriteString(fmt.Sprintf("%s: %s\n", key, value))
				}
				sb.WriteString("```\n\n")
			}

			if ep.Body != nil && *ep.Body != "" {
				sb.WriteString("**Request Body:**\n\n")
				sb.WriteString("```json\n")
				sb.WriteString(*ep.Body)
				sb.WriteString("\n```\n\n")
			}

			sb.WriteString("---\n\n")
		}
	}

	return sb.String()
}

// generateOpenAPIDoc generates OpenAPI 3.0 JSON documentation
func (h *Handler) generateOpenAPIDoc(project *model.Project, collections []model.Collection) string {
	// Simplified OpenAPI 3.0 structure
	var sb strings.Builder

	projectName := project.Name
	if projectName == "" {
		projectName = "API Documentation"
	}

	sb.WriteString(`{
  "openapi": "3.0.0",
  "info": {
    "title": "`)
	sb.WriteString(projectName)
	sb.WriteString(`",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "https://api.example.com"
    }
  ],
  "paths": {
`)

	// Process each collection and endpoint
	for _, col := range collections {
		// Get endpoints for this collection
		endpoints, _ := h.store.GetEndpointsByCollectionID(col.ID)

		for _, ep := range endpoints {
			// Extract path from URL (remove domain)
			path := ep.URL
			if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
				if idx := strings.Index(path, "//"); idx != -1 {
					rest := path[idx+2:]
					if slashIdx := strings.Index(rest, "/"); slashIdx != -1 {
						path = rest[slashIdx:]
					}
				}
			}

			sb.WriteString(fmt.Sprintf("    %s: {\n", path))
			sb.WriteString(fmt.Sprintf("      %s: {\n", strings.ToLower(ep.Method)))
			sb.WriteString(`        "summary": "` + ep.Name + `",
        "description": "` + ep.Description + `",
        "responses": {
          "200": {
            "description": "Success"
          }
        }
      }
    }
`)
		}
	}

	sb.WriteString(`  }
}`)

	return sb.String()
}

// generateHTMLDoc generates HTML documentation
func (h *Handler) generateHTMLDoc(project *model.Project, collections []model.Collection) string {
	md := h.generateMarkdownDoc(project, collections)

	// Simple Markdown to HTML conversion
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API Documentation</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; line-height: 1.6; }
        h1 { border-bottom: 2px solid #333; padding-bottom: 10px; }
        h2 { color: #333; margin-top: 30px; }
        h3 { color: #555; }
        pre { background: #f4f4f4; padding: 10px; border-radius: 4px; overflow-x: auto; }
        code { background: #f4f4f4; padding: 2px 6px; border-radius: 3px; }
        hr { border: none; border-top: 1px solid #ddd; margin: 20px 0; }
    </style>
</head>
<body>
` + md + `
</body>
</html>`

	return html
}

// ExportPostmanCollection generates Postman Collection v2.1
func (h *Handler) ExportPostman(c *gin.Context) {
	userID := c.GetString("user_id")
	projectID := c.Param("pid")

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

	// Get all collections
	collections, _ := h.store.GetCollectionsByProjectID(projectID)

	// Build Postman collection
	collection := fmt.Sprintf(`{
  "info": {
    "name": "%s",
    "description": "%s",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
`,
		project.Name,
		project.Description,
	)

	// Add each collection as a folder
	for i, col := range collections {
		endpoints, _ := h.store.GetEndpointsByCollectionID(col.ID)

		collection += fmt.Sprintf(`    {
      "name": "%s",
      "item": [
`, col.Name)

		for j, ep := range endpoints {
			comma := ""
			if j < len(endpoints)-1 {
				comma = ","
			}

			collection += fmt.Sprintf(`        {
          "name": "%s",
          "request": {
            "method": "%s",
            "header": [],
            "url": {
              "raw": "%s"
            }
          }
        }%s
`,
			ep.Name,
			ep.Method,
			ep.URL,
			comma,
			)
		}

		collection += `      ]
    }`

		if i < len(collections)-1 {
			collection += ","
		}
	}

	collection += `
  ]
}
`

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", `attachment; filename="${project.Name}.postman_collection.json"`)
	c.String(200, collection)
}
