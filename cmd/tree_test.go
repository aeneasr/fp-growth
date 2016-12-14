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
						&FPTreeNode{Item: A, Count: 1, Children: []*FPTreeNode{&FPTreeNode{Item: B, Count: 1, Children: []*FPTreeNode{}}}},
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
						&FPTreeNode{Item: B, Count: 2, Children: []*FPTreeNode{&FPTreeNode{Item: A, Count: 1, Children: []*FPTreeNode{}}}},
						// &FPTreeNode{Item: B, Count: 2, Children: []*FPTreeNode{}},
						//
						// &FPTreeNode{Item: A, Count: 1, Children: []*FPTreeNode{&FPTreeNode{Item: B, Count: 1, Children: []*FPTreeNode{}}}},
						// &FPTreeNode{Item: B, Count: 1, Children: []*FPTreeNode{}},
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
						&FPTreeNode{Item: A, Count: 2, Children: []*FPTreeNode{&FPTreeNode{Item: B, Count: 1, Children: []*FPTreeNode{}}}},
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
							&FPTreeNode{Item: B, Count: 2, Children: []*FPTreeNode{
								&FPTreeNode{Item: C, Count: 2, Children: []*FPTreeNode{}},
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
							&FPTreeNode{Item: C, Count: 4, Children: []*FPTreeNode{
								&FPTreeNode{Item: B, Count: 3, Children: []*FPTreeNode{
									&FPTreeNode{Item: D, Count: 1, Children: []*FPTreeNode{
										&FPTreeNode{Item: E, Count: 1, Children: []*FPTreeNode{}},
									}},
								}},
								&FPTreeNode{Item: E, Count: 1, Children: []*FPTreeNode{
									&FPTreeNode{Item: F, Count: 1, Children: []*FPTreeNode{}},
								}},
							}},
							&FPTreeNode{Item: D, Count: 3, Children: []*FPTreeNode{
								&FPTreeNode{Item: E, Count: 1, Children: []*FPTreeNode{}},
								&FPTreeNode{Item: F, Count: 1, Children: []*FPTreeNode{}},
							}},
						}},
						&FPTreeNode{Item: B, Count: 2, Children: []*FPTreeNode{
							&FPTreeNode{Item: D, Count: 1, Children: []*FPTreeNode{
								&FPTreeNode{Item: E, Count: 1, Children: []*FPTreeNode{}},
							}},
							&FPTreeNode{Item: E, Count: 1, Children: []*FPTreeNode{}},
						}},
						&FPTreeNode{Item: C, Count: 2, Children: []*FPTreeNode{
							&FPTreeNode{Item: F, Count: 1, Children: []*FPTreeNode{}},
						}},
					},
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			h := NewHeadTable(c.db, c.minSup)
			t.Logf("Head table: %v", h)
			ordered := OrderItems(c.db, h)

			a := NewFPTree(ordered)
			assert.True(t, treeEquals(t, c.e.Root, a.Root))
			// clean up parent things, somehow compare without order?
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

func treeEquals(t *testing.T, expect, actual *FPTreeNode) bool {
	if expect.Item != actual.Item || expect.Count != actual.Count || len(expect.Children) != len(actual.Children) {
		t.Logf("reason=notequal expected=%v+ actual=%v+", expect, actual)
		return false
	}

	for _, v1 := range expect.Children {
		var found bool
		for _, v2 := range actual.Children {
			if v2.Item == v1.Item {
				if !treeEquals(t, v1, v2) {
					t.Logf("reason=branchfail expected=%v actual=%v", expect, actual)
					return false
				}
				found = true
			}
		}
		if !found {
			t.Logf("reason=nonexistent search=%v+ expected=%v actual=%v", v1, expect.Children, actual.Children)
			return false
		}
	}

	return true
}

//
//func TestNewFPTree(t *testing.T) {
//	head := NewHeadTable(exDB, 2)
//	ordered := OrderItems(exDB, head)
//	e := FPTree{
//		Root: &FPTreeNode{
//			Children: []*FPTreeNode{
//				&FPTreeNode{Item: C, Count: 2, Children: []*FPTreeNode{
//					&FPTreeNode{Item: F, Count: 1},
//				}},
//				&FPTreeNode{Item: A, Count: 7, Children: []*FPTreeNode{
//					&FPTreeNode{Item: C, Count: 4, Children: []*FPTreeNode{
//						&FPTreeNode{Item: E, Count: 1, Children: []*FPTreeNode{
//							&FPTreeNode{Item: F, Count: 1},
//						}},
//						&FPTreeNode{Item: B, Count: 3, Children: []*FPTreeNode{
//							&FPTreeNode{Item: D, Count: 1, Children: []*FPTreeNode{
//								&FPTreeNode{Item: E, Count: 1},
//							}},
//						}},
//					}},
//					&FPTreeNode{Item: D, Count: 3, Children: []*FPTreeNode{
//						&FPTreeNode{Item: E, Count: 1},
//						&FPTreeNode{Item: F, Count: 1},
//					}},
//				}},
//				&FPTreeNode{Item: B, Count: 2, Children: []*FPTreeNode{
//					&FPTreeNode{Item: F, Count: 1},
//					&FPTreeNode{Item: E, Count: 1, Children: []*FPTreeNode{
//						&FPTreeNode{Item: E, Count: 1},
//					}},
//				}},
//			},
//		},
//	}
//
//	a := NewFPTree(ordered)
//	assert.EqualValues(t, e, a)
//}
