package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

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
}

func show_income_for_sym(db *sql.DB, sym string) {
	if sym == "" {
		log.Panic("can't happen")
	}
	fmt.Printf("showing income for %s\n", sym)
}
