package cmd

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestConstructSupportCountTable(t *testing.T) {
	for k, c := range []struct {
		bs []Items
		e  SupportCountTable
	}{
		{
			bs:        []Items{
				Items{A},
				Items{A, B, E                                },
				Items{B, D                                },
				Items{B, C                                },
				Items{A, C                                },
				Items{A, B, D                                },
				Items{A, B, C, E                                },
				Items{A, B, C                                },
				Items{A, D, E                                },
				Items{A, F                                },
				Items{B, C, G                                },
				Items{C, E, F                                },
			},
			e: SupportCountTable{
				A: {A: 0, B:4, C:3, D:2, E:3, F:1, G:0},
				B: {A: 0, B:0, C:4, D:2, E:2, F:0, G:1},
				C: {A: 0, B:0, C:0, D:0, E:2, F:1, G:1},
				D: {A: 0, B:0, C:0, D:0, E:1, F:0, G:0},
				E: {A: 0, B:0, C:0, D:0, E:0, F:1, G:0},
				F: {A: 0, B:0, C:0, D:0, E:0, F:0, G:0},
				G: {A: 0, B:0, C:0, D:0, E:0, F:0, G:0},
			},
		},
		{
			bs:        []Items{
				Items{C},
				Items{A, C, E},
				Items{A, D},
			},
			e: SupportCountTable{
				A: {A: 0, C:1, D:1, E:1},
				C: {A: 0, C:0, D:0, E:1},
				D: {A: 0, C:0, D:0, E:0},
				E: {A: 0, C:0, D:0, E:0},
			},
		},
		{
			bs:        []Items(exDB),
			e: SupportCountTable{
				A: {A: 0, B:3, C:3, D:4, E:3, F:0, G:1},
				B: {A: 0, B:0, C:1, D:1, E:1, F:0, G:0},
				C: {A: 1, B:2, C:0, D:1, E:2, F:0, G:0},
				D: {A: 0, B:1, C:0, D:0, E:2, F:0, G:0},
				E: {A: 0, B:2, C:0, D:1, E:0, F:0, G:0},
				F: {A: 2, B:0, C:2, D:1, E:1, F:0, G:0},
				G: {A: 0, B:0, C:0, D:1, E:1, F:0, G:0},
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			out := ConstructSupportCountTable(c.bs)
			assert.Equal(t, len(c.e), len(out))
			for i, cells := range c.e {
				//assert.Equal(t, cells, out[i], "lookup=%d",i)
				for j, exp := range cells {
					assert.EqualValues(t, exp, out.Get(i, j), "i=%d j=%d", i, j)
				}
			}
		})
	}
}

func TestConstructConditionalSupportCountTables(t *testing.T) {
	for k, c := range []struct {
		bs ConditionalPatternBases
		e  map[int]SupportCountTable
	}{
		{
			bs: g_cbp[4],
			e: map[int]SupportCountTable{
				F: {
					//A: {A: 0, B:0, C:1, D:1, E:1, F:0, G:0},
					//B: {A: 0, B:0, C:0, D:0, E:0, F:0, G:0},
					C: {A: 1, C:0, D:0, E:0},
					D: {A: 1, C:0, D:0, E:0},
					E: {A: 1, C:1, D:0, E:0},
					A: {A: 0, C:0, D:0, E:0},
					//G: {A: 0, B:0, C:0, D:0, E:0, F:0, G:0},
				},
				//E: {
				//	D: {A: 2, B:2, C:1, D:0, E:0, F:0, G:0},
				//	B: {A: 1, B:0, C:1, D:0, E:0, F:0, G:0},
				//	C: {A: 2, B:0, C:0, D:0, E:0, F:0, G:0},
				//	//A: {A: 0, B:0, C:0, D:0, E:0, F:0, G:0},
				//	//D: {A: 0, B:0, C:0, D:0, E:0, F:0, G:0},
				//	//D: {A: 0, B:0, C:0, D:0, E:0, F:0, G:0},
				//	//D: {A: 0, B:0, C:0, D:0, E:0, F:0, G:0},
				//},
				//D: {
				//	B: {A: 1, B:0, C:1, D:0, E:0, F:0, G:0},
				//	C: {A: 1, B:0, C:0, D:0, E:0, F:0, G:0},
				//
				//},
				//B: {
				//	C: {A: 3, B:0, C:0, D:0, E:0, F:0, G:0},
				//},
				//C: {},
				//A: {},
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			out := ConstructConditionalSupportCountTables(c.bs)
			//assert.Equal(t, len(c.e), len(out))

			for f := range c.e {
				for i, cells := range c.e[f] {
					//assert.Equal(t, cells, out[f][i], "lookup=%d", i)
					for j, exp := range cells {
						assert.EqualValues(t, exp, out[f].Get(i, j), "f=%d i=%d j=%d", f, i, j)
					}
				}
			}
		})
	}
}

func TestConstructImprovedConditionalHeadTables(t *testing.T) {
	for k, c := range []struct {
		ds     DataSet
		minSup int
		e      ConditionalPatternBases
	}{
		{
			minSup: 2,
			ds: exDB,
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
			h := NewHeadTable(c.ds, c.minSup)
			OrderItems(c.ds, h)
			fpt := NewFPTree(c.ds, &h)
			counts := ConstructSupportCountTable(c.ds)
			cpbs := MineConditionalPatternBases(fpt, h)
			tables := ConstructImprovedConditionalHeadTables(cpbs, counts, c.minSup)
			res := OrderConditionalPatternBases(cpbs, tables)
			assert.EqualValues(t, c.e, res)
		})
	}
}
