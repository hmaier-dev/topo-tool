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

func init() {

}

func main() {
	conn = dbConnect()
	// thx for this, chat-gpt
	// Register the route for serving the CSS file
	// Use-Case for this:
	// You can get content from "<public_path>/<content>" by removing
	// <public_path> from <content> (with StripPrefix), and setting FileServer to your <private_path>
	// e.g. http.Handle("/my_public_path/", http.StripPrefix("/my_public_path/", http.FileServer(http.Dir("secret_location"))))
	static := http.FileServer(http.Dir("static"))
	js := http.FileServer(http.Dir("js"))
	http.Handle("/static/", http.StripPrefix("/static/", static))
	http.Handle("/js/", http.StripPrefix("/js/", js))
	// Register function to "/"
	http.HandleFunc("/", indexHandler)
	fmt.Println("Server is starting...")
	err := http.ListenAndServe("0.0.0.0:8181", nil)
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
	wd, err := os.Getwd()
	var static = filepath.Join(wd, "static")
	var index = filepath.Join(static, "index.html")
	var table = filepath.Join(static, "table.html")
	if err != nil {
		log.Fatal("cannot get working directory", err)
	}

	// Check if the request is an OPTIONS preflight request
	if r.Method == "OPTIONS" {
		// Respond with status OK to preflight requests
		w.WriteHeader(http.StatusOK)
		return
	}

	// POST-Request
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse POST-Request", http.StatusBadRequest)
		}
		fmt.Println("POST-Request received...")
		hostname := r.Form.Get("hostname")
		tableData := searchHostname(hostname)
		if tableData == nil || len(tableData) == 0 { // returning if there is no sufficient tableData
			return
		}
		tmpl, err := template.ParseFiles(table)
		err = tmpl.Execute(w, struct{ Data []Row }{Data: tableData}) // write response to w
		if err != nil {
			log.Fatal("problem with executing the template ", err)
		}
	} else {
		// Loading the index.html without any data
		tmpl, err := template.ParseFiles(index)
		err = tmpl.Execute(w, nil) // write response to w
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal("problem with parsing the index template ", err)
		}

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
	//var host = "db"
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
		log.Fatal("problems with query", err)
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

func searchHostname(h string) []Row {
	var switchNames = getQuery(conn, "SHOW TABLES;")
	for i, sw := range switchNames {
		if sw[0] == "clearpass" { // clearpass is not a switch
			switchNames = append(switchNames[:i], switchNames[i+1:]...) // cut out clearpass
		}
	}
	// Get Query Response
	var resp [][]string
	var totalResp [][]string
	for _, sw := range switchNames {
		n := sw[0]
		query := "Select * From `" + n + "` WHERE hostname LIKE '%" + h + "%';"
		resp = getQuery(conn, query) // response from a single table
		for _, entry := range resp {
			entry = append(entry, n) // append the slice with the SwitchName
			totalResp = append(totalResp, entry)
		}
	}

	return makeTableStruct(totalResp)
}

// Row single switch-table row
// Note: to export(public) a variable, it must begin with an Uppercase Letter
type Row struct {
	Id            string
	InterfaceName string
	Mac           string
	Hostname      string
	Ip            string
	Vlan          string
	Stack         string
	InterfaceNum  string
	SwitchName    string
}

// Converting the database-output
// into a usable []struct-format
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
			SwitchName:    entry[8],
		}
		table = append(table, row)
	}
	//fmt.Printf("%v", table)
	return table
}

//func requestRefresh(w http.ResponseWriter) { // from Python!
//	// w.Header().Set("Content-Type", "text/plain")
//
//	url := "localhost:8181"
//	response, err := net.Dial("tcp", url)
//	if err != nil {
//		log.Fatal("Cannot connect to: ", url, " Error: ", err)
//	}
//	defer response.Close()
//
//	buffer := make([]byte, 64) // Hardcoding the message-size to 64 byte
//	//allData := []byte{}
//	for {
//		// n is the amount of data left in response
//		n, err := response.Read(buffer) // read into buffer
//		fmt.Printf("n: %v \n", n)
//		if err != nil {
//			fmt.Println("n: ", n, " Problem reading to buffer... ", err)
//			break
//		}
//		if n == 0 {
//			break
//		}
//		msg := "<p>" + string(buffer[:n]) + "</p>"
//		fmt.Fprintf(w, msg)
//
//		//fmt.Printf("%v \n", string(buffer))
//		//allData = append(allData, buffer[:n]...)
//	}
//	// fmt.Printf("%v \n", n)
//	//fmt.Printf("%v \n", string(allData))
//}
