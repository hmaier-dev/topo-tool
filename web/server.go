package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
)
import _ "github.com/go-sql-driver/mysql"

var conn *sql.DB

func main() {

	var _ = "azt-4816"

	conn = dbConnect()
	// var switchNames []string = switchTableNames()
	// for _, sw := range switchNames {
	// selectRowByHostname(sw, hostname)
	// }

	getQuery(conn, "select * from `sw_b`;")

	conn.Close()
}

type switchRow struct {
	id            int
	interfaceName string
	mac           string
	hostname      string
	vlan          int
	stack         int
	interfaceNum  int
}

func dbConnect() *sql.DB {
	var user = "www-data"
	var password = "password123"
	var host = "localhost"
	var port = "3306"
	var dbName = "topology-tool"
	var source = user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName
	conn, err := sql.Open("mysql", source)
	if err != nil {
		fmt.Println("sql.open", err)
		os.Exit(1)
	}
	return conn
}

func getQuery(conn *sql.DB, query string) []string {
	rows, err := conn.Query(query)
	if err != nil {
		fmt.Printf("problems with query: %v \n", err)
	}
	cols, _ := rows.Columns()
	// [row][values] -> e.g. row: [[value][value][value]]

	rawResult := make([][]byte, len(cols))
	// singleRow := make([]string, len(cols))
	// .Scan() needs []any as result type
	dest := make([]interface{}, len(cols))

	for i := range cols {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {
		rows.Scan(dest...)
		for i, elem := range dest {
			fmt.Printf("%v -> %v \n", i, elem)
		}

	}

	return make([]string, 0)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	fmt.Println(r)
}
