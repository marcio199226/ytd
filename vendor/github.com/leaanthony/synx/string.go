package synx

// String represents a string type
type String struct {
	locker
	value string
}

// NewString creates a new wrapper for a string value
func NewString(value string) *String {
	return &String{
		locker: newLock(),
		value:  value,
	}
}

// SetValue sets the value
func (s *String) SetValue(value string) {
	s.Lock()
	s.value = value
	s.Unlock()
}

// GetValue returns the value originally given
func (s *String) GetValue() (value string) {
	s.Lock()
	value = s.value
	s.Unlock()
	return
}
