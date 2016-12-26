package cmd

type ConditionalPatternBases []ConditionalPatternBase

type ConditionalPatternBase struct {
	Prefix ConditionalItem
	Bases  [][]ConditionalItem
}

type ConditionalItem struct {
	Item  int
	Count int
}

func mineParents(p *FPTreeNode, max int) []ConditionalItem {
	var pb = []ConditionalItem{}
	for p != nil {
		count := p.Count
		if count > max {
			count = max
		}

		pb = append(pb, ConditionalItem{Item: p.Item, Count: count})
		p = p.Parent
	}
	return pb
}

func MineConditionalPatternBases(t FPTree, ht HeadTable) ConditionalPatternBases {
	base := ConditionalPatternBases{}

	for i := len(ht) - 1; i >= 0; i-- {
		l := ht[i].Link
		pb := ConditionalPatternBase{
			Prefix: ConditionalItem{
				Item: l.Item,
				Count: ht[i].Count,
			},
			Bases: [][]ConditionalItem{[]ConditionalItem{}},
		}
		pb.Bases[0] = mineParents(l.Parent, ht[i].Count)

		branch := 1
		l = l.Link
		for l != nil {
			count := l.Count
			if count > pb.Prefix.Count {
				count = pb.Prefix.Count
			}

			pb.Bases = append(pb.Bases, []ConditionalItem{{
				Item: l.Item,
				Count: count,
			}})

			pb.Bases[branch] = mineParents(l.Parent, ht[i].Count)
			branch++
			l = l.Link
		}

		base = append(base, pb)
	}

	return base
}