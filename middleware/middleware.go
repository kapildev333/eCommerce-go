package middleware

import (
	config "eCommerce-go/utils"
	logger "eCommerce-go/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

const (
	AuthorizationHeaderKey  = "Authorization"
	AuthorizationTypeBearer = "Bearer"
	AuthorizationPayloadKey = "authorization_payload" // Key for storing claims in context
	UserIDKey               = "user_id"               // Key for storing userID in context
)

// AuthMiddleware creates a Gin middleware for JWT authentication
func AuthMiddleware() gin.HandlerFunc {
	log := logger.GetLogger()
	log.With("component", "auth_middleware")
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeaderKey)
		if len(authHeader) == 0 {
			err := errors.New("authorization header is not provided")
			log.Warn("Auth middleware failed", "path", c.Request.URL.Path, "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			log.Warn("Auth middleware failed", "path", c.Request.URL.Path, "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != strings.ToLower(AuthorizationTypeBearer) {
			err := fmt.Errorf("unsupported authorization type %s", authType)
			log.Warn("Auth middleware failed", "path", c.Request.URL.Path, "type", authType, "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		accessToken := fields[1]
		token, err := TokenServiceInstance.ValidateToken(accessToken, config.Env.AccessTokenSecret)
		if err != nil {
			log.Warn("Auth middleware failed: Invalid token", "path", c.Request.URL.Path, "error", err)
			// Check for specific JWT errors like expiry
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has expired"})
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			}
			return
		}

		// Extract claims (we stored Claims struct during creation)
		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			log.Error("Auth middleware failed: Invalid token claims", "path", c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		// Set user information in the context for downstream handlers
		c.Set(AuthorizationPayloadKey, claims) // Store full claims if needed
		c.Set(UserIDKey, claims.UserID)        // Store UserID directly for convenience

		log.Debug("Auth middleware success", "path", c.Request.URL.Path, "userID", claims.UserID)
		c.Next() // Proceed to the next handler
	}
}
