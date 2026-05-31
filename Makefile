generate:
	go run github.com/99designs/gqlgen generate
	go generate ./...

sqlc:
	sqlc generate

run-postgres:
	docker-compose -f docker-compose.postgres.yml up -d --build

run-inmemory:
	docker-compose -f docker-compose.inmemory.yml up -d --build

stop:
	docker-compose -f docker-compose.inmemory.yml down -v
	docker-compose -f docker-compose.postgres.yml down -v

test:
	go test -v -count=1 ./...
