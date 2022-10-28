package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	sql_db, _ := sql.Open("sqlite3", "/home/tgphelps/src/go/etfs/etfs.db")
	defer sql_db.Close()

main_loop:
	for {
		f := get_command()
		if len(f) > 1 {
			f[1] = strings.ToUpper(f[1])
		}
		fmt.Println("cmd: ", f)
		switch f[0] {
		case "q":
			break main_loop
		case "h":
			do_help()
		case "l":
			do_list(sql_db)
		case "i":
			do_insert(f, sql_db)
		case "u":
			do_update(f, sql_db)
		case "a":
			do_activate(f, sql_db)
		case "d":
			do_delete(f, sql_db)
		case "v":
			do_view(f, sql_db)
		default:
			do_help()
		}
	}
}

func get_command() []string {
	cmd := read_command("cmd >")
	f := strings.Fields(cmd)
	if len(f) == 0 {
		return []string{"h"}
	} else {
		return f
	}
}

func read_command(prompt string) string {
	fmt.Print(prompt)
	r := bufio.NewReader(os.Stdin)
	cmd, err := r.ReadString(('\n'))
	if err == io.EOF {
		return "q"
	}
	return strings.TrimSpace(cmd)
}

func do_help() {
	fmt.Println("commands: l(ist) i(nsert) u(update) d(delete) v(iew) a(activate) q(uit)")
	fmt.Println("i <sym> <name>")
	fmt.Println("u <sym> <name>")
	fmt.Println("d <sym>")
	fmt.Println("v <sym>")
	fmt.Println("a <sym> <0/1>")
}

func do_list(db *sql.DB) {
	row, err := db.Query(("select symbol, name, active from etf order by symbol"))
	check(err)
	defer row.Close()
	for row.Next() {
		var symbol string
		var text string
		var active int
		row.Scan(&symbol, &text, &active)
		fmt.Printf("%5s %30s   %d\n", symbol, text, active)
	}
}

func do_insert(f []string, db *sql.DB) {
	stmt, err := db.Prepare("insert into etf (symbol, name, active) values (?, ?, ?)")
	check(err)
	stmt.Exec(f[1], strings.Join(f[2:], " "), 0)
}

func do_update(f []string, db *sql.DB) {
	stmt, err := db.Prepare("update etf set name = ? where symbol = ?")
	check(err)
	stmt.Exec(strings.Join(f[2:], " "), f[1])
}

func do_activate(f []string, db *sql.DB) {
	stmt, err := db.Prepare("update etf set active = ? where symbol = ?")
	check(err)
	flag, _ := strconv.Atoi(f[2])
	stmt.Exec(flag, f[1])
}

func do_delete(f []string, db *sql.DB) {
	stmt, err := db.Prepare("delete from etf where symbol = ?")
	check(err)
	stmt.Exec(f[1])
}

func do_view(f []string, db *sql.DB) {
	row, err := db.Query("select symbol, name, active from etf where symbol == ?", f[1])
	check(err)
	for row.Next() {
		var symbol string
		var text string
		var active int
		row.Scan(&symbol, &text, &active)
		fmt.Printf("%4s %30ss   %d\n", symbol, text, active)
	}
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}
