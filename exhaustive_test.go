package tidbtest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCNum(t *testing.T) {
	s, r := cNum(0, 2, []int{1, 2, 3, 4})
	assert.Equal(t, [][]int{
		{1, 2},
		{1, 3},
		{1, 4},
		{2, 3},
		{2, 4},
		{3, 4},
	}, s)
	assert.Equal(t, [][]int{
		{3, 4},
		{2, 4},
		{2, 3},
		{1, 4},
		{1, 3},
		{1, 2},
	}, r)
}

func TestShuffle(t *testing.T) {
	e := &exhaustive{}
	turns := e.Shuffle(map[string]int{
		"A": 1,
		"B": 1,
		"C": 2,
	})
	assert.Equal(t, [][]string{
		{"A", "B", "C", "C"},
		{"A", "C", "B", "C"},
		{"A", "C", "C", "B"},

		{"B", "A", "C", "C"},
		{"C", "A", "B", "C"},
		{"C", "A", "C", "B"},

		{"B", "C", "A", "C"},
		{"C", "B", "A", "C"},
		{"C", "C", "A", "B"},

		{"B", "C", "C", "A"},
		{"C", "B", "C", "A"},
		{"C", "C", "B", "A"},
	}, turns)
}
