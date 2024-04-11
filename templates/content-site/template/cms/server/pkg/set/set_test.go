package set

import (
	"cms/pkg/asserts"
	"testing"
)

func TestNew(t *testing.T) {
	type myInt int64
	asserts.Equals(t, set[int]{1: {}, 2: {}, 3: {}}, New([]int{1, 2, 3}))
	asserts.Equals(t, set[string]{"a": {}, "b": {}, "c": {}}, New([]string{"a", "b", "c"}))
	asserts.Equals(t, set[myInt]{1: {}, 2: {}, 3: {}}, New([]myInt{1, 2, 3}))
}

func TestSet_Has(t *testing.T) {
	s := New([]int{1, 2, 3})
	asserts.Equals(t, true, s.Has(1))
	asserts.Equals(t, false, s.Has(4))
}

func TestSet_SymmetricDifference(t *testing.T) {
	s1 := New([]int{1, 2, 3})
	s2 := New([]int{3, 4, 5})

	toAdd, toRemove := s1.Difference(s2)
	asserts.EqualsIgnoreOrder(t, []int{4, 5}, toAdd)
	asserts.EqualsIgnoreOrder(t, []int{1, 2}, toRemove)
}
