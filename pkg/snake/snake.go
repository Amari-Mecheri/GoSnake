package snake

import (
	"errors"
	"gosnake/pkg/common"
)

// Snaker is the snake interface
type Snaker interface {
	Size() (size int, err error)
	SetDirection(direction common.Direction)
	Position() (position common.Position, err error)
	NextMove() (nextPosition common.Position, err error)
	MoveTo(newPosition common.Position) (theTail common.Position, err error)
	GrowTo(newPosition common.Position) (err error)
}

type snake struct {
	body      []common.Position
	direction common.Direction
}

// ErrNoSnakeBody is a custom error type
var ErrNoSnakeBody = errors.New("the snake has no body")

// New returns an instance of snake
func New() Snaker {
	return new(snake)
}

func (aSnake *snake) Size() (size int, err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	return len(aSnake.body), nil
}

func (aSnake *snake) SetDirection(direction common.Direction) {
	aSnake.direction = direction
}

func (aSnake *snake) Position() (position common.Position, err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	var size int
	size, err = aSnake.Size()
	if err != nil { // should not occur
		return position, err
	}
	if size > 0 {
		return aSnake.body[size-1], nil
	}

	return position, ErrNoSnakeBody
}

func (aSnake *snake) NextMove() (nextPosition common.Position, err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	nextPosition, err = aSnake.Position()
	if err != nil {
		return nextPosition, err
	}
	nextPosition.X += aSnake.direction.DX
	nextPosition.Y += aSnake.direction.DY
	return nextPosition, nil
}

func (aSnake *snake) MoveTo(newPosition common.Position) (theTail common.Position, err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	var size int
	size, err = aSnake.Size()
	if err != nil { // should not occur
		return theTail, err
	}
	if size > 0 {
		// Removes the tail
		theTail = aSnake.body[0]
		aSnake.body = aSnake.body[1:]
		// Inserts the head at new position
		aSnake.body = append(aSnake.body, newPosition)
		// return the tail
		return theTail, nil
	}

	return theTail, ErrNoSnakeBody
}

func (aSnake *snake) GrowTo(newPosition common.Position) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	// Doesn't remove the tail
	// Inserts the head at new position
	aSnake.body = append(aSnake.body, newPosition)
	return nil
}
