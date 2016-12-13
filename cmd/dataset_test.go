package cmd

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestOrderItems(t *testing.T) {
	head := NewHeadTable(exDB, 2)
	a := OrderItems(exDB, head)
	e := DataSet{
		Items{A,C,B,D,E},
		Items{A,C,B},
		Items{A,D,F},
		Items{A,D,E},
		Items{B,D,E},
		Items{A,C,E,F},
		Items{B,E},
		Items{A,C,B},
		Items{A,D},
		Items{C,F},
		Items{C},
	}
	assert.EqualValues(t, e, a)
}