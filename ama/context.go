package ama

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSON 数据
type H map[string]interface{}

// 根据请求 *http.Request，构造响应 http.ResponseWriter
// 一个完整的响应，需要考虑消息头(Header)和消息体(Body)，而 Header 包含了状态码(StatusCode)，消息类型(ContentType)等

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

	// 解析到的参数
	// {lang: "go"} {filepath: "css/style.css"}
	Params map[string]string

	// 响应信息
	StatusCode int

	// 中间件
	handlers []HandlerFunc
	// 记录当前执行到第几个中间件
	index int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req: req,
		Path: req.URL.Path,
		Method: req.Method,
		index: -1,
	}
}

// 中间件等待用户自己定义的 Handler 处理结束后，再做一些额外的操作
// 当在中间件中调用 Next 方法时，控制权交给了下一个中间件，直到调用到最后一个中间件，
// 然后再从后往前返回，调用每个中间件在 Next 方法之后定义的部分。
/*
	func A(c *Context) {
		part1
		c.Next()
		part2
	}
	func B(c *Context) {
		part3
		c.Next()
		part4
	}
	假设我们应用了中间件 A 和 B，和路由映射的 Handler。c.handlers是这样的[A, B, Handler]，c.index初始化为-1。调用c.Next()
	最终的顺序是part1 -> part3 -> Handler -> part 4 -> part2
*/
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// 获取路由参数
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
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