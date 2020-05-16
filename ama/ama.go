package ama

import (
	"fmt"
	"net/http"
)

// 定义路由映射的处理方法为 handler 类型，
type HandlerFunc func(http.ResponseWriter, *http.Request)

// 实现 http Handler 接口
type Engine struct{
	// 路由映射 map
	router map[string]HandlerFunc
}

// 构造 Engine
func New() *Engine {
	return &Engine{
		router: make(map[string]HandlerFunc),
	}
}

// 添加路由
// method：请求方法
// pattern：静态路由地址
// handler：路由映射的处理方法
func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	// key 由请求方法和静态路由地址构成，如 GET-/
	key := method + "-" + pattern
	engine.router[key] = handler
}

// GET 请求
// pattern：静态路由地址
// handler：路由映射的处理方法
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRouter("GET", pattern, handler)
}

// POST 请求
// pattern：静态路由地址
// handler：路由映射的处理方法
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRouter("POST", pattern, handler)
}

// 启动 http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 重写 Handler ServerHTTP 方法，处理 HTTP 请求。解析请求的路径，查找路由映射表，
// 如果查到，就执行注册的处理方法。如果查不到，就返回 404 NOT FOUND
// w：构造针对请求的响应
// req：包含 HTTP 请求的索引信息，请求地址、Header、Body等
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path

	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		_, err := fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
		if err != nil {
			fmt.Println("fprintf err:", err)
		}
	}
}