package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

// є дерево, яке записане строкою у вигляді: (a, b), (b, c), (a, d), (d, i)
// написати скріпт який приймає строку дерева, і перефразовує в строку дерева іншого формату: (a(b(c)d(i))).
type nodeOfTree struct {
	value string
	nodes []*nodeOfTree
}

func (n *nodeOfTree) buildTreeString() string {
	if len(n.nodes) == 0 {
		return n.value
	}

	var child string
	for _, v := range n.nodes {
		child += v.buildTreeString()
	}

	n.value = fmt.Sprintf("%s(%s)", n.value, child)

	return n.value
}

func getRootNodeFromArray(array [][]string) (*nodeOfTree, error) {
	mapNodes := make(map[string]*nodeOfTree)

	var root *nodeOfTree

	for _, node := range array {
		parent := node[0]
		child := node[1]

		_, ok := mapNodes[parent]
		if !ok {
			if _, ok = mapNodes[child]; !ok && root != nil {
				return nil, errors.New("your tree has more than one root")
			}

			mapNodes[parent] = &nodeOfTree{
				value: parent,
			}
			root = mapNodes[parent]
		}

		if ok {
			if _, ok = mapNodes[child]; ok {
				return nil, errors.New("deadlock occurred in the tree")
			}
		}

		_, ok = mapNodes[child]
		if !ok {
			mapNodes[child] = &nodeOfTree{
				value: child,
			}
			parentNode := mapNodes[parent]
			parentNode.nodes = append(parentNode.nodes, mapNodes[child])
		}

		if ok {
			parentNode := mapNodes[parent]
			parentNode.nodes = append(parentNode.nodes, mapNodes[child])
			root = parentNode
		}
	}

	return root, nil
}

func splitTreeStringToArray(s string) [][]string {
	var nodeArray [][]string

	for _, v := range strings.Split(s, "),(") {
		node := strings.Split(strings.Trim(v, "()"), ",")
		nodeArray = append(nodeArray, node)
	}

	return nodeArray
}

func main() {
	s := "(a,b),(b,c),(a,d),(d,i)"
	nodeArray := splitTreeStringToArray(s)

	rootNode, err := getRootNodeFromArray(nodeArray)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("(%s)\n", rootNode.buildTreeString())
}
