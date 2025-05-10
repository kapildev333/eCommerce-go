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
)

func getAllProductsHandler(c *gin.Context) {
	_ = middleware.CheckUserExist(c)
	var products []models.Product
	query := `SELECT id, name, description, sku, price, compare_at_price, cost_price, slug, 
              weight, weight_unit, dimensions, is_taxable, tax_code, is_digital, 
              is_published, featured, meta_title, meta_description, updated_at, created_at 
              FROM products`

	// Handle optional filtering
	categoryID := c.Query("category_id")
	if categoryID != "" {
		query = `SELECT p.id, p.name, p.description, p.sku, p.price, p.compare_at_price, 
                 p.cost_price, p.slug, p.weight, p.weight_unit, p.dimensions, p.is_taxable, 
                 p.tax_code, p.is_digital, p.is_published, p.featured, p.meta_title, 
                 p.meta_description, p.updated_at, p.created_at 
                 FROM products p
                 JOIN product_categories pc ON p.id = pc.product_id
                 WHERE pc.category_id = $1`
	}

	var rows *sql.Rows
	var err error

	if categoryID != "" {
		rows, err = db.DB.Query(query, categoryID)
	} else {
		rows, err = db.DB.Query(query)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error retrieving products",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		var description, sku, weightUnit, metaTitle, metaDescription sql.NullString
		var compareAtPrice, costPrice, weight sql.NullFloat64
		var dimensions []byte

		if err := rows.Scan(
			&product.ID, &product.Name, &description, &sku, &product.Price,
			&compareAtPrice, &costPrice, &product.Slug, &weight, &weightUnit,
			&dimensions, &product.IsTaxable, &product.TaxCode, &product.IsDigital,
			&product.IsPublished, &product.Featured, &metaTitle, &metaDescription,
			&product.UpdatedAt, &product.CreatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				StatusCode: http.StatusInternalServerError,
				Error:      true,
				Message:    "Error scanning product",
			})
			return
		}

		// Convert nullable types
		if description.Valid {
			product.Description = description.String
		}
		if sku.Valid {
			product.SKU = sku.String
		}
		if compareAtPrice.Valid {
			product.CompareAtPrice = &compareAtPrice.Float64
		}
		if costPrice.Valid {
			product.CostPrice = &costPrice.Float64
		}
		if weight.Valid {
			product.Weight = &weight.Float64
		}
		if weightUnit.Valid {
			product.WeightUnit = weightUnit.String
		}
		if metaTitle.Valid {
			product.MetaTitle = metaTitle.String
		}
		if metaDescription.Valid {
			product.MetaDescription = metaDescription.String
		}

		// Load product images
		imagesQuery := `SELECT id, product_id, url, alt_text, position, created_at 
                       FROM product_images WHERE product_id = $1 ORDER BY position`
		imgRows, err := db.DB.Query(imagesQuery, product.ID)
		if err == nil {
			defer imgRows.Close()
			for imgRows.Next() {
				var img models.ProductImage
				var altText sql.NullString
				if err := imgRows.Scan(&img.ID, &img.ProductID, &img.URL, &altText,
					&img.Position, &img.CreatedAt); err == nil {
					if altText.Valid {
						img.AltText = altText.String
					}
					product.Images = append(product.Images, img)
				}
			}
		}

		// Load product categories
		categoriesQuery := `SELECT c.id, c.name, c.description, c.parent_id, c.slug, 
                           c.image_url, c.meta_title, c.meta_description, c.is_active, 
                           c.updated_at, c.created_at 
                           FROM categories c
                           JOIN product_categories pc ON c.id = pc.category_id
                           WHERE pc.product_id = $1`
		catRows, err := db.DB.Query(categoriesQuery, product.ID)
		if err == nil {
			defer catRows.Close()
			for catRows.Next() {
				var cat models.Category
				var catDesc, imageURL, catMetaTitle, catMetaDescription sql.NullString
				var parentID sql.NullInt64

				if err := catRows.Scan(
					&cat.ID, &cat.Name, &catDesc, &parentID, &cat.Slug,
					&imageURL, &catMetaTitle, &catMetaDescription, &cat.IsActive,
					&cat.UpdatedAt, &cat.CreatedAt,
				); err == nil {
					if catDesc.Valid {
						cat.Description = catDesc.String
					}
					if parentID.Valid {
						pid := int(parentID.Int64)
						cat.ParentID = &pid
					}
					if imageURL.Valid {
						cat.ImageURL = imageURL.String
					}
					if catMetaTitle.Valid {
						cat.MetaTitle = catMetaTitle.String
					}
					if catMetaDescription.Valid {
						cat.MetaDescription = catMetaDescription.String
					}
					product.Categories = append(product.Categories, cat)
				}
			}
		}

		products = append(products, product)
	}

	c.JSON(http.StatusOK, utils.Response{
		StatusCode: http.StatusOK,
		Error:      false,
		Message:    "Products retrieved",
		Data:       gin.H{"products": products},
	})
}

func getProductHandler(c *gin.Context) {
	_ = middleware.CheckUserExist(c)
	productID := c.DefaultQuery("id", "")

	// Get basic product info
	query := `SELECT id, name, description, sku, price, compare_at_price, cost_price, slug, 
              weight, weight_unit, dimensions, is_taxable, tax_code, is_digital, is_published, 
              featured, meta_title, meta_description, updated_at, created_at 
              FROM products WHERE id = $1`

	row := db.DB.QueryRow(query, productID)

	var product models.Product
	var description, sku, weightUnit, metaTitle, metaDescription sql.NullString
	var compareAtPrice, costPrice, weight sql.NullFloat64
	var dimensions []byte

	if err := row.Scan(
		&product.ID, &product.Name, &description, &sku, &product.Price,
		&compareAtPrice, &costPrice, &product.Slug, &weight, &weightUnit,
		&dimensions, &product.IsTaxable, &product.TaxCode, &product.IsDigital,
		&product.IsPublished, &product.Featured, &metaTitle, &metaDescription,
		&product.UpdatedAt, &product.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, utils.Response{
				StatusCode: http.StatusNotFound,
				Error:      true,
				Message:    "Product not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, utils.Response{
				StatusCode: http.StatusInternalServerError,
				Error:      true,
				Message:    "Error retrieving product",
			})
		}
		return
	}

	// Convert nullable types
	if description.Valid {
		product.Description = description.String
	}
	if sku.Valid {
		product.SKU = sku.String
	}
	if compareAtPrice.Valid {
		product.CompareAtPrice = &compareAtPrice.Float64
	}
	if costPrice.Valid {
		product.CostPrice = &costPrice.Float64
	}
	if weight.Valid {
		product.Weight = &weight.Float64
	}
	if weightUnit.Valid {
		product.WeightUnit = weightUnit.String
	}
	if metaTitle.Valid {
		product.MetaTitle = metaTitle.String
	}
	if metaDescription.Valid {
		product.MetaDescription = metaDescription.String
	}

	// Load product images
	imagesQuery := `SELECT id, product_id, url, alt_text, position, created_at 
                   FROM product_images WHERE product_id = $1 ORDER BY position`
	imgRows, err := db.DB.Query(imagesQuery, product.ID)
	if err == nil {
		defer imgRows.Close()
		for imgRows.Next() {
			var img models.ProductImage
			var altText sql.NullString
			if err := imgRows.Scan(&img.ID, &img.ProductID, &img.URL, &altText,
				&img.Position, &img.CreatedAt); err == nil {
				if altText.Valid {
					img.AltText = altText.String
				}
				product.Images = append(product.Images, img)
			}
		}
	}

	// Load product categories
	categoriesQuery := `SELECT c.id, c.name, c.description, c.parent_id, c.slug, 
                       c.image_url, c.meta_title, c.meta_description, c.is_active, 
                       c.updated_at, c.created_at 
                       FROM categories c
                       JOIN product_categories pc ON c.id = pc.category_id
                       WHERE pc.product_id = $1`
	catRows, err := db.DB.Query(categoriesQuery, product.ID)
	if err == nil {
		defer catRows.Close()
		for catRows.Next() {
			var cat models.Category
			var catDesc, imageURL, catMetaTitle, catMetaDescription sql.NullString
			var parentID sql.NullInt64

			if err := catRows.Scan(
				&cat.ID, &cat.Name, &catDesc, &parentID, &cat.Slug,
				&imageURL, &catMetaTitle, &catMetaDescription, &cat.IsActive,
				&cat.UpdatedAt, &cat.CreatedAt,
			); err == nil {
				if catDesc.Valid {
					cat.Description = catDesc.String
				}
				if parentID.Valid {
					pid := int(parentID.Int64)
					cat.ParentID = &pid
				}
				if imageURL.Valid {
					cat.ImageURL = imageURL.String
				}
				if catMetaTitle.Valid {
					cat.MetaTitle = catMetaTitle.String
				}
				if catMetaDescription.Valid {
					cat.MetaDescription = catMetaDescription.String
				}
				product.Categories = append(product.Categories, cat)
			}
		}
	}

	// Load inventory information
	inventoryQuery := `SELECT id, product_id, variant_id, quantity, low_stock_threshold, 
                      reserved_quantity, warehouse_id, location, updated_at, created_at 
                      FROM inventory 
                      WHERE product_id = $1 AND variant_id IS NULL`
	invRow := db.DB.QueryRow(inventoryQuery, product.ID)

	var inventory models.Inventory
	var variantID, lowStockThreshold, warehouseID sql.NullInt64
	var location sql.NullString

	err = invRow.Scan(
		&inventory.ID, &inventory.ProductID, &variantID, &inventory.Quantity,
		&lowStockThreshold, &inventory.ReservedQuantity, &warehouseID,
		&location, &inventory.UpdatedAt, &inventory.CreatedAt,
	)

	if err == nil {
		if lowStockThreshold.Valid {
			lst := int(lowStockThreshold.Int64)
			inventory.LowStockThreshold = &lst
		}
		if warehouseID.Valid {
			wid := int(warehouseID.Int64)
			inventory.WarehouseID = &wid
		}
		if location.Valid {
			inventory.Location = location.String
		}
		product.Inventory = &inventory
	}

	c.JSON(http.StatusOK, utils.Response{
		StatusCode: http.StatusOK,
		Error:      false,
		Message:    "Product retrieved",
		Data:       gin.H{"product": product},
	})
}

func createProductHandler(c *gin.Context) {
	_ = middleware.CheckUserExist(c)
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Error:      true,
			Message:    "Invalid input: " + err.Error(),
		})
		return
	}

	// Start a transaction
	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error starting transaction",
		})
		return
	}
	defer tx.Rollback()

	// Insert product
	productQuery := `INSERT INTO products 
                    (name, description, sku, price, compare_at_price, cost_price, slug, 
                    weight, weight_unit, is_taxable, tax_code, is_digital, is_published, 
                    featured, meta_title, meta_description) 
                    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) 
                    RETURNING id, created_at, updated_at`

	err = tx.QueryRow(
		productQuery,
		product.Name, product.Description, product.SKU, product.Price,
		product.CompareAtPrice, product.CostPrice, product.Slug,
		product.Weight, product.WeightUnit, product.IsTaxable,
		product.TaxCode, product.IsDigital, product.IsPublished,
		product.Featured, product.MetaTitle, product.MetaDescription,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error creating product: " + err.Error(),
		})
		return
	}

	// Insert categories if provided
	if len(product.Categories) > 0 {
		categoryQuery := `INSERT INTO product_categories (product_id, category_id, is_primary) 
                         VALUES ($1, $2, $3)`

		for i, category := range product.Categories {
			isPrimary := i == 0 // First category is primary
			_, err = tx.Exec(categoryQuery, product.ID, category.ID, isPrimary)
			if err != nil {
				c.JSON(http.StatusInternalServerError, utils.Response{
					StatusCode: http.StatusInternalServerError,
					Error:      true,
					Message:    "Error linking category: " + err.Error(),
				})
				return
			}
		}
	}

	// Insert inventory record
	if product.Inventory != nil {
		inventoryQuery := `INSERT INTO inventory 
                          (product_id, quantity, low_stock_threshold, reserved_quantity, 
                          warehouse_id, location) 
                          VALUES ($1, $2, $3, $4, $5, $6) 
                          RETURNING id, created_at, updated_at`

		var inv models.Inventory
		inv.ProductID = &product.ID
		inv.Quantity = product.Inventory.Quantity
		inv.LowStockThreshold = product.Inventory.LowStockThreshold
		inv.ReservedQuantity = product.Inventory.ReservedQuantity
		inv.WarehouseID = product.Inventory.WarehouseID
		inv.Location = product.Inventory.Location

		err = tx.QueryRow(
			inventoryQuery,
			product.ID, inv.Quantity, inv.LowStockThreshold,
			inv.ReservedQuantity, inv.WarehouseID, inv.Location,
		).Scan(&inv.ID, &inv.CreatedAt, &inv.UpdatedAt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				StatusCode: http.StatusInternalServerError,
				Error:      true,
				Message:    "Error creating inventory: " + err.Error(),
			})
			return
		}

		product.Inventory = &inv
	}

	// Insert images if provided
	if len(product.Images) > 0 {
		imageQuery := `INSERT INTO product_images (product_id, url, alt_text, position) 
                      VALUES ($1, $2, $3, $4) 
                      RETURNING id, created_at`

		for i, image := range product.Images {
			err = tx.QueryRow(
				imageQuery,
				product.ID, image.URL, image.AltText, i,
			).Scan(&product.Images[i].ID, &product.Images[i].CreatedAt)

			if err != nil {
				c.JSON(http.StatusInternalServerError, utils.Response{
					StatusCode: http.StatusInternalServerError,
					Error:      true,
					Message:    "Error adding product image: " + err.Error(),
				})
				return
			}

			product.Images[i].ProductID = product.ID
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error committing transaction: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, utils.Response{
		StatusCode: http.StatusCreated,
		Error:      false,
		Message:    "Product created successfully",
		Data:       gin.H{"product": product},
	})
}

func updateProductHandler(c *gin.Context) {
	_ = middleware.CheckUserExist(c)
	productID := c.DefaultQuery("id", "")

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Error:      true,
			Message:    "Invalid input: " + err.Error(),
		})
		return
	}

	// Start a transaction
	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error starting transaction",
		})
		return
	}
	defer tx.Rollback()

	// Update product
	productQuery := `UPDATE products SET 
                    name = $1, description = $2, sku = $3, price = $4, 
                    compare_at_price = $5, cost_price = $6, slug = $7, 
                    weight = $8, weight_unit = $9, is_taxable = $10, 
                    tax_code = $11, is_digital = $12, is_published = $13, 
                    featured = $14, meta_title = $15, meta_description = $16,
                    updated_at = NOW()
                    WHERE id = $17 
                    RETURNING updated_at`

	err = tx.QueryRow(
		productQuery,
		product.Name, product.Description, product.SKU, product.Price,
		product.CompareAtPrice, product.CostPrice, product.Slug,
		product.Weight, product.WeightUnit, product.IsTaxable,
		product.TaxCode, product.IsDigital, product.IsPublished,
		product.Featured, product.MetaTitle, product.MetaDescription,
		productID,
	).Scan(&product.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, utils.Response{
				StatusCode: http.StatusNotFound,
				Error:      true,
				Message:    "Product not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, utils.Response{
				StatusCode: http.StatusInternalServerError,
				Error:      true,
				Message:    "Error updating product: " + err.Error(),
			})
		}
		return
	}

	// Update product categories if provided
	if len(product.Categories) > 0 {
		// First delete existing relationships
		_, err = tx.Exec("DELETE FROM product_categories WHERE product_id = $1", productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				StatusCode: http.StatusInternalServerError,
				Error:      true,
				Message:    "Error updating categories: " + err.Error(),
			})
			return
		}

		// Then insert new relationships
		categoryQuery := `INSERT INTO product_categories (product_id, category_id, is_primary) 
                         VALUES ($1, $2, $3)`

		for i, category := range product.Categories {
			isPrimary := i == 0 // First category is primary
			_, err = tx.Exec(categoryQuery, productID, category.ID, isPrimary)
			if err != nil {
				c.JSON(http.StatusInternalServerError, utils.Response{
					StatusCode: http.StatusInternalServerError,
					Error:      true,
					Message:    "Error linking category: " + err.Error(),
				})
				return
			}
		}
	}

	// Update inventory if provided
	if product.Inventory != nil {
		// Check if inventory record exists
		var invID int
		err := tx.QueryRow("SELECT id FROM inventory WHERE product_id = $1 AND variant_id IS NULL", productID).Scan(&invID)

		if err != nil && errors.Is(err, sql.ErrNoRows) {
			// Create new inventory record
			inventoryQuery := `INSERT INTO inventory 
                              (product_id, quantity, low_stock_threshold, reserved_quantity, 
                              warehouse_id, location) 
                              VALUES ($1, $2, $3, $4, $5, $6) 
                              RETURNING id, created_at, updated_at`

			err = tx.QueryRow(
				inventoryQuery,
				productID, product.Inventory.Quantity, product.Inventory.LowStockThreshold,
				product.Inventory.ReservedQuantity, product.Inventory.WarehouseID, product.Inventory.Location,
			).Scan(&product.Inventory.ID, &product.Inventory.CreatedAt, &product.Inventory.UpdatedAt)
		} else if err == nil {
			// Update existing inventory record
			inventoryQuery := `UPDATE inventory SET 
                              quantity = $1, low_stock_threshold = $2, 
                              reserved_quantity = $3, warehouse_id = $4, 
                              location = $5, updated_at = NOW() 
                              WHERE id = $6 
                              RETURNING updated_at`

			err = tx.QueryRow(
				inventoryQuery,
				product.Inventory.Quantity, product.Inventory.LowStockThreshold,
				product.Inventory.ReservedQuantity, product.Inventory.WarehouseID,
				product.Inventory.Location, invID,
			).Scan(&product.Inventory.UpdatedAt)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				StatusCode: http.StatusInternalServerError,
				Error:      true,
				Message:    "Error updating inventory: " + err.Error(),
			})
			return
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error committing transaction: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		StatusCode: http.StatusOK,
		Error:      false,
		Message:    "Product updated successfully",
		Data:       gin.H{"product": product},
	})
}

func deleteProductHandler(c *gin.Context) {
	_ = middleware.CheckUserExist(c)
	productID := c.DefaultQuery("id", "")

	// Start a transaction
	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error starting transaction",
		})
		return
	}
	defer tx.Rollback()

	// Check if product exists
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", productID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error checking product existence",
		})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, utils.Response{
			StatusCode: http.StatusNotFound,
			Error:      true,
			Message:    "Product not found",
		})
		return
	}

	// Delete related records
	_, err = tx.Exec("DELETE FROM product_categories WHERE product_id = $1", productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error deleting product categories",
		})
		return
	}

	_, err = tx.Exec("DELETE FROM product_images WHERE product_id = $1", productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error deleting product images",
		})
		return
	}

	_, err = tx.Exec("DELETE FROM inventory WHERE product_id = $1", productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error deleting product inventory",
		})
		return
	}

	// Delete the product
	_, err = tx.Exec("DELETE FROM products WHERE id = $1", productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error deleting product",
		})
		return
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      true,
			Message:    "Error committing transaction",
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		StatusCode: http.StatusOK,
		Error:      false,
		Message:    "Product deleted successfully",
	})
}

func ConfigProductController(group *gin.RouterGroup) {
	productRoutes := group.Group("products")
	productRoutes.GET("/", getAllProductsHandler)
	productRoutes.GET("/getProduct", getProductHandler)
	productRoutes.POST("/addProduct", createProductHandler)
	productRoutes.PUT("/updateProduct", updateProductHandler)
	productRoutes.DELETE("/deleteProduct", deleteProductHandler)
}
