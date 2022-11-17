package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       // root for each HTTP method
	handlers map[string]HandlerFunc // key is method-pattern
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	segments := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range segments {
		if item == "" {
			continue
		}
		parts = append(parts, item)
		if item[0] == '*' {
			break
		}
	}
	return parts
}

func (r *router) getRoute(method string, path string) (node *node, params map[string]string) {
	searchParts := parsePattern(path)
	params = make(map[string]string)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	node = root.search(searchParts, 0)

	if node != nil {
		parts := parsePattern(node.pattern)
		for idx, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[idx]
			} else if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[idx:], "/")
				// this must the last part, we could not to break
				break
			}
		}
		return node, params
	}

	return nil, nil
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %6s - %s", method, pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	parts := parsePattern(pattern)
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	node, params := r.getRoute(c.Method, c.Path)

	if node != nil {
		c.Params = params
		key := c.Method + "-" + node.pattern
		r.handlers[key](c)
	} else {
		c.Status(http.StatusNotFound)
		c.String("404 NOT FOUND: %s\n", c.Path)
	}
}
