package cmd

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestNewFPTree(t *testing.T) {
	for k, c := range []struct {
		db     DataSet
		minSup float32
		e      FPTree
	}{
		{
			db: DataSet{
				Items{A},
				Items{B},
			},
			minSup: 1,
			e: FPTree{
				Root: &FPTreeNode{
					Children: []*FPTreeNode{
						&FPTreeNode{Item: A, Count: 1, Children: []*FPTreeNode{}},
						&FPTreeNode{Item: B, Count: 1, Children: []*FPTreeNode{}},
					},
				},
			},
		},
		{
			db: DataSet{
				Items{A, B},
			},
			minSup: 1,
			e: FPTree{
				Root: &FPTreeNode{
					Children: []*FPTreeNode{
						&FPTreeNode{Item: A, Count: 1, Children: []*FPTreeNode{&FPTreeNode{Item: B, Parent: &FPTreeNode{Item: A, Count: 1}, Count: 1, Children: []*FPTreeNode{}}}},
					},
				},
			},
		},
		{
			db: DataSet{
				Items{A, B},
				Items{B},
			},
			minSup: 1,
			e: FPTree{
				Root: &FPTreeNode{
					Children: []*FPTreeNode{
						&FPTreeNode{Item: B, Count: 2, Children: []*FPTreeNode{&FPTreeNode{Item: A, Parent: &FPTreeNode{Item: B, Count: 2}, Count: 1, Children: []*FPTreeNode{}}}},
					},
				},
			},
		},
		{
			db: DataSet{
				Items{A},
				Items{A, B},
				Items{B},
			},
			minSup: 1,
			e: FPTree{
				Root: &FPTreeNode{
					Children: []*FPTreeNode{
						&FPTreeNode{Item: A, Count: 2, Children: []*FPTreeNode{&FPTreeNode{Item: B, Parent: &FPTreeNode{Item: A, Count: 2}, Count: 1, Children: []*FPTreeNode{}, Link: &FPTreeNode{Item: B, Count: 1}}}},
						&FPTreeNode{Item: B, Count: 1, Children: []*FPTreeNode{}},
					},
				},
			},
		},
		{
			db: DataSet{
				Items{A, B, C},
				Items{A, B, C},
			},
			minSup: 1,
			e: FPTree{
				Root: &FPTreeNode{
					Children: []*FPTreeNode{
						&FPTreeNode{Item: A, Count: 2, Children: []*FPTreeNode{
							&FPTreeNode{Item: B, Parent: &FPTreeNode{Item: A, Count: 2}, Count: 2, Children: []*FPTreeNode{
								&FPTreeNode{Item: C, Parent: &FPTreeNode{Item: B, Count: 2}, Count: 2, Children: []*FPTreeNode{}},
							}},
						}},
					},
				},
			},
		},
		{
			db: exDB,
			minSup: 2,
			e: FPTree{
				Root: &FPTreeNode{
					Children: []*FPTreeNode{
						&FPTreeNode{Item: A, Count: 7, Children: []*FPTreeNode{
							&FPTreeNode{Item: C, Count: 4, Parent: &FPTreeNode{Item: A, Count: 7}, Link: &FPTreeNode{Item: C, Count: 2}, Children: []*FPTreeNode{
								&FPTreeNode{Item: B, Count: 3, Parent: &FPTreeNode{Item: C, Count: 4}, Link: &FPTreeNode{Item: B, Count: 2}, Children: []*FPTreeNode{
									&FPTreeNode{Item: D, Count: 1, Parent: &FPTreeNode{Item: B, Count: 3}, Link: &FPTreeNode{Item: D, Count: 3}, Children: []*FPTreeNode{
										&FPTreeNode{Item: E, Count: 1, Parent: &FPTreeNode{Item: D, Count: 1}, Link: &FPTreeNode{Item: E, Count: 1}, Children: []*FPTreeNode{}},
									}},
								}},
								&FPTreeNode{Item: E, Count: 1, Parent: &FPTreeNode{Item: C, Count: 4}, Link: &FPTreeNode{Item: E, Count: 1}, Children: []*FPTreeNode{
									&FPTreeNode{Item: F, Count: 1, Parent: &FPTreeNode{Item: E, Count: 1}, Link: &FPTreeNode{Item: F, Count: 1}, Children: []*FPTreeNode{}},
								}},
							}},
							&FPTreeNode{Item: D, Count: 3, Parent: &FPTreeNode{Item: A, Count: 7}, Link: &FPTreeNode{Item: D, Count: 1}, Children: []*FPTreeNode{
								&FPTreeNode{Item: E, Count: 1, Parent: &FPTreeNode{Item: D, Count: 3}, Link: &FPTreeNode{Item: E, Count: 1}, Children: []*FPTreeNode{}},
								&FPTreeNode{Item: F, Count: 1, Parent: &FPTreeNode{Item: D, Count: 3}, Link: &FPTreeNode{Item: F, Count: 1}, Children: []*FPTreeNode{}},
							}},
						}},
						&FPTreeNode{Item: B, Count: 2, Children: []*FPTreeNode{
							&FPTreeNode{Item: D, Count: 1, Parent: &FPTreeNode{Item: B, Count: 2}, Children: []*FPTreeNode{
								&FPTreeNode{Item: E, Count: 1, Parent: &FPTreeNode{Item: D, Count: 1}, Link: &FPTreeNode{Item: E, Count: 1}, Children: []*FPTreeNode{}},
							}},
							&FPTreeNode{Item: E, Count: 1, Parent: &FPTreeNode{Item: B, Count: 2}, Children: []*FPTreeNode{}},
						}},
						&FPTreeNode{Item: C, Count: 2, Children: []*FPTreeNode{
							&FPTreeNode{Item: F, Count: 1, Parent: &FPTreeNode{Item: C, Count: 2}, Children: []*FPTreeNode{}},
						}},
					},
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			h := NewHeadTable(c.db, c.minSup)
			ordered := OrderItems(c.db, h)

			a := NewFPTree(ordered, &h)
			assert.True(t, treeEquals(t, c.e.Root, a.Root))

			for _, r := range h {
				assert.NotNil(t, r.Link)
			}
		})
	}
}

func TestOnlyOneBranch(t *testing.T) {
	for k, c := range []struct {
		Root *FPTreeNode
		e    bool
	}{
		{
			Root: &FPTreeNode{Item: A, Count: 2, Children: []*FPTreeNode{}},
			e: true,
		},
		{
			Root: &FPTreeNode{Item: A, Count: 2, Children: []*FPTreeNode{&FPTreeNode{Item: A, Count: 2, Children: []*FPTreeNode{&FPTreeNode{Item: A, Count: 2, Children: []*FPTreeNode{}}}}}},
			e: true,
		},
		{
			Root: &FPTreeNode{
				Item: A, Count: 2, Children: []*FPTreeNode{
					&FPTreeNode{
						Item: A, Count: 2, Children: []*FPTreeNode{
							&FPTreeNode{Item: A, Count: 2, Children: []*FPTreeNode{}},
							&FPTreeNode{Item: A, Count: 2, Children: []*FPTreeNode{}},
						},
					},
				},
			},
			e: false,
		},
		{
			Root: &FPTreeNode{
				Item: A, Count: 2, Children: []*FPTreeNode{
					&FPTreeNode{
						Item: A, Count: 2, Children: []*FPTreeNode{
							&FPTreeNode{Item: A, Count: 2, Children: []*FPTreeNode{
								&FPTreeNode{Item: A, Count: 2, Children: []*FPTreeNode{}},
								&FPTreeNode{Item: A, Count: 2, Children: []*FPTreeNode{}},
							}},
						},
					},
				},
			},
			e: false,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.Equal(t, c.e, c.Root.OnlyOneBranch())
		})
	}
}

func TestTreeEquals(t *testing.T) {
	assert.True(t, treeEquals(t, &FPTreeNode{
		Children: []*FPTreeNode{
			&FPTreeNode{Item: A, Count: 2, Children: []*FPTreeNode{
				&FPTreeNode{Item: B, Count: 2, Children: []*FPTreeNode{}},
				&FPTreeNode{Item: C, Count: 2, Children: []*FPTreeNode{}},
			}},
		},
	}, &FPTreeNode{
		Children: []*FPTreeNode{
			&FPTreeNode{Item: A, Count: 2, Children: []*FPTreeNode{
				&FPTreeNode{Item: C, Count: 2, Children: []*FPTreeNode{}},
				&FPTreeNode{Item: B, Count: 2, Children: []*FPTreeNode{}},
			}},
		},
	}))
	assert.False(t, treeEquals(t, &FPTreeNode{
		Children: []*FPTreeNode{
			&FPTreeNode{Item: A, Count: 2, Children: []*FPTreeNode{
				&FPTreeNode{Item: B, Count: 2, Children: []*FPTreeNode{}},
				&FPTreeNode{Item: C, Count: 2, Children: []*FPTreeNode{}},
			}},
		},
	}, &FPTreeNode{
		Children: []*FPTreeNode{
			&FPTreeNode{Item: A, Count: 2, Children: []*FPTreeNode{
				&FPTreeNode{Item: C, Count: 2, Children: []*FPTreeNode{}},
				&FPTreeNode{Item: D, Count: 2, Children: []*FPTreeNode{}},
			}},
		},
	}))

}

func TestMinePatterns(t *testing.T) {
	for k, c := range []struct {
		p ConditionalItem
		i *FPTreeNode
		e []Pattern
	}{
		{
			p: ConditionalItem{Item: B, Count: 3},
			i: &FPTreeNode{
				Item: A,
				Count: 2,
				Children: []*FPTreeNode{},
			},
			e: []Pattern{
				{
					Count: 3,
					Pattern: Items{B},
				},
				{
					Count: 2,
					Pattern: Items{A, B},
				},
			},
		},
		{
			p: ConditionalItem{Item: D, Count: 3},
			i: &FPTreeNode{
				Item: F,
				Count: 3,
				Children: []*FPTreeNode{
					&FPTreeNode{
						Item: C,
						Count: 3,
						Children: []*FPTreeNode{
							&FPTreeNode{
								Item: A,
								Count: 3,
								Children: []*FPTreeNode{},
							},
						},
					},
				},
			},
			e: []Pattern{
				{
					Count: 3,
					Pattern: Items{D},
				},
				{
					Count: 3,
					Pattern: Items{F,D},
				},
				{
					Count: 3,
					Pattern: Items{C,D},
				},
				{
					Count: 3,
					Pattern: Items{A,D},
				},
				{
					Count: 3,
					Pattern: Items{F,C,D},
				},
				{
					Count: 3,
					Pattern: Items{F,A,D},
				},
				{
					Count: 3,
					Pattern: Items{C,A,D},
				},
				{
					Count: 3,
					Pattern: Items{F,C,A,D},
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			out := c.i.MinePatterns(c.p)
			assert.EqualValues(t, c.e, out)
		})
	}
}

func treeEquals(t *testing.T, expect, actual *FPTreeNode) bool {
	if expect.Item != actual.Item || expect.Count != actual.Count || len(expect.Children) != len(actual.Children) {
		t.Logf("reason=notequal\nexpected=%v+\nactual  =%v+", expect, actual)
		return false
	}

	if expect.Link == nil {
		if actual.Link != nil {
			t.Logf("reason=unexpectedlink\nexpected=%v+\nactual  =%v+", expect, actual)
			return false
		}
	} else {
		if actual.Link == nil {
			t.Logf("reason=expectedlink\nexpected=%v+\nactual  =%v+", expect, actual)
			return false
		}

		if actual.Link.Item != expect.Link.Item || actual.Link.Count != expect.Link.Count {
			t.Logf("reason=wronglink\nexpected=%v+\nactual  =%v+", expect, actual)
			return false
		}
	}

	if expect.Parent == nil {
		if actual.Parent != nil {
			t.Logf("reason=unexparent\nexpected=%v\nactual  =%v", expect, actual)
			return false
		}
	} else {
		if actual.Parent == nil {
			t.Logf("reason=parentmissing\nexpected=%v\nactual  =%v", expect, actual)
			return false
		}

		if actual.Parent.Item != expect.Parent.Item || actual.Parent.Count != expect.Parent.Count {
			t.Logf("reason=wrongparent\nexpected=%v\nactual  =%v", expect, actual)
			return false
		}
	}

	for _, v1 := range expect.Children {
		var found bool
		for _, v2 := range actual.Children {
			if v2.Item == v1.Item {
				if !treeEquals(t, v1, v2) {
					t.Logf("reason=branchfail\nexpected=%v\nactual  =%v", expect, actual)
					return false
				}
				found = true
			}
		}
		if !found {
			t.Logf("reason=nonexistent\nsearch=%v+\nexpected=%v\nactual=%v", v1, expect.Children, actual.Children)
			return false
		}
	}

	return true
}
