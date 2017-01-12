package cmd

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
	"time"
	"encoding/json"
)

var g_cbp = []ConditionalPatternBases{
	ConditionalPatternBases{
		{
			Prefix: ConditionalItem{Item: B, Count: 1},
			Bases: [][]ConditionalItem{{}},
		},
		{
			Prefix: ConditionalItem{Item: A, Count: 1},
			Bases: [][]ConditionalItem{{}},
		},
	},
	ConditionalPatternBases{
		{
			Prefix: ConditionalItem{Item: B, Count: 1},
			Bases: [][]ConditionalItem{
				{ConditionalItem{Item: A, Count: 1}},
			},
		},
		{
			Prefix: ConditionalItem{Item: A, Count: 1},
			Bases: [][]ConditionalItem{{}},
		},
	},
	ConditionalPatternBases{
		{
			Prefix: ConditionalItem{Item: B, Count: 2},
			Bases: [][]ConditionalItem{
				{ConditionalItem{Item: A, Count: 1}},
			},
		},
		{
			Prefix: ConditionalItem{Item: A, Count: 2},
			Bases: [][]ConditionalItem{{}},
		},
	},
	ConditionalPatternBases{
		{
			Prefix: ConditionalItem{Item: C, Count: 2},
			Bases: [][]ConditionalItem{
				{
					ConditionalItem{Item: B, Count: 2},
					ConditionalItem{Item: A, Count: 2},
				},
			},
		},
		{
			Prefix: ConditionalItem{Item: B, Count: 2},
			Bases: [][]ConditionalItem{
				{
					ConditionalItem{Item: A, Count: 2},
				},
			},
		},
		{
			Prefix: ConditionalItem{Item: A, Count: 2},
			Bases: [][]ConditionalItem{{}                                        },
		},
	},
	ConditionalPatternBases{
		{
			Prefix: ConditionalItem{Item: F, Count: 3},
			Bases: [][]ConditionalItem{
				{
					ConditionalItem{Item: D, Count: 1},
					ConditionalItem{Item: A, Count: 1},
				},
				{
					ConditionalItem{Item: E, Count: 1},
					ConditionalItem{Item: C, Count: 1},
					ConditionalItem{Item: A, Count: 1},
				},
				{
					ConditionalItem{Item: C, Count: 1},
				},
			},
		},
		{
			Prefix: ConditionalItem{Item: E, Count: 5},
			Bases: [][]ConditionalItem{
				{
					ConditionalItem{Item: D, Count: 1},
					ConditionalItem{Item: B, Count: 1},
					ConditionalItem{Item: C, Count: 1},
					ConditionalItem{Item: A, Count: 1},
				},
				{
					ConditionalItem{Item: D, Count: 1},
					ConditionalItem{Item: A, Count: 1},
				},
				{
					ConditionalItem{Item: D, Count: 1},
					ConditionalItem{Item: B, Count: 1},
				},
				{
					ConditionalItem{Item: C, Count: 1},
					ConditionalItem{Item: A, Count: 1},
				},
				{
					ConditionalItem{Item: B, Count: 1},
				},
			},
		},
		{
			Prefix: ConditionalItem{Item: D, Count: 5},
			Bases: [][]ConditionalItem{
				{
					ConditionalItem{Item: B, Count: 1},
					ConditionalItem{Item: C, Count: 1},
					ConditionalItem{Item: A, Count: 1},
				},
				{
					ConditionalItem{Item: A, Count: 3},
				},
				{
					ConditionalItem{Item: B, Count: 1},
				},
			},
		},
		{
			Prefix: ConditionalItem{Item: B, Count: 5},
			Bases: [][]ConditionalItem{
				{
					ConditionalItem{Item: C, Count: 3},
					ConditionalItem{Item: A, Count: 3},
				},
			},
		},
		{
			Prefix: ConditionalItem{Item: C, Count: 6},
			Bases: [][]ConditionalItem{
				{ConditionalItem{Item: A, Count: 4}},
			},
		},
		{
			Prefix: ConditionalItem{Item: A, Count: 7},
			Bases: [][]ConditionalItem{{}},
		},
	},
}

func TestMineConditionalPatternBases(t *testing.T) {
	for k, c := range []struct {
		db     DataSet
		minSup float32
		e      ConditionalPatternBases
	}{
		{
			db: DataSet{
				Items{A},
				Items{B},
			},
			minSup: 1,
			e: g_cbp[0],
		},
		{
			db: DataSet{
				Items{A, B},
			},
			minSup: 1,
			e:  g_cbp[1],
		},
		{
			db: DataSet{
				Items{A},
				Items{A, B},
				Items{B},
			},
			minSup: 1,
			e:g_cbp[2],
		},
		{
			db: DataSet{
				Items{A, B, C},
				Items{A, B, C},
			},
			minSup: 1,
			e: g_cbp[3],
		},
		{
			db: exDB,
			minSup: 2,
			e: g_cbp[4],
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			h := NewHeadTable(c.db, c.minSup)
			ordered := OrderItems(c.db, h)
			fpt := NewFPTree(ordered, &h)
			cbp := MineConditionalPatternBases(fpt, h)

			assert.EqualValues(t, c.e, cbp)
		})
	}

}

func TestBenchmarkMining(t *testing.T) {
	t.SkipNow()

	times := map[string]float64{}
	minSups := []float32{
		0.5,
		0.4,
		0.03,
		0.02,
		0.01,
		0.005,
		0.001,
		0.0001,
		0.00001,
	}
	txss := []int{
		//10,
		//100,
		//10000,
		1000,
	}
	uss := []int{
		100,
	}

	var dbs = map[float32]map[int]map[int]DataSet{}
	for _, minSup := range minSups {
		for _, txs := range txss {
			for _, us := range uss {
				dbs[minSup] = make(map[int]map[int]DataSet)
				dbs[minSup][txs] = make(map[int]DataSet)
				dbs[minSup][txs][us], _ = generateDatabase(txs, us)
			}
		}
	}

	for _, minSup := range minSups {
		for _, txs := range txss {
			for _, us := range uss {
				d := fmt.Sprintf("minsup=%f/transactions=%d/items=%d/minItems=%f", minSup, txs, us, float32(txs) * minSup)
				db := dbs[minSup][txs][us]
				t.Run(d, func(t *testing.T) {
					start := time.Now()
					h := NewHeadTable(db, float32(txs) * minSup)
					ordered := OrderItems(db, h)
					fpt := NewFPTree(ordered, &h)
					_ = MineConditionalPatternBases(fpt, h)
					end := time.Now()
					t.Logf("Iteration took %f seconds", end.Sub(start).Seconds())
					times[d] = end.Sub(start).Seconds()
				})
			}
		}
	}
	out, _ := json.MarshalIndent(times, "", "\t")
	t.Logf("%s", string(out))
}

func TestConstructConditionalHeadTables(t *testing.T) {
	for k, c := range []struct {
		db     DataSet
		minSup int
		bs     ConditionalPatternBases
		e      []ConditionalHeadTable
	}{
		{
			minSup: 1,
			bs: g_cbp[0],
			e: []ConditionalHeadTable{
				{
					Prefix: ConditionalItem{Item: B, Count: 1},
					HeadTable: HeadTable{},
				},
				{
					Prefix: ConditionalItem{Item: A, Count: 1},
					HeadTable: HeadTable{},
				},
			},
		},
		{
			minSup: 1,
			bs: g_cbp[1],
			e: []ConditionalHeadTable{
				{
					Prefix: ConditionalItem{Item: B, Count: 1},
					HeadTable: HeadTable{
						{Item: A, Count: 1},
					},
				},
				{
					Prefix: ConditionalItem{Item: A, Count: 1},
					HeadTable: HeadTable{},
				},
			},
		},
		{
			minSup: 1,
			bs: g_cbp[2],
			e: []ConditionalHeadTable{
				{
					Prefix: ConditionalItem{Item: B, Count: 2},
					HeadTable: HeadTable{
						{Item: A, Count: 1},
					},
				},
				{
					Prefix: ConditionalItem{Item: A, Count: 2},
					HeadTable: HeadTable{},
				},
			},
		},
		{
			minSup: 1,
			bs: g_cbp[3],
			e: []ConditionalHeadTable{
				{
					Prefix: ConditionalItem{Item: C, Count: 2},
					HeadTable: HeadTable{
						{Item: A, Count: 2},
						{Item: B, Count: 2},
					},
				},
				{
					Prefix: ConditionalItem{Item: B, Count: 2},
					HeadTable: HeadTable{
						{Item: A, Count: 2},
					},
				},
				{
					Prefix: ConditionalItem{Item: A, Count: 2},
					HeadTable: HeadTable{},
				},
			},
		},
		{
			minSup: 2,
			bs: g_cbp[4],
			e: []ConditionalHeadTable{
				{
					Prefix: ConditionalItem{Item: F, Count: 3},
					HeadTable: HeadTable{
						{Item: A, Count: 2},
						{Item: C, Count: 2},
					},
				},
				{
					Prefix: ConditionalItem{Item: E, Count: 5},
					HeadTable: HeadTable{
						{Item: A, Count: 3},
						{Item: B, Count: 3},
						{Item: D, Count: 3},
						{Item: C, Count: 2},
					},
				},
				{
					Prefix: ConditionalItem{Item: D, Count: 5},
					HeadTable: HeadTable{
						{Item: A, Count: 4},
						{Item: B, Count: 2},
					},
				},
				{
					Prefix: ConditionalItem{Item: B, Count: 5},
					HeadTable: HeadTable{
						{Item: A, Count: 3},
						{Item: C, Count: 3},
					},
				},
				{
					Prefix: ConditionalItem{Item: C, Count: 6},
					HeadTable: HeadTable{
						{Item: A, Count: 4},
					},
				},
				{
					Prefix: ConditionalItem{Item: A, Count: 7},
					HeadTable: HeadTable{},
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			tables := ConstructConditionalHeadTables(c.bs, c.minSup)
			assert.EqualValues(t, c.e, tables)
		})
	}
}

func TestOrderConditionalPatternBases(t *testing.T) {
	for k, c := range []struct {
		db     DataSet
		minSup int
		bs     ConditionalPatternBases
		e      ConditionalPatternBases
	}{
		{
			minSup: 1,
			bs: g_cbp[0],
			e: g_cbp[0],
		},
		{
			minSup: 1,
			bs: g_cbp[1],
			e: g_cbp[1],
		},
		{
			minSup: 1,
			bs: g_cbp[2],
			e: g_cbp[2],
		},
		{
			minSup: 1,
			bs: g_cbp[3],
			e:
			ConditionalPatternBases{
				{
					Prefix: ConditionalItem{Item: C, Count: 2},
					Bases: [][]ConditionalItem{
						{
							ConditionalItem{Item: A, Count: 2},
							ConditionalItem{Item: B, Count: 2},
						},
					},
				},
				{
					Prefix: ConditionalItem{Item: B, Count: 2},
					Bases: [][]ConditionalItem{
						{
							ConditionalItem{Item: A, Count: 2},
						},
					},
				},
				{
					Prefix: ConditionalItem{Item: A, Count: 2},
					Bases: [][]ConditionalItem{{}                                        },
				},
			},
		},
		{
			minSup: 2,
			bs: g_cbp[4],
			e: ConditionalPatternBases{
				{
					Prefix: ConditionalItem{Item: F, Count: 3},
					Bases: [][]ConditionalItem{
						{
							ConditionalItem{Item: A, Count: 1},
						},
						{
							ConditionalItem{Item: A, Count: 1},
							ConditionalItem{Item: C, Count: 1},
						},
						{
							ConditionalItem{Item: C, Count: 1},
						},
					},
				},
				{
					Prefix: ConditionalItem{Item: E, Count: 5},
					Bases: [][]ConditionalItem{
						{
							ConditionalItem{Item: A, Count: 1},
							ConditionalItem{Item: B, Count: 1},
							ConditionalItem{Item: D, Count: 1},
							ConditionalItem{Item: C, Count: 1},
						},
						{
							ConditionalItem{Item: A, Count: 1},
							ConditionalItem{Item: D, Count: 1},
						},
						{
							ConditionalItem{Item: B, Count: 1},
							ConditionalItem{Item: D, Count: 1},
						},
						{
							ConditionalItem{Item: A, Count: 1},
							ConditionalItem{Item: C, Count: 1},
						},
						{
							ConditionalItem{Item: B, Count: 1},
						},
					},
				},
				{
					Prefix: ConditionalItem{Item: D, Count: 5},
					Bases: [][]ConditionalItem{
						{
							ConditionalItem{Item: A, Count: 1},
							ConditionalItem{Item: B, Count: 1},
						},
						{
							ConditionalItem{Item: A, Count: 3},
						},
						{
							ConditionalItem{Item: B, Count: 1},
						},
					},
				},
				{
					Prefix: ConditionalItem{Item: B, Count: 5},
					Bases: [][]ConditionalItem{
						{
							ConditionalItem{Item: A, Count: 3},
							ConditionalItem{Item: C, Count: 3},
						},
					},
				},
				{
					Prefix: ConditionalItem{Item: C, Count: 6},
					Bases: [][]ConditionalItem{
						{ConditionalItem{Item: A, Count: 4}},
					},
				},
				{
					Prefix: ConditionalItem{Item: A, Count: 7},
					Bases: [][]ConditionalItem{{}},
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			ht := ConstructConditionalHeadTables(c.bs, c.minSup)
			res := OrderConditionalPatternBases(c.bs, ht)
			assert.EqualValues(t, c.e, res)
		})
	}

}

func TestConstructConditionalFPTrees(t *testing.T) {
	for k, c := range []struct {
		db     DataSet
		minSup int
		bs     ConditionalPatternBases
		e      []ConditionalFPTree
	}{
		{
			minSup: 1,
			bs: g_cbp[0],
			e: []ConditionalFPTree{
				{
					Prefix: ConditionalItem{Item: B, Count: 1},
					Tree: FPTree{Root: &FPTreeNode{}},
				},
				{
					Prefix: ConditionalItem{Item: A, Count: 1},
					Tree: FPTree{Root: &FPTreeNode{}},
				},
			},
		},
		{
			minSup: 1,
			bs: g_cbp[1],
			e: []ConditionalFPTree{
				{
					Prefix: ConditionalItem{Item: B, Count: 1},
					Tree: FPTree{
						Root: &FPTreeNode{
							Children: []*FPTreeNode{
								&FPTreeNode{Item: A, Count: 1, Children: []*FPTreeNode{}},
							},
						},
					},
				},
				{
					Prefix: ConditionalItem{Item: A, Count: 1},
					Tree: FPTree{Root: &FPTreeNode{}                                },
				},
			},
		},
		{
			minSup: 1,
			bs: g_cbp[2],
			e: []ConditionalFPTree{
				{
					Prefix: ConditionalItem{Item: B, Count: 2},
					Tree: FPTree{
						Root: &FPTreeNode{
							Children: []*FPTreeNode{
								&FPTreeNode{Item: A, Count: 1, Children: []*FPTreeNode{}},
							},
						},
					},
				},
				{
					Prefix: ConditionalItem{Item: A, Count: 2},
					Tree: FPTree{Root: &FPTreeNode{}                                },
				},
			},
		},
		{
			minSup: 1,
			bs: g_cbp[3],
			e: []ConditionalFPTree{
				{
					Prefix: ConditionalItem{Item: C, Count: 2},
					Tree: FPTree{
						Root: &FPTreeNode{
							Children: []*FPTreeNode{
								&FPTreeNode{
									Item: B, Count: 2, Children: []*FPTreeNode{
										&FPTreeNode{Item: A,  Parent: &FPTreeNode{Item: B, Count: 2}, Count: 2, Children: []*FPTreeNode{}},
									},
								},
							},
						},
					},
				},/*
				{
					Prefix: ConditionalItem{Item: B, Count: 2},
					Tree: FPTree{
						Root: &FPTreeNode{
							Children: []*FPTreeNode{
								&FPTreeNode{Item: A, Count: 2, Children: []*FPTreeNode{}},
							},
						},
					},
				},
				{
					Prefix: ConditionalItem{Item: A, Count: 2},
					Tree: FPTree{Root: &FPTreeNode{}                                },
				},*/
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			ht := ConstructConditionalHeadTables(c.bs, c.minSup)
			bs := OrderConditionalPatternBases(c.bs, ht)
			trees := ConstructConditionalFPTrees(bs, ht)
			assert.Equal(t, len(c.e), len(trees))
			for i, tree := range trees {
				if i >= len(c.e) {
					assert.True(t, false)
					return
				}
				assert.Equal(t, c.e[i].Prefix, tree.Prefix)
				assert.True(t, treeEquals(t, c.e[i].Tree.Root, tree.Tree.Root))
			}
		})
	}

}