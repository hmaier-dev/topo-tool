package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)
import _ "github.com/go-sql-driver/mysql"

var conn *sql.DB

func main() {

	conn = dbConnect()

	//var switchNames [][]string = getQuery(conn, "SHOW TABLES;")
	//// clearpass is not a switch
	//for i, sw := range switchNames {
	//	if sw[0] == "clearpass" {
	//		switchNames = append(switchNames[:i], switchNames[i+1:]...) // cut out clearpass
	//	}
	//}

	// thx for this, chat-gpt
	// Register the route for serving the CSS file
	// Use-Case for this:
	// You can get content from "<public_path>/<content>" by removing
	// <public_path> from <content> (with StripPrefix), and setting FileServer to your <private_path>
	// e.g. http.Handle("/my_public_path/", http.StripPrefix("/my_public_path/", http.FileServer(http.Dir("secret_location"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	// Register function to "/"
	http.HandleFunc("/", indexHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("cannot listen and server", err)
	}

	err = conn.Close()
	if err != nil {
		log.Fatal("cannot close the database-connection ", err)
	}
}

// ----------------------------------------------------
// Web-Server section
//
// ----------------------------------------------------

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Println("POST received...")
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse POST-Request", http.StatusBadRequest)
		}

		hostname := r.Form.Get("hostname")
		fmt.Printf("%v", hostname)
	}

	dbSlice := getQuery(conn, "SELECT * FROM `sw_c-sued`;")
	table := makeTableStruct(dbSlice)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	//var path = filepath.Join(wd, "web", "static", "index.html")
	var path = filepath.Join(wd, "static", "index.html")
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		log.Fatal("problem with parsing the template ", err)
	}
	err = tmpl.Execute(w, struct{ Data []Row }{Data: table}) // write response to w
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("problem with executing the template ", err)
	}

}

// ----------------------------------------------------
// Section for DB-fun
//
// ----------------------------------------------------

func dbConnect() *sql.DB {
	var user = "www-data"
	var password = "password123"
	var host = "localhost"
	var port = "3306"
	var dbName = "topology-tool"
	var source = user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName
	conn, err := sql.Open("mysql", source)
	if err != nil {
		log.Fatal("sql.open", err)
	}
	return conn
}

func getQuery(conn *sql.DB, query string) [][]string {
	rows, err := conn.Query(query)
	if err != nil {
		fmt.Printf("problems with query: %v \n", err)
	}
	cols, _ := rows.Columns()

	rawResult := make([][]byte, len(cols)) // [row][values] -> e.g. row: [[value][value][value]]
	dest := make([]interface{}, len(cols)) // .Scan() needs []any as result type
	allRows := make([][]string, 0)
	for i := range cols {
		dest[i] = &rawResult[i] // mapping dest indices to byte slice
	}
	for rows.Next() {
		err := rows.Scan(dest...)
		if err != nil {
			log.Fatal("problems scanning the database", err)
		}
		singleRow := make([]string, len(cols))
		for i, raw := range rawResult {
			singleRow[i] = string(raw) // from byte to string
			//fmt.Printf("%v -> %v \n", i, singleRow)
		}
		allRows = append(allRows, singleRow)
	}
	return allRows
}

// single switch-table row
// Note: to export(public) a variable, it must begin with a Uppercase Letter
type Row struct {
	Id            string
	InterfaceName string
	Mac           string
	Hostname      string
	Ip            string
	Vlan          string
	Stack         string
	InterfaceNum  string
}

func makeTableStruct(array [][]string) []Row {
	var table []Row
	for _, entry := range array {
		row := Row{
			Id:            entry[0],
			InterfaceName: entry[1],
			Mac:           entry[2],
			Hostname:      entry[3],
			Ip:            entry[4],
			Vlan:          entry[5],
			Stack:         entry[6],
			InterfaceNum:  entry[7],
		}
		table = append(table, row)
	}
	//fmt.Printf("%v", table)
	return table
}
