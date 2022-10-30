package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type stock_data struct {
	symbol   string
	shares   int
	basis    float64
	ave_cost float64
	price    float64
	gain     float64
	pct_gain float64
}

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

	var stocks []stock_data
	var syms []string
	for row.Next() {
		var stock stock_data
		row.Scan(&stock.symbol, &stock.shares, &stock.basis, &stock.ave_cost)
		stocks = append(stocks, stock)
		syms = append(syms, stock.symbol)
	}
	stock_list := strings.Join(syms, ",")
	get_price_data(stock_list)
	fmt.Println("stock list:", stock_list)
	for _, s := range stocks {
		s.price = get_closing_price(s.symbol)
		s.gain = (s.price - s.ave_cost) * float64(s.shares)
		s.pct_gain = 100 * (s.gain / s.basis)

		fmt.Printf("%-5s %4d %8.2f %7.2f %7.2f % 10.2f % 7.2f%%\n", s.symbol, s.shares, s.basis, s.ave_cost, s.price, s.gain, s.pct_gain)

	}
}

func build_portfolio(db *sql.DB) {
	log.Panic("build not implemented")
}

func get_closing_price(sym string) float64 {
	return 100.0
}

func get_price_data(syms string) {

}
