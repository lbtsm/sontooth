package gee

import (
	"encoding/json"
	"net/http"
)

/*
设计context的必要性
1、	src内的http serveHttp方法的参数，参数设计的很细，每次编写返回信息时，总是需要显式的设置header，对于少量接口来说，还算是方便，但是
	现在企业业务都是需要大量的接口去支持，故设计一个context（框架）便于编写易读的代码（将header的设置，隐藏在context｜框架内）
2、 便于以后方便的扩展，使接入链路追踪、接口监控更加方便
3、	支持动态路由，例如：/get/:name （动态路由这个比较鸡肋）
*/
type Context struct {
	writer       http.ResponseWriter
	request      *http.Request
	Path, Method string
	StatusCode   int
	UrlParams    map[string]string
	// 自定义的中间件
	handlers []HandlerFunc
	index    int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		writer:   w,
		request:  r,
		Path:     r.URL.Path,
		Method:   r.Method,
		handlers: make([]HandlerFunc, 0),
		index:    -1,
	}
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.writer.WriteHeader(c.StatusCode)
}

func (c *Context) SetHeader(key, value string) {
	c.writer.Header().Set(key, value)
}

func (c *Context) Query(key string) string {
	return c.request.URL.Query().Get(key)
}

func (c *Context) PostForm(key string) string {
	return c.request.FormValue(key)
}

func (c *Context) Json(data interface{}) {
	c.Status(http.StatusOK)
	c.SetHeader("Content-Type", "application/json")
	encoder := json.NewEncoder(c.writer)
	if err := encoder.Encode(data); err != nil {
		// todo ready 500 func
		http.Error(c.writer, err.Error(), 500)
	}
}

func (c *Context) HTML(html string) {
	c.Status(http.StatusOK)
	c.SetHeader("Content-Type", "text/html")
	if _, err := c.writer.Write([]byte(html)); err != nil {
		// todo ready 500 func
		http.Error(c.writer, err.Error(), 500)
	}
}

func (c *Context) String(data string) (int, error) {
	return c.writer.Write([]byte(data))
}

func (c *Context) UrlParam(key string) (bool, string) {
	value, ok := c.UrlParams[key]
	return ok, value
}

func (c *Context) Next() {
	c.index++
	for ; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index](c)
	}
}
