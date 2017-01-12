package cmd

import (
	"fmt"
)

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
	var parent = "none"
	var link = "none"
	if f.Parent != nil {
		parent = fmt.Sprintf("%d:%d", f.Parent.Item, f.Parent.Count)
	}
	if f.Link != nil {
		link = fmt.Sprintf("%d:%d", f.Link.Item, f.Link.Count)
	}
	return fmt.Sprintf("{item=%d count=%d parent=%s link=%s children=%v}", f.Item, f.Count, parent, link,f.Children)
}

func NewFPTree(ordered DataSet, ht *HeadTable) FPTree {
	t := FPTree{
		Root: &FPTreeNode{
			Children: []*FPTreeNode{},
		},
	}

	links := map[int][]*FPTreeNode{}
	for _, tx := range ordered {
		buildTree(tx, t.Root, nil, links)
	}

	for item, l := range links {
		ht.SetLink(item, l[0])
		if len(l) == 1 {
			continue
		}

		for k, n := range l[0:len(l) - 1] {
			n.Link = l[k+1]
		}
	}

	return t
}

func buildTree(items Items, root *FPTreeNode, parent *FPTreeNode, links map[int][]*FPTreeNode) {
	current := root
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
				Parent: parent,
				Children: []*FPTreeNode{},
				Link: nil,
			}
			current.Children = append(current.Children, nn)
			current = nn

			links[item] = append(links[item], nn)
		}
		parent = current
	}
}