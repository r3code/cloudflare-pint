package utils

import "github.com/prometheus/prometheus/promql/parser"

// Node is used to turn the parsed PromQL query expression into a tree.
// This allows us to walk the tree up & down and look for either parents
// or children of specific type. Which is useful if you, for example,
// want to check if all vector selectors are wrapped inside function
// calls etc.
type Node struct {
	Parent   *Node
	Expr     parser.Node
	children []Node
}

// Tree takes a parsed PromQL node and turns it into a Node
// instance with parent and children populated.
func Tree(expr parser.Node, parent *Node) Node {
	n := Node{
		Parent: parent,
		Expr:   expr,
	}
	for _, child := range parser.Children(expr) {
		n.children = append(n.children, Tree(child, &n))
	}
	return n
}

// WalkUpExpr allows to iterate a promQLNode node looking for
// parents of specific type.
// Prometheus parser returns interfaces which makes it more difficult
// to figure out what kind of node we're dealing with, hence this
// helper takes a type parameter it tries to cast.
// It starts by checking the node passed to it and then walks
// up by visiting all parent nodes.
func WalkUpExpr[T parser.Node](node *Node) (nodes []*Node) {
	if node == nil {
		return nodes
	}
	if _, ok := node.Expr.(T); ok {
		nodes = append(nodes, node)
	}
	if node.Parent != nil {
		nodes = append(nodes, WalkUpExpr[T](node.Parent)...)
	}
	return nodes
}

// WalkDownExpr works just like WalkUpExpr but it walks the tree
// down, visiting all children.
// It also starts by checking the node passed to it before walking
// down the tree.
func WalkDownExpr[T parser.Node](node *Node) (nodes []*Node) {
	if _, ok := node.Expr.(T); ok {
		nodes = append(nodes, node)
	}
	for _, child := range node.children {
		nodes = append(nodes, WalkDownExpr[T](&child)...)
	}
	return nodes
}

// WalkUpParent works like WalkUpExpr but checks the parent
// (if present) instead of the node itself.
// It returns the nodes where the parent is of given type.
func WalkUpParent[T parser.Node](node *Node) (nodes []*Node) {
	if node == nil || node.Parent == nil {
		return nodes
	}
	if _, ok := node.Parent.Expr.(T); ok {
		nodes = append(nodes, node)
	}
	if node.Parent != nil {
		nodes = append(nodes, WalkUpParent[T](node.Parent)...)
	}
	return nodes
}
