.PHONY: deploy

run:
	go run server.go

db:
	cat db/schema.sql | sqlite3 db/gohire.db