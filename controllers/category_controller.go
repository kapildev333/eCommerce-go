package controllers

import (
	"database/sql"
	"eCommerce-go/db"
	"eCommerce-go/middleware"
	"eCommerce-go/models"
	"eCommerce-go/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// getAllCategoriesHandler retrieves all product categories
func getAllCategoriesHandler(c *gin.Context) {
	_ = middleware.CheckUserExist(c)
	query := `
		SELECT id, name, description, parent_id, slug, image_url, 
		       meta_title, meta_description, is_active, updated_at, created_at 
		FROM categories 
		WHERE is_active = true
		ORDER BY name
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error retrieving categories",
		})
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		var description, imageURL, metaTitle, metaDescription sql.NullString
		var parentID sql.NullInt64

		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&description,
			&parentID,
			&category.Slug,
			&imageURL,
			&metaTitle,
			&metaDescription,
			&category.IsActive,
			&category.UpdatedAt,
			&category.CreatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				StatusCode: http.StatusInternalServerError,
				Error:      true,
				Message:    "Error scanning category",
			})
			return
		}

		// Convert nullable fields
		if description.Valid {
			category.Description = description.String
		}
		if parentID.Valid {
			parentIDInt := int(parentID.Int64)
			category.ParentID = &parentIDInt
		}
		if imageURL.Valid {
			category.ImageURL = imageURL.String
		}
		if metaTitle.Valid {
			category.MetaTitle = metaTitle.String
		}
		if metaDescription.Valid {
			category.MetaDescription = metaDescription.String
		}

		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, utils.Response{
		StatusCode: http.StatusOK,
		Error:      false,
		Message:    "Categories retrieved successfully",
		Data:       gin.H{"categories": categories},
	})
}

// getCategoryHandler retrieves a specific category by ID
func getCategoryHandler(c *gin.Context) {
	_ = middleware.CheckUserExist(c)
	categoryID := c.DefaultQuery("id", "")

	query := `
		SELECT id, name, description, parent_id, slug, image_url, 
		       meta_title, meta_description, is_active, updated_at, created_at 
		FROM categories 
		WHERE id = $1
	`
	row := db.DB.QueryRow(query, categoryID)

	var category models.Category
	var description, imageURL, metaTitle, metaDescription sql.NullString
	var parentID sql.NullInt64

	if err := row.Scan(
		&category.ID,
		&category.Name,
		&description,
		&parentID,
		&category.Slug,
		&imageURL,
		&metaTitle,
		&metaDescription,
		&category.IsActive,
		&category.UpdatedAt,
		&category.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, utils.Response{
				StatusCode: http.StatusNotFound,
				Error:      true,
				Message:    "Category not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, utils.Response{
				StatusCode: http.StatusInternalServerError,
				Error:      true,
				Message:    "Error retrieving category",
			})
		}
		return
	}

	// Convert nullable fields
	if description.Valid {
		category.Description = description.String
	}
	if parentID.Valid {
		parentIDInt := int(parentID.Int64)
		category.ParentID = &parentIDInt
	}
	if imageURL.Valid {
		category.ImageURL = imageURL.String
	}
	if metaTitle.Valid {
		category.MetaTitle = metaTitle.String
	}
	if metaDescription.Valid {
		category.MetaDescription = metaDescription.String
	}

	c.JSON(http.StatusOK, utils.Response{
		StatusCode: http.StatusOK,
		Error:      false,
		Message:    "Category retrieved successfully",
		Data:       gin.H{"category": category},
	})
}

// createCategoryHandler creates a new category
func createCategoryHandler(c *gin.Context) {
	_ = middleware.CheckUserExist(c)
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Error:      true,
			Message:    "Invalid input",
		})
		return
	}

	query := `
		INSERT INTO categories (name, description, parent_id, slug, image_url, 
						   meta_title, meta_description, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := db.DB.QueryRow(
		query,
		category.Name,
		category.Description,
		category.ParentID,
		category.Slug,
		category.ImageURL,
		category.MetaTitle,
		category.MetaDescription,
		category.IsActive,
	).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error creating category: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, utils.Response{
		StatusCode: http.StatusCreated,
		Error:      false,
		Message:    "Category created successfully",
		Data:       gin.H{"category": category},
	})
}

// updateCategoryHandler updates an existing category
func updateCategoryHandler(c *gin.Context) {
	_ = middleware.CheckUserExist(c)
	categoryID := c.DefaultQuery("id", "")
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Error:      true,
			Message:    "Invalid input",
		})
		return
	}

	query := `
		UPDATE categories
		SET name = $1, description = $2, parent_id = $3, slug = $4, 
			image_url = $5, meta_title = $6, meta_description = $7, 
			is_active = $8, updated_at = NOW()
		WHERE id = $9
		RETURNING updated_at
	`

	err := db.DB.QueryRow(
		query,
		category.Name,
		category.Description,
		category.ParentID,
		category.Slug,
		category.ImageURL,
		category.MetaTitle,
		category.MetaDescription,
		category.IsActive,
		categoryID,
	).Scan(&category.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, utils.Response{
				StatusCode: http.StatusNotFound,
				Error:      true,
				Message:    "Category not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, utils.Response{
				StatusCode: http.StatusInternalServerError,
				Error:      true,
				Message:    "Error updating category: " + err.Error(),
			})
		}
		return
	}

	category.ID, _ = strconv.Atoi(categoryID)

	c.JSON(http.StatusOK, utils.Response{
		StatusCode: http.StatusOK,
		Error:      false,
		Message:    "Category updated successfully",
		Data:       gin.H{"category": category},
	})
}

// deleteCategoryHandler deletes a category
func deleteCategoryHandler(c *gin.Context) {
	_ = middleware.CheckUserExist(c)
	categoryID := c.DefaultQuery("id", "")

	// First check if category has children
	var childCount int
	childQuery := "SELECT COUNT(*) FROM categories WHERE parent_id = $1"
	if err := db.DB.QueryRow(childQuery, categoryID).Scan(&childCount); err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error checking category children",
		})
		return
	}

	if childCount > 0 {
		c.JSON(http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Error:      true,
			Message:    "Cannot delete category with subcategories",
		})
		return
	}

	// Then check if products are associated with this category
	var productCount int
	productQuery := "SELECT COUNT(*) FROM product_categories WHERE category_id = $1"
	if err := db.DB.QueryRow(productQuery, categoryID).Scan(&productCount); err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error checking category products",
		})
		return
	}

	if productCount > 0 {
		c.JSON(http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Error:      true,
			Message:    "Cannot delete category with associated products",
		})
		return
	}

	// If no children or products, proceed with deletion
	query := "DELETE FROM categories WHERE id = $1"
	result, err := db.DB.Exec(query, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error deleting category",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, utils.Response{
			StatusCode: http.StatusNotFound,
			Error:      true,
			Message:    "Category not found",
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		StatusCode: http.StatusOK,
		Error:      false,
		Message:    "Category deleted successfully",
	})
}

// ConfigCategoryController sets up category routes
func ConfigCategoryController(group *gin.RouterGroup) {
	categoryRoutes := group.Group("categories")
	categoryRoutes.GET("/", getAllCategoriesHandler)
	categoryRoutes.GET("/getCategory", getCategoryHandler)
	categoryRoutes.POST("/addCategory", createCategoryHandler)
	categoryRoutes.PUT("/updateCategory", updateCategoryHandler)
	categoryRoutes.DELETE("/deleteCategory", deleteCategoryHandler)
}
