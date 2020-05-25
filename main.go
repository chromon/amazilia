package main

import (
	"amazilia/ama"
	"fmt"
	"net/http"
)

func main() {
	r := ama.New()

	r.GET("/index", func(c *ama.Context) {
		c.HTML(http.StatusOK, "<h2>Hello Ama</h2>")
	})

	// curl "http://localhost:8080/hello?name=hi"
	// hello hi, you're at /hello
	//r.GET("/hello", func(c *ama.Context) {
	//	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	//})

	// curl "http://localhost:8080/hello/abc"
	// hello abc, you're at /hello/abc
	//r.GET("/hello/:name", func(c *ama.Context) {
	//	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	//})

	// curl "http://localhost:8080/login" -X POST -d 'username=haha&password=hehe'
	// {"password":"hehe","username":"haha"}
	//r.POST("/login", func(c *ama.Context) {
	//	c.JSON(http.StatusOK, ama.H{
	//		"username": c.PostForm("username"),
	//		"password": c.PostForm("password"),
	//	})
	//})

	g1 := r.Group("/g1")
	// http://localhost:8080/g1
	g1.GET("/", func(c *ama.Context) {
		c.HTML(http.StatusOK, "<h2>Hello Group</h2>")
	})
	// http://localhost:8080/g1/hello?name=hehe
	g1.GET("/hello", func(c *ama.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	g2 := r.Group("/g2")
	// http://localhost:8080/g2/hello/hehe
	g2.GET("/hello/:name", func(c *ama.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})
	// curl "http://localhost:8080/g2/login" -X POST -d 'username=haha&password=hehe'
	g2.POST("/login", func(c *ama.Context) {
		c.JSON(http.StatusOK, ama.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	err := r.Run(":8080")
	if err != nil {
		fmt.Println("Run server err:", err)
	}
}

