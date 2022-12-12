package main

import (
	"io"
	"log"
	"net/http"
)
func main() {
	pingHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "hello world!\n")
	}
	http.HandleFunc("/ping", pingHandler)
	log.Fatal(http.ListenAndServe(":8080",nil))
}