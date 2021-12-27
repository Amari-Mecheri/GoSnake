package gamestate

import (
	"errors"
	"gosnake/pkg/common"
	"gosnake/pkg/gameboard"
)

// GameStater is the gameState interface
type GameStater interface {
	InitBoard(size common.Size) (err error)
	CreateObjects() (listSprite []common.Sprite, err error)
	Start()
	Play() (listSprite []common.Sprite, err error)
	GameInProgress() bool
	SetGameInProgress(bool)
	Dirty() bool
	HighScore() int
	Score() int
	Round() int
	MoveLeft()
	MoveRight()
	MoveDown()
	MoveUp()
	BoardSize() common.Size
	SnakePosition() (position common.Position, err error)
	SnakeSize() (size int, err error)
}

type gameState struct {
	gameInProgress bool
	round          int
	score          int
	highScore      int
	dirty          bool
	gameboard.GameBoarder
}

// ErrInvalidBoardReference is a custom error thrown when the board object is nil
var ErrInvalidBoardReference = errors.New("The board object is nil")

var (
	goLeft common.Direction = common.Direction{
		DX: -1,
		DY: 0,
	}
	goRight common.Direction = common.Direction{
		DX: 1,
		DY: 0,
	}
	goUp common.Direction = common.Direction{
		DX: 0,
		DY: -1,
	}
	goDown common.Direction = common.Direction{
		DX: 0,
		DY: 1,
	}
)

// New returns an instance of gameState
func New() GameStater {
	var aGameState gameState
	aGameState.GameBoarder = gameboard.New()
	return &aGameState
}

func (aGameState *gameState) InitBoard(size common.Size) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	aGameState.GameBoarder = gameboard.New()

	if err = aGameState.InitGameBoard(size); err != nil {
		return err
	}
	aGameState.dirty = false
	return nil
}

func (aGameState *gameState) CreateObjects() (listSprite []common.Sprite, err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	if aGameState.GameBoarder == nil {
		return listSprite, ErrInvalidBoardReference
	}
	position := common.Position{
		X: aGameState.BoardSize().Width / 2,
		Y: aGameState.BoardSize().Height / 2,
	}
	snake, err := aGameState.CreateSnake(position, goRight)
	if err != nil {
		return nil, err
	}
	candy, err := aGameState.CreateCandy()
	return []common.Sprite{snake, candy}, err
}

func (aGameState *gameState) Start() {
	aGameState.gameInProgress = true
	aGameState.score = 0
	aGameState.round = 0
	aGameState.dirty = true
}

func (aGameState *gameState) Play() (listSprite []common.Sprite, err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	if aGameState.GameBoarder == nil {
		return listSprite, ErrInvalidBoardReference
	}

	//Plays a round
	aGameState.round++

	//Move the snake
	oldValue, spriteList, err := aGameState.MoveSnake()
	if err != nil {
		aGameState.gameInProgress = false
		return spriteList, err
	}
	//Game over?
	if aGameState.IsSnakePart(oldValue) {
		aGameState.gameInProgress = false
		return spriteList, nil
	}

	//Ate a candy?
	if aGameState.IsCandy(oldValue) {
		//Remove the candy since it's been eaten
		aGameState.RemoveCandy()
		//updates the score
		aGameState.score++
		//updates the highscore
		if aGameState.score > aGameState.highScore {
			aGameState.highScore = aGameState.score
		}
	}

	//No more candies?
	if !aGameState.CandyAlive() {
		sprite, err := aGameState.CreateCandy()
		if err != nil {
			return nil, err
		}
		spriteList = append(spriteList, sprite)
	}

	return spriteList, nil
}

func (aGameState *gameState) GameInProgress() bool {
	return aGameState.gameInProgress
}

func (aGameState *gameState) SetGameInProgress(val bool) {
	aGameState.gameInProgress = val
}

func (aGameState *gameState) Dirty() bool {
	return aGameState.dirty
}

func (aGameState *gameState) HighScore() int {
	return aGameState.highScore
}

func (aGameState *gameState) Score() int {
	return aGameState.score
}

func (aGameState *gameState) Round() int {
	return aGameState.round
}

func (aGameState *gameState) MoveLeft() {
	aGameState.SetSnakeDirection(goLeft)
}
func (aGameState *gameState) MoveRight() {
	aGameState.SetSnakeDirection(goRight)
}
func (aGameState *gameState) MoveDown() {
	aGameState.SetSnakeDirection(goDown)
}
func (aGameState *gameState) MoveUp() {
	aGameState.SetSnakeDirection(goUp)
}
