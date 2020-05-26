package ama

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		// 开始计时
		t := time.Now()
		// 处理请求
		c.Next()
		// 计算消耗时间
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func LoggerForG2() HandlerFunc {
	return func(c *Context) {
		// 起始时间
		t := time.Now()
		// 发生错误
		c.Fail(500, "Internal Server Error")
		// 计算消耗时间
		log.Printf("[%d] %s in %v for group g2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}