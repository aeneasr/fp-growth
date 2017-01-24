package cmd

import (
	"sort"
)

type HeadTable []HeadTableRow

type HeadTableRow struct {
	Item  int
	Count int
	Link  *FPTreeNode
}

func NewHeadTable(db DataSet, minSup int) HeadTable {
	ic := map[int]int{}
	for _, tx := range db {
		for _, i := range tx {
			ic[i] = ic[i] + 1
		}
	}

	for k, c := range ic {
		if c < int(minSup) {
			delete(ic, k)
		}
	}

	pl := make(HeadTable, len(ic))
	i := 0
	for k, v := range ic {
		pl[i] = HeadTableRow{Item: k, Count: v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func (p HeadTable) Len() int {
	return len(p)
}
func (p HeadTable) Less(i, j int) bool {
	if p[i].Count == p[j].Count {
		return p[i].Item > p[j].Item
	}
	return p[i].Count < p[j].Count
}
func (p HeadTable) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p HeadTable) SetLink(id int, link *FPTreeNode) {
	for n, i := range p {
		if i.Item != id {
			continue;
		}

		if i.Link == nil {
			i.Link = link
			p[n] = i
		}
	}
}

func (p HeadTable) Get(id int) HeadTableRow {
	for _, i := range p {
		if i.Item == id {
			return i
		}
	}

	return HeadTableRow{Count: -1}
}
func (p HeadTable) GetPosition(id int) int {
	for k, i := range p {
		if i.Item == id {
			return k
		}
	}

	return -1
}