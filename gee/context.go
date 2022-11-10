package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	// response info
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:     w,
		Req:        req,
		Path:       req.URL.Path,
		Method:     req.Method,
		StatusCode: http.StatusOK,
	}
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Form(key string) string {
	return c.Req.FormValue(key)
}

// call it before write body
func (c *Context) Status(code int) {
	c.StatusCode = code
}

// call it before write body
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// after calling, do not touch context
func (c *Context) String(format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Writer.WriteHeader(c.StatusCode)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// after calling, do not touch context
func (c *Context) JSON(obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Writer.WriteHeader(c.StatusCode)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// after calling, do not touch context
func (c *Context) Data(data []byte) {
	c.Writer.WriteHeader(c.StatusCode)
	c.Writer.Write(data)
}

// after calling, do not touch context
func (c *Context) HTML(html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(c.StatusCode)
	c.Writer.Write([]byte(html))
}
