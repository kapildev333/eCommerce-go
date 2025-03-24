package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"time"
)

var secretKey = []byte(os.Getenv("SECRET_KEY"))

func GenerateToken(userID uint, username, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userID,
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return token, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: http.StatusUnauthorized,
				Message:    "Authorization token not provided",
			})
			c.Abort()
			return
		}

		token, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: http.StatusUnauthorized,
				Error:      true,
				Message:    "Invalid or expired token",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: http.StatusUnauthorized,
				Error:      true,
				Message:    "Invalid token claims",
			})
			c.Abort()
			return
		}

		// Check if the token is expired
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				// Generate a new token if expired
				newToken, err := GenerateToken(uint(claims["id"].(float64)), claims["username"].(string), claims["email"].(string))
				if err != nil {
					c.JSON(http.StatusInternalServerError, Response{
						StatusCode: http.StatusInternalServerError,
						Error:      true,
						Message:    "Error generating new token",
					})
					c.Abort()
					return
				}
				c.Header("Authorization", newToken)
			}
		}

		// Set user information in context
		c.Set("userID", claims["id"])
		c.Set("username", claims["username"])
		c.Set("email", claims["email"])

		c.Next()
	}
}
