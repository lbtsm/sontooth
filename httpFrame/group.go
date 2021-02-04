package gee

import (
	"net/http"
	"path"
)

type RouterGroup struct {
	basePath   string        // api分组前缀
	middleware []HandlerFunc // 中间间支持，这个是group最重要的一个功能，todo 在中间件中，有些请求其实是没有必要去进行中间件函数执行的，例如健康检查
	engine     *Engine
}

func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	g := &RouterGroup{
		basePath:   rg.engine.basePath + prefix,
		middleware: make([]HandlerFunc, 0),
		engine:     rg.engine,
	}
	rg.engine.groups = append(rg.engine.groups, g)
	return g
}

func (rg *RouterGroup) addRoute(method, route string, handler HandlerFunc) {
	fullPath := path.Join(rg.basePath, route)
	rg.engine.router.addRoute(method, fullPath, handler)
}

func (rg *RouterGroup) Get(path string, handler HandlerFunc) {
	rg.addRoute(http.MethodGet, path, handler)
}

func (rg *RouterGroup) Post(path string, handler HandlerFunc) {
	rg.addRoute(http.MethodPost, path, handler)
}

func (rg *RouterGroup) Delete(path string, handler HandlerFunc) {
	rg.addRoute(http.MethodDelete, path, handler)
}

func (rg *RouterGroup) Put(path string, handler HandlerFunc) {
	rg.addRoute(http.MethodPut, path, handler)
}

func (rg *RouterGroup) Option(path string, handler HandlerFunc) {
	rg.addRoute(http.MethodOptions, path, handler)
}

func (rg *RouterGroup) Use(middleware ...HandlerFunc) {
	rg.middleware = append(rg.middleware, middleware...)
}
