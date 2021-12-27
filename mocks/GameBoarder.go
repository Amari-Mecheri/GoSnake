// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import common "gosnake/pkg/common"

import mock "github.com/stretchr/testify/mock"

// GameBoarder is an autogenerated mock type for the GameBoarder type
type GameBoarder struct {
	mock.Mock
}

// BoardSize provides a mock function with given fields:
func (_m *GameBoarder) BoardSize() common.Size {
	ret := _m.Called()

	var r0 common.Size
	if rf, ok := ret.Get(0).(func() common.Size); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(common.Size)
	}

	return r0
}

// CandyAlive provides a mock function with given fields:
func (_m *GameBoarder) CandyAlive() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// CandyPosition provides a mock function with given fields:
func (_m *GameBoarder) CandyPosition() common.Position {
	ret := _m.Called()

	var r0 common.Position
	if rf, ok := ret.Get(0).(func() common.Position); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(common.Position)
	}

	return r0
}

// CreateCandy provides a mock function with given fields:
func (_m *GameBoarder) CreateCandy() (common.Sprite, error) {
	ret := _m.Called()

	var r0 common.Sprite
	if rf, ok := ret.Get(0).(func() common.Sprite); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(common.Sprite)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateSnake provides a mock function with given fields: position, direction
func (_m *GameBoarder) CreateSnake(position common.Position, direction common.Direction) (common.Sprite, error) {
	ret := _m.Called(position, direction)

	var r0 common.Sprite
	if rf, ok := ret.Get(0).(func(common.Position, common.Direction) common.Sprite); ok {
		r0 = rf(position, direction)
	} else {
		r0 = ret.Get(0).(common.Sprite)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Position, common.Direction) error); ok {
		r1 = rf(position, direction)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InitGameBoard provides a mock function with given fields: size
func (_m *GameBoarder) InitGameBoard(size common.Size) error {
	ret := _m.Called(size)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Size) error); ok {
		r0 = rf(size)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsCandy provides a mock function with given fields: ch
func (_m *GameBoarder) IsCandy(ch rune) bool {
	ret := _m.Called(ch)

	var r0 bool
	if rf, ok := ret.Get(0).(func(rune) bool); ok {
		r0 = rf(ch)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsSnakePart provides a mock function with given fields: ch
func (_m *GameBoarder) IsSnakePart(ch rune) bool {
	ret := _m.Called(ch)

	var r0 bool
	if rf, ok := ret.Get(0).(func(rune) bool); ok {
		r0 = rf(ch)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MoveSnake provides a mock function with given fields:
func (_m *GameBoarder) MoveSnake() (rune, []common.Sprite, error) {
	ret := _m.Called()

	var r0 rune
	if rf, ok := ret.Get(0).(func() rune); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(rune)
	}

	var r1 []common.Sprite
	if rf, ok := ret.Get(1).(func() []common.Sprite); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]common.Sprite)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// RandomFreePosition provides a mock function with given fields:
func (_m *GameBoarder) RandomFreePosition() (common.Position, error) {
	ret := _m.Called()

	var r0 common.Position
	if rf, ok := ret.Get(0).(func() common.Position); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(common.Position)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveCandy provides a mock function with given fields:
func (_m *GameBoarder) RemoveCandy() {
	_m.Called()
}

// SetSnakeDirection provides a mock function with given fields: direction
func (_m *GameBoarder) SetSnakeDirection(direction common.Direction) {
	_m.Called(direction)
}

// SnakePosition provides a mock function with given fields:
func (_m *GameBoarder) SnakePosition() (common.Position, error) {
	ret := _m.Called()

	var r0 common.Position
	if rf, ok := ret.Get(0).(func() common.Position); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(common.Position)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SnakeSize provides a mock function with given fields:
func (_m *GameBoarder) SnakeSize() (int, error) {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}