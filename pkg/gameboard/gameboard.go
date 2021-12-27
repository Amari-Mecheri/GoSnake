package gameboard

import (
	"crypto/rand"
	"errors"
	"gosnake/pkg/candy"
	"gosnake/pkg/common"
	"gosnake/pkg/snake"
	"math/big"
)

// Objects' body representation
const (
	FreeSpace rune = ' '
	SnakePart rune = 'S'
	CandyBody rune = '*'
)

// Defines custom errors
var (
	ErrInvalidSnakeReference = errors.New("The snake object is nil")
	ErrInvalidCandyReference = errors.New("The candy object is nil")
	ErrInvalidSize           = errors.New("Invalid board size")
	ErrInvalidPosition       = errors.New("Invalid position")
)

// GameBoarder is the interface defining gameBoard exported methods
type GameBoarder interface {
	InitGameBoard(size common.Size) (err error)
	BoardSize() common.Size
	IsSnakePart(ch rune) bool
	SetSnakeDirection(direction common.Direction)
	SnakeSize() (size int, err error)
	MoveSnake() (oldValue rune, listSprite []common.Sprite, err error)
	CreateSnake(position common.Position,
		direction common.Direction) (sprite common.Sprite, err error)
	SnakePosition() (position common.Position, err error)
	IsCandy(ch rune) bool
	CandyPosition() common.Position
	CandyAlive() bool
	RemoveCandy()
	CreateCandy() (sprite common.Sprite, err error)
	RandomFreePosition() (position common.Position, err error)
}

// gameBoard defines the properties of a game board
type gameBoard struct {
	size        common.Size
	board       [][]rune
	movingSnake snake.Snaker
	candy       candy.Candyer
}

// New returns an instance of gameBoard
func New() GameBoarder {
	var aGameBoard gameBoard
	aGameBoard.movingSnake = snake.New()
	aGameBoard.candy = candy.New()
	return &aGameBoard
}

func (aGameBoard *gameBoard) InitGameBoard(size common.Size) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	if err := aGameBoard.createBoard(size); err != nil {
		return err
	}
	if err := aGameBoard.clearBoard(); err != nil {
		return err // Shouldn't happen
	}

	return nil
}

func (aGameBoard *gameBoard) createBoard(size common.Size) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	// creates a slice of slices (2d slice)
	if size.Width < 0 || size.Height < 0 {
		return ErrInvalidSize
	}
	aGameBoard.board = make([][]rune, size.Width)
	for i := range aGameBoard.board {
		aGameBoard.board[i] = make([]rune, size.Height)
	}
	aGameBoard.size = size
	return nil
}

func (aGameBoard *gameBoard) clearBoard() (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	// fills the gameBoard with FreeSpaces
	for i := range aGameBoard.board {
		for j := range aGameBoard.board[i] {
			aGameBoard.board[i][j] = FreeSpace
		}
	}

	return nil
}

func (aGameBoard *gameBoard) BoardSize() common.Size {
	return aGameBoard.size
}

func (aGameBoard *gameBoard) IsSnakePart(ch rune) bool {
	return ch == SnakePart
}

func (aGameBoard *gameBoard) IsCandy(ch rune) bool {
	return ch == CandyBody
}

func (aGameBoard *gameBoard) CreateSnake(position common.Position,
	direction common.Direction) (sprite common.Sprite, err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	// Creates and position the snake
	aGameBoard.movingSnake = snake.New()
	aGameBoard.movingSnake.SetDirection(direction)
	if err = aGameBoard.movingSnake.GrowTo(position); err != nil {
		return sprite, err // Shouldn't happen
	}

	// Writes the snake to the board
	if err = aGameBoard.setCell(position, SnakePart); err != nil {
		return sprite, err // We actually want to return a default sprite
	}

	return common.Sprite{
		Value:    SnakePart,
		Position: position,
	}, nil
}

func (aGameBoard *gameBoard) SnakePosition() (position common.Position, err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	return aGameBoard.movingSnake.Position()
}

func (aGameBoard *gameBoard) CandyPosition() common.Position {
	return aGameBoard.candy.Position()
}

func (aGameBoard *gameBoard) CandyAlive() bool {
	return aGameBoard.candy.Alive()
}

func (aGameBoard *gameBoard) RemoveCandy() {
	aGameBoard.candy.Remove()
}

func (aGameBoard *gameBoard) CreateCandy() (sprite common.Sprite, err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	// Gets a free position
	position, err := aGameBoard.RandomFreePosition()
	if err != nil {
		return sprite, err
	}

	// Creates a candy
	aGameBoard.candy = candy.New()
	aGameBoard.candy.Init(position)

	// Sets the candy on the board
	if err = aGameBoard.setCell(position, CandyBody); err != nil {
		return sprite, err
	}

	return common.Sprite{
		Value:    CandyBody,
		Position: position,
	}, nil
}

func (aGameBoard *gameBoard) RandomFreePosition() (position common.Position, err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	if aGameBoard.size.Width <= 0 || aGameBoard.size.Height <= 0 {
		return position, ErrInvalidSize
	}
	// Randomly defines the candy position
	maxW := aGameBoard.size.Width
	maxH := aGameBoard.size.Height

	rndX, err := random(maxW)
	if err != nil {
		return position, err
	}

	rndY, err := random(maxH)
	if err != nil {
		return position, err
	}

	// loops for a free spot
	val, err := aGameBoard.cell(common.Position{
		X: rndX,
		Y: rndY,
	})
	if err != nil {
		return position, err
	}

	for val != FreeSpace {
		rndX, err = random(maxW)
		if err != nil {
			break
		}

		rndY, err = random(maxH)
		if err != nil {
			break
		}

		val, err = aGameBoard.cell(common.Position{
			X: rndX,
			Y: rndY,
		})
		if err != nil {
			break
		}
	}

	return common.Position{
		X: rndX,
		Y: rndY,
	}, err
}

func random(max int) (rnd int, err error) {
	aBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err == nil {
		rnd = int(aBig.Int64())
	}

	return rnd, err
}

func (aGameBoard *gameBoard) SetSnakeDirection(direction common.Direction) {
	aGameBoard.movingSnake.SetDirection(direction)
}

func (aGameBoard *gameBoard) SnakeSize() (size int, err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	if aGameBoard.movingSnake == nil {
		return 0, ErrInvalidSnakeReference
	}

	return aGameBoard.movingSnake.Size()
}

func (aGameBoard *gameBoard) MoveSnake() (oldValue rune, listSprite []common.Sprite, err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	// Asks the snake where it wants to move
	requestedPosition, err := aGameBoard.movingSnake.NextMove()
	if err != nil {
		return oldValue, listSprite, err
	}

	// Translates the requested position inside the board
	actualPosition := aGameBoard.translatePosition(requestedPosition)

	// Gets the content at the actual position
	oldValue, err = aGameBoard.cell(actualPosition)
	if err != nil {
		return oldValue, listSprite, err
	}

	// Call the actual move (or growth)
	listSprite, err = aGameBoard.actualMove(actualPosition, oldValue)
	// returns the old content and the list of sprites
	return oldValue, listSprite, err
}

func (aGameBoard *gameBoard) actualMove(position common.Position, oldValue rune) (
	listSprite []common.Sprite, err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	// has the snake eaten a candy?
	if oldValue == CandyBody {
		// Grow the snake
		err = aGameBoard.movingSnake.GrowTo(position)
		if err != nil {
			return nil, err
		}
		// update the board
		err = aGameBoard.setCell(position, SnakePart)
		// Only the head of the snake is updated
		return []common.Sprite{
			{
				Value:    SnakePart,
				Position: position,
			},
		}, err
	}

	// Move the snake
	oldTail, err := aGameBoard.movingSnake.MoveTo(position)
	if err != nil {
		return nil, err
	}

	// update the board with the new head
	err = aGameBoard.setCell(position, SnakePart)
	if err != nil {
		return nil, err
	}

	// Remove the tail
	err = aGameBoard.setCell(oldTail, FreeSpace)
	// The head and the tail have been updated
	return []common.Sprite{
		{
			Value:    SnakePart,
			Position: position,
		},
		{
			Value:    FreeSpace,
			Position: oldTail,
		},
	}, err
}

func (aGameBoard *gameBoard) cell(position common.Position) (value rune, err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	if !aGameBoard.checkPosition(position) {
		return value, ErrInvalidPosition
	}

	return aGameBoard.board[position.X][position.Y], nil
}

func (aGameBoard *gameBoard) setCell(position common.Position, value rune) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	if !aGameBoard.checkPosition(position) {
		return ErrInvalidPosition
	}

	aGameBoard.board[position.X][position.Y] = value
	return nil
}

func (aGameBoard *gameBoard) checkPosition(position common.Position) bool {
	if position.X < 0 || position.X >= aGameBoard.size.Width {
		return false
	}
	if position.Y < 0 || position.Y >= aGameBoard.size.Height {
		return false
	}

	return true
}

func (aGameBoard *gameBoard) translatePosition(requestedPosition common.Position) common.Position {
	// The position is kept inside the board
	// If it gets out one side, it enters the other side
	translatedPosition := requestedPosition
	maxW := aGameBoard.size.Width - 1
	maxH := aGameBoard.size.Height - 1

	if translatedPosition.X < 0 {
		translatedPosition.X = maxW
	}
	if translatedPosition.X > maxW {
		translatedPosition.X = 0
	}
	if translatedPosition.Y < 0 {
		translatedPosition.Y = maxH
	}
	if translatedPosition.Y > maxH {
		translatedPosition.Y = 0
	}

	return translatedPosition
}
