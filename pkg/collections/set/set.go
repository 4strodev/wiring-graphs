package set

type Set[T comparable] map[T]struct{}

func New[T comparable]() Set[T] {
	return make(Set[T], 0)
}

func (s Set[T]) Add(val T) {
	s[val] = struct{}{}
}

func (s Set[T]) Has(val T) bool {
	_, ok := s[val]
	return ok
}

func (s Set[T]) Remove(val T) {
	delete(s, val)
}
