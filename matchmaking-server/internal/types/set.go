package types

type Set[T comparable] map[T]struct{}

func (s Set[T]) Has(v T) bool {
	_, found := s[v]

	return found
}

func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}
