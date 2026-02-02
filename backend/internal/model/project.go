package model

import "time"

type Project struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	UserID      string    `json:"user_id" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=500"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=500"`
}

type Collection struct {
	ID          string    `json:"id" db:"id"`
	ProjectID   string    `json:"project_id" db:"project_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	ParentID    *string   `json:"parent_id" db:"parent_id"`
	SortOrder   int       `json:"sort_order" db:"sort_order"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateCollectionRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	Description string  `json:"description" binding:"max=500"`
	ParentID    *string `json:"parent_id"`
}

type UpdateCollectionRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	Description string  `json:"description" binding:"max=500"`
	ParentID    *string `json:"parent_id"`
}

// Endpoint represents an API endpoint with detailed field information
type Endpoint struct {
	ID           string            `json:"id" db:"id"`
	CollectionID string            `json:"collection_id" db:"collection_id"`
	Name         string            `json:"name" db:"name"`
	Method       string            `json:"method" db:"method"`
	URL          string            `json:"url" db:"url"`
	Headers      map[string]string `json:"headers" db:"headers"`
	Body         *string           `json:"body" db:"body"`
	Description  string            `json:"description" db:"description"`
	SortOrder    int               `json:"sort_order" db:"sort_order"`
	// Detailed API field information
	RequestParams  []APIParam `json:"request_params" db:"request_params"`
	RequestBody    *APIBody   `json:"request_body" db:"request_body"`
	ResponseParams []APIParam `json:"response_params" db:"response_params"`
	ResponseBody   *APIBody   `json:"response_body" db:"response_body"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// APIParam represents a parameter in request or response
type APIParam struct {
	Name        string      `json:"name" db:"name"`
	Type        string      `json:"type" db:"type"`          // string, number, boolean, object, array
	ParamType   string      `json:"param_type" db:"param_type"` // path, query, header, body
	Required    bool        `json:"required" db:"required"`
	Description string      `json:"description" db:"description"`
	DefaultValue interface{} `json:"default_value" db:"default_value"`
	Children    []APIParam  `json:"children" db:"children"` // For nested object properties
}

// APIBody represents request/response body structure
type APIBody struct {
	Type       string      `json:"type" db:"type"`           // json, form-data, raw, xml
	DataType   string      `json:"data_type" db:"data_type"` // object, array, string, etc.
	Schema     []APIParam  `json:"schema" db:"schema"`       // Body field definitions
	Example    interface{} `json:"example" db:"example"`
	JSONSchema string      `json:"json_schema" db:"json_schema"` // JSON Schema string
}

type CreateEndpointRequest struct {
	Name         string            `json:"name" binding:"required,min=1,max=100"`
	Method       string            `json:"method" binding:"required,oneof=GET POST PUT PATCH DELETE OPTIONS HEAD"`
	URL          string            `json:"url" binding:"required"`
	Headers      map[string]string `json:"headers"`
	Body         *string           `json:"body"`
	Description  string            `json:"description" binding:"max=1000"`
}

type UpdateEndpointRequest struct {
	Name         string            `json:"name" binding:"required,min=1,max=100"`
	Method       string            `json:"method" binding:"required,oneof=GET POST PUT PATCH DELETE OPTIONS HEAD"`
	URL          string            `json:"url" binding:"required"`
	Headers      map[string]string `json:"headers"`
	Body         *string           `json:"body"`
	Description  string            `json:"description" binding:"max=1000"`
}

type Environment struct {
	ID          string            `json:"id" db:"id"`
	ProjectID   string            `json:"project_id" db:"project_id"`
	Name        string            `json:"name" db:"name"`
	Variables   map[string]string `json:"variables" db:"variables"`
	IsDefault   bool              `json:"is_default" db:"is_default"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`
}

type CreateEnvironmentRequest struct {
	Name      string            `json:"name" binding:"required,min=1,max=50"`
	Variables map[string]string `json:"variables"`
	IsDefault bool              `json:"is_default"`
}

type UpdateEnvironmentRequest struct {
	Name      string            `json:"name" binding:"required,min=1,max=50"`
	Variables map[string]string `json:"variables"`
	IsDefault bool              `json:"is_default"`
}

// YAPI related types

type YAPIProject struct {
	ID          int    `json:"_id"`
	Name        string `json:"name"`
	BaseURL     string `json:"baseurl"`
	UID         int    `json:"uid"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	ContentType string `json:"contentType"`
}

type YAPIInterface struct {
	ID          int                  `json:"_id"`
	ProjectID   int                  `json:"project_id"`
	CatID       int                  `json:"catid"`
	Path        string               `json:"path"`
	Title       string               `json:"title"`
	Method      string               `json:"method"`
	Status      string               `json:"status"` // undone, done
	ReqParams   []YAPIParam          `json:"req_params"`
	ReqBodyOther YAPIRequestBodyOther  `json:"req_body_other"`
	ResBody     string               `json:"res_body"`
	ReqBodyType  string               `json:"req_body_type"`
	ReqHeaders   []YAPIHeader         `json:"req_headers"`
	Description string               `json:"desc"`
	UpdatedAt   int64                `json:"up_time"`
}

type YAPIParam struct {
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Required   string `json:"required"` // "1" for required
	Type       string `json:"type"`     // string, number, boolean, etc.
}

type YAPIRequestBodyOther struct {
	Type    string        `json:"type"` // json, form, text
	Schema  []YAPIParam   `json:"schema"`
}

type YAPIHeader struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Desc   string `json:"desc"`
	Required string `json:"required"`
}

type YAPICategory struct {
	ID      int    `json:"_id"`
	Name    string `json:"name"`
	ProjectID int  `json:"project_id"`
	Order   int    `json:"order"`
}

type YAPIResponse struct {
	Errcode int         `json:"errcode"`
	Data    interface{} `json:"data"`
	ErrMsg  string      `json:"errmsg"`
}

type YAPISyncConfig struct {
	ProjectID    string `json:"project_id" binding:"required"`
	YAPIUrl      string `json:"yapi_url" binding:"required"`      // e.g., http://yapi.example.com
	YAPIToken    string `json:"yapi_token" binding:"required"`    // YAPI project token
	YAPIProjectID int   `json:"yapi_project_id" binding:"required"` // YAPI project ID
}
