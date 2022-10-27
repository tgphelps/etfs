package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	sql_db, _ := sql.Open("sqlite3", "/home/tgphelps/src/go/etfs/etfs.db")
	defer sql_db.Close()

	for {
		f := get_command()
		if len(f) > 1 {
			f[1] = strings.ToUpper(f[1])
		}
		fmt.Println(f)
		if f[0] == "q" {
			break
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
	cmd, _ := r.ReadString(('\n'))
	return strings.TrimSpace(cmd)
}
