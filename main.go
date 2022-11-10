package main

import (
	"example.com/gee"
)

// // Engine is the uni handler for all requests
// type Engine struct{}

func main() {
	r := gee.New()

	r.GET("/", func(c *gee.Context) {
		c.HTML("<h1> Hello Gee</h1>")
	})

	r.GET("/hello", func(c *gee.Context) {
		// expect /hello?name=geekandrew
		c.String("hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(gee.H{
			"username": c.Form("username"),
			"password": c.Form("password"),
		})
	})

	r.Run(":9999")
}
