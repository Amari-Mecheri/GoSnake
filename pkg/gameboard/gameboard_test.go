package gameboard

import (
	"gosnake/mocks"
	"gosnake/pkg/candy"
	"gosnake/pkg/common"
	"gosnake/pkg/snake"
	"gosnake/testdata"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGameBoard_createBoard(t *testing.T) {
	type fields struct {
		size        common.Size
		board       [][]rune
		movingSnake snake.Snaker
		candy       candy.Candyer
	}
	type args struct {
		size common.Size
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantSize common.Size
		wantErr  bool
	}{
		{
			name: "TestEmptyBoardTo0,0",
			fields: fields{
				size:  testdata.Size0_0,
				board: nil,
			},
			args: args{
				size: testdata.Size0_0,
			},
			wantSize: testdata.Size0_0,
			wantErr:  false,
		},
		{
			name: "TestEmptyBoardTo18,19",
			fields: fields{
				size:  testdata.Size0_0,
				board: nil,
			},
			args: args{
				size: common.Size{
					Width:  18,
					Height: 19,
				},
			},
			wantSize: common.Size{
				Width:  18,
				Height: 19,
			},
			wantErr: false,
		},
		{
			name: "TestEmptyBoardTo18,-1",
			fields: fields{
				size:  testdata.Size0_0,
				board: nil,
			},
			args: args{
				size: common.Size{
					Width:  18,
					Height: -1,
				},
			},
			wantSize: testdata.Size0_0,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameBoard := &gameBoard{
				size:        tt.fields.size,
				board:       tt.fields.board,
				movingSnake: tt.fields.movingSnake,
				candy:       tt.fields.candy,
			}
			err := aGameBoard.createBoard(tt.args.size)
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr)
		})
	}
}

func TestGameBoard_clearBoard(t *testing.T) {
	type fields struct {
		size        common.Size
		board       [][]rune
		movingSnake snake.Snaker
		candy       candy.Candyer
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestEmptyBoard",
			fields: fields{
				size:  testdata.Size0_0,
				board: nil,
			},
			wantErr: false,
		},
		{
			name: "TestBoard3,3",
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3),
			},
			wantErr: false,
		},
		{
			name: "TestBoard3,3WrongSize",
			fields: fields{
				size:  testdata.Size4_4,
				board: testdata.Duplicate(testdata.Board3_3),
			},
			wantErr: false,
		},
		{
			name: "TestBoard3,3WrongSizeWrongAllocation",
			fields: fields{
				size:  testdata.Size4_4,
				board: testdata.Duplicate(testdata.Board3_3WrongAllocation),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameBoard := &gameBoard{
				size:        tt.fields.size,
				board:       tt.fields.board,
				movingSnake: tt.fields.movingSnake,
				candy:       tt.fields.candy,
			}
			err := aGameBoard.clearBoard()
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr)
		})
	}
}

func TestGameBoard_InitGameBoard(t *testing.T) {
	type fields struct {
		size        common.Size
		board       [][]rune
		movingSnake snake.Snaker
		candy       candy.Candyer
	}
	type args struct {
		size common.Size
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantErrType error
	}{
		{
			name: "TestNegativeSize",
			args: args{
				size: testdata.SizeMinus1_Minus1,
			},
			wantErr:     true,
			wantErrType: ErrInvalidSize,
		},
		{
			name:    "TestEmptyBoard",
			wantErr: false,
		},
		{
			name: "TestBoard3,3",
			args: args{
				size: testdata.Size3_3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameBoard := &gameBoard{
				size:        tt.fields.size,
				board:       tt.fields.board,
				movingSnake: tt.fields.movingSnake,
				candy:       tt.fields.candy,
			}
			err := aGameBoard.InitGameBoard(tt.args.size)
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr)
			if gotErr {
				require.ErrorIs(t, err, tt.wantErrType)
			}
		})
	}
}

func TestGameBoard_BoardSize(t *testing.T) {
	type fields struct {
		size        common.Size
		board       [][]rune
		movingSnake snake.Snaker
		candy       candy.Candyer
	}
	tests := []struct {
		name     string
		fields   fields
		wantSize common.Size
	}{
		{
			name: "TestEmptyBoard",
			fields: fields{
				size:  testdata.Size0_0,
				board: nil,
			},
			wantSize: testdata.Size0_0,
		},
		{
			name: "TestBoard3,3",
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3),
			},
			wantSize: testdata.Size3_3,
		},
		{
			name: "TestBoard3,3WrongSizeWrongAllocation",
			fields: fields{
				size:  testdata.Size4_4,
				board: testdata.Duplicate(testdata.Board3_3WrongAllocation),
			},
			wantSize: testdata.Size4_4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameBoard := &gameBoard{
				size:        tt.fields.size,
				board:       tt.fields.board,
				movingSnake: tt.fields.movingSnake,
				candy:       tt.fields.candy,
			}
			gotSize := aGameBoard.BoardSize()
			require.Equal(t, tt.wantSize, gotSize)
		})
	}
}

func TestGameBoard_IsSnakePart(t *testing.T) {
	type fields struct {
		size        common.Size
		board       [][]rune
		movingSnake snake.Snaker
		candy       candy.Candyer
	}
	type args struct {
		ch rune
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "CandyBody",
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3),
			},
			args: args{
				ch: CandyBody,
			},
			want: false,
		},
		{
			name: "SnakePart",
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3),
			},
			args: args{
				ch: SnakePart,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameBoard := &gameBoard{
				size:        tt.fields.size,
				board:       tt.fields.board,
				movingSnake: tt.fields.movingSnake,
				candy:       tt.fields.candy,
			}
			got := aGameBoard.IsSnakePart(tt.args.ch)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGameBoard_IsCandy(t *testing.T) {
	type fields struct {
		size        common.Size
		board       [][]rune
		movingSnake snake.Snaker
		candy       candy.Candyer
	}
	type args struct {
		ch rune
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "CandyBody",
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3),
			},
			args: args{
				ch: CandyBody,
			},
			want: true,
		},
		{
			name: "SnakePart",
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3),
			},
			args: args{
				ch: SnakePart,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameBoard := &gameBoard{
				size:        tt.fields.size,
				board:       tt.fields.board,
				movingSnake: tt.fields.movingSnake,
				candy:       tt.fields.candy,
			}
			got := aGameBoard.IsCandy(tt.args.ch)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGameBoard_CreateSnake(t *testing.T) {
	type fields struct {
		size        common.Size
		board       [][]rune
		movingSnake snake.Snaker
		candy       candy.Candyer
	}
	type args struct {
		position  common.Position
		direction common.Direction
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantSprite  common.Sprite
		wantErrType error
		wantErr     bool
	}{
		{
			name: "TestNilSnake",
			fields: fields{
				size:  testdata.Size0_0,
				board: nil,
			},
			wantErrType: ErrInvalidPosition,
			wantErr:     true, // The snake is created but the board is empty
		},
		{
			name: "TestEmptyBoardPos0,0",
			fields: fields{
				size:        testdata.Size0_0,
				board:       nil,
				movingSnake: snake.New(),
			},
			args: args{
				position:  testdata.Position0_0,
				direction: testdata.Direction0_0,
			},
			wantErrType: ErrInvalidPosition,
			wantErr:     true, // There is a snake but the board is empty
		},
		{
			name: "TestBoard3_3",
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3),
			},
			args: args{
				position:  testdata.Position0_0,
				direction: testdata.Direction0_0,
			},
			wantSprite: common.Sprite{
				Value:    SnakePart,
				Position: testdata.Position0_0,
			},
			wantErr: false,
		},
		{
			name: "TestBoard3_3Pos4_3",
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3),
			},
			args: args{
				position:  testdata.Position4_3,
				direction: testdata.Direction0_0,
			},
			wantErrType: ErrInvalidPosition,
			wantErr:     true, // the position is out of board's ranges
		},
		{
			name: "TestBoard3_3Pos3_4",
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3),
			},
			args: args{
				position:  testdata.Position3_4,
				direction: testdata.Direction0_0,
			},
			wantErrType: ErrInvalidPosition,
			wantErr:     true, // the position is out of board's ranges
		},
		{
			name: "TestBoard3_3Pos-1,-1",
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3),
			},
			args: args{
				position: testdata.PositionMinus1_Minus1,
			},
			wantErrType: ErrInvalidPosition,
			wantErr:     true, // the position is out of board's ranges
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameBoard := &gameBoard{
				size:        tt.fields.size,
				board:       tt.fields.board,
				movingSnake: tt.fields.movingSnake,
				candy:       tt.fields.candy,
			}
			gotSprite, err := aGameBoard.CreateSnake(tt.args.position, tt.args.direction)
			gotErr := (err != nil)
			if gotErr && tt.wantErrType != nil {
				require.ErrorIs(t, err, tt.wantErrType)
			}
			require.Equal(t, tt.wantErr, gotErr, "%w", err)
			require.Equal(t, tt.wantSprite, gotSprite)
		})
	}
}

func TestGameBoard_SnakePosition(t *testing.T) {
	// At this point this method just returns snake.Position
	// In this test the returned values should be the mocked values
	type fields struct {
		size        common.Size
		board       [][]rune
		movingSnake snake.Snaker
		candy       candy.Candyer
	}
	tests := []struct {
		name         string
		fields       fields
		mockPosition common.Position
		wantPosition common.Position
		wantErrType  error
		wantErr      bool
	}{
		{
			name:         "TestEmptyBoardPos4,3",
			mockPosition: testdata.Position4_3,
			wantPosition: testdata.Position4_3,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameBoard := &gameBoard{
				size:        tt.fields.size,
				board:       tt.fields.board,
				movingSnake: tt.fields.movingSnake,
				candy:       tt.fields.candy,
			}
			aSnake := &mocks.Snaker{}
			aSnake.On("Position").Return(tt.mockPosition, nil)
			aGameBoard.movingSnake = aSnake
			gotPosition, err := aGameBoard.SnakePosition()
			gotErr := (err != nil)
			if gotErr && tt.wantErrType != nil {
				require.ErrorIs(t, err, tt.wantErrType)
			}
			require.Equal(t, tt.wantErr, gotErr, "%w", err)
			require.Equal(t, tt.wantPosition, gotPosition)
		})
	}
}

func TestGameBoard_CandyPosition(t *testing.T) {
	// At this point this method just returns candy.Position
	// In this test the returned value should be the mocked value
	type fields struct {
		size        common.Size
		board       [][]rune
		movingSnake snake.Snaker
		candy       candy.Candyer
	}
	tests := []struct {
		name         string
		fields       fields
		mockPosition common.Position
		wantPosition common.Position
	}{
		{
			name:         "TestEmptyBoardPos4,3",
			mockPosition: testdata.Position4_3,
			wantPosition: testdata.Position4_3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameBoard := &gameBoard{
				size:        tt.fields.size,
				board:       tt.fields.board,
				movingSnake: tt.fields.movingSnake,
				candy:       tt.fields.candy,
			}
			aCandy := &mocks.Candyer{}
			aCandy.On("Position").Return(tt.mockPosition)
			aGameBoard.candy = aCandy
			gotPosition := aGameBoard.CandyPosition()
			require.Equal(t, tt.wantPosition, gotPosition)

		})
	}
}

func TestGameBoard_CandyAlive(t *testing.T) {
	// At this point this method just returns candy.CandyAlive
	// In this test the returned value should be the mocked value
	type fields struct {
		size        common.Size
		board       [][]rune
		movingSnake snake.Snaker
		candy       candy.Candyer
	}
	tests := []struct {
		name       string
		fields     fields
		mockStatus bool
		wantStatus bool
	}{
		{
			name:       "TestEmptyBoardAlive",
			mockStatus: true,
			wantStatus: true,
		},
		{
			name:       "TestEmptyBoardAlive",
			mockStatus: false,
			wantStatus: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameBoard := &gameBoard{
				size:        tt.fields.size,
				board:       tt.fields.board,
				movingSnake: tt.fields.movingSnake,
				candy:       tt.fields.candy,
			}
			aCandy := &mocks.Candyer{}
			aCandy.On("Alive").Return(tt.mockStatus)
			aGameBoard.candy = aCandy
			gotStatus := aGameBoard.CandyAlive()
			require.Equal(t, tt.wantStatus, gotStatus)
		})
	}
}

func TestGameBoard_CreateCandy(t *testing.T) {
	type fields struct {
		size        common.Size
		board       [][]rune
		movingSnake snake.Snaker
		candy       candy.Candyer
	}
	tests := []struct {
		name        string
		fields      fields
		wantSprite  common.Sprite
		wantErrType error
		wantErr     bool
	}{
		{
			name: "TestNilCandy",
			fields: fields{
				size:  testdata.Size0_0,
				board: nil,
			},
			wantErrType: ErrInvalidSize,
			wantErr:     true, // The candy is created but the board is empty
		},
		{
			// Randomly positions a candy
			// We provide a board with only 1 free spot
			name: "TestBoard3_3OneFreeSpot",
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3_OneFreeSpotPos1_1),
				candy: candy.New(),
			},
			wantSprite: common.Sprite{
				Value:    CandyBody,
				Position: testdata.Position1_1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameBoard := &gameBoard{
				size:        tt.fields.size,
				board:       tt.fields.board,
				movingSnake: tt.fields.movingSnake,
				candy:       tt.fields.candy,
			}
			gotSprite, err := aGameBoard.CreateCandy()
			gotErr := (err != nil)
			if gotErr && tt.wantErrType != nil {
				require.ErrorIs(t, err, tt.wantErrType)
			}
			require.Equal(t, tt.wantErr, gotErr, "%w", err)
			require.Equal(t, tt.wantSprite, gotSprite)
		})
	}
}

func TestGameBoard_RandomFreePosition(t *testing.T) {
	// To test we just a provide a board with 1 free spot
	type fields struct {
		size        common.Size
		board       [][]rune
		movingSnake snake.Snaker
		candy       candy.Candyer
	}
	tests := []struct {
		name         string
		fields       fields
		wantPosition common.Position
		wantErr      bool
	}{
		{
			name: "TestBoard3_3_OneFreeSpot",
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3_OneFreeSpotPos1_1),
			},
			wantPosition: testdata.Position1_1,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameBoard := &gameBoard{
				size:        tt.fields.size,
				board:       tt.fields.board,
				movingSnake: tt.fields.movingSnake,
				candy:       tt.fields.candy,
			}
			gotPosition, err := aGameBoard.RandomFreePosition()
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr, "%w", err)
			require.Equal(t, tt.wantPosition, gotPosition)
		})
	}
}

func TestGameBoard_SnakeSize(t *testing.T) {
	type fields struct {
		size        common.Size
		board       [][]rune
		movingSnake snake.Snaker
		candy       candy.Candyer
	}
	tests := []struct {
		name     string
		fields   fields
		mockSize int
		wantSize int
		wantErr  bool
	}{
		{
			name: "TestEmptySnake",
			fields: fields{
				movingSnake: snake.New(),
			},
			wantSize: 0,
			wantErr:  false,
		},
		{
			name:     "TestSnakeSize5",
			mockSize: 5,
			wantSize: 5,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameBoard := &gameBoard{
				size:        tt.fields.size,
				board:       tt.fields.board,
				movingSnake: tt.fields.movingSnake,
				candy:       tt.fields.candy,
			}
			if aGameBoard.movingSnake == nil {
				// we shall mock a snake
				aSnake := &mocks.Snaker{}
				aSnake.On("Size").Return(tt.mockSize, nil)
				aGameBoard.movingSnake = aSnake
			}
			gotSize, err := aGameBoard.SnakeSize()
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr, "%w", err)
			require.Equal(t, tt.wantSize, gotSize)
		})
	}
}

func TestGameBoard_MoveSnake(t *testing.T) {
	type fields struct {
		size        common.Size
		board       [][]rune
		movingSnake snake.Snaker
		candy       candy.Candyer
	}
	tests := []struct {
		name           string
		fields         fields
		mockNextMove   common.Position
		mockOldTail    common.Position
		wantOldValue   rune
		wantListSprite []common.Sprite
		wantTypeErr    error
		wantErr        bool
	}{
		{
			name: "TestEmptySnake", // The snake has no body, it can't be moved
			fields: fields{
				size:        testdata.Size3_3,
				board:       testdata.Duplicate(testdata.Board3_3),
				movingSnake: snake.New(),
			},
			wantTypeErr: snake.ErrNoSnakeBody,
			wantErr:     true,
		},
		{
			name:         "TestEmptyBoardSnakeMove1,1", // There is no board, no move possible
			mockNextMove: testdata.Position1_1,
			wantTypeErr:  ErrInvalidPosition,
			wantOldValue: 0,
			wantErr:      true,
		},
		{
			name: "TestSnakePos0,0Move1,1", // The snake moves to 1,1 from 0,0
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3),
			},
			mockNextMove: testdata.Position1_1,
			mockOldTail:  testdata.Position0_0,
			wantOldValue: testdata.Duplicate(testdata.Board3_3)[1][1],
			wantListSprite: []common.Sprite{
				{
					Value:    SnakePart,
					Position: testdata.Position1_1,
				},
				{
					Value:    FreeSpace,
					Position: testdata.Position0_0,
				},
			},
			wantErr: false,
		},
		{
			name: "TestSnakePos0,0GrowTo1,1", // The snake eats a candy and grow from 0,0 to 1,1
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3Candy1_1),
			},
			mockNextMove: testdata.Position1_1,
			mockOldTail:  testdata.Position0_0,
			wantOldValue: testdata.Duplicate(testdata.Board3_3Candy1_1)[1][1],
			wantListSprite: []common.Sprite{
				{
					Value:    SnakePart,
					Position: testdata.Position1_1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameBoard := &gameBoard{
				size:        tt.fields.size,
				board:       tt.fields.board,
				movingSnake: tt.fields.movingSnake,
				candy:       tt.fields.candy,
			}
			if aGameBoard.movingSnake == nil {
				// we shall mock a snake
				aSnake := &mocks.Snaker{}
				aSnake.On("NextMove").Return(tt.mockNextMove, nil)
				aSnake.On("GrowTo", tt.mockNextMove).Return(nil)
				aSnake.On("MoveTo", tt.mockNextMove).Return(tt.mockOldTail, nil)
				aGameBoard.movingSnake = aSnake
			}
			gotOldValue, gotListSprite, err := aGameBoard.MoveSnake()
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr, "%w", err)
			if tt.wantTypeErr != nil {
				require.ErrorIs(t, err, tt.wantTypeErr)
			}
			require.Equal(t, tt.wantOldValue, gotOldValue)
			require.Equal(t, tt.wantListSprite, gotListSprite)
		})
	}
}

func TestGameBoard_actualMove(t *testing.T) {
	type fields struct {
		size        common.Size
		board       [][]rune
		movingSnake snake.Snaker
		candy       candy.Candyer
	}
	type args struct {
		position common.Position
		oldValue rune
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		mockNextMove   common.Position
		mockOldTail    common.Position
		wantTypeErr    error
		wantListSprite []common.Sprite
		wantErr        bool
	}{
		{
			name: "TestEmptySnake", // The snake has no body, it can't be moved
			fields: fields{
				size:        testdata.Size3_3,
				board:       testdata.Duplicate(testdata.Board3_3),
				movingSnake: snake.New(),
			},
			args: args{
				position: testdata.Position1_1,
				oldValue: FreeSpace,
			},
			wantTypeErr: snake.ErrNoSnakeBody,
			wantErr:     true,
		},
		{
			name: "TestEmptyBoardSnakeMove1,1", // There is no board, no move possible
			args: args{
				position: testdata.Position1_1,
				oldValue: FreeSpace,
			},
			mockNextMove: testdata.Position1_1,
			wantTypeErr:  ErrInvalidPosition,
			wantErr:      true,
		},
		{
			name: "TestSnakePos0,0Move1,1", // The snake moves to 1,1 from 0,0
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3),
			},
			args: args{
				position: testdata.Position1_1,
				oldValue: FreeSpace,
			},
			mockNextMove: testdata.Position1_1,
			mockOldTail:  testdata.Position0_0,
			wantListSprite: []common.Sprite{
				{
					Value:    SnakePart,
					Position: testdata.Position1_1,
				},
				{
					Value:    FreeSpace,
					Position: testdata.Position0_0,
				},
			},
			wantErr: false,
		},
		{
			name: "TestSnakePos0,0GrowTo1,1", // The snake eats a candy and grow from 0,0 to 1,1
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3Candy1_1),
			},
			args: args{
				position: testdata.Position1_1,
				oldValue: CandyBody,
			},
			mockNextMove: testdata.Position1_1,
			mockOldTail:  testdata.Position0_0,
			wantListSprite: []common.Sprite{
				{
					Value:    SnakePart,
					Position: testdata.Position1_1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameBoard := &gameBoard{
				size:        tt.fields.size,
				board:       tt.fields.board,
				movingSnake: tt.fields.movingSnake,
				candy:       tt.fields.candy,
			}
			if aGameBoard.movingSnake == nil {
				// we shall mock a snake
				aSnake := &mocks.Snaker{}
				aSnake.On("NextMove").Return(tt.mockNextMove, nil)
				aSnake.On("GrowTo", tt.mockNextMove).Return(nil)
				aSnake.On("MoveTo", tt.mockNextMove).Return(tt.mockOldTail, nil)
				aGameBoard.movingSnake = aSnake
			}
			gotListSprite, err := aGameBoard.actualMove(tt.args.position, tt.args.oldValue)

			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr, "%w", err)
			if tt.wantTypeErr != nil {
				require.ErrorIs(t, err, tt.wantTypeErr)
			}
			require.Equal(t, tt.wantListSprite, gotListSprite)
		})
	}
}

func TestGameBoard_cell(t *testing.T) {
	type fields struct {
		size        common.Size
		board       [][]rune
		movingSnake snake.Snaker
		candy       candy.Candyer
	}
	type args struct {
		position common.Position
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantValue   rune
		wantTypeErr error
		wantErr     bool
	}{
		{
			name: "TestEmptyBoard", // There is no board, error returned
			args: args{
				position: testdata.Position1_1,
			},
			wantValue:   0,
			wantTypeErr: ErrInvalidPosition,
			wantErr:     true,
		},
		{
			name: "TestBoard3_3Pos0_0", // Return a value
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3),
			},
			args: args{
				position: testdata.Position0_0,
			},
			wantValue: testdata.Duplicate(testdata.Board3_3)[0][0],
			wantErr:   false,
		},
		{
			name: "TestBoard3_3Candy1_1Pos1_1", // Return a candy
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3Candy1_1),
			},
			args: args{
				position: testdata.Position1_1,
			},
			wantValue: CandyBody,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameBoard := &gameBoard{
				size:        tt.fields.size,
				board:       tt.fields.board,
				movingSnake: tt.fields.movingSnake,
				candy:       tt.fields.candy,
			}
			gotValue, err := aGameBoard.cell(tt.args.position)

			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr, "%w", err)
			if tt.wantTypeErr != nil {
				require.ErrorIs(t, err, tt.wantTypeErr)
			}
			require.Equal(t, tt.wantValue, gotValue)
		})
	}
}

func TestGameBoard_setCell(t *testing.T) {
	type fields struct {
		size        common.Size
		board       [][]rune
		movingSnake snake.Snaker
		candy       candy.Candyer
	}
	type args struct {
		position common.Position
		value    rune
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantTypeErr error
		wantErr     bool
	}{
		{
			name: "TestEmptyBoard", // There is no board, error returned
			args: args{
				position: testdata.Position1_1,
				value:    'Y',
			},
			wantTypeErr: ErrInvalidPosition,
			wantErr:     true,
		},
		{
			name: "TestBoard3_3Pos0_0",
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3),
			},
			args: args{
				position: testdata.Position0_0,
				value:    'Z',
			},
			wantErr: false,
		},
		{
			name: "TestBoard3_3Candy1_1Pos1_1",
			fields: fields{
				size:  testdata.Size3_3,
				board: testdata.Duplicate(testdata.Board3_3Candy1_1),
			},
			args: args{
				position: testdata.Position1_1,
				value:    'X',
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameBoard := &gameBoard{
				size:        tt.fields.size,
				board:       tt.fields.board,
				movingSnake: tt.fields.movingSnake,
				candy:       tt.fields.candy,
			}
			err := aGameBoard.setCell(tt.args.position, tt.args.value)
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr, "%w", err)
			if tt.wantTypeErr != nil {
				require.ErrorIs(t, err, tt.wantTypeErr)
			}
		})
	}
}

func TestGameBoard_translatePosition(t *testing.T) {
	type fields struct {
		size        common.Size
		board       [][]rune
		movingSnake snake.Snaker
		candy       candy.Candyer
	}
	type args struct {
		requestedPosition common.Position
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantPosition common.Position
	}{
		{
			name: "TestEmptyBoard", // The size of the board is 0,0
			args: args{
				requestedPosition: testdata.Position0_0,
			},
			wantPosition: testdata.Position0_0,
		},
		{
			name: "TestSize3_3Pos1,1",
			fields: fields{
				size: testdata.Size3_3,
			},
			args: args{
				requestedPosition: testdata.Position1_1,
			},
			wantPosition: testdata.Position1_1,
		},
		{
			name: "TestSize4_4Pos4,3",
			fields: fields{
				size: testdata.Size4_4,
			},
			args: args{
				requestedPosition: testdata.Position4_3,
			},
			wantPosition: testdata.Position0_3,
		},
		{
			name: "TestSize4_4Pos3,4",
			fields: fields{
				size: testdata.Size4_4,
			},
			args: args{
				requestedPosition: testdata.Position3_4,
			},
			wantPosition: testdata.Position3_0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameBoard := &gameBoard{
				size:        tt.fields.size,
				board:       tt.fields.board,
				movingSnake: tt.fields.movingSnake,
				candy:       tt.fields.candy,
			}
			gotPosition := aGameBoard.translatePosition(tt.args.requestedPosition)
			require.Equal(t, tt.wantPosition, gotPosition)
		})
	}
}

func TestNew(t *testing.T) {
	var wantType *gameBoard
	var got = New()
	require.IsType(t, wantType, got)
}

func Test_random(t *testing.T) {
	type args struct {
		max int
	}
	tests := []struct {
		name    string
		args    args
		wantRnd int
		wantErr bool
	}{
		{
			name: "TestRandomMax1",
			args: args{
				max: 1,
			},
			wantRnd: 0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRnd, err := random(tt.args.max)
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr)
			require.Equal(t, tt.wantRnd, gotRnd)
		})
	}
}
