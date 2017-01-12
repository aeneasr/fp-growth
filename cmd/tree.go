package cmd

import (
	"fmt"
)

type FPTree struct {
	Root *FPTreeNode
}

type FPTreeNode struct {
	Item     int
	Count    int
	Link     *FPTreeNode
	Parent   *FPTreeNode
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
	return fmt.Sprintf("{item=%d count=%d parent=%s link=%s children=%v}", f.Item, f.Count, parent, link, f.Children)
}

func (f *FPTreeNode) OnlyOneBranch() bool {
	if len(f.Children) > 1 {
		return false
	}

	for _, c := range f.Children {
		if !c.OnlyOneBranch() {
			return false
		}
	}
	return true
}

type Pattern struct {
	Pattern Items
	Count   int
}

func (f *FPTreeNode) MinePatterns(p ConditionalItem) []Pattern {
	if !f.OnlyOneBranch() {
		return nil
	}

	var list []ConditionalItem
	pa := f
	for {
		list = append(list, ConditionalItem{
			Item: pa.Item,
			Count: pa.Count,
		})
		if len(pa.Children) == 0 {
			break;
		}
		pa = pa.Children[0]
	}

	comb := func(n, m int, emit func([]int)) {
		s := make([]int, m)
		last := m - 1
		var rc func(int, int)
		rc = func(i, next int) {
			for j := next; j < n; j++ {
				s[i] = j
				if i == last {
					emit(s)
				} else {
					rc(i + 1, j + 1)
				}
			}
			return
		}
		rc(0, 0)
	}

	res := []Pattern{
		{
			Pattern: Items{p.Item},
			Count: p.Count,
		},
	}
	for i := 1; i <= len(list); i++ {
		comb(len(list), i, func(c []int) {
			pattern := Pattern{
				Count: -1,
			}
			for _, v := range c {
				pattern.Pattern = append(pattern.Pattern, list[v].Item)
				if list[v].Count < pattern.Count || pattern.Count == -1 {
					pattern.Count = list[v].Count
				}
			}
			pattern.Pattern = append(pattern.Pattern, p.Item)
			res = append(res, pattern)
		})
	}

	return res
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
			n.Link = l[k + 1]
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