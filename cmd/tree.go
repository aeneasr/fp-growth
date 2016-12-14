package cmd

import "fmt"

type FPTree struct {
	Root *FPTreeNode
}

type FPTreeNode struct {
	Item    int
	Count   int
	Link    *FPTreeNode
	Parent  *FPTreeNode
	Children []*FPTreeNode
}

func (f *FPTreeNode) String() string {
	var parent = 0
	if f.Parent != nil {
		parent = f.Parent.Item
	}
	return fmt.Sprintf("{item=%d count=%d link=%v parent=%d children=%v}", f.Item, f.Count, f.Link, parent ,f.Children)
}

func NewFPTree(ordered DataSet) FPTree {
	t := FPTree{
		Root: &FPTreeNode{
			Children: []*FPTreeNode{},
		},
	}

	for _, tx := range ordered {
		buildBranch(tx, t.Root)
	}
	return t
}

func buildBranch(items Items, n *FPTreeNode) {
	current := n
	var c *FPTreeNode
	for _, item := range items {
		var found bool
		for _, c = range current.Children {
			if c.Item == item {
				c.Count++
				current = c
				found = true
				break;
			}
		}
		if !found {
			nn := &FPTreeNode{
				Item: item,
				Count: 1,
				Parent: n,
				Children: []*FPTreeNode{},
			}
			current.Children = append(current.Children, nn)
			current = nn
		}
	}
}