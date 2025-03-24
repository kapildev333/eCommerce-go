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