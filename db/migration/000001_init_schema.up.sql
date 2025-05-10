CREATE TABLE "users"
(
    id         SERIAL PRIMARY KEY,
    email      VARCHAR(255) UNIQUE NOT NULL,
    username   VARCHAR(255)        NOT NULL,
    password   VARCHAR(255)        NOT NULL,
    updated_at TIMESTAMPTZ         NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ         NOT NULL DEFAULT now()
);

CREATE TABLE "shipping_addresses"
(
    id             SERIAL PRIMARY KEY,
    user_id        INTEGER REFERENCES users (id) ON DELETE CASCADE,
    address_line_1 VARCHAR(255) NOT NULL,
    address_line_2 VARCHAR(255),
    city           VARCHAR(255) NOT NULL,
    state          VARCHAR(255),
    postal_code    VARCHAR(20)  NOT NULL,
    country        VARCHAR(255) NOT NULL,
    is_default     BOOLEAN      NOT NULL DEFAULT FALSE,
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE "user_payments"
(
    id                 SERIAL PRIMARY KEY,
    user_id            INTEGER REFERENCES users (id) ON DELETE CASCADE,
    payment_method     VARCHAR(50)    NOT NULL,                                              -- e.g., credit_card, paypal, bank_transfer
    transaction_id     VARCHAR(255) UNIQUE,                                                  -- Unique ID from payment gateway
    amount             DECIMAL(10, 2) NOT NULL,
    currency           VARCHAR(3)     NOT NULL DEFAULT 'USD',                                -- ISO 4217 currency code
    payment_status     VARCHAR(50)    NOT NULL,                                              -- e.g., pending, completed, failed, refunded
    payment_date       TIMESTAMPTZ    NOT NULL DEFAULT now(),
    billing_address_id INTEGER        REFERENCES shipping_addresses (id) ON DELETE SET NULL, -- Optional billing address
    description        TEXT,                                                                 -- Optional description of the payment
    updated_at         TIMESTAMPTZ    NOT NULL DEFAULT now(),
    created_at         TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE TABLE "categories"
(
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    parent_id       INTEGER REFERENCES categories (id) ON DELETE SET NULL,
    slug            VARCHAR(100) NOT NULL UNIQUE,
    image_url       VARCHAR(255),
    meta_title      VARCHAR(100),
    meta_description VARCHAR(255),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE "products"
(
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    sku             VARCHAR(50) UNIQUE,
    price           DECIMAL(10, 2) NOT NULL,
    compare_at_price DECIMAL(10, 2),
    cost_price      DECIMAL(10, 2),
    slug            VARCHAR(255) NOT NULL UNIQUE,
    weight          DECIMAL(8, 2),
    weight_unit     VARCHAR(10) DEFAULT 'kg',
    dimensions      JSON, -- {length, width, height}
    is_taxable      BOOLEAN DEFAULT TRUE,
    tax_code        VARCHAR(50),
    is_digital      BOOLEAN DEFAULT FALSE,
    is_published    BOOLEAN NOT NULL DEFAULT FALSE,
    featured        BOOLEAN NOT NULL DEFAULT FALSE,
    meta_title      VARCHAR(100),
    meta_description VARCHAR(255),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE "product_categories"
(
    id              SERIAL PRIMARY KEY,
    product_id      INTEGER NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    category_id     INTEGER NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    is_primary      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(product_id, category_id)
);

CREATE TABLE "product_images"
(
    id              SERIAL PRIMARY KEY,
    product_id      INTEGER NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    url             VARCHAR(255) NOT NULL,
    alt_text        VARCHAR(255),
    position        INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE "product_variants"
(
    id              SERIAL PRIMARY KEY,
    product_id      INTEGER NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    sku             VARCHAR(50) UNIQUE,
    name            VARCHAR(255),
    price           DECIMAL(10, 2),
    compare_at_price DECIMAL(10, 2),
    cost_price      DECIMAL(10, 2),
    weight          DECIMAL(8, 2),
    dimensions      JSON, -- {length, width, height}
    is_default      BOOLEAN NOT NULL DEFAULT FALSE,
    options         JSONB, -- e.g., {"color": "red", "size": "XL"}
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE "product_attributes"
(
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL UNIQUE,
    display_name    VARCHAR(100) NOT NULL,
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE "product_attribute_values"
(
    id              SERIAL PRIMARY KEY,
    attribute_id    INTEGER NOT NULL REFERENCES product_attributes (id) ON DELETE CASCADE,
    value           VARCHAR(100) NOT NULL,
    display_value   VARCHAR(100) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(attribute_id, value)
);

CREATE TABLE "inventory"
(
    id              SERIAL PRIMARY KEY,
    product_id      INTEGER REFERENCES products (id) ON DELETE SET NULL,
    variant_id      INTEGER REFERENCES product_variants (id) ON DELETE SET NULL,
    quantity        INTEGER NOT NULL DEFAULT 0,
    low_stock_threshold INTEGER DEFAULT 5,
    reserved_quantity INTEGER NOT NULL DEFAULT 0,
    warehouse_id    INTEGER, -- If you implement multiple warehouses later
    location        VARCHAR(100), -- Shelf/bin reference
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT inventory_product_variant_check CHECK ((product_id IS NOT NULL AND variant_id IS NULL) OR (variant_id IS NOT NULL))
);

CREATE TABLE "inventory_movements"
(
    id              SERIAL PRIMARY KEY,
    inventory_id    INTEGER NOT NULL REFERENCES inventory (id) ON DELETE CASCADE,
    quantity_change INTEGER NOT NULL, -- Positive for additions, negative for removals
    reference_type  VARCHAR(50) NOT NULL, -- e.g., "order", "manual_adjustment", "return", "restock"
    reference_id    VARCHAR(100), -- ID of the related entity (order_id, etc.)
    notes           TEXT,
    created_by      INTEGER REFERENCES users (id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE "orders"
(
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER REFERENCES users (id) ON DELETE SET NULL,
    order_number    VARCHAR(50) UNIQUE NOT NULL,
    status          VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, processing, completed, cancelled, refunded
    subtotal        DECIMAL(10, 2) NOT NULL,
    tax_amount      DECIMAL(10, 2) NOT NULL DEFAULT 0,
    shipping_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    total_amount    DECIMAL(10, 2) NOT NULL,
    currency        VARCHAR(3) NOT NULL DEFAULT 'USD',
    shipping_address_id INTEGER REFERENCES shipping_addresses (id) ON DELETE SET NULL,
    payment_id      INTEGER REFERENCES user_payments (id) ON DELETE SET NULL,
    notes           TEXT,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE "order_items"
(
    id              SERIAL PRIMARY KEY,
    order_id        INTEGER NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    product_id      INTEGER REFERENCES products (id) ON DELETE SET NULL,
    variant_id      INTEGER REFERENCES product_variants (id) ON DELETE SET NULL,
    quantity        INTEGER NOT NULL,
    unit_price      DECIMAL(10, 2) NOT NULL,
    subtotal        DECIMAL(10, 2) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT order_items_product_variant_check CHECK ((product_id IS NOT NULL AND variant_id IS NULL) OR (variant_id IS NOT NULL))
);

-- Create indexes for better performance
CREATE INDEX idx_product_categories_product_id ON product_categories(product_id);
CREATE INDEX idx_product_categories_category_id ON product_categories(category_id);
CREATE INDEX idx_inventory_product_id ON inventory(product_id);
CREATE INDEX idx_inventory_variant_id ON inventory(variant_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);
CREATE INDEX idx_order_items_variant_id ON order_items(variant_id);