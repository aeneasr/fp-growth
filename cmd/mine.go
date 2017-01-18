package cmd

import "sort"

//import "fmt"

//import "sort"

func Mine(db DataSet, minSup int) []Pattern {
	h := NewHeadTable(db, minSup)
	OrderItems(db, h)
	fpt := NewFPTree(db, &h)

	return MineAllPatterns(fpt, nil, minSup, h)
}

func MineAllPatterns(fpt FPTree, patterns []Pattern, minSup int, h HeadTable) []Pattern {
	cpbs := MineConditionalPatternBases(fpt, h)
	cht := ConstructConditionalHeadTables(cpbs, minSup)
	ocbps := OrderConditionalPatternBases(cpbs, cht)
	trees := ConstructConditionalFPTrees(ocbps, cht)

	for i, tree := range trees {
		if tree.Tree.Root.OnlyOneBranch() {
			patterns = append(patterns, tree.MineFrequentPatterns()...)
		} else {
			patterns = append(patterns, MineAllPatterns(tree.Tree, patterns, minSup, cht[i].HeadTable)...)
		}
	}
	return patterns
}

func MineImproved(db DataSet, minSup int) []Pattern {
	h := NewHeadTable(db, minSup)
	OrderItems(db, h)
	fpt := NewFPTree(db, &h)
	counts := ConstructSupportCountTable(db)

	return MineAllPatternsImproved(fpt, nil, minSup, h, counts)
}

func MineAllPatternsImproved(fpt FPTree, patterns []Pattern, minSup int, h HeadTable, counts SupportCountTable) []Pattern {
	cpbs := MineConditionalPatternBases(fpt, h)

	cht := ConstructImprovedConditionalHeadTables(cpbs, counts, minSup)
	ocbps := OrderConditionalPatternBases(cpbs, cht)
	condcs := ConstructConditionalSupportCountTables(ocbps)
	trees := ConstructConditionalFPTrees(ocbps, cht)

	for _, tree := range trees {
		if tree.Tree.Root.OnlyOneBranch() {
			patterns = append(patterns, tree.MineFrequentPatterns()...)
		} else {
			for _, ht := range cht {
				if ht.Prefix.Item == tree.Prefix.Item {
					patterns = append(patterns, MineAllPatternsImproved(tree.Tree, patterns, minSup, ht.HeadTable, condcs[tree.Prefix.Item])...)
				}
			}
		}
	}
	return patterns
}

type SupportCountTable map[int]map[int]int

func (s SupportCountTable) Get(i, j int) int {
	if cells, ok := s[i]; !ok {
		return -1
	} else {
		if count, ok := cells[j]; !ok {
			return -1
		} else {
			return count
		}
	}
}

func ConstructConditionalSupportCountTables(bases ConditionalPatternBases) map[int]SupportCountTable {
	var scts = map[int]SupportCountTable{}

	for _, base := range bases {
		if scts[base.Prefix.Item] == nil {
			scts[base.Prefix.Item] = SupportCountTable{}
		}
		var observed = map[int]bool{}

		sct := scts[base.Prefix.Item]
		for _, list := range base.Bases {
			for i := range list {
				observed[list[i].Item] = true
				if sct[list[i].Item] == nil {
					sct[list[i].Item] = map[int]int{}
				}
				for j := i + 1; j < len(list); j++ {
					observed[list[j].Item] = true
					sct[list[i].Item][list[j].Item] += list[j].Count
					//fmt.Printf("%d:%d i=%d j=%d count=%d\n", i, j, list[i].Item, list[j].Item, sct[list[i].Item][list[j].Item])
				}
			}
		}

		for o := range observed {
			if _, ok := sct[o]; !ok {
				sct[o] = map[int]int{}
			}
			for i := range sct {
				if _, ok := sct[i][o]; !ok {
					sct[i][o] = 0
				}
			}
		}
	}

	return scts
}

func ConstructSupportCountTable(items []Items) SupportCountTable {
	var sct = SupportCountTable{}
	var observed = map[int]bool{}

	for _, list := range items {
		for i := range list {
			observed[list[i]] = true
			if sct[list[i]] == nil {
				sct[list[i]] = map[int]int{}
			}
			for j := i + 1; j < len(list); j++ {
				sct[list[i]][list[j]] += 1
				observed[list[j]] = true
				//fmt.Printf("%d:%d i=%d j=%d count=%d\n", i, j, list[i].Item, list[j].Item, sct[list[i].Item][list[j].Item])
			}
		}
	}

	for o := range observed {
		if _, ok := sct[o]; !ok {
			sct[o] = map[int]int{}
		}
		for i := range sct {
			if _, ok := sct[i][o]; !ok {
				sct[i][o] = 0
			}
		}
	}

	return sct
}

func ConstructImprovedConditionalHeadTables(bs ConditionalPatternBases, scts  SupportCountTable, minSup int) ConditionalHeadTables {
	var tables ConditionalHeadTables
	for _, base := range bs {
		pl := HeadTable{}
		i := base.Prefix.Item
		for j := range scts[i] {
			count := scts.Get(i, j) + scts.Get(j, i)
			if count < minSup {
				continue
			}
			pl = append(pl, HeadTableRow{Item: j, Count:count})
		}

		sort.Sort(sort.Reverse(pl))
		tables = append(tables, ConditionalHeadTable{
			Prefix: base.Prefix,
			HeadTable: pl,
		})
	}

	return tables
}