package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	sql_db, err := sql.Open("sqlite3", "/home/tgphelps/src/go/etfs/etfs.db")
	check(err)
	defer sql_db.Close()

	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "build" {
			build_portfolio(sql_db)
		}
	}
	show_portfolio(sql_db)
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func show_portfolio(db *sql.DB) {
	row, err := db.Query(("select symbol, shares, basis, basis/shares as ave_cost from portfolio order by symbol"))
	check(err)

	var symbol string
	var shares int
	var basis float32
	var ave_cost float32
	for row.Next() {
		row.Scan(&symbol, &shares, &basis, &ave_cost)
		fmt.Printf("%5s %4d %7.2f %7.2f\n", symbol, shares, basis, ave_cost)
	}
}

func build_portfolio(db *sql.DB) {
	log.Panic("build not implemented")
}
