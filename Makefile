createdb:
	docker exec -it my-postgres createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it my-postgres dropdb --username=root --owner=root simple_bank
postgres_init:
	docker run --name my-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16.1-alpine3.19	
postgres_start:
	docker start my-postgres
postgres_stop:
	docker stop my-postgres
migrate_up:
	migrate -path scripts/db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrate_down:
	migrate -path scripts/db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test_db:
	go test -v -cover ./test/db/

.PHONY: createdb dropdb postgres_init postgres_start postgres_stop migrate_up migrate_down sqlc