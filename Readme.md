# eCommerce-go

This is an eCommerce application built with Go, Gin framework, and PostgreSQL. The application provides functionalities for managing users, shipping addresses, and user payments.

## Table of Contents

- [Installation](#installation)
- [Database Migrations](#database-migrations)
- [API Endpoints](#api-endpoints)
- [Project Structure](#project-structure)
- [License](#license)

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/kapildev333/eCommerce-go.git
    cd eCommerce-go
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

3. Set up the PostgreSQL database and update the connection string in the `Makefile` and `db/db.go`.

## Database Migrations

To manage database migrations, the project uses the `migrate` tool. Here are the commands to manage migrations:

- Create a new migration:

    ```sh
    make create_schema
    ```

- Apply all pending migrations:

    ```sh
    make migrateup
    ```

- Roll back the most recent migration:

    ```sh
    make migratedown
    ```

- Force the migration version:

    ```sh
    make force_version
    ```

## API Endpoints

### Address Endpoints

- **Get All Addresses**

    ```http
    GET /address/getAllAddress
    ```

- **Add Address**

    ```http
    POST /address/addAddress
    ```

- **Get Address**

    ```http
    GET /address/getAddress
    ```

### Payment Endpoints

- **Submit Payment**

    ```http
    POST /payments/submitPayment
    ```

- **Get User Payment History**

    ```http
    GET /payments/getUserPaymentHistory
    ```

## Project Structure

```plaintext
eCommerce-go/
├── controllers/
│   ├── auth_controller.go
│   ├── address_controller.go
│   ├── payment_controller.go
├── db/
│   ├── db.go
│   ├── migration/
│       ├── 000001_init_schema.up.sql
│       ├── 000001_init_schema.down.sql
├── engine/
│   ├── app_engine.go
│   ├── app_routes.go
├── models/
│   ├── auth_model.go
│   ├── shipping_address.go
│   ├── user_payments.go
├── utils/
│   ├── common_response.go
│   ├── jwt_handler.go    
├── main.go
├── Makefile
├── go.mod
├── go.sum