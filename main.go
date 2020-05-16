package main

import (
	"amazilia/ama"
	"fmt"
	"net/http"
)

func main() {
	r := ama.New()

	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		_, err := fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
		if err != nil {
			fmt.Println("URL.Path err:", err)
		}
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			_, err := fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
			if err != nil {
				fmt.Println("Header err:", err)
			}
		}

	})

	err := r.Run(":8080")
	if err != nil {
		fmt.Println("Run server err:", err)
	}
}

