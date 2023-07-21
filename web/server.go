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

	var hostname = "azt-4816"

	conn = dbConnect()
	var switchNames []string = switchTableNames()
	for _, sw := range switchNames {
		selectRowByHostname(sw, hostname)
	}
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

func selectAllFromClearpass() {
	query, err := conn.Prepare("SELECT *  FROM `clearpass`;")
	rows, err := query.Query()
	if err != nil {
		fmt.Println("query.Query", err)
	}
	// iterator which yields the next row
	for rows.Next() {
		// declare variables to which Scan() can save the data to
		var (
			id       int
			hostname string
			mac      string
		)
		err := rows.Scan(&id, &hostname, &mac)
		if err != nil {
			return
		}
		fmt.Printf("%v, %v, %v \n", id, hostname, mac)
	}
}

func selectRowByHostname(switchName string, hostname string) {
	stmt, _ := conn.Prepare("SELECT id, interface_name, mac, hostname, vlan, stack , interface_num FROM `?` WHERE hostname = ? ;")
	rows, err := stmt.Query(switchName, hostname)

	if err == nil {
		fmt.Println("stmt.Query", err)
	}
	//columns, _ := rows.ColumnTypes()
	//response_array := make([]string, len(columns))
	//iterator which yields the next row
	for rows.Next() {
		// declare variables to which Scan() can save the data to
		var sr switchRow
		rows.Scan(&sr.id, &sr.interfaceName, &sr.mac, &sr.hostname, &sr.vlan, &sr.stack, &sr.interfaceNum)
		fmt.Printf("%v", sr)

	}
}

func switchTableNames() []string {
	query, _ := conn.Prepare("SHOW TABLES;")
	rows, _ := query.Query()
	columns, _ := rows.ColumnTypes()
	switchNames := make([]string, len(columns))
	var tableName string
	for rows.Next() {
		rows.Scan(&tableName)
		if tableName != "clearpass" {
			switchNames = append(switchNames, tableName)
		}
	}
	return switchNames
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	fmt.Println(r)

}
