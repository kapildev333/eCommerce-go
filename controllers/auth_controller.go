package controllers

import (
	"database/sql"
	"eCommerce-go/db"
	"eCommerce-go/models"
	utils "eCommerce-go/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func signUpHandler(c *gin.Context) {
	var user models.User
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

	token, err := utils.GenerateToken(user.ID, user.Username, user.Email)
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
		Data:       gin.H{"user": user, "token": token},
	})
}

func signInToAccount(c *gin.Context) {
	var user models.User
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

	token, err := utils.GenerateToken(userDB.ID, userDB.Username, userDB.Email)
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
		Data:       gin.H{"user": userDB, "token": token},
	})
}
func ConfigAuthController(group *gin.RouterGroup) {
	accounts := group.Group("account")
	accounts.POST("/createAccount", signUpHandler)
	accounts.POST("/signing", signInToAccount)
}
