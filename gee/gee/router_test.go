package gee

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name/:gender", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(
		parsePattern("/p/:name"), []string{"p", ":name"},
	)
	if !ok {
		t.Fatal("111 test parsePattern failed")
	}
	ok = ok && reflect.DeepEqual(
		parsePattern("/p/*"), []string{"p", "*"},
	)
	if !ok {
		t.Fatal("222 test parsePattern failed")
	}
	ok = ok && reflect.DeepEqual(
		parsePattern("/p/*name/*"), []string{"p", "*name"},
	)
	if !ok {
		t.Fatal("333 test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, params := r.getRoute("GET", "/hello/andrew")

	if n == nil {
		t.Fatal("nil should't be returned")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("pattern is wrong")
	}

	if params["name"] != "andrew" {
		t.Fatal("path params \"name\" is wrong")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, params["name"])
}
