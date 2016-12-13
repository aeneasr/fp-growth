package cmd

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestNewFPTree(t *testing.T) {

	for k, c := range []struct {
		db DataSet
		minSup float32
		e  FPTree
	}{
		{
			db: DataSet{
				Items{A},
			},
			minSup: 1,
			e: FPTree{
				Root: &FPTreeNode{
					Children: []*FPTreeNode{
						&FPTreeNode{Item: A, Count: 1, Children: []*FPTreeNode{}},
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
			db: DataSet{
				Items{A, B, C},
				Items{A, B, C},
				Items{A, B, D},
			},
			minSup: 1,
			e: FPTree{
				Root: &FPTreeNode{
					Children: []*FPTreeNode{
						&FPTreeNode{Item: A, Count: 3, Children: []*FPTreeNode{
							&FPTreeNode{Item: B, Count: 3, Children: []*FPTreeNode{
								&FPTreeNode{Item: C, Count: 2, Children: []*FPTreeNode{}},
								&FPTreeNode{Item: D, Count: 1, Children: []*FPTreeNode{}},
							}},
						}},
					},
				},
			},
		},

	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			h := NewHeadTable(exDB, c.minSup)
			ordered := OrderItems(c.db, h)

			a := NewFPTree(ordered)
			assert.EqualValues(t, c.e, a)

			// clean up parent things, somehow compare without order?
		})
	}
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
