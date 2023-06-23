postgres:
	docker run --name rezy -p 5432:5432 -d -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret postgres:12-alpine

createdb:
	docker exec -it rezy createdb --username=root --owner=root rezydb

dropdb:
	docker exec -it rezy dropdb rezydb

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/rezydb?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/rezydb?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

mock:
	~/go/bin/mockgen -package mockdb -destination db/mock/store.go github.com/adam-macioszek/rezy/db/sqlc Store 


.PHONY: postgres createdb dropdb migrateup migratedown test server