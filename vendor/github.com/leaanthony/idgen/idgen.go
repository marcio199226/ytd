package idgen

import (
	"errors"
)

//Generator generates unique IDs (uint)
type Generator struct {
	counter uint
	usedIDs map[uint]struct{}
	maximum uint
}

//New creates a new ID (uint) Generator
func New() *Generator {
	return &Generator{
		usedIDs: make(map[uint]struct{}),
		maximum: ^uint(0), // MaxUInt
	}
}

// NewWithMaximum creates a new ID (uint) Generator with
// a maximum number of unique IDs that may be generated
func NewWithMaximum(maximum uint) *Generator {
	return &Generator{
		usedIDs: make(map[uint]struct{}),
		maximum: maximum,
	}
}

// NewID returns an ID that has not been previously issued.
// Returns an error if all available IDs have been issued.
func (g *Generator) NewID() (uint, error) {

	startID := g.counter
	for {
		_, exists := g.usedIDs[g.counter]
		if !exists {
			g.usedIDs[g.counter] = struct{}{}
			return g.counter, nil
		}
		g.counter++
		// Wrap at max size
		if g.counter >= g.maximum {
			g.counter = 0
		}
		// Check if we have looped
		if g.counter == startID {
			return 0, errors.New("maximum number of unique IDs generated")
		}
	}
}

// ReleaseID frees up the given ID so that it may be reused
func (g *Generator) ReleaseID(ID uint) {
	delete(g.usedIDs, ID)
}
