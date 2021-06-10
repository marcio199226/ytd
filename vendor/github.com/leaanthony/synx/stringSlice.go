package synx

// StringSlice represents a string slice type
type StringSlice struct {
	locker
	value []string
}

// NewStringSlice creates a new wrapper for a string slice
func NewStringSlice(value []string) *StringSlice {
	return &StringSlice{
		locker: newLock(),
		value:  value,
	}
}

// SetValue sets the value
func (s *StringSlice) SetValue(value []string) {
	s.Lock()
	s.value = value
	s.Unlock()
}

// GetValue returns the value originally given
func (s *StringSlice) GetValue() (value []string) {
	s.Lock()
	value = s.value
	s.Unlock()
	return
}

// GetElement returns the value at the given element index
func (s *StringSlice) GetElement(index int) (value string) {
	s.Lock()
	value = s.value[index]
	s.Unlock()
	return
}

// Length returns the number of elements in the slice
func (s *StringSlice) Length() (length int) {
	return len(s.GetValue())
}
