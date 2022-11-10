package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	etfs "tgphelps.com/etfs/pkg"
)

func main() {
	db, err := sql.Open("sqlite3", "/home/tgphelps/src/go/etfs/etfs.db")
	check(err)
	defer db.Close()

	tx, err := db.Begin()
	check(err)
	tx.Exec("delete from capital_gain")
	hld := make(map[string]etfs.Holding)
	row, err := tx.Query("select event_type, symbol, shares, date, amount " +
		"from event order by id")
	check(err)
	for row.Next() {
		process(row, tx, hld)
	}
	tx.Commit()
}

func process(row *sql.Rows, tx *sql.Tx, hld map[string]etfs.Holding) {
	var event string
	var symbol string
	var shares int
	var date string
	var amount float64
	row.Scan(&event, &symbol, &shares, &date, &amount)
	// fmt.Printf("date: %s event: %s symbol: %s\n", date, event, symbol)
	if _, found := hld[symbol]; !found {
		hld[symbol] = etfs.Holding{Sym: symbol}
		// fmt.Println("added key:", symbol)
	}
	switch event {
	case "BUY":
		h := hld[symbol]
		h.Buy_shares(shares, amount)
		hld[symbol] = h
	case "SELL":
		h := hld[symbol]
		basis := h.Ave_cost
		profit := h.Sell_shares(shares, amount)
		hld[symbol] = h
		fmt.Printf("%s Sold %d shares of %s. Profit = %6.2f\n", date, shares, symbol, profit)
		_, err := tx.Exec("insert into capital_gain (date, symbol, shares, basis, amount) values (?,?,?,?,?)",
			date, symbol, shares, basis, amount)
		check(err)
	default:
		log.Panic("cannot happen")
	}
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}
