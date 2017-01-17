package cmd

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