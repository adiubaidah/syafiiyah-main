DB_URL=postgresql://postgres:root@localhost:5432/as_syafiiyah?sslmode=disable

sqlcinit:
	docker run --rm -v "%cd%:/src" -w /src sqlc/sqlc init
sqlc:
	docker run --rm -v "%cd%:/src" -w /src sqlc/sqlc generate

proto:
	protoc \
  --proto_path=./internal/constant/proto \
  --go_out=paths=source_relative:./platform/protobuf \
  --go-grpc_out=paths=source_relative:./platform/protobuf \
  ./internal/constant/proto/*.proto

start_postgres:
	docker container start postgres
postgres:
	docker exec -it postgres psql -U postgres -d as_syafiiyah

redis:
	docker exec -it redis redis-cli

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrate_force:
	migrate -path db/migration -database $(DB_URL) -verbose force $(version)

migrate_up:
	migrate -path db/migration -database $(DB_URL) -verbose up

migrate_down:
	migrate -path db/migration -database $(DB_URL) -verbose down

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

build_windows:
	go build -o main.exe cmd/rest/main.go

test:
	go test -v ./... -cover