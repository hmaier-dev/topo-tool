package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	// fileServer := http.FileServer(http.Dir("./static"))
	// http.Handle("/", fileServer)
	http.HandleFunc("/", indexHandler)
	fmt.Printf("Server at port 8080...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	for _, element := range [100]int{} {
		fmt.Fprintf(w, "%d", element)
		time.Sleep(1000000)
	}

}
