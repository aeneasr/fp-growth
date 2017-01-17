package cmd

import (
	"sort"
)

type DataSet []Items

type Items []int


func (s Items) Len() int      { return len(s) }
func (s Items) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type OrderItemsByHeaderTableWrapper struct{
	Items
	H HeadTable
}

func (s OrderItemsByHeaderTableWrapper) Less(i, j int) bool {
	return s.H.GetPosition(s.Items[i]) <s.H.GetPosition(s.Items[j])
}


func OrderItems(d DataSet, h HeadTable) {
	var o OrderItemsByHeaderTableWrapper
	o.H = h
	for x, items := range d {
		o.Items = items
		sort.Sort(o)
		d[x] = filterItems(items, func (i int) bool {
			return h.Get(i) != nil
		})
	}
}

func filterItems(s Items, fn func(int) bool) (p Items) {
	for _, v := range s {
		if fn(v) {
			p = append(p, v)
		}
	}
	return p
}