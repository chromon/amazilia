package ama

import (
	"log"
	"net/http"
	"strings"
)

// 路由相关
type router struct {
	// 存储每种请求方式的Trie 树根节点
	// key - roots["GET"] roots["POST"]
	roots map[string]*node

	// 存储每种请求方式的 HandlerFunc
	// key - handlers["GET-/p/:lang/doc"]
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots: make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 解析路由
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 注册路由规则，映射 handler
func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)

	parts := parsePattern(pattern)

	// key 由请求方法和静态路由地址构成，如 GET-/
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}

	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// 获取路由
func (r *router) getRouter(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	// 解析 : 和 * 两种匹配符的参数，返回一个 map
	// /p/go/doc 匹配到 /p/:lang/doc，解析结果为：{lang: "go"}
	// /static/css/style.css 匹配到 /static/*filepath，解析结果为：{filepath: "css/style.css"}
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// 解析请求的路径，查找路由映射表，匹配路由规则，查找到对应的 handler
// 如果查到，就执行注册的处理方法。如果查不到，就返回 404 NOT FOUND
func (r *router) handle(c *Context) {
	n, params := r.getRouter(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}