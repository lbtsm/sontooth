package gee

import (
	"strings"
)

type node struct {
	path     string  // 完整url路径
	part     string  // 部分路径
	children []*node // part 路径下的子节点
	isWild   bool    // 是否精准匹配，含有: & * 时，为true
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchAllChild(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(path string, parts []string, index int) {
	if len(parts) == index {
		n.path = path
		return
	}

	part := parts[index]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*', // 插入路由仅在项目初始化使用，如考虑性能，可考虑其他方式对比
		}
		n.children = append(n.children, child) // 只需要在child==nil的情况下添加
	}
	child.insert(path, parts, index+1)
}

func (n *node) search(parts []string, index int) *node {
	// 递归出口
	if len(parts) == index || strings.HasPrefix(n.part, "*") {
		if n.path == "" {
			return nil
		}
		return n
	}

	part := parts[index]
	children := n.matchAllChild(part)
	for _, child := range children {
		result := child.search(parts, index+1)
		if result != nil {
			return result
		}
	}

	return nil
}
