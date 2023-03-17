package skip_list_test

import (
	"testing"

	skip_list "github.com/ForeverSRC/noelle/skip-list"
	"github.com/ForeverSRC/noelle/skip-list/level"
	"github.com/stretchr/testify/assert"
)

type TestLevelGenerator struct {
	index    int
	sequence []int
}

func (t *TestLevelGenerator) GetLevel() int {
	lvl := t.sequence[t.index%len(t.sequence)]
	t.index++
	return lvl
}

func (t *TestLevelGenerator) MaxLevel() int {
	return 6
}

func TestInsertAndDelete(t *testing.T) {
	elems := []skip_list.Element[string, int]{
		{Key: "a", Value: 3},
		{Key: "b", Value: 6},
		{Key: "d", Value: 9},
		{Key: "c", Value: 7},
		{Key: "e", Value: 12},
		{Key: "f", Value: 17},
		{Key: "g", Value: 19},
		{Key: "i", Value: 25},
		{Key: "j", Value: 26},
		{Key: "h", Value: 21},
		{Key: "f", Value: 18},
	}

	levels := []int{1, 4, 2, 1, 1, 2, 1, 3, 1, 1}

	expects := [][]skip_list.Element[string, int]{
		{
			{Key: "a", Value: 3},
			{Key: "b", Value: 6},
			{Key: "c", Value: 7},
			{Key: "d", Value: 9},
			{Key: "e", Value: 12},
			{Key: "f", Value: 18},
			{Key: "g", Value: 19},
			{Key: "h", Value: 21},
			{Key: "i", Value: 25},
			{Key: "j", Value: 26},
		},
		{
			{Key: "b", Value: 6},
			{Key: "d", Value: 9},
			{Key: "f", Value: 18},
			{Key: "i", Value: 25},
		},
		{
			{Key: "b", Value: 6},
			{Key: "i", Value: 25},
		},
		{
			{Key: "b", Value: 6},
		},
	}

	lg := &TestLevelGenerator{
		index:    0,
		sequence: levels,
	}

	sl := skip_list.NewSkipList[string, int](lg)
	for _, e := range elems {
		sl.Insert(e.Key, e.Value)
	}

	assert.Equal(t, 4, sl.GetLevel())
	for i, exps := range expects {
		lvl := i + 1
		res := sl.LevelIteration(lvl)
		assert.Equalf(t, exps, res, "not equal, level: %d", lvl)
	}

	sl.Delete("b")
	assert.Equal(t, 3, sl.GetLevel())
	expectsAfterDelete := [][]skip_list.Element[string, int]{
		{
			{Key: "a", Value: 3},
			{Key: "c", Value: 7},
			{Key: "d", Value: 9},
			{Key: "e", Value: 12},
			{Key: "f", Value: 18},
			{Key: "g", Value: 19},
			{Key: "h", Value: 21},
			{Key: "i", Value: 25},
			{Key: "j", Value: 26},
		},
		{
			{Key: "d", Value: 9},
			{Key: "f", Value: 18},
			{Key: "i", Value: 25},
		},
		{
			{Key: "i", Value: 25},
		},
	}

	for i, exps := range expectsAfterDelete {
		lvl := i + 1
		res := sl.LevelIteration(lvl)
		assert.Equalf(t, exps, res, "not equal, level: %d", lvl)
	}

}

func TestWork(t *testing.T) {
	cases := []struct {
		op            string
		elem          skip_list.Element[string, int]
		expected      int
		expectedFound bool
	}{
		{
			op:            "search",
			elem:          skip_list.Element[string, int]{Key: "ab"},
			expectedFound: false,
		},
		{
			op:   "insert",
			elem: skip_list.Element[string, int]{Key: "ab", Value: 10},
		},
		{
			op:            "search",
			elem:          skip_list.Element[string, int]{Key: "ab"},
			expected:      10,
			expectedFound: true,
		},
		{
			op:   "insert",
			elem: skip_list.Element[string, int]{Key: "abc", Value: 110},
		},
		{
			op:   "insert",
			elem: skip_list.Element[string, int]{Key: "ad", Value: 120},
		},
		{
			op:            "search",
			elem:          skip_list.Element[string, int]{Key: "ad"},
			expected:      120,
			expectedFound: true,
		},
		{
			op:   "insert",
			elem: skip_list.Element[string, int]{Key: "abc", Value: 19},
		},
		{
			op:            "search",
			elem:          skip_list.Element[string, int]{Key: "abc"},
			expected:      19,
			expectedFound: true,
		},
		{
			op:   "delete",
			elem: skip_list.Element[string, int]{Key: "abc"},
		},
		{
			op:            "search",
			elem:          skip_list.Element[string, int]{Key: "abc"},
			expectedFound: false,
		},
	}

	sl := skip_list.NewSkipList[string, int](level.NewRandomLevelGenerator())
	for _, c := range cases {
		switch c.op {
		case "insert":
			sl.Insert(c.elem.Key, c.elem.Value)
		case "search":
			res, found := sl.Search(c.elem.Key)
			assert.Equal(t, c.expectedFound, found)
			if c.expectedFound {
				assert.Equal(t, c.expected, res)
			}
		case "delete":
			sl.Delete(c.elem.Key)
		}
	}

}
