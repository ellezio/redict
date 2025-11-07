package redict

type string_ struct {
	inner []byte
}

func newStrings() *string_ {
	return &string_{}
}

func (s *string_) set(value []byte) {
	s.inner = make([]byte, len(value))
	copy(s.inner, value)
}

func (s *string_) get() []byte {
	return s.inner
}
