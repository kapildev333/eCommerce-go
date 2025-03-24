create_schema:
	migrate create -ext sql -dir db/migration -seq init_schema

migrateup:
	migrate -path db/migration -database "postgres://kapil:Kapil333@localhost:5432/eCommerce?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://kapil:Kapil333@localhost:5432/eCommerce?sslmode=disable" -verbose down

force_version:
	migrate -path db/migration -database "postgres://kapil:Kapil333@localhost:5432/eCommerce?sslmode=disable" -verbose force 1