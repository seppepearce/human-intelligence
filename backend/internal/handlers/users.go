package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"human-intelligence/internal/database"
	"human-intelligence/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	db *database.DB
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Bio      string `json:"bio"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	Username string  `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	Bio      string  `json:"bio,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
}

type AuthResponse struct {
	User         *models.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
}

func NewUserHandler(db *database.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Check if user already exists
	var existingUser models.User
	query := "SELECT id FROM users WHERE email = $1 OR username = $2"
	err := h.db.QueryRow(query, req.Email, req.Username).Scan(&existingUser.ID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this email or username already exists",
		})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create user
	user := models.User{
		ID:        uuid.New(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Bio:       req.Bio,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  true,
	}

	query = `
		INSERT INTO users (id, username, email, password_hash, bio, created_at, updated_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, username, email, bio, avatar, created_at, updated_at, is_active
	`

	err = h.db.QueryRow(query, user.ID, user.Username, user.Email, user.Password,
		user.Bio, user.CreatedAt, user.UpdatedAt, user.IsActive).Scan(
		&user.ID, &user.Username, &user.Email, &user.Bio, &user.Avatar,
		&user.CreatedAt, &user.UpdatedAt, &user.IsActive)

	if err != nil {
		log.Printf("Failed to create user - DB error: %v", err)
		log.Printf("Query: %s", query)
		log.Printf("Values: ID=%v, Username=%s, Email=%s, Bio=%s", user.ID, user.Username, user.Email, user.Bio)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Generate tokens
	accessToken, refreshToken, expiresAt, err := h.generateTokens(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate tokens",
		})
		return
	}

	// Clear password before returning
	user.Password = ""

	c.JSON(http.StatusCreated, AuthResponse{
		User:         &user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Find user
	var user models.User
	query := `
		SELECT id, username, email, password_hash, bio, avatar, created_at, updated_at, is_active
		FROM users
		WHERE email = $1 AND is_active = true
	`

	err := h.db.QueryRow(query, req.Email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.Bio,
		&user.Avatar, &user.CreatedAt, &user.UpdatedAt, &user.IsActive)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Generate tokens
	accessToken, refreshToken, expiresAt, err := h.generateTokens(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate tokens",
		})
		return
	}

	// Clear password before returning
	user.Password = ""

	c.JSON(http.StatusOK, AuthResponse{
		User:         &user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	})
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	// Parse and validate refresh token
	token, err := jwt.ParseWithClaims(req.RefreshToken, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(getJWTSecret()), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid refresh token",
		})
		return
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token claims",
		})
		return
	}

	// Get user from database
	var user models.User
	query := `
		SELECT id, username, email, bio, avatar, created_at, updated_at, is_active
		FROM users
		WHERE id = $1 AND is_active = true
	`

	err = h.db.QueryRow(query, claims.UserID).Scan(
		&user.ID, &user.Username, &user.Email, &user.Bio,
		&user.Avatar, &user.CreatedAt, &user.UpdatedAt, &user.IsActive)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	// Generate new tokens
	accessToken, refreshToken, expiresAt, err := h.generateTokens(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate tokens",
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		User:         &user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	user, err := h.getUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Username != "" {
		// Check if username is already taken
		var existingID uuid.UUID
		checkQuery := "SELECT id FROM users WHERE username = $1 AND id != $2"
		err := h.db.QueryRow(checkQuery, req.Username, userID).Scan(&existingID)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Username already taken",
			})
			return
		} else if err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Database error",
			})
			return
		}

		setParts = append(setParts, "username = $"+strconv.Itoa(argIndex))
		args = append(args, req.Username)
		argIndex++
	}

	if req.Bio != "" {
		setParts = append(setParts, "bio = $"+strconv.Itoa(argIndex))
		args = append(args, req.Bio)
		argIndex++
	}

	if req.Avatar != nil && *req.Avatar != "" {
		setParts = append(setParts, "avatar = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Avatar)
		argIndex++
	}

	if len(setParts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No fields to update",
		})
		return
	}

	// Add updated_at
	setParts = append(setParts, "updated_at = $"+strconv.Itoa(argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add WHERE clause
	args = append(args, userID)

	query := "UPDATE users SET " + strings.Join(setParts, ", ") + " WHERE id = $" + strconv.Itoa(argIndex)

	_, err := h.db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update profile",
		})
		return
	}

	// Return updated user
	user, err := h.getUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch updated user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := h.getUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// Helper methods
func (h *UserHandler) generateTokens(user *models.User) (string, string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)            // 24 hours for access token
	refreshExpiresAt := now.Add(7 * 24 * time.Hour) // 7 days for refresh token

	// Access token claims
	accessClaims := &models.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.ID.String(),
		},
	}

	// Refresh token claims
	refreshClaims := &models.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.ID.String(),
		},
	}

	// Generate tokens
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessTokenString, err := accessToken.SignedString([]byte(getJWTSecret()))
	if err != nil {
		return "", "", time.Time{}, err
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(getJWTSecret()))
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessTokenString, refreshTokenString, expiresAt, nil
}

func (h *UserHandler) getUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, username, email, bio, avatar, created_at, updated_at, is_active
		FROM users
		WHERE id = $1 AND is_active = true
	`

	err := h.db.QueryRow(query, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.Bio,
		&user.Avatar, &user.CreatedAt, &user.UpdatedAt, &user.IsActive)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (h *UserHandler) getUserIDFromContext(c *gin.Context) uuid.UUID {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil
	}

	if id, ok := userID.(uuid.UUID); ok {
		return id
	}

	return uuid.Nil
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Use a default secret for development
		return "your_super_secret_jwt_key_change_this_in_production"
	}
	return secret
}
