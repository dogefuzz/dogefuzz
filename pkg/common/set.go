package common

type void = struct{}

type set[T comparable] struct {
	internal map[T]void
}

func NewSet[T comparable](args ...T) *set[T] {
	internal := make(map[T]void)
	for _, arg := range args {
		internal[arg] = void{}
	}

	return &set[T]{
		internal: internal,
	}
}

func (s *set[T]) Add(args ...T) {
	for _, arg := range args {
		s.internal[arg] = void{}
	}
}

func (s *set[T]) Has(value T) bool {
	if _, ok := s.internal[value]; ok {
		return true
	}
	return false
}
