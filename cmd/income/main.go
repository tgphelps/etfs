package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var show_all bool
var only_sym string

func main() {
	flag.BoolVar(&show_all, "l", false, "list income for all symbolss")
	flag.StringVar(&only_sym, "s", "", "show income for this symbol")
	flag.Parse()
	// fmt.Printf("list_all: %v\n", list_all)
	// fmt.Printf("only_sym: %v\n", only_sym)
	if (!show_all && only_sym == "") ||
		(show_all && only_sym != "") {
		log.Panicln("Use exactly one of -l and -s")
	}
	sql_db, err := sql.Open("sqlite3", "/home/tgphelps/src/go/etfs/etfs.db")
	check(err)
	defer sql_db.Close()

	if show_all {
		show_all_income(sql_db)
	} else {
		show_income_for_sym(sql_db, only_sym)
	}
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func show_all_income(db *sql.DB) {
	fmt.Println("showing all income")

	row, err := db.Query("select date, symbol, shares, basis, amount from capital_gain order by date, symbol")
	check(err)
	var total_gain float64

	fmt.Println("   DATE    SYM  SHRS   BUY->SELL       DIFF    GAIN")
	fmt.Println("---------- ---- ---- ---------------- ------- ------")
	for row.Next() {
		var date string
		var symbol string
		var shares int
		var basis float64
		var amount float64
		row.Scan(&date, &symbol, &shares, &basis, &amount)

		diff := amount - basis
		gain := float64(shares) * diff
		total_gain += gain

		fmt.Printf("%10s %-4s %4d (%6.2f->%6.2f) %7.2f %7.2f\n",
			date, symbol, shares, basis, amount, diff, gain)
	}
	fmt.Printf("Total gain: %7.2f\n", total_gain)
}

func show_income_for_sym(db *sql.DB, sym string) {
	if sym == "" {
		log.Panic("can't happen")
	}
	sym = strings.ToUpper(sym)
	row, err := db.Query("select date, shares, basis, amount from capital_gain where symbol = ? order by date",
		sym)
	check(err)
	fmt.Println("   DATE    SYM  SHRS   BUY->SELL       DIFF    GAIN")
	fmt.Println("---------- ---- ---- ---------------- ------- -------")

	for row.Next() {
		var date string
		var shares int
		var basis float64
		var amount float64
		row.Scan(&date, &shares, &basis, &amount)
		diff := amount - basis
		gain := float64(shares) * diff
		fmt.Printf("%10s %-4s %4d (%6.2f->%6.2f) %7.2f %7.2f\n",
			date, sym, shares, basis, amount, diff, gain)
	}
}
