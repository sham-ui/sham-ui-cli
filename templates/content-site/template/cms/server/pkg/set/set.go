package set

type set[T comparable] map[T]struct{}

func (s set[T]) Has(v T) bool {
	_, ok := s[v]
	return ok
}

// Difference returns the elements that are in the set but not in the other set.
//
// Parameters:
// - other: the other set to compare against.
//
// Returns:
// - toAdd: a slice of elements that are in the other set but not in the set.
// - toRemove: a slice of elements that are in the set but not in the other set.
//
//nolint:nonamedreturns
func (s set[T]) Difference(other set[T]) (toAdd []T, toRemove []T) {
	for v := range other {
		if !s.Has(v) {
			toAdd = append(toAdd, v)
		}
	}
	for v := range s {
		if !other.Has(v) {
			toRemove = append(toRemove, v)
		}
	}
	return
}

func New[T comparable](s []T) set[T] {
	m := make(set[T], len(s))
	for _, v := range s {
		m[v] = struct{}{}
	}
	return m
}
