package ama

import (
	"log"
	"net/http"
)

// 路由相关
type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

// 添加路由
func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)

	// key 由请求方法和静态路由地址构成，如 GET-/
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// 解析请求的路径，查找路由映射表，
// 如果查到，就执行注册的处理方法。如果查不到，就返回 404 NOT FOUND
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}