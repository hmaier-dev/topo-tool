package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
)
import _ "github.com/go-sql-driver/mysql"

func main() {

	dbConnect()

	//// fileServer := http.FileServer(http.Dir("./static"))
	//// http.Handle("/", fileServer)
	//http.HandleFunc("/", indexHandler)
	//fmt.Printf("Server at port 8080...\n")
	//if err := http.ListenAndServe(":8080", nil); err != nil {
	//	log.Fatal(err)
	//}

}

func dbConnect() {
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

	query, err := conn.Prepare("SELECT *  FROM `clearpass`;")
	rows, err := query.Query()
	if err != nil {
		fmt.Println("query.Query", err)
	}
	defer conn.Close()
	colNames, err := rows.Columns()
	colTypes, err := rows.ColumnTypes()
	fmt.Println(colNames, colTypes)
	for i, _ := range colTypes {
		fmt.Println(colTypes[i].ScanType())
	}

	//for rows.Next() {
	//	var id int16
	//	rows.Scan(&id)
	//	fmt.Println("resp", id)
	//}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	fmt.Println(r)

}
