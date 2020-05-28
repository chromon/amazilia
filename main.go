package main

import (
	"amazilia/ama"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age int8
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := ama.New()

	// 全局中间件
	r.Use(ama.Logger())
	r.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	r.LoadHTMLGlob("templates/html/*")
	// 访问 localhost:8080/assets/css/style.css
	// 返回 ./templates/static/css/style.css
	r.Static("/assets", "templates/static")

	//// -- basic router --
	//r.GET("/index", func(c *ama.Context) {
	//	c.HTML(http.StatusOK, "<h2>Hello Ama</h2>", nil)
	//})

	//// curl "http://localhost:8080/hello?name=hi"
	//// hello hi, you're at /hello
	//r.GET("/hello", func(c *ama.Context) {
	//	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	//})

	//// curl "http://localhost:8080/hello/abc"
	//// hello abc, you're at /hello/abc
	//r.GET("/hello/:name", func(c *ama.Context) {
	//	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	//})

	//// curl "http://localhost:8080/login" -X POST -d 'username=haha&password=hehe'
	//// {"password":"hehe","username":"haha"}
	//r.POST("/login", func(c *ama.Context) {
	//	c.JSON(http.StatusOK, ama.H{
	//		"username": c.PostForm("username"),
	//		"password": c.PostForm("password"),
	//	})
	//})

	// -- router group --
	//g1 := r.Group("/g1")
	//// http://localhost:8080/g1
	//g1.GET("/", func(c *ama.Context) {
	//	c.HTML(http.StatusOK, "<h2>Hello Group</h2>", nil)
	//})
	//// http://localhost:8080/g1/hello?name=hehe
	//g1.GET("/hello", func(c *ama.Context) {
	//	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	//})

	//g2 := r.Group("/g2")
	//g2.Use(ama.LoggerForG2())
	//// http://localhost:8080/g2/hello/hehe
	//g2.GET("/hello/:name", func(c *ama.Context) {
	//	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	//})
	//// curl "http://localhost:8080/g2/login" -X POST -d 'username=haha&password=hehe'
	//g2.POST("/login", func(c *ama.Context) {
	//	c.JSON(http.StatusOK, ama.H{
	//		"username": c.PostForm("username"),
	//		"password": c.PostForm("password"),
	//	})
	//})

	// -- render template --
	r.GET("/", func (c *ama.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	stu1 := &student{Name: "ellery1", Age: 18}
	stu2 := &student{Name: "ellery2", Age: 18}
	r.GET("/students", func(c *ama.Context) {
		c.HTML(http.StatusOK, "test.tmpl", ama.H{
			"title":  "ama",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *ama.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", ama.H{
			"title": "ama",
			"now":   time.Date(2020, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	err := r.Run(":8080")
	if err != nil {
		fmt.Println("Run server err:", err)
	}
}

