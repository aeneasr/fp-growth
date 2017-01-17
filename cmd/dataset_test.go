package cmd

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestOrderItems(t *testing.T) {
	db := DataSet{
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

	head := NewHeadTable(db, 2)
	OrderItems(db, head)
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
	assert.EqualValues(t, e, db)
}