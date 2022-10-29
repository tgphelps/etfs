
default: fetch

fetch:
	go build -o bin ./cmd/fetch

manage:
	go build -o bin ./cmd/manage

portfolio:
	go build -o bin ./cmd/portfolio

clean:
	rm bin/*

sql:
	sqlite3 etfs.db
