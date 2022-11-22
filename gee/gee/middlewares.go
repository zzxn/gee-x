package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

const MAX_TRACEBACK = 10

func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Status(http.StatusInternalServerError)
				c.String("Server Internal Error")
			}
		}()
        c.Next()
	}
}

//print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip trace, defer func in Recovery and Recovery
    m := n
    if m > MAX_TRACEBACK {
        m = MAX_TRACEBACK
    }

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
    for _, pc := range pcs[:m] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc) // file and line for the start of the func
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
    if m < n {
        str.WriteString(fmt.Sprintf("\n\t...%d more lines", n - m))
    }
	return str.String()
}
