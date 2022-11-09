package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
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

type event struct {
	event_type string
	symbol     string
	shares     int
	date       string
	amount     float64
}

func main() {
	var testing bool
	flag.BoolVar(&testing, "t", false, "turn on testing code")
	flag.Parse()

	sql_db, err := sql.Open("sqlite3", "/home/tgphelps/src/go/etfs/etfs.db")
	check(err)
	defer sql_db.Close()

	args := flag.Args()
	// fmt.Println("args:", args)
	if len(args) > 0 && args[0] == "build" {
		build_portfolio(sql_db, testing)
	}
	show_portfolio(sql_db, testing)
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func show_portfolio(db *sql.DB, testing bool) {
	var suffix string
	if testing {
		suffix = "2"
		fmt.Println("testing...")
	} else {
		suffix = ""
	}
	sql := fmt.Sprintf("select symbol, shares, basis, basis/shares as ave_cost "+
		"from portfolio%s order by symbol", suffix)
	row, err := db.Query(sql)
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
	if len(syms) == 0 {
		log.Panic("No portfolio rows.")
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
		fmt.Printf("%-6s %6d %8.2f %8.2f %6.2f % 8.2f % 6.2f%%\n",
			s.symbol, s.shares, s.basis, s.ave_cost, s.price, s.gain, s.pct_gain)
	}
	fmt.Println("")
	fmt.Printf("Total basis: %10.2f\n", total_basis)
	fmt.Printf("Total gain: %8.2f %6.2f\n", total_gain, (100 * total_gain / total_basis))
}

func build_portfolio(db *sql.DB, testing bool) {
	// Read all buy/sell events and build portfolio array.

	// The 'order by id' insures that the events are in chronological order.
	row, err := db.Query("select event_type, symbol, shares, date, amount from event order by id")
	check(err)

	var e event
	p := make(map[string]etfs.Holding)

	for row.Next() {
		row.Scan(&e.event_type, &e.symbol, &e.shares, &e.date, &e.amount)
		// fmt.Println("event for: ", e.symbol)
		if _, found := p[e.symbol]; !found {
			p[e.symbol] = etfs.Holding{Sym: e.symbol}
			// fmt.Println("added key:", e.symbol)
		}
		switch e.event_type {
		case "BUY":
			// fmt.Println("buy ", e.symbol)
			h := p[e.symbol]
			h.Buy_shares(e.shares, e.amount)
			p[e.symbol] = h
		case "SELL":
			// fmt.Println("sell ", e.symbol)
			h := p[e.symbol]
			h.Sell_shares(e.shares, e.amount)
			p[e.symbol] = h
		default:
			log.Panic("Bad event type:", e.event_type, e.symbol)
		}
	}
	if testing {
		fmt.Println("testing...")
		for h := range p {
			if p[h].Shares > 0 {
				fmt.Println(p[h].Sym, p[h].Shares, p[h].Total_cost, p[h].Ave_cost)
			}
		}
	}

	var tbl_name string
	if testing {
		tbl_name = "portfolio2"
	} else {
		tbl_name = "portfolio"
	}

	del_sql := fmt.Sprintf("delete from %s", tbl_name)
	tx, err := db.Begin()
	check(err)
	_, err = tx.Exec(del_sql)
	check(err)
	// fmt.Println("truncated table: ", tbl_name)

	insert_holdings(tx, tbl_name, p)
	err = tx.Commit()
	check(err)
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

func insert_holdings(tx *sql.Tx, tbl string, p map[string]etfs.Holding) {
	// fmt.Println("insert holdings...")
	for h := range p {
		if p[h].Shares > 0 {
			sql := fmt.Sprintf("insert into %s (symbol, shares, basis) values (?, ?, ?)", tbl)
			_, err := tx.Exec(sql, p[h].Sym, p[h].Shares, p[h].Total_cost)
			check(err)
		}
	}
}
