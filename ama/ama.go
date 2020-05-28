package ama

import (
	"html/template"
	"net/http"
	"strings"
)

// 定义路由映射的处理方法为 handler 类型，
type HandlerFunc func(*Context)

// 实现 http Handler 接口
type Engine struct{
	*RouterGroup

	// 路由映射 map
	router *router

	// 存储所有分组
	groups []*RouterGroup

	// html 模板渲染
	// 将所有的模板加载进内存
	htmlTemplates *template.Template
	// 所有自定义模板渲染函数
	funcMap template.FuncMap
}

// 构造 Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}

	return engine
}

// 添加路由
// method：请求方法
// pattern：静态路由地址
// handler：路由映射的处理方法
func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	engine.router.addRouter(method, pattern, handler)
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

// 设置自定义渲染函数
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// 设置加载模板的方法
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(
		template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// 启动 http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 重写 Handler ServerHTTP 方法，处理 HTTP 请求。
// w：构造针对请求的响应
// req：包含 HTTP 请求的索引信息，请求地址、Header、Body等
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// 接收到一个具体请求时，要判断该请求适用于哪些中间件
	var middleWares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middleWares = append(middleWares, group.middleWares...)
		}
	}

	c := newContext(w, req)
	c.handlers = middleWares
	c.engine = engine
	engine.router.handle(c)
}