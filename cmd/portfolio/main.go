package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	etfs "tgphelps.com/etfs/pkg"
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
	var total_gain float64
	var total_basis float64
	for row.Next() {
		var stock stock_data
		row.Scan(&stock.symbol, &stock.shares, &stock.basis, &stock.ave_cost)
		stocks = append(stocks, stock)
		syms = append(syms, stock.symbol)
	}
	stock_list := strings.Join(syms, ",")
	get_price_data(stock_list)
	// fmt.Println("stock list:", stock_list)
	fmt.Println("")
	fmt.Println("SYMBOL SHARES  BASIS   AVE COST PRICE    GAIN   %% GAIN")
	fmt.Println("------ ------ -------- -------- ------ -------- -------")
	for _, s := range stocks {
		s.price = get_closing_price(s.symbol)
		s.gain = (s.price - s.ave_cost) * float64(s.shares)
		s.pct_gain = 100 * (s.gain / s.basis)
		total_gain += s.gain
		total_basis += s.basis
		fmt.Printf("%-6s %6d %8.2f %8.2f %6.2f % 8.2f % 6.2f%%\n", s.symbol, s.shares, s.basis, s.ave_cost, s.price, s.gain, s.pct_gain)
	}
	fmt.Println("")
	fmt.Printf("Total basis: %10.2f\n", total_basis)
	fmt.Printf("Total gain: %8.2f %6.2f\n", total_gain, (100 * total_gain / total_basis))
}

func build_portfolio(db *sql.DB) {
	log.Panic("build not implemented")
}

func get_closing_price(symbol string) float64 {
	for sym, price := range prices {
		if symbol == sym {
			return price
		}
	}
	log.Panic("symbol not in price map")
	return 0.0
}

var prices map[string]float64

func get_price_data(syms string) {
	body := etfs.Fetch_stock_data((syms))
	data := map[string]map[string][]map[string]interface{}{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("json.Unmarshal error:")
		fmt.Println(err)
		fmt.Println(body)
		panic("STOP")
	}
	// fmt.Println(data)
	prices = make(map[string]float64)
	for _, sd := range data["quoteResponse"]["result"] {
		symbol := sd["symbol"].(string)
		last := sd["regularMarketPrice"].(float64)
		prices[symbol] = last
	}
	// fmt.Println(prices)
}
