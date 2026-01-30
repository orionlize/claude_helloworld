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
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at" db:"updated_at"`
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
