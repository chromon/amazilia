package main

import (
	"amazilia/ama"
	"fmt"
	"net/http"
)

func main() {
	r := ama.New()

	r.GET("/", func(c *ama.Context) {
		c.HTML(http.StatusOK, "<h2>Hello Amazilia</h2>")
	})

	r.GET("/hello", func(c *ama.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *ama.Context) {
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

