package gee

import "strings"

type node struct {
	pattern  string // (only for valid node) full path to this node, e.g. /p/:lang
	part     string // the last part of path
	children []*node
	wild     bool // if wild is true, then it can match any part
}

// matchChild find the first matched child
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.wild {
			return child
		}
	}
	return nil
}

// matchChildren find all matched children
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.wild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// insert insert the parts[height:] to this node
// pattern == "/part[0]/.../part[n - 1]"
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, wild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// search search parts[height:] in this node, only returns valid node
func (n *node) search(parts []string, height int) *node {
    // here n.part could be empty, so n.part[0] is invalid
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			// not a valid node
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		var result *node = child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
