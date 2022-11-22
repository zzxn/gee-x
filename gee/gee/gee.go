package gee

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type RouterGroup struct {
	prefix      string // prefix which this group handle，a group's prefix cannot be a prefix of another group's prefix, e.g /v1 and /v1/hello
	middlewares []HandlerFunc
	engine      *Engine // a group share engine instance with its parent
}

// Engine implements the interface of ServeHTTP
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // store all groups under it
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
    engine.Use(Logger())
    engine.Use(Recovery())

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

// Use if defined to add middleware to the group.
// the first defined middleware in the first defined will be first executed.
// user handler itself will be executed at last.
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// implement http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
			// multiple parent-child group could be matched
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
