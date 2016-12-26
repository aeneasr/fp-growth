package cmd

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

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
			e: ConditionalPatternBases{
				{
					Prefix: ConditionalItem{Item: B, Count: 1},
					Bases: [][]ConditionalItem{{}},
				},
				{
					Prefix: ConditionalItem{Item: A, Count: 1},
					Bases: [][]ConditionalItem{{}},
				},
			},
		},
		{
			db: DataSet{
				Items{A, B},
			},
			minSup: 1,
			e: ConditionalPatternBases{
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
		},
		{
			db: DataSet{
				Items{A},
				Items{A, B},
				Items{B},
			},
			minSup: 1,
			e: ConditionalPatternBases{
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
		},
	}{
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			h := NewHeadTable(c.db, c.minSup)
			ordered := OrderItems(c.db, h)

			fpt := NewFPTree(ordered, &h)
			cbp := MineConditionalPatternBases(fpt, h)
			assert.EqualValues(t, c.e, cbp)
		})
	}

}
