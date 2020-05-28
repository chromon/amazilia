package ama

import (
	"log"
	"net/http"
	"path"
)

// 路由分组控制
type RouterGroup struct {
	// 分组前缀
	prefix string

	// 中间件
	middleWares []HandlerFunc

	// 分组嵌套的父分组
	parent *RouterGroup

	// 所有分组共享一个 engine 实例
	engine *Engine
}

// 创建一个新 RouterGroup
func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	group := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		engine: engine,
	}

	engine.groups = append(engine.groups, group)
	return group
}

// 添加路由
func (g *RouterGroup) addRouter(method string, comp string, handler HandlerFunc) {
	pattern := g.prefix + comp
	log.Printf("GroupRouter %4s - %s", method, pattern)
	g.engine.router.addRouter(method, pattern, handler)
}

// GET 请求
func (g *RouterGroup) GET(pattern string, handler HandlerFunc) {
	g.addRouter("GET", pattern, handler)
}

// POST 请求
func (g *RouterGroup) POST(pattern string, handler HandlerFunc) {
	g.addRouter("POST", pattern, handler)
}

// 添加中间件到群组
func (g *RouterGroup) Use(middleWares ...HandlerFunc) {
	g.middleWares = append(g.middleWares, middleWares...)
}

// 创建静态文件处理 handler
func (g *RouterGroup) createStaticHandler(relativePath string, fileSys http.FileSystem) HandlerFunc {
	absolutePath := path.Join(g.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fileSys))
	return func (c *Context) {
		file := c.Param("filepath")
		// 判断文件是否存在，或是否有权限访问
		if _, err := fileSys.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// 将本地文件映射到路由
func (g *RouterGroup) Static(relativePath string, root string) {
	handler := g.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	g.GET(urlPattern, handler)
}




