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