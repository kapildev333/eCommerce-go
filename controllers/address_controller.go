package controllers

import (
	"eCommerce-go/db"
	"eCommerce-go/models"
	"eCommerce-go/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getAllShippingAddresses(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.Response{
			StatusCode: http.StatusUnauthorized,
			Error:      true,
			Message:    "Unauthorized",
		})
		return
	}

	var addresses []models.ShippingAddress
	query := "SELECT id, user_id, address_line_1, address_line_2, city, state, postal_code, country, is_default, updated_at, created_at FROM shipping_addresses WHERE user_id = $1"
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error retrieving addresses",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var address models.ShippingAddress
		if err := rows.Scan(&address.ID, &address.UserID, &address.AddressLine1, &address.AddressLine2, &address.City, &address.State, &address.PostalCode, &address.Country, &address.IsDefault, &address.UpdatedAt, &address.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				StatusCode: http.StatusInternalServerError,
				Error:      true,
				Message:    "Error scanning address",
			})
			return
		}
		addresses = append(addresses, address)
	}

	c.JSON(http.StatusOK, utils.Response{
		StatusCode: http.StatusOK,
		Error:      false,
		Message:    "Addresses retrieved",
		Data:       gin.H{"addresses": addresses},
	})
}

func addAddress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.Response{
			StatusCode: http.StatusUnauthorized,
			Error:      true,
			Message:    "Unauthorized",
		})
		return
	}
	var address models.ShippingAddress
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Error:      true,
			Message:    "Invalid input",
		})
		return
	}
	query := "INSERT INTO shipping_addresses (user_id, address_line_1, address_line_2, city, state, postal_code, country, is_default) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	err := db.DB.QueryRow(query, userID, address.AddressLine1, address.AddressLine2, address.City, address.State, address.PostalCode, address.Country, address.IsDefault).Scan(&address.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error inserting address",
		})
		return
	}
	c.JSON(http.StatusCreated, utils.Response{
		StatusCode: http.StatusCreated,
		Error:      false,
		Message:    "Address created",
		Data:       gin.H{"address": address},
	})

}
func ConfigAddressController(group *gin.RouterGroup) {
	accounts := group.Group("address")
	accounts.GET("/getAllAddress", getAllShippingAddresses)
	accounts.POST("/addAddress", addAddress)
}
