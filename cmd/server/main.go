package main

import (
	"fmt"
	"log"
	"net/http"
)

func health(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s --> %s", req.Method, req.URL.Path)
	fmt.Fprintf(w, "ok\n")
}

func main() {

	http.HandleFunc("/health", health)

	http.ListenAndServe(":8090", nil)
}
