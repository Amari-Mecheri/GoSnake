package candy

import (
	"gosnake/pkg/common"
)

// Candyer is a candy interface
type Candyer interface {
	Remove()
	Init(newPosition common.Position)
	Position() common.Position
	Alive() bool
}

// candy has the properties of a candy
type candy struct {
	alive    bool //false by default
	position common.Position
}

// New returns an instance of candy
func New() Candyer {
	return new(candy)
}

// Remove sets alive to false
func (aCandy *candy) Remove() {
	aCandy.alive = false
}

// Init set the position of the candy and alive to true
func (aCandy *candy) Init(newPosition common.Position) {
	aCandy.position = newPosition
	aCandy.alive = true
}

// Position returns the candy position
func (aCandy *candy) Position() common.Position {
	return aCandy.position
}

// Alive returns the candy status: eaten or not
func (aCandy *candy) Alive() bool {
	return aCandy.alive
}
