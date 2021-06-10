package synx

// Bool represents a boolean value
type Bool struct {
	locker
	value bool
}

// NewBool creates a new wrapper type for a Bool value
func NewBool(value bool) *Bool {
	return &Bool{
		locker: newLock(),
		value:  value,
	}
}

// SetValue sets the value to the given value
func (s *Bool) SetValue(value bool) {
	s.Lock()
	s.value = value
	s.Unlock()
}

// GetValue returns the current value
func (s *Bool) GetValue() (value bool) {
	s.Lock()
	value = s.value
	s.Unlock()
	return
}
