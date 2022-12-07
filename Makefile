
default: fetch

fetch:
	go build -o bin ./cmd/fetch

manage:
	go build -o bin ./cmd/manage

mktrans:
	go build -o bin ./cmd/mktrans
portfolio:
	go build -o bin ./cmd/portfolio

income:
	go build -o bin ./cmd/income

build_gains:
	go build -o bin ./cmd/build_gains

clean:
	rm bin/*

sql:
	sqlite3 etfs.db
