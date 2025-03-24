package controllers

import (
	"eCommerce-go/db"
	"eCommerce-go/models"
	"eCommerce-go/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func submitPaymentHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.Response{
			StatusCode: http.StatusUnauthorized,
			Error:      true,
			Message:    "Unauthorized",
		})
		return
	}
	var payment models.UserPayment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Error:      true,
			Message:    "Invalid input",
		})
		return
	}

	query := "INSERT INTO user_payments (user_id, payment_method, transaction_id, amount, currency, payment_status, billing_address_id, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	err := db.DB.QueryRow(query, userID, payment.PaymentMethod, payment.TransactionID, payment.Amount, payment.Currency, payment.PaymentStatus, payment.BillingAddressID, payment.Description).Scan(&payment.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error inserting payment",
		})
		return
	}
	c.JSON(http.StatusCreated, utils.Response{
		StatusCode: http.StatusCreated,
		Error:      false,
		Message:    "Payment submitted",
		Data:       gin.H{"payment": payment},
	})
}
func getPaymentHistoryHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.Response{
			StatusCode: http.StatusUnauthorized,
			Error:      true,
			Message:    "Unauthorized",
		})
		return
	}
	var paymentHistory []models.UserPayment
	query := "SELECT *FROM user_payments WHERE user_id = $1"
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error retrieving user's payment history",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var payment models.UserPayment
		if err := rows.Scan(&payment.ID, &payment.UserID, &payment.PaymentMethod, &payment.TransactionID, &payment.Amount, &payment.Currency, &payment.PaymentStatus, &payment.PaymentDate, &payment.BillingAddressID, &payment.Description, &payment.UpdatedAt, &payment.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				StatusCode: http.StatusInternalServerError,
				Error:      true,
				Message:    "Error scanning payment",
			})
			return
		}
		paymentHistory = append(paymentHistory, payment)
	}
	c.JSON(http.StatusOK, utils.Response{
		StatusCode: http.StatusOK,
		Error:      false,
		Message:    "History of user payments fetched successfully",
		Data:       gin.H{"payments": paymentHistory},
	})
}
func ConfigPaymentController(group *gin.RouterGroup) {
	accounts := group.Group("payments")
	accounts.POST("/submitPayment", submitPaymentHandler)
	accounts.GET("/getUserPaymentHistory", getPaymentHistoryHandler)
}
