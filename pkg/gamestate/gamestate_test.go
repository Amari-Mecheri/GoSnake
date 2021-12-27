package gamestate

import (
	"gosnake/mocks"
	"gosnake/pkg/common"
	"gosnake/pkg/gameboard"
	"gosnake/testdata"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGameState_InitBoard(t *testing.T) {
	type fields struct {
		gameInProgress bool
		round          int
		score          int
		highScore      int
		dirty          bool
		GameBoard      gameboard.GameBoarder
	}
	type args struct {
		size common.Size
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantMock    bool
		mockErrType error
		wantErrType error
		wantErr     bool
	}{
		{
			name:    "TestNilBoard", // The function will create an instance
			wantErr: false,
		},
		{
			name: "TestNegativeBoardSize",
			args: args{
				size: testdata.SizeMinus1_Minus1,
			},
			wantMock:    true,
			mockErrType: gameboard.ErrInvalidSize, // A negative size will fail
			wantErrType: gameboard.ErrInvalidSize,
			wantErr:     true,
		},
		{
			name:     "TestInitOk", // A size of 0 won't fail.
			wantMock: true,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameState := &gameState{
				gameInProgress: tt.fields.gameInProgress,
				round:          tt.fields.round,
				score:          tt.fields.score,
				highScore:      tt.fields.highScore,
				dirty:          tt.fields.dirty,
				GameBoarder:    tt.fields.GameBoard,
			}
			if tt.wantMock {
				aGameBoard := &mocks.GameBoarder{}
				aGameBoard.On("InitGameBoard", tt.args.size).Return(tt.mockErrType)
				aGameState.GameBoarder = aGameBoard
			}
			err := aGameState.InitBoard(tt.args.size)
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr)
			if tt.wantErrType != nil {
				require.ErrorIs(t, err, tt.wantErrType)
			}
		})
	}
}

func TestGameState_CreateObjects(t *testing.T) {
	type fields struct {
		gameInProgress bool
		round          int
		score          int
		highScore      int
		dirty          bool
		GameBoarder    gameboard.GameBoarder
	}
	tests := []struct {
		name              string
		fields            fields
		wantMock          bool
		mockBoardSize     common.Size
		mockSnakePosition common.Position
		mockCandyPosition common.Position
		mockErr           error
		wantListSprite    []common.Sprite
		wantErrType       error
		wantErr           bool
	}{
		{
			name:        "TestNilBoard",
			wantMock:    false,
			wantErrType: ErrInvalidBoardReference,
			wantErr:     true,
		},
		{
			name:          "TestEmptyBoard",
			wantMock:      true,
			mockBoardSize: testdata.Size0_0,
			mockErr:       gameboard.ErrInvalidPosition,
			wantErrType:   gameboard.ErrInvalidPosition,
			wantErr:       true,
		},
		{
			name:              "TestBoard3_3",
			wantMock:          true,
			mockBoardSize:     testdata.Size3_3,
			mockSnakePosition: testdata.Position1_1,
			mockCandyPosition: testdata.Position0_0,
			mockErr:           nil,
			wantListSprite: []common.Sprite{
				{
					Value:    gameboard.SnakePart,
					Position: testdata.Position1_1,
				},
				{
					Value:    gameboard.CandyBody,
					Position: testdata.Position0_0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameState := &gameState{
				gameInProgress: tt.fields.gameInProgress,
				round:          tt.fields.round,
				score:          tt.fields.score,
				highScore:      tt.fields.highScore,
				dirty:          tt.fields.dirty,
				GameBoarder:    tt.fields.GameBoarder,
			}
			if tt.wantMock {
				aGameBoard := &mocks.GameBoarder{}
				aGameBoard.On("BoardSize").Return(tt.mockBoardSize)
				aGameBoard.On("CreateSnake", tt.mockSnakePosition, goRight).Return(
					common.Sprite{
						Value:    gameboard.SnakePart,
						Position: tt.mockSnakePosition,
					},
					tt.mockErr,
				)
				aGameBoard.On("CreateCandy").Return(
					common.Sprite{
						Value:    gameboard.CandyBody,
						Position: tt.mockCandyPosition,
					},
					tt.mockErr,
				)
				aGameState.GameBoarder = aGameBoard
			}
			gotListSprite, err := aGameState.CreateObjects()
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr)
			if tt.wantErrType != nil {
				require.ErrorIs(t, err, tt.wantErrType)
			}
			require.Equal(t, tt.wantListSprite, gotListSprite)
		})
	}
}

func TestGameState_Start(t *testing.T) {
	type fields struct {
		gameInProgress bool
		round          int
		score          int
		highScore      int
		dirty          bool
		GameBoarder    gameboard.GameBoarder
	}
	tests := []struct {
		name               string
		fields             fields
		wantGameInProgress bool
		wantRound          int
		wantScore          int
		wantHighestScore   int
		wantDirty          bool
	}{
		{
			name: "Start",
			fields: fields{
				gameInProgress: false,
				round:          100,
				score:          15,
				highScore:      17,
				dirty:          false,
			},
			wantGameInProgress: true,
			wantRound:          0,
			wantScore:          0,
			wantHighestScore:   17,
			wantDirty:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameState := &gameState{
				gameInProgress: tt.fields.gameInProgress,
				round:          tt.fields.round,
				score:          tt.fields.score,
				highScore:      tt.fields.highScore,
				dirty:          tt.fields.dirty,
				GameBoarder:    tt.fields.GameBoarder,
			}
			aGameState.Start()
			require.Equal(t, tt.wantGameInProgress, aGameState.GameInProgress())
			require.Equal(t, tt.wantRound, aGameState.Round())
			require.Equal(t, tt.wantScore, aGameState.Score())
			require.Equal(t, tt.wantHighestScore, aGameState.HighScore())
			require.Equal(t, tt.wantDirty, aGameState.Dirty())
		})
	}
}

func TestGameState_Play(t *testing.T) {
	type fields struct {
		gameInProgress bool
		round          int
		score          int
		highScore      int
		dirty          bool
		GameBoard      gameboard.GameBoarder
	}
	tests := []struct {
		name               string
		fields             fields
		wantMock           bool
		mockBoardSize      common.Size
		mockOldValue       rune
		mockIsSnakePart    bool
		mockIsCandyBody    bool
		mockIsCandyAlive   bool
		mockSnakePosition  common.Position
		mockCandyPosition  common.Position
		mockListSprite     []common.Sprite
		mockErr            error
		wantGameInProgress bool
		wantListSprite     []common.Sprite
		wantScore          int
		wantHighScore      int
		wantErrType        error
		wantErr            bool
	}{
		{
			name:        "TestNilBoard",
			wantErrType: ErrInvalidBoardReference,
			wantErr:     true,
		},
		{
			name: "TestInvalidPosition",
			fields: fields{
				gameInProgress: true,
			},
			wantMock:           true,
			mockErr:            gameboard.ErrInvalidPosition,
			wantGameInProgress: false,
			wantErrType:        gameboard.ErrInvalidPosition,
			wantErr:            true,
		},
		{
			name: "TestMoveOK",
			fields: fields{
				gameInProgress: true,
			},
			wantMock:          true,
			mockBoardSize:     testdata.Size0_0,
			mockOldValue:      0,
			mockSnakePosition: testdata.Position0_0,
			mockCandyPosition: testdata.Position1_1,
			mockIsSnakePart:   false,
			mockIsCandyBody:   false,
			mockIsCandyAlive:  true,
			mockListSprite: []common.Sprite{
				{
					Value:    gameboard.SnakePart,
					Position: testdata.Position1_1,
				},
				{
					Value:    gameboard.FreeSpace,
					Position: testdata.Position0_0,
				},
			},
			wantGameInProgress: true,
			wantListSprite: []common.Sprite{
				{
					Value:    gameboard.SnakePart,
					Position: testdata.Position1_1,
				},
				{
					Value:    gameboard.FreeSpace,
					Position: testdata.Position0_0,
				},
			},
			wantErr: false,
		},
		{
			name: "TestGameOver",
			fields: fields{
				gameInProgress: true,
			},
			wantMock:          true,
			mockBoardSize:     testdata.Size0_0,
			mockOldValue:      0,
			mockSnakePosition: testdata.Position0_0,
			mockCandyPosition: testdata.Position1_1,
			mockIsSnakePart:   true, // The snake ate itself
			mockIsCandyBody:   false,
			mockIsCandyAlive:  true,
			mockListSprite: []common.Sprite{
				{
					Value:    gameboard.SnakePart,
					Position: testdata.Position1_1,
				},
				{
					Value:    gameboard.FreeSpace,
					Position: testdata.Position0_0,
				},
			},
			wantGameInProgress: false, // Then the game is over
			wantListSprite: []common.Sprite{
				{
					Value:    gameboard.SnakePart,
					Position: testdata.Position1_1,
				},
				{
					Value:    gameboard.FreeSpace,
					Position: testdata.Position0_0,
				},
			},
			wantErr: false,
		},
		{
			name: "TestEatTheCandyScore1HighSCore10",
			fields: fields{
				gameInProgress: true,
				score:          1,
				highScore:      10,
			},
			wantMock:          true,
			wantScore:         2,
			wantHighScore:     10,
			mockBoardSize:     testdata.Size0_0,
			mockOldValue:      0,
			mockSnakePosition: testdata.Position0_0,
			mockCandyPosition: testdata.Position1_1,
			mockIsSnakePart:   false,
			mockIsCandyBody:   true,  // The snake ate the candy
			mockIsCandyAlive:  false, // No more candies, a new one shall be generated
			mockListSprite: []common.Sprite{
				{
					Value:    gameboard.SnakePart,
					Position: testdata.Position1_1,
				},
				{
					Value:    gameboard.FreeSpace,
					Position: testdata.Position0_0,
				},
			},
			wantGameInProgress: true,
			wantListSprite: []common.Sprite{
				{
					Value: gameboard.SnakePart,
					Position: common.Position{
						X: 1,
						Y: 1,
					},
				},
				{
					Value: gameboard.FreeSpace,
					Position: common.Position{
						X: 0,
						Y: 0,
					},
				},
				{
					Value: gameboard.CandyBody,
					Position: common.Position{
						X: 1,
						Y: 1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "TestEatTheCandyScore10HighSCore10",
			fields: fields{
				gameInProgress: true,
				score:          10,
				highScore:      10,
			},
			wantMock:          true,
			wantScore:         11,
			wantHighScore:     11,
			mockBoardSize:     testdata.Size0_0,
			mockOldValue:      0,
			mockSnakePosition: testdata.Position0_0,
			mockCandyPosition: testdata.Position1_1,
			mockIsSnakePart:   false,
			mockIsCandyBody:   true,  // The snake ate the candy
			mockIsCandyAlive:  false, // No more candies, a new one shall be generated
			mockListSprite: []common.Sprite{
				{
					Value:    gameboard.SnakePart,
					Position: testdata.Position1_1,
				},
				{
					Value:    gameboard.FreeSpace,
					Position: testdata.Position0_0,
				},
			},
			wantGameInProgress: true,
			wantListSprite: []common.Sprite{
				{
					Value: gameboard.SnakePart,
					Position: common.Position{
						X: 1,
						Y: 1,
					},
				},
				{
					Value: gameboard.FreeSpace,
					Position: common.Position{
						X: 0,
						Y: 0,
					},
				},
				{
					Value: gameboard.CandyBody,
					Position: common.Position{
						X: 1,
						Y: 1,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameState := &gameState{
				gameInProgress: tt.fields.gameInProgress,
				round:          tt.fields.round,
				score:          tt.fields.score,
				highScore:      tt.fields.highScore,
				dirty:          tt.fields.dirty,
				GameBoarder:    tt.fields.GameBoard,
			}
			if tt.wantMock {
				aGameBoard := &mocks.GameBoarder{}
				aGameBoard.On("BoardSize").Return(tt.mockBoardSize)
				aGameBoard.On("MoveSnake").Return(
					tt.mockOldValue,
					tt.mockListSprite,
					tt.mockErr,
				)
				aGameBoard.On("IsSnakePart", tt.mockOldValue).Return(tt.mockIsSnakePart)
				aGameBoard.On("IsCandy", tt.mockOldValue).Return(tt.mockIsCandyBody)
				aGameBoard.On("CandyAlive").Return(tt.mockIsCandyAlive)
				aGameBoard.On("CreateCandy").Return(
					common.Sprite{
						Value:    gameboard.CandyBody,
						Position: tt.mockCandyPosition,
					},
					tt.mockErr,
				)
				aGameBoard.On("RemoveCandy").Return()
				aGameState.GameBoarder = aGameBoard
			}
			gotListSprite, err := aGameState.Play()
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr, err)
			if tt.wantErrType != nil {
				require.ErrorIs(t, err, tt.wantErrType)
			}
			require.Equal(t, tt.wantListSprite, gotListSprite, "gotListSprite")
			gotScore := aGameState.Score()
			gotHighScore := aGameState.HighScore()
			require.Equal(t, tt.wantScore, gotScore, "gotScore")
			require.Equal(t, tt.wantHighScore, gotHighScore, "gotHighScore")
		})
	}
}

func TestGameState_GameInProgress(t *testing.T) {
	type fields struct {
		gameInProgress bool
		round          int
		score          int
		highScore      int
		dirty          bool
		GameBoard      gameboard.GameBoarder
	}
	tests := []struct {
		name               string
		fields             fields
		wantGameInProgress bool
	}{
		{
			name: "GameInProgressFalse",
			fields: fields{
				gameInProgress: false,
			},
			wantGameInProgress: false,
		},
		{
			name: "GameInProgressTrue",
			fields: fields{
				gameInProgress: true,
			},
			wantGameInProgress: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameState := &gameState{
				gameInProgress: tt.fields.gameInProgress,
				round:          tt.fields.round,
				score:          tt.fields.score,
				highScore:      tt.fields.highScore,
				dirty:          tt.fields.dirty,
				GameBoarder:    tt.fields.GameBoard,
			}
			gotGameInProgress := aGameState.GameInProgress()
			require.Equal(t, tt.wantGameInProgress, gotGameInProgress, "gotGameInProgress")
		})
	}
}

func TestGameState_Dirty(t *testing.T) {
	type fields struct {
		gameInProgress bool
		round          int
		score          int
		highScore      int
		dirty          bool
		GameBoard      gameboard.GameBoarder
	}
	tests := []struct {
		name      string
		fields    fields
		wantDirty bool
	}{
		{
			name: "DirtyFalse",
			fields: fields{
				dirty: false,
			},
			wantDirty: false,
		},
		{
			name: "DirtyTrue",
			fields: fields{
				dirty: true,
			},
			wantDirty: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameState := &gameState{
				gameInProgress: tt.fields.gameInProgress,
				round:          tt.fields.round,
				score:          tt.fields.score,
				highScore:      tt.fields.highScore,
				dirty:          tt.fields.dirty,
				GameBoarder:    tt.fields.GameBoard,
			}
			gotDirty := aGameState.Dirty()
			require.Equal(t, tt.wantDirty, gotDirty, "gotDirty")
		})
	}
}

func TestGameState_HighestScore(t *testing.T) {
	type fields struct {
		gameInProgress bool
		round          int
		score          int
		highScore      int
		dirty          bool
		GameBoard      gameboard.GameBoarder
	}
	tests := []struct {
		name          string
		fields        fields
		wantHighscore int
	}{
		{
			name: "HighScore",
			fields: fields{
				highScore: 10,
			},
			wantHighscore: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameState := &gameState{
				gameInProgress: tt.fields.gameInProgress,
				round:          tt.fields.round,
				score:          tt.fields.score,
				highScore:      tt.fields.highScore,
				dirty:          tt.fields.dirty,
				GameBoarder:    tt.fields.GameBoard,
			}
			gotHighSCore := aGameState.HighScore()
			require.Equal(t, tt.wantHighscore, gotHighSCore, "gotHighSCore")
		})
	}
}

func TestGameState_Score(t *testing.T) {
	type fields struct {
		gameInProgress bool
		round          int
		score          int
		highScore      int
		dirty          bool
		GameBoard      gameboard.GameBoarder
	}
	tests := []struct {
		name      string
		fields    fields
		wantScore int
	}{
		{
			name: "Score",
			fields: fields{
				score: 7,
			},
			wantScore: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameState := &gameState{
				gameInProgress: tt.fields.gameInProgress,
				round:          tt.fields.round,
				score:          tt.fields.score,
				highScore:      tt.fields.highScore,
				dirty:          tt.fields.dirty,
				GameBoarder:    tt.fields.GameBoard,
			}
			gotScore := aGameState.Score()
			require.Equal(t, tt.wantScore, gotScore, "gotScore")
		})
	}
}

func TestGameState_Round(t *testing.T) {
	type fields struct {
		gameInProgress bool
		round          int
		score          int
		highScore      int
		dirty          bool
		GameBoard      gameboard.GameBoarder
	}
	tests := []struct {
		name      string
		fields    fields
		wantRound int
	}{
		{
			name: "Round",
			fields: fields{
				round: 17,
			},
			wantRound: 17,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGameState := &gameState{
				gameInProgress: tt.fields.gameInProgress,
				round:          tt.fields.round,
				score:          tt.fields.score,
				highScore:      tt.fields.highScore,
				dirty:          tt.fields.dirty,
				GameBoarder:    tt.fields.GameBoard,
			}
			gotRound := aGameState.Round()
			require.Equal(t, tt.wantRound, gotRound, "gotRound")
		})
	}
}

func TestNew(t *testing.T) {
	var wantType *gameState
	var got = New()
	require.IsType(t, wantType, got)
}
