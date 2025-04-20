package controllers

import (
	"context"
	"database/sql"
	"eCommerce-go/db"
	redis "eCommerce-go/db"
	"eCommerce-go/middleware"
	"eCommerce-go/models"
	"eCommerce-go/utils"
	config "eCommerce-go/utils"
	logger "eCommerce-go/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func signUpHandler(c *gin.Context) {
	var user models.User
	log := logger.GetLogger()
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Error:      true,
			Message:    "Invalid input",
		})
		return
	}

	if user.Email == "" || user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Error:      true,
			Message:    "Missing fields",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error hashing password",
		})
		return
	}
	user.Password = string(hashedPassword)

	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"
	err = db.DB.QueryRow(query, user.Username, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error saving user",
		})
		return
	}

	userID := strconv.FormatInt(int64(user.ID), 10)
	td, err := middleware.TokenServiceInstance.GenerateTokens(userID)
	if err != nil {
		log.Error("Signup failed: Could not generate tokens", "userID", userID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	err = redis.RedisStoreInstance.StoreRefreshToken(c, userID, td.RefreshUUID, config.Env.RefreshTokenLifespan)
	if err != nil {
		log.Error("Signup failed: Could not store refresh token", "refreshUUID", td.RefreshUUID, "userID", userID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	log.Info("Login successful, tokens generated", "userID", userID, "accessUUID", td.AccessUUID, "refreshUUID", td.RefreshUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error generating token",
		})
		return
	}

	c.JSON(http.StatusCreated, utils.Response{
		StatusCode: http.StatusCreated,
		Error:      false,
		Message:    "User created",
		Data:       td,
	})
}

func signInToAccount(c *gin.Context) {
	var user models.User
	log := logger.GetLogger()
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Error:      true,
			Message:    "Invalid input",
		})
		return
	}

	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Error:      true,
			Message:    "Missing fields",
		})
		return
	}

	var userDB models.User
	query := "SELECT id, username, email, password FROM users WHERE email = $1"
	err := db.DB.QueryRow(query, user.Email).Scan(&userDB.ID, &userDB.Username, &userDB.Email, &userDB.Password)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, utils.Response{
			StatusCode: http.StatusNotFound,
			Error:      true,
			Message:    "User not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error retrieving user",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, utils.Response{
			StatusCode: http.StatusUnauthorized,
			Error:      true,
			Message:    "Invalid password",
		})
		return
	}

	userID := strconv.FormatInt(int64(userDB.ID), 10)
	td, err := middleware.TokenServiceInstance.GenerateTokens(userID)
	if err != nil {
		log.Error("Login failed: Could not generate tokens", "userID", userID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	err = redis.RedisStoreInstance.StoreRefreshToken(c, userID, td.RefreshUUID, config.Env.RefreshTokenLifespan)
	if err != nil {
		log.Error("Login failed: Could not store refresh token", "refreshUUID", td.RefreshUUID, "userID", userID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	log.Info("Login successful, tokens generated", "userID", userID, "accessUUID", td.AccessUUID, "refreshUUID", td.RefreshUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error generating token",
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		StatusCode: http.StatusOK,
		Error:      false,
		Message:    "User signed in",
		Data:       gin.H{"user": userDB, "token": td},
	})
}

// RefreshTokenRequest - Structure for refresh token request body
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Refresh godoc
// @Summary Refresh access token
// @Description Provides a new access token using a valid refresh token.
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body RefreshTokenRequest true "Refresh Token"
// @Success 200 {object} map[string]string "New access token"
// @Failure 400 {object} gin.H "Invalid input"
// @Failure 401 {object} gin.H "Invalid or expired refresh token"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /refresh [post]
func refresh(c *gin.Context) {
	var req RefreshTokenRequest
	log := logger.GetLogger()
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("Token refresh failed: Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 1. Validate the refresh token structure and signature
	token, err := middleware.TokenServiceInstance.ValidateToken(req.RefreshToken, config.Env.RefreshTokenSecret)
	if err != nil {
		log.Warn("Token refresh failed: Invalid refresh token", "error", err)
		// Distinguish expired from other invalid errors if needed
		if errors.Is(err, jwt.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token has expired"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		}
		return
	}

	// 2. Extract claims (ensure it has RefreshUUID and UserID)
	refreshClaims, err := middleware.ExtractRefreshClaims(token)
	if err != nil {
		log.Warn("Token refresh failed: Could not extract refresh claims", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token claims"})
		return
	}

	// 3. Validate against Redis store
	ctx := context.Background() // Use request context
	storedUserID, err := redis.RedisStoreInstance.ValidateRefreshToken(ctx, refreshClaims.RefreshUUID)
	if err != nil {
		// Error already logged in ValidateRefreshToken
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	// Optional: Check if the UserID from the token matches the one stored in Redis
	if storedUserID != refreshClaims.UserID {
		log.Error("Token refresh failed: UserID mismatch", "tokenUserID", refreshClaims.UserID, "redisUserID", storedUserID, "refreshUUID", refreshClaims.RefreshUUID)
		// Security measure: If mismatch, invalidate the token in Redis immediately
		_ = redis.RedisStoreInstance.DeleteRefreshToken(ctx, refreshClaims.RefreshUUID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token (user mismatch)"})
		return
	}

	// --- Refresh Token Rotation (Optional but Recommended) ---
	// For better security, invalidate the old refresh token and issue a new one along with the new access token.
	// If you don't rotate, skip the deletion and generation of a new refresh token.

	// 4. Delete the old refresh token from Redis
	err = redis.RedisStoreInstance.DeleteRefreshToken(ctx, refreshClaims.RefreshUUID)
	if err != nil {
		// Log the error but might proceed if deletion failed, though it's not ideal
		log.Error("Token refresh: Failed to delete old refresh token, proceeding...", "refreshUUID", refreshClaims.RefreshUUID, "error", err)
		// Depending on policy, you might want to return an error here instead.
	}
	// --- End Rotation Step ---

	// 5. Generate *new* tokens (both access and potentially refresh if rotating)
	newTd, err := middleware.TokenServiceInstance.GenerateTokens(refreshClaims.UserID) // Generate for the validated user
	if err != nil {
		log.Error("Token refresh failed: Could not generate new tokens", "userID", refreshClaims.UserID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new tokens"})
		return
	}

	// 6. Store the *new* refresh token details in Redis (only if rotating refresh tokens)
	err = redis.RedisStoreInstance.StoreRefreshToken(ctx, refreshClaims.UserID, newTd.RefreshUUID, config.Env.RefreshTokenLifespan)
	if err != nil {
		// This is critical. If storing the new refresh token fails, the user might be locked out after the new access token expires.
		log.Error("Token refresh critical error: Could not store new refresh token", "newRefreshUUID", newTd.RefreshUUID, "userID", refreshClaims.UserID, "error", err)
		// You might want to try deleting the newly generated access token info or handle this state carefully.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store new refresh token state"})
		return
	}
	// --- End Rotation Step ---

	log.Info("Token refresh successful", "userID", refreshClaims.UserID, "oldRefreshUUID", refreshClaims.RefreshUUID, "newAccessUUID", newTd.AccessUUID, "newRefreshUUID", newTd.RefreshUUID)

	// 7. Send back the new tokens
	c.JSON(http.StatusOK, gin.H{
		"access_token":  newTd.AccessToken,
		"refresh_token": newTd.RefreshToken,
	})
}

func ConfigAuthController(group *gin.RouterGroup) {
	accounts := group.Group("account")
	accounts.POST("/createAccount", signUpHandler)
	accounts.POST("/signing", signInToAccount)
	accounts.POST("/refresh", refresh)
}
