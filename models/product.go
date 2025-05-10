package models

import "time"

// Category represents a product category
type Category struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	ParentID        *int      `json:"parent_id,omitempty"`
	Slug            string    `json:"slug"`
	ImageURL        string    `json:"image_url,omitempty"`
	MetaTitle       string    `json:"meta_title,omitempty"`
	MetaDescription string    `json:"meta_description,omitempty"`
	IsActive        bool      `json:"is_active"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedAt       time.Time `json:"created_at"`
}

// Dimensions represents product dimensions
type Dimensions struct {
	Length float64 `json:"length,omitempty"`
	Width  float64 `json:"width,omitempty"`
	Height float64 `json:"height,omitempty"`
}

// Product represents a product in the system
type Product struct {
	ID              int         `json:"id"`
	Name            string      `json:"name"`
	Description     string      `json:"description,omitempty"`
	SKU             string      `json:"sku,omitempty"`
	Price           float64     `json:"price"`
	CompareAtPrice  *float64    `json:"compare_at_price,omitempty"`
	CostPrice       *float64    `json:"cost_price,omitempty"`
	Slug            string      `json:"slug"`
	Weight          *float64    `json:"weight,omitempty"`
	WeightUnit      string      `json:"weight_unit,omitempty"`
	Dimensions      *Dimensions `json:"dimensions,omitempty"`
	IsTaxable       bool        `json:"is_taxable"`
	TaxCode         string      `json:"tax_code,omitempty"`
	IsDigital       bool        `json:"is_digital"`
	IsPublished     bool        `json:"is_published"`
	Featured        bool        `json:"featured"`
	MetaTitle       string      `json:"meta_title,omitempty"`
	MetaDescription string      `json:"meta_description,omitempty"`
	UpdatedAt       time.Time   `json:"updated_at"`
	CreatedAt       time.Time   `json:"created_at"`

	// Relations
	Categories []Category       `json:"categories,omitempty"`
	Images     []ProductImage   `json:"images,omitempty"`
	Variants   []ProductVariant `json:"variants,omitempty"`
	Inventory  *Inventory       `json:"inventory,omitempty"`
}

// ProductCategory represents the many-to-many relationship between products and categories
type ProductCategory struct {
	ID         int       `json:"id"`
	ProductID  int       `json:"product_id"`
	CategoryID int       `json:"category_id"`
	IsPrimary  bool      `json:"is_primary"`
	CreatedAt  time.Time `json:"created_at"`
}

// ProductImage represents an image associated with a product
type ProductImage struct {
	ID        int       `json:"id"`
	ProductID int       `json:"product_id"`
	URL       string    `json:"url"`
	AltText   string    `json:"alt_text,omitempty"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
}

// ProductVariant represents a variant of a product
type ProductVariant struct {
	ID             int               `json:"id"`
	ProductID      int               `json:"product_id"`
	SKU            string            `json:"sku,omitempty"`
	Name           string            `json:"name,omitempty"`
	Price          *float64          `json:"price,omitempty"`
	CompareAtPrice *float64          `json:"compare_at_price,omitempty"`
	CostPrice      *float64          `json:"cost_price,omitempty"`
	Weight         *float64          `json:"weight,omitempty"`
	Dimensions     *Dimensions       `json:"dimensions,omitempty"`
	IsDefault      bool              `json:"is_default"`
	Options        map[string]string `json:"options,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`

	// Relations
	Inventory *Inventory `json:"inventory,omitempty"`
}

// ProductAttribute represents a product attribute (like color, size, etc.)
type ProductAttribute struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`

	// Relations
	Values []ProductAttributeValue `json:"values,omitempty"`
}

// ProductAttributeValue represents a value for a product attribute
type ProductAttributeValue struct {
	ID           int       `json:"id"`
	AttributeID  int       `json:"attribute_id"`
	Value        string    `json:"value"`
	DisplayValue string    `json:"display_value"`
	CreatedAt    time.Time `json:"created_at"`
}

// Inventory represents the stock level of a product or variant
type Inventory struct {
	ID                int       `json:"id"`
	ProductID         *int      `json:"product_id,omitempty"`
	VariantID         *int      `json:"variant_id,omitempty"`
	Quantity          int       `json:"quantity"`
	LowStockThreshold *int      `json:"low_stock_threshold,omitempty"`
	ReservedQuantity  int       `json:"reserved_quantity"`
	WarehouseID       *int      `json:"warehouse_id,omitempty"`
	Location          string    `json:"location,omitempty"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedAt         time.Time `json:"created_at"`

	// Relations
	Movements []InventoryMovement `json:"movements,omitempty"`
}

// InventoryMovement represents a change in inventory
type InventoryMovement struct {
	ID             int       `json:"id"`
	InventoryID    int       `json:"inventory_id"`
	QuantityChange int       `json:"quantity_change"`
	ReferenceType  string    `json:"reference_type"`
	ReferenceID    string    `json:"reference_id,omitempty"`
	Notes          string    `json:"notes,omitempty"`
	CreatedBy      *int      `json:"created_by,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// Order represents a customer order
type Order struct {
	ID                int       `json:"id"`
	UserID            *int      `json:"user_id,omitempty"`
	OrderNumber       string    `json:"order_number"`
	Status            string    `json:"status"`
	Subtotal          float64   `json:"subtotal"`
	TaxAmount         float64   `json:"tax_amount"`
	ShippingAmount    float64   `json:"shipping_amount"`
	DiscountAmount    float64   `json:"discount_amount"`
	TotalAmount       float64   `json:"total_amount"`
	Currency          string    `json:"currency"`
	ShippingAddressID *int      `json:"shipping_address_id,omitempty"`
	PaymentID         *int      `json:"payment_id,omitempty"`
	Notes             string    `json:"notes,omitempty"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedAt         time.Time `json:"created_at"`

	// Relations
	Items           []OrderItem      `json:"items,omitempty"`
	ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
	Payment         *UserPayment     `json:"payment,omitempty"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	ProductID *int      `json:"product_id,omitempty"`
	VariantID *int      `json:"variant_id,omitempty"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
	Subtotal  float64   `json:"subtotal"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Product *Product        `json:"product,omitempty"`
	Variant *ProductVariant `json:"variant,omitempty"`
}
