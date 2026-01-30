package handler

import (
	"context"
	"fmt"
	"time"

	"apihub/internal/config"
	"apihub/internal/model"
	"apihub/pkg/auth"
	"apihub/pkg/response"
	"apihub/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func New(db *pgxpool.Pool, cfg *config.Config) *Handler {
	return &Handler{db: db, cfg: cfg}
}

func (h *Handler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, 500, "Failed to hash password")
		return
	}

	// Generate UUID
	userID := generateUUID()

	// Insert user
	ctx := context.Background()
	query := `
		INSERT INTO users (id, username, email, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, created_at, updated_at
	`

	var user model.User
	err = h.db.QueryRow(ctx, query, userID, req.Username, req.Email, string(hashedPassword)).Scan(
		&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		logger.Error("Failed to create user: " + err.Error())
		response.Error(c, 500, "Failed to create user")
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Username, h.cfg.JWT.Secret, h.cfg.JWT.Expiration)
	if err != nil {
		response.Error(c, 500, "Failed to generate token")
		return
	}

	response.Success(c, 201, model.AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	ctx := context.Background()
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users WHERE email = $1
	`

	var user model.User
	err := h.db.QueryRow(ctx, query, req.Email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			response.Error(c, 401, "Invalid credentials")
		} else {
			response.Error(c, 500, "Database error")
		}
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		response.Error(c, 401, "Invalid credentials")
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Username, h.cfg.JWT.Secret, h.cfg.JWT.Expiration)
	if err != nil {
		response.Error(c, 500, "Failed to generate token")
		return
	}

	response.Success(c, 200, model.AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 401, "Unauthorized")
		return
	}

	ctx := context.Background()
	query := `
		SELECT id, username, email, created_at, updated_at
		FROM users WHERE id = $1
	`

	var user model.User
	err := h.db.QueryRow(ctx, query, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		response.Error(c, 404, "User not found")
		return
	}

	// Generate new token
	token, err := auth.GenerateToken(user.ID, user.Username, h.cfg.JWT.Secret, h.cfg.JWT.Expiration)
	if err != nil {
		response.Error(c, 500, "Failed to generate token")
		return
	}

	response.Success(c, 200, model.AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 401, "Unauthorized")
		return
	}

	ctx := context.Background()
	query := `
		SELECT id, username, email, created_at, updated_at
		FROM users WHERE id = $1
	`

	var user model.User
	err := h.db.QueryRow(ctx, query, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		response.Error(c, 404, "User not found")
		return
	}

	response.Success(c, 200, user)
}

func (h *Handler) UpdateCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 401, "Unauthorized")
		return
	}

	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Email    string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	ctx := context.Background()
	query := `
		UPDATE users
		SET username = $1, email = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, username, email, created_at, updated_at
	`

	var user model.User
	err := h.db.QueryRow(ctx, query, req.Username, req.Email, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		response.Error(c, 500, "Failed to update user")
		return
	}

	response.Success(c, 200, user)
}

func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
