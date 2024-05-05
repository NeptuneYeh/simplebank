createdb:
	docker exec -it my-postgres createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it my-postgres dropdb --username=root --owner=root simple_bank
postgres_init:
	docker run --name my-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16.1-alpine3.19	
postgres_start:
	docker start simple_bank-my-postgres-1
postgres_stop:
	docker-compose -p simple_bank stop my-postgres
migrate_up:
	migrate -path scripts/db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrate_down:
	migrate -path scripts/db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test_db:
	go test -v -cover ./test/db/
server:
	go run ./cmd/main.go
mock_db:
	mockgen -package mockdb -destination internal/infrastructure/database/postgres/mock/store.go github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc Store
dev_service_up:
	docker compose -p simple_bank up -d
test_all:
	go test -coverpkg=./... -coverprofile=coverage.out ./test/...
proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
        --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
        proto/*.proto

.PHONY: createdb dropdb postgres_init postgres_start postgres_stop migrate_up migrate_down sqlc test server mock_db dev_service_up proto