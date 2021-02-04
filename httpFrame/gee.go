package gee

import (
	"net/http"
	"strings"
)

type HandlerFunc func(ctx *Context)

type Engine struct {
	*RouterGroup // 匿名继承
	router       *router
	groups       []*RouterGroup // store all routeGroup
}

func New() *Engine {
	e := &Engine{
		router: newRouter(),
		groups: make([]*RouterGroup, 0),
	}
	e.RouterGroup = &RouterGroup{
		basePath:   "/", // 这里减少使用group时，
		engine:     e,
		middleware: make([]HandlerFunc, 0),
	}
	return e
}

func (e *Engine) Run(port string) error {
	return http.ListenAndServe(port, e)
}

/*
重写 ServeHTTP 的必要性，其实也是写这个框架的原因
为什么要写这个框架，那肯定是go src内的http server，有什么缺陷
*/
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r) // gin框架在这里创建context是通过sync.Pool中获取的
	// 根结点的路由集合，绑定的中间件
	middles := make([]HandlerFunc, len(e.RouterGroup.middleware))
	copy(middles, e.RouterGroup.middleware)
	// 组织各种中间件
	for _, group := range e.groups {
		// 匹配对应的 group
		if strings.HasPrefix(r.URL.Path, group.basePath) {
			middles = append(middles, group.middleware...)
		}
	}
	c.handlers = middles
	e.router.handle(c)
}

/*
	MethodHead    = "HEAD"
	MethodPatch   = "PATCH" // RFC 5789
	MethodConnect = "CONNECT"
	MethodTrace   = "TRACE"
*/
