package cmd

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestGenerateDatabase(t *testing.T) {
	db, err := generateDatabase(100, 10)
	require.Nil(t, err)
	assert.Len(t, db, 100)
	for _, items := range db {
		assert.True(t, len(items) < 11)
		assert.True(t, len(items) > 0)
		for _, item := range items {
			assert.NotEmpty(t, item)
			// t.Logf("tx=%d i=%d v=%d", y, x, item)
		}
	}
}