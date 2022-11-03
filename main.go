package main

import (
	"fmt"
	"net/http"
    "gee"
)

// // Engine is the uni handler for all requests
// type Engine struct{}

func main() {
    var counter int
    r := gee.New()
    
    r.GET("/", func (w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
        for k, v := range req.Header {
            fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
        }
    })

    r.GET("/count", func (w http.ResponseWriter, req *http.Request) {
        counter += 1 
        // counter is forever 1, because this func will be called multiple times
        fmt.Fprintf(w, "counter = %d\n", counter)
    })

    r.Run(":9999")
}
