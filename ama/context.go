package ama

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSON 数据
type H map[string]interface{}

// 设计上下文 Context，封装 Request 和 Response ，
// 提供对 JSON、HTML 等返回类型的支持
type Context struct {

	// 响应
	Writer http.ResponseWriter

	// 请求
	Req *http.Request

	// 请求信息
	Path string
	Method string

	// 响应信息
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req: req,
		Path: req.URL.Path,
		Method: req.Method,
	}
}

// 获取 form values
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 获取 value
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 响应状态
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 设置 header
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 构造字符串响应内容
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	_, err := c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
	if err != nil {
		fmt.Println("writer write err:", err)
	}
}

// 构造 JSON 响应内容
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// 构造 Data 响应内容
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	_, err := c.Writer.Write(data)
	if err != nil {
		fmt.Println("writer write err:", err)
	}
}

// 构造 HTML 响应内容
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	_, err := c.Writer.Write([]byte(html))
	if err != nil {
		fmt.Println("writer write err:", err)
	}
}