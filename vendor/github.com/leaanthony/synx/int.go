package synx

// Int represents an int type
type Int struct {
	locker
	value int
}

// NewInt creates a new wrapper type for an int value
func NewInt(value int) *Int {
	return &Int{
		locker: newLock(),
		value:  value,
	}
}

// SetValue sets the value
func (s *Int) SetValue(value int) {
	s.Lock()
	s.value = value
	s.Unlock()
}

// GetValue returns the value originally given
func (s *Int) GetValue() (value int) {
	s.Lock()
	value = s.value
	s.Unlock()
	return
}
