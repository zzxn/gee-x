package main

import (
	"log"
	"time"

	"example.com/gee"
)

// // Engine is the uni handler for all requests
// type Engine struct{}

func main() {
	r := gee.New()
	r.Use(gee.Logger())

	r.GET("/index", func(c *gee.Context) {
		c.HTML("<h1> GEE INDEX </h1>")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML("<h1> Hello Gee</h1>")
		})

		v1.GET("/hello", func(c *gee.Context) {
			// expect /hello?name=geekandrew
			c.String("hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	v2.Use(func(c *gee.Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	})
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String("hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		v2.GET("/assets/*filepath", func(c *gee.Context) {
			c.JSON(gee.H{"filepath": c.Param("filepath")})
		})

		v2.POST("/login", func(c *gee.Context) {
			c.JSON(gee.H{
				"username": c.Form("username"),
				"password": c.Form("password"),
			})
		})
	}

	r.Run(":9999")
}
