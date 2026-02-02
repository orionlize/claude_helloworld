package handler

import (
	"time"

	"apihub/internal/config"
	"apihub/internal/model"
	"apihub/pkg/auth"
	"apihub/pkg/response"
	"apihub/pkg/logger"
	"apihub/pkg/store"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db    *pgxpool.Pool
	cfg   *config.Config
	store store.Store
}

func New(db *pgxpool.Pool, cfg *config.Config, store store.Store) *Handler {
	return &Handler{db: db, cfg: cfg, store: store}
}

func (h *Handler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	// Check if user already exists
	existingUser, _ := h.store.GetUserByEmail(req.Email)
	if existingUser != nil {
		response.Error(c, 400, "User already exists")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, 500, "Failed to hash password")
		return
	}

	// Create user
	user := &model.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Username: req.Username,
	}

	if err := h.store.CreateUser(user); err != nil {
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

	// Clear password before returning
	user.Password = ""

	response.Success(c, 201, model.AuthResponse{
		Token: token,
		User:  *user,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	// Get user by email
	user, err := h.store.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		response.Error(c, 401, "Invalid credentials")
		return
	}

	// Verify password using bcrypt only
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Error(c, 401, "Invalid credentials")
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Username, h.cfg.JWT.Secret, h.cfg.JWT.Expiration)
	if err != nil {
		response.Error(c, 500, "Failed to generate token")
		return
	}

	// Clear password before returning
	user.Password = ""

	response.Success(c, 200, model.AuthResponse{
		Token: token,
		User:  *user,
	})
}

func (h *Handler) GetCurrentUser(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 401, "Unauthorized")
		return
	}

	userID, ok := userIDStr.(string)
	if !ok {
		response.Error(c, 401, "Invalid user ID type")
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil || user == nil {
		response.Error(c, 404, "User not found")
		return
	}

	// Clear password before returning
	user.Password = ""

	response.Success(c, 200, user)
}

func generateUUID() string {
	return time.Now().Format("20060102150405") + randomString(6)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
