postgres:
	docker run --name postgres12 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=gobank -p 5432:5432 -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root gobank

dropdb:
	docker exec -it postgres12 dropdb gobank

migrateup:
	migrate -path db/migrations -database "postgresql://root:gobank@localhost:5432/gobank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:gobank@localhost:5432/gobank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test