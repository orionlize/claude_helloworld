package yapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"apihub/internal/model"
)

// Client provides methods to interact with YAPI Open API
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NewClient creates a new YAPI client
func NewClient(baseURL, token string) *Client {
	// Ensure baseURL doesn't have trailing slash
	baseURL = strings.TrimSuffix(baseURL, "/")
	// Add /api if not present
	if !strings.HasSuffix(baseURL, "/api") {
		baseURL += "/api"
	}

	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetInterfaces fetches all interfaces from a YAPI project with pagination support
func (c *Client) GetInterfaces(projectID int) ([]model.YAPIInterface, error) {
	// Try the category-based approach first (more reliable for YAPI)
	interfaces, err := c.getInterfacesByCategory(projectID)
	if err == nil && len(interfaces) > 0 {
		return interfaces, nil
	}

	// Fallback to get_list with pagination
	var allInterfaces []model.YAPIInterface
	page := 1
	limit := 100 // YAPI default limit per page

	for {
		url := fmt.Sprintf("%s/interface/get_list?project_id=%d&token=%s&page=%d&limit=%d", c.BaseURL, projectID, c.Token, page, limit)

		resp, err := c.HTTPClient.Get(url)
		if err != nil {
			// If get_list also fails, return the error from category method
			return nil, fmt.Errorf("failed to fetch interfaces: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		var yapiResp struct {
			Errcode int    `json:"errcode"`
			ErrMsg  string `json:"errmsg"`
			Data    struct {
				List []model.YAPIInterface `json:"list"`
				Count int                  `json:"count"` // Total count
			} `json:"data"`
		}

		if err := json.Unmarshal(body, &yapiResp); err != nil {
			// Try category method as fallback
			return c.getInterfacesByCategory(projectID)
		}

		if yapiResp.Errcode != 0 {
			// Try alternative API endpoint if get_list fails
			return c.getInterfacesByCategory(projectID)
		}

		allInterfaces = append(allInterfaces, yapiResp.Data.List...)

		// Check if we've fetched all interfaces
		if len(yapiResp.Data.List) == 0 || len(allInterfaces) >= yapiResp.Data.Count {
			break
		}

		page++
	}

	return allInterfaces, nil
}

// getInterfacesByCategory fetches interfaces by category (alternative method)
func (c *Client) getInterfacesByCategory(projectID int) ([]model.YAPIInterface, error) {
	// First get all categories
	categories, err := c.GetCategories(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	var allInterfaces []model.YAPIInterface
	for _, cat := range categories {
		url := fmt.Sprintf("%s/interface/list_cat?catid=%d&token=%s", c.BaseURL, cat.ID, c.Token)

		resp, err := c.HTTPClient.Get(url)
		if err != nil {
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		// YAPI list_cat returns: { "errcode": 0, "data": { "list": [...], "count": N, "total": M } }
		var yapiResp struct {
			Errcode int `json:"errcode"`
			ErrMsg  string `json:"errmsg"`
			Data    struct {
				List []model.YAPIInterface `json:"list"`
				Count int                  `json:"count"`
			} `json:"data"`
		}

		if err := json.Unmarshal(body, &yapiResp); err != nil {
			continue
		}

		if yapiResp.Errcode == 0 && len(yapiResp.Data.List) > 0 {
			allInterfaces = append(allInterfaces, yapiResp.Data.List...)
		}
	}

	if len(allInterfaces) == 0 {
		return nil, fmt.Errorf("no interfaces found, please check your YAPI configuration and token permissions")
	}

	return allInterfaces, nil
}

// GetCategories fetches all categories from a YAPI project
func (c *Client) GetCategories(projectID int) ([]model.YAPICategory, error) {
	url := fmt.Sprintf("%s/interface/getCatMenu?project_id=%d&token=%s", c.BaseURL, projectID, c.Token)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var yapiResp struct {
		Errcode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		Data    []model.YAPICategory `json:"data"`
	}

	if err := json.Unmarshal(body, &yapiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if yapiResp.Errcode != 0 {
		return nil, fmt.Errorf("YAPI error: %s", yapiResp.ErrMsg)
	}

	return yapiResp.Data, nil
}

// GetProject fetches YAPI project information
func (c *Client) GetProject(projectID int) (*model.YAPIProject, error) {
	url := fmt.Sprintf("%s/project/get?project_id=%d&token=%s", c.BaseURL, projectID, c.Token)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var yapiResp struct {
		Errcode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		Data    *model.YAPIProject `json:"data"`
	}

	if err := json.Unmarshal(body, &yapiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if yapiResp.Errcode != 0 {
		return nil, fmt.Errorf("YAPI error: %s", yapiResp.ErrMsg)
	}

	return yapiResp.Data, nil
}

// TestConnection tests the YAPI connection
func (c *Client) TestConnection(projectID int) error {
	_, err := c.GetProject(projectID)
	return err
}

// ConvertYAPIParamToAPIParam converts YAPI param to internal APIParam
func ConvertYAPIParamToAPIParam(yParam model.YAPIParam, paramType string) model.APIParam {
	required := yParam.Required == "1"

	return model.APIParam{
		Name:        yParam.Name,
		Type:        yParam.Type,
		ParamType:   paramType,
		Required:    required,
		Description: yParam.Desc,
	}
}

// ConvertYAPISchemaToAPIParams converts YAPI schema to APIParams with nested support
func ConvertYAPISchemaToAPIParams(schema []model.YAPIParam) []model.APIParam {
	var params []model.APIParam
	for _, s := range schema {
		param := model.APIParam{
			Name:        s.Name,
			Type:        s.Type,
			ParamType:   "body",
			Required:    s.Required == "1",
			Description: s.Desc,
		}
		params = append(params, param)
	}
	return params
}

// ParseYAPIResponse attempts to parse YAPI response body to extract field information
func ParseYAPIResponse(resBody string) ([]model.APIParam, *model.APIBody) {
	// Try to parse as JSON to extract schema
	var jsonData interface{}
	if err := json.Unmarshal([]byte(resBody), &jsonData); err != nil {
		return nil, nil
	}

	// Extract schema from parsed JSON
	schema := extractSchemaFromJSON(jsonData)
	body := &model.APIBody{
		Type:       "json",
		DataType:   guessDataType(jsonData),
		Schema:     schema,
		JSONSchema: resBody,
	}

	return schema, body
}

func extractSchemaFromJSON(data interface{}) []model.APIParam {
	var params []model.APIParam

	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			param := model.APIParam{
				Name:        key,
				Type:        guessType(val),
				ParamType:   "body",
				Description: "",
			}
			// Recursively extract nested object properties
			if nestedObj, ok := val.(map[string]interface{}); ok {
				param.Children = extractSchemaFromJSON(nestedObj)
			}
			params = append(params, param)
		}
	case []interface{}:
		if len(v) > 0 {
			// For arrays, use first element as schema reference
			params = extractSchemaFromJSON(v[0])
		}
	}

	return params
}

func guessType(val interface{}) string {
	if val == nil {
		return "string"
	}
	switch val.(type) {
	case bool:
		return "boolean"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return "number"
	case map[string]interface{}:
		return "object"
	case []interface{}:
		return "array"
	default:
		return "string"
	}
}

func guessDataType(data interface{}) string {
	switch data.(type) {
	case map[string]interface{}:
		return "object"
	case []interface{}:
		return "array"
	default:
		return "string"
	}
}

// ConvertYAPIInterfaceToEndpoint converts YAPI interface to internal Endpoint model
func ConvertYAPIInterfaceToEndpoint(yInt model.YAPIInterface) (*model.Endpoint, error) {
	endpoint := &model.Endpoint{
		Name:        yInt.Title,
		Method:      yInt.Method,
		URL:         yInt.Path,
		Description: yInt.Description,
	}

	// Convert request parameters (path and query params)
	var reqParams []model.APIParam
	for _, p := range yInt.ReqParams {
		paramType := "query"
		if strings.Contains(yInt.Path, ":"+p.Name) || strings.Contains(yInt.Path, "{"+p.Name+"}") {
			paramType = "path"
		}
		reqParams = append(reqParams, ConvertYAPIParamToAPIParam(p, paramType))
	}
	endpoint.RequestParams = reqParams

	// Convert headers
	headers := make(map[string]string)
	for _, h := range yInt.ReqHeaders {
		headers[h.Name] = h.Value
	}
	endpoint.Headers = headers

	// Convert request body
	if yInt.ReqBodyOther.Type != "" && len(yInt.ReqBodyOther.Schema) > 0 {
		reqBody := &model.APIBody{
			Type:     yInt.ReqBodyOther.Type,
			DataType: "object",
			Schema:   ConvertYAPISchemaToAPIParams(yInt.ReqBodyOther.Schema),
		}
		endpoint.RequestBody = reqBody
	}

	// Parse response body
	if yInt.ResBody != "" {
		respParams, respBody := ParseYAPIResponse(yInt.ResBody)
		endpoint.ResponseParams = respParams
		endpoint.ResponseBody = respBody
	}

	return endpoint, nil
}

// ValidateConfig validates YAPI sync configuration
func ValidateConfig(config model.YAPISyncConfig) error {
	if config.ProjectID == "" {
		return errors.New("project_id is required")
	}
	if config.YAPIUrl == "" {
		return errors.New("yapi_url is required")
	}
	if config.YAPIToken == "" {
		return errors.New("yapi_token is required")
	}
	if config.YAPIProjectID == 0 {
		return errors.New("yapi_project_id is required")
	}
	return nil
}
