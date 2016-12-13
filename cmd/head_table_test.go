package cmd

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

const A = 1
const B = 2
const C = 3
const D = 4
const E = 5
const F = 6
const G = 7

var exDB = DataSet{
	Items{A,B,C,D,E},
	Items{A,C,B},
	Items{F,A,D},
	Items{A,G,E,D},
	Items{D,E,B},
	Items{F,A,C,E},
	Items{E,B},
	Items{C,A,B},
	Items{A,D},
	Items{F,C},
	Items{C},
}

var exCount = map[int]int{
	A: 7,
	B: 5,
	C: 6,
	D: 5,
	E: 5,
	F: 3,
	G: 1,
}

var exOrder = []int{A, C, B, D, E, F}

func TestNewHeadTable(t *testing.T) {
	expected := HeadTable{
		HeadTableRow{
			Item: A,
			Count: 7,
		},
		HeadTableRow{
			Item: C,
			Count: 6,
		},
		HeadTableRow{
			Item: B,
			Count: 5,
		},
		HeadTableRow{
			Item: D,
			Count: 5,
		},
		HeadTableRow{
			Item: E,
			Count: 5,
		},
		HeadTableRow{
			Item: F,
			Count: 3,
		},
	}
	actual := NewHeadTable(exDB, 2)
	assert.EqualValues(t, expected, actual)
}
