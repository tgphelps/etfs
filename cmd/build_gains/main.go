package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "/home/tgphelps/src/go/etfs/etfs.db")
	check(err)
	defer db.Close()

	tx, err := db.Begin()
	check(err)
	tx.Exec("delete from capital_gain")
	row, err := tx.Query("select event_type, symbol, shares, date, amount " +
		"from event order by id")
	for row.Next() {
		fmt.Printf("row: %v\n", row)
	}
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}
