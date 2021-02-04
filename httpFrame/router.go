package gee

import (
	"net/http"
	"strings"
)

/*
拆出一个router结构体，方便后续的路由动作，例如支持动态路由
*/
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc // gin框架对应的value是handlerFunc的切片
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *router) parsePath(path string) []string {
	parts := strings.Split(path, "/")
	if parts[0] == "*" {
		return parts
	}

	ret := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}

		ret = append(ret, part)
	}

	return ret
}

func (r *router) addRoute(method, path string, handler HandlerFunc) {
	key := method + "-" + path
	r.handlers[key] = handler
	// 添加到tree🌲中
	parts := r.parsePath(path)
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(path, parts, 0)
}

// example return map[string]string {"lang":go}
func (r *router) getRoute(method, path string) (*node, map[string]string) {
	var (
		receiveParts = r.parsePath(path)
		params       = make(map[string]string)
		root, ok     = r.roots[method]
	)
	if !ok {
		return nil, nil
	}

	target := root.search(receiveParts, 0)
	if target == nil {
		return nil, nil
	}
	// 根据找到
	targetParts := r.parsePath(target.path)
	for idx, tp := range targetParts {
		if tp[0] == ':' {
			params[tp[1:]] = receiveParts[idx]
		}
		// 文件路径
		if tp[0] == '*' && len(tp) > 1 {
			params[tp[1:]] = strings.Join(receiveParts[idx:], "/")
			break
		}
	}

	return target, params
}

func (r *router) handle(c *Context) {
	target, params := r.getRoute(c.Method, c.Path)
	if target != nil {
		c.UrlParams = params
		// todo 这里这种拼接字符串的方式，是否会影响效率
		key := c.Method + "-" + c.Path
		c.handlers = append(c.handlers, r.handlers[key])
		c.Next()
	} else {
		// todo ready 404 handler
		c.Status(http.StatusNotFound)
	}
}
