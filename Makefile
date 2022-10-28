
default: fetch

fetch:
	go build -o bin ./cmd/fetch

manage:
	go build -o bin ./cmd/manage

clean:
	rm bin/*

sql:
	sqlite3 etfs.db
