package main

// mktrans: Insert buy/sell record into the 'event' table.
//
// Usage:
//    mktrans TYPE SYMBOL SHARES AMOUNT DATE
//
// TYPE must be buy or sell
// DATE MUST be yyyy-mm-dd

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 5 {
		log.Panic("need exactly 5 arguments")
	}
	sql_db, err := sql.Open("sqlite3", "/home/tgphelps/src/go/etfs/etfs.db")
	check(err)
	defer sql_db.Close()
	xactn, symbol, shares, amount, date := convert_args(flag.Args())
	fmt.Println(xactn, symbol, shares, amount, date)

	stmt, err := sql_db.Prepare("insert into event (event_type, symbol, shares, date, amount) values (?, ?, ?, ?, ?)")
	check(err)
	_, err = stmt.Exec(xactn, symbol, shares, date, amount)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func convert_args(args []string) (string, string, int, float64, string) {
	xactn := strings.ToUpper((args[0]))
	symbol := strings.ToUpper(args[1])
	shares, err := strconv.Atoi(args[2])
	check(err)
	amount, err := strconv.ParseFloat(args[3], 64)
	check(err)
	date := args[4]

	if xactn != "BUY" && xactn != "SELL" {
		log.Panic("arg 1 must be BUY or SELL")
	}
	year := date[0:4]
	month := date[5:7]
	day := date[8:10]
	_, err = strconv.Atoi(year)
	check(err)
	_, err = strconv.Atoi(month)
	check(err)
	_, err = strconv.Atoi(day)
	check(err)
	fmt.Println("yyyy-mm-dd", year, month, day)
	return xactn, symbol, shares, amount, date
}
