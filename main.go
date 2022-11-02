package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
    counter := 0
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	})
    http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
        counter += 1
        fmt.Fprintf(w, "counter = %d\n", counter)
    }) 
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
