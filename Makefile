postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=USER -e POSTGRES_PASSWORD=PASSWORD -d postgres:14.11-bookworm

createdb:
	docker exec -it postgres14 createdb --username=USER --owner=USER eCommerce

dropdb:
	docker exec -it postgres14 dropdb eCommerce

create_schema:
	migrate create -ext sql -dir db/migration -seq init_schema

migrateup:
	migrate -path db/migration -database "postgresql://USER:PASSWORD@localhost:5432/eCommerce?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://USER:PASSWORD@localhost:5432/eCommerce?sslmode=disable" -verbose down

force_version:
	migrate -path db/migration -database "postgres://USER:PASSWORD@localhost:5432/eCommerce?sslmode=disable" -verbose force 1