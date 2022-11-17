package gee

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type RouterGroup struct {
    prefix string // prefix which this group handle
    middlewares []HandlerFunc
    engine *Engine // a group share engine instance with its parent
}

// Engine implements the interface of ServeHTTP
type Engine struct {
    *RouterGroup
	router *router
    groups []*RouterGroup  // store all groups under it
}

func New() *Engine {
    engine := &Engine{router: newRouter()}
    engine.RouterGroup = &RouterGroup{engine: engine}
    engine.groups = []*RouterGroup{engine.RouterGroup}
    return engine
}

func (this *RouterGroup) Group(prefix string) *RouterGroup {
    engine := this.engine
    if !strings.HasPrefix(prefix, "/") {
        log.Panicf("group prefix must start with /, but found %s\n", prefix)
    }
    newGroup := &RouterGroup{
        prefix: this.prefix + prefix, // so prefix must start with "/"
        engine: engine,
    }
    engine.groups = append(engine.groups, newGroup)
    return newGroup
}

func (this *RouterGroup) addRoute(method string, path string, handler HandlerFunc) {
    if !strings.HasPrefix(path, "/") {
        log.Panicf("router path must start with /, but found %s\n", path)
    }
    pattern := this.prefix + path
    this.engine.router.addRoute(method, pattern, handler)
}

func (this *RouterGroup) GET(pattern string, handler HandlerFunc) {
	this.addRoute("GET", pattern, handler)
}

func (this *RouterGroup) POST(pattern string, handler HandlerFunc) {
	this.addRoute("POST", pattern, handler)
}

// implement http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
