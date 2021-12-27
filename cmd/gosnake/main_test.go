// amari.mecheri@gmail.com
//
// Architecture of GoSnake
//
// By choice, the package "main" is made of functions, variables are local to the main function and are sent as arguments where needed.
//
// uimanager, gamestate, gameboard, snake and candy packages export Interfaces, Methods implementing the corresponding interfaces and New() functions
// They also define a struct with the needed properties. These structs are not exported.
// The New() function of each package is used to get an instance of the corresponding object.
// The New function also performs instantiation of the childs by calling the respective New() functions.
// This allows the implementation of the mechanics of "classes" with private members, instantiated via a mandatory constructor/fabric (New function)
//
// inheritance (is a) is implemented by using an anonymous field of type child
// composition (has) is implemented by using a named field of type child
//
// => In the case of inheritance, exported methods of the child can be used and exported by the parent as if they were methods of the parent.
// => In the case of composition methods of the child are only accessible via the corresponding named field.
//
// main function
// 		| (var) ->	gamestate
// 						| (is a) -> gameboard
// 										| (has) -> snake
// 	 									| (has) -> candy
// 		| (var) ->	uimanager
// 						| (has) -> gocui
//
// uimanager encapsulates functions from Gocui https://github.com/jroimartin/gocui
//
// The main package controls the gamestate, creates the views layouts
// and updates them to reflect the state via uimanager
//
// The game logic is held by gamestate. The Play() methods plays a round.
//
// The rounds are called in a loop controlled by tickers at intervals
// The errors from the routines are channelled back to the main function
//
// In order to manage the keys pressed, the bindings are all affected to the same eventHandler which selects the appropriate action
// Since the eventHandler will only receive the key pressed as parameter, a closure is used to allow access to the main parameters
//
// Most functions and methods start with a defer common.ErrorWrapper
// Which handles unexpected panics and wraps any error with the function/method name
// This allows identification of the function where the error occurred
// as well as the chain of the functions called (since each parent function calls common.ErrorWrapper too).
// Where needed custom errors are created and later identified with errors Is().
// In the case of panics the content is checked, the panic can be an error or a string (cf common.ErrorWrapper test)
// The errors and unexpected panics can be controlled and logged if needed
// For example, if an error/panic occurs in the routines gameEngine or gameOverAnim, a red alert is displayed without terminating the program.
package main

import (
	"errors"
	"gosnake/mocks"
	"gosnake/pkg/common"
	"gosnake/pkg/gameboard"
	"gosnake/pkg/gamestate"
	"gosnake/pkg/uimanager"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_chunks(t *testing.T) {
	type args struct {
		str      string
		lenChunk int
		nbChunks int
	}
	tests := []struct {
		name       string
		args       args
		wantChunks []string
		wantErr    bool
	}{
		{
			name:       "EmptyString",
			wantChunks: []string{""},
		},
		{
			name: "EmptyStringChunks3",
			args: args{
				nbChunks: 3,
			},
			wantChunks: []string{"", "", ""},
		},
		{
			name: "EmptyStringChunks3Len5",
			args: args{
				lenChunk: 5,
				nbChunks: 3,
			},
			wantChunks: []string{"", "", ""},
		},
		{
			name: "String10Chunks3Len5",
			args: args{
				str:      "0123456789",
				lenChunk: 5,
				nbChunks: 3,
			},
			wantChunks: []string{"01234", "56789", ""},
		},
		{
			name: "String10Chunks3Len3",
			args: args{
				str:      "0123456789",
				lenChunk: 3,
				nbChunks: 3,
			},
			wantChunks: []string{"012", "345", "678", "9"},
		},
		{
			name: "String10Chunks4Len3",
			args: args{
				str:      "0123456789",
				lenChunk: 3,
				nbChunks: 4,
			},
			wantChunks: []string{"012", "345", "678", "9"},
		},
		{
			name: "String10Chunks3Len-5",
			args: args{
				str:      "0123456789",
				lenChunk: -5,
				nbChunks: 4,
			},
			wantChunks: []string{"0123456789", "", "", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChunks, err := chunks(tt.args.str, tt.args.lenChunk, tt.args.nbChunks)
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr)
			require.Equal(t, tt.wantChunks, gotChunks)
		})
	}
}

func Test_handleRoutineError(t *testing.T) {
	type args struct {
		userInterface uimanager.UIManagerer
		errChan       chan error
		err           *error
		origin        string
	}
	tests := []struct {
		name                    string
		args                    args
		argErr                  error // The error we provide to handleRoutineError
		mockDisplayRedLayoutErr error // The error returned by uimanager
		wantErrType             error // The returned error
		wantErr                 bool
	}{
		{
			name: "TestArgError",
			args: args{
				userInterface: &mocks.UIManagerer{},
				errChan:       make(chan error),
			},
			argErr:                  errors.New("ArgError"),
			mockDisplayRedLayoutErr: nil,
			wantErrType:             errors.New("ArgError"),
			wantErr:                 true,
		},
		{
			name: "TestArgErrorMockError",
			args: args{
				userInterface: &mocks.UIManagerer{},
				errChan:       make(chan error),
			},
			argErr:                  errors.New("ArgError"),
			mockDisplayRedLayoutErr: errors.New("ErrorView"),
			wantErrType:             errors.New("ErrorView"),
			wantErr:                 true,
		},
		{
			name: "TestOK",
			args: args{
				userInterface: &mocks.UIManagerer{},
				errChan:       make(chan error),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aUI := &mocks.UIManagerer{}
			aUI.On("DisplayRedLayout", errorViewTitle, mock.Anything).Return(tt.mockDisplayRedLayoutErr)
			tt.args.userInterface = aUI
			tt.args.err = &tt.argErr
			go handleRoutineError(tt.args.userInterface, tt.args.errChan, tt.args.err, tt.args.origin)
			err := <-tt.args.errChan
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr)
			if gotErr && tt.wantErrType != nil {
				require.Contains(t, err.Error(), tt.wantErrType.Error())
			}
		})
	}
}

func Test_gameOverAnim(t *testing.T) {
	type args struct {
		userInterface uimanager.UIManagerer
		scrollOver    *bool
		errChan       chan error
	}
	tests := []struct {
		name                    string
		args                    args
		mockUpdateLnErr         error
		mockDisplayRedLayoutErr error
		wantErrType             error
		wantErr                 bool
	}{
		{
			name: "TestMockNoError",
			args: args{
				userInterface: &mocks.UIManagerer{},
				errChan:       make(chan error),
				scrollOver:    new(bool),
			},
			mockUpdateLnErr:         nil,
			mockDisplayRedLayoutErr: nil,
			wantErr:                 false,
		},
		{
			name: "TestMockUpdateError",
			args: args{
				userInterface: &mocks.UIManagerer{},
				errChan:       make(chan error),
				scrollOver:    new(bool),
			},
			mockUpdateLnErr:         errors.New("UpdateError"),
			mockDisplayRedLayoutErr: nil,
			wantErrType:             errors.New("UpdateError"),
			wantErr:                 true,
		},
		{
			name: "TestDisplayError",
			args: args{
				userInterface: &mocks.UIManagerer{},
				errChan:       make(chan error),
				scrollOver:    new(bool),
			},
			// To have a DisplayError error we have to have an error in the first place
			// so let's use UpdateError to provide an error
			mockUpdateLnErr:         errors.New("UpdateError"),
			mockDisplayRedLayoutErr: errors.New("DisplayError"),
			wantErrType:             errors.New("DisplayError"),
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aUI := &mocks.UIManagerer{}
			aUI.On("UpdateLn", messageViewTitle, mock.Anything).Return(tt.mockUpdateLnErr)
			aUI.On("DisplayRedLayout", errorViewTitle, mock.Anything).Return(tt.mockDisplayRedLayoutErr)
			tt.args.userInterface = aUI
			go gameOverAnim(tt.args.userInterface, tt.args.scrollOver, tt.args.errChan)
			err := <-tt.args.errChan
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr, err)
			if gotErr && tt.wantErrType != nil {
				require.Contains(t, err.Error(), tt.wantErrType.Error())
			}
		})
	}
}

func Test_gameEngine(t *testing.T) {
	type args struct {
		gameState     gamestate.GameStater
		userInterface uimanager.UIManagerer
		//gameInProgress *bool
		errChan chan error
	}
	tests := []struct {
		name                    string
		args                    args
		mockGameInProgess       bool // has to be set to false otherwise the routine won't end
		MockSpriteList          []common.Sprite
		mockSnakeSize           int
		mockHighScore           int
		mockRound               int
		mockScore               int
		mockSnakeSizeErr        error
		mockSnakePosition       error
		MockPlayErr             error
		mockSetViewErr          error
		mockSetViewLayoutErr    error
		mockUpdateLnErr         error
		mockDisplayRedLayoutErr error
		mockUpdateErr           error
		wantErrType             error
		wantErr                 bool
	}{
		{
			name: "TestMockNoError",
			args: args{
				gameState:     &mocks.GameStater{},
				userInterface: &mocks.UIManagerer{},
				errChan:       make(chan error),
				//gameInProgress: new(bool),
			},
			mockGameInProgess: false,
			MockSpriteList: []common.Sprite{
				{
					Value: gameboard.SnakePart,
					Position: common.Position{
						X: 0,
						Y: 0,
					},
				},
			},
			mockSnakeSize:        0,
			mockHighScore:        0,
			mockRound:            0,
			mockScore:            0,
			MockPlayErr:          nil,
			mockSetViewErr:       nil,
			mockSetViewLayoutErr: nil,
			wantErrType:          nil,
			wantErr:              false,
		},
		{
			name: "TestMockNoError_negatives",
			args: args{
				gameState:     &mocks.GameStater{},
				userInterface: &mocks.UIManagerer{},
				errChan:       make(chan error),
				//gameInProgress: new(bool),
			},
			mockGameInProgess: false,
			MockSpriteList: []common.Sprite{
				{
					Value: gameboard.SnakePart,
					Position: common.Position{
						X: 0,
						Y: 0,
					},
				},
			},
			mockSnakeSize:        -1,
			mockHighScore:        -1,
			mockRound:            -1,
			mockScore:            -1,
			MockPlayErr:          nil,
			mockSetViewErr:       nil,
			mockSetViewLayoutErr: nil,
			wantErrType:          nil,
			wantErr:              false,
		},
		{
			name: "TestPlayError",
			args: args{
				gameState:     &mocks.GameStater{},
				userInterface: &mocks.UIManagerer{},
				errChan:       make(chan error),
				//gameInProgress: new(bool),
			},
			mockGameInProgess: false,
			MockSpriteList: []common.Sprite{
				{
					Value: gameboard.SnakePart,
					Position: common.Position{
						X: 0,
						Y: 0,
					},
				},
			},
			mockSnakeSize:        5,
			mockHighScore:        17,
			mockRound:            18,
			mockScore:            29,
			MockPlayErr:          errors.New("PlayError"),
			mockSetViewErr:       nil,
			mockSetViewLayoutErr: nil,
			wantErrType:          errors.New("PlayError"),
			wantErr:              true,
		},
		{
			name: "TestSetViewError",
			args: args{
				gameState:     &mocks.GameStater{},
				userInterface: &mocks.UIManagerer{},
				errChan:       make(chan error),
				//gameInProgress: new(bool),
			},
			mockGameInProgess: false,
			MockSpriteList: []common.Sprite{
				{
					Value: gameboard.SnakePart,
					Position: common.Position{
						X: 0,
						Y: 0,
					},
				},
			},
			mockSnakeSize:        5,
			mockHighScore:        17,
			mockRound:            18,
			mockScore:            29,
			MockPlayErr:          nil,
			mockSetViewErr:       errors.New("SetViewError"),
			mockSetViewLayoutErr: nil,
			wantErrType:          errors.New("SetViewError"),
			wantErr:              true,
		},
		{
			name: "TestSetViewLayoutError",
			args: args{
				gameState:     &mocks.GameStater{},
				userInterface: &mocks.UIManagerer{},
				errChan:       make(chan error),
				//gameInProgress: new(bool),
			},
			mockGameInProgess: false,
			MockSpriteList: []common.Sprite{
				{
					Value: gameboard.SnakePart,
					Position: common.Position{
						X: 0,
						Y: 0,
					},
				},
			},
			mockSnakeSize:        5,
			mockHighScore:        17,
			mockRound:            18,
			mockScore:            29,
			MockPlayErr:          nil,
			mockSetViewErr:       nil,
			mockSetViewLayoutErr: errors.New("SetViewLayoutError"),
			wantErrType:          errors.New("SetViewLayoutError"),
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGmState := &mocks.GameStater{}
			aGmState.On("Play").Return(tt.MockSpriteList, tt.MockPlayErr)
			aGmState.On("GameInProgress").Return(tt.mockGameInProgess)
			aGmState.On("SetGameInProgress", false).Return()
			aGmState.On("BoardSize").Return(common.Size{
				Width:  0,
				Height: 0,
			})
			aGmState.On("SnakePosition").Return(
				common.Position{
					X: 0,
					Y: 0,
				}, nil)
			aGmState.On("SnakeSize").Return(tt.mockSnakeSize, nil)
			aGmState.On("Round").Return(tt.mockRound)
			aGmState.On("Score").Return(tt.mockScore)
			aGmState.On("HighScore").Return(tt.mockHighScore)
			tt.args.gameState = aGmState
			aUI := &mocks.UIManagerer{}

			aUI.On("SetView", scoreViewTitle, mock.Anything).Return(tt.mockSetViewErr)
			aUI.On("SetViewLayout", scoreViewTitle, mock.Anything).Return(tt.mockSetViewLayoutErr)
			aUI.On("UpdateLn", messageViewTitle, mock.Anything).Return(tt.mockUpdateLnErr)
			aUI.On("DisplayRedLayout", errorViewTitle, mock.Anything).Return(tt.mockDisplayRedLayoutErr)
			aUI.On("Update", boardViewTitle, mock.Anything).Return(tt.mockUpdateErr)
			tt.args.userInterface = aUI

			go gameEngine(tt.args.gameState, tt.args.userInterface, tt.args.errChan)
			err := <-tt.args.errChan
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr, err)
			if gotErr && tt.wantErrType != nil {
				require.Contains(t, err.Error(), tt.wantErrType.Error())
			}
			require.Equal(t, tt.mockGameInProgess, tt.args.gameState.GameInProgress())
		})
	}
}

func Test_startGame(t *testing.T) {
	type args struct {
		gameState     gamestate.GameStater
		userInterface uimanager.UIManagerer
		scrollOver    *bool
		errChn        *error
	}
	tests := []struct {
		name                    string
		args                    args
		mockGameInProgess       bool // has to be set to false otherwise the routine won't end
		MockSpriteList          []common.Sprite
		mockSnakeSize           int
		mockHighScore           int
		mockRound               int
		mockScore               int
		mockSnakeSizeErr        error
		mockSnakePosition       error
		MockPlayErr             error
		mockSetViewErr          error
		mockSetViewLayoutErr    error
		mockUpdateLnErr         error
		mockDisplayRedLayoutErr error
		mockUpdateErr           error
		wantErrType             error
		wantErr                 bool
	}{
		{
			name: "TestMockNoError",
			args: args{
				gameState:     &mocks.GameStater{},
				userInterface: &mocks.UIManagerer{},
				scrollOver:    new(bool),
				errChn:        new(error),
			},
			mockGameInProgess: false,
			MockSpriteList: []common.Sprite{
				{
					Value: gameboard.SnakePart,
					Position: common.Position{
						X: 0,
						Y: 0,
					},
				},
			},
			mockSnakeSize:        0,
			mockHighScore:        0,
			mockRound:            0,
			mockScore:            0,
			MockPlayErr:          nil,
			mockSetViewErr:       nil,
			mockSetViewLayoutErr: nil,
			wantErrType:          nil,
			wantErr:              false,
		},
		{
			name: "TestMockSetViewLayoutError",
			args: args{
				gameState:     &mocks.GameStater{},
				userInterface: &mocks.UIManagerer{},
				scrollOver:    new(bool),
				errChn:        new(error),
			},
			mockGameInProgess: false,
			MockSpriteList: []common.Sprite{
				{
					Value: gameboard.SnakePart,
					Position: common.Position{
						X: 0,
						Y: 0,
					},
				},
			},
			mockSnakeSize:        0,
			mockHighScore:        0,
			mockRound:            0,
			mockScore:            0,
			MockPlayErr:          nil,
			mockSetViewErr:       nil,
			mockSetViewLayoutErr: errors.New("SetViewLayoutError"),
			wantErrType:          errors.New("SetViewLayoutError"),
			wantErr:              true,
		},
		{
			name: "TestMockSetViewError",
			args: args{
				gameState:     &mocks.GameStater{},
				userInterface: &mocks.UIManagerer{},
				scrollOver:    new(bool),
				errChn:        new(error),
			},
			mockGameInProgess: false,
			MockSpriteList: []common.Sprite{
				{
					Value: gameboard.SnakePart,
					Position: common.Position{
						X: 0,
						Y: 0,
					},
				},
			},
			mockSnakeSize:        0,
			mockHighScore:        0,
			mockRound:            0,
			mockScore:            0,
			MockPlayErr:          nil,
			mockSetViewErr:       errors.New("SetViewError"),
			mockSetViewLayoutErr: nil,
			wantErrType:          errors.New("SetViewError"),
			wantErr:              true,
		},
		{
			name: "TestMockUpdateLnError",
			args: args{
				gameState:     &mocks.GameStater{},
				userInterface: &mocks.UIManagerer{},
				scrollOver:    new(bool),
				errChn:        new(error),
			},
			mockGameInProgess: false,
			MockSpriteList: []common.Sprite{
				{
					Value: gameboard.SnakePart,
					Position: common.Position{
						X: 0,
						Y: 0,
					},
				},
			},
			mockSnakeSize:           0,
			mockHighScore:           0,
			mockRound:               0,
			mockScore:               0,
			mockSnakeSizeErr:        nil,
			mockSnakePosition:       nil,
			MockPlayErr:             nil,
			mockSetViewErr:          nil,
			mockSetViewLayoutErr:    nil,
			mockUpdateLnErr:         errors.New("UpdateLnError"),
			mockDisplayRedLayoutErr: nil,
			mockUpdateErr:           nil,
			wantErrType:             errors.New("UpdateLnError"),
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aGmState := &mocks.GameStater{}
			aGmState.On("Start").Return()
			aGmState.On("Play").Return(tt.MockSpriteList, tt.MockPlayErr)
			aGmState.On("GameInProgress").Return(tt.mockGameInProgess)
			aGmState.On("SetGameInProgress", false).Return()
			aGmState.On("BoardSize").Return(common.Size{
				Width:  0,
				Height: 0,
			})
			aGmState.On("SnakePosition").Return(
				common.Position{
					X: 0,
					Y: 0,
				}, nil)
			aGmState.On("SnakeSize").Return(tt.mockSnakeSize, nil)
			aGmState.On("Round").Return(tt.mockRound)
			aGmState.On("Score").Return(tt.mockScore)
			aGmState.On("HighScore").Return(tt.mockHighScore)
			tt.args.gameState = aGmState
			aUI := &mocks.UIManagerer{}

			aUI.On("SetView", scoreViewTitle, mock.Anything).Return(tt.mockSetViewErr)
			aUI.On("SetViewLayout", scoreViewTitle, mock.Anything).Return(tt.mockSetViewLayoutErr)
			aUI.On("UpdateLn", messageViewTitle, mock.Anything).Return(tt.mockUpdateLnErr)
			aUI.On("DisplayRedLayout", errorViewTitle, mock.Anything).Return(tt.mockDisplayRedLayoutErr)
			aUI.On("Update", boardViewTitle, mock.Anything).Return(tt.mockUpdateErr)
			tt.args.userInterface = aUI
			// To check that we got out the routines we set errChn to an error
			*tt.args.errChn = errors.New("Start")
			*tt.args.scrollOver = true
			//refreshInterval = time.Millisecond
			err := startGame(tt.args.gameState, tt.args.userInterface, tt.args.scrollOver, tt.args.errChn)
			// Wait for any error to be received from the channel
			time.Sleep(1 * time.Second)
			// checks/waits for the gameOverAnim routine to terminate
			for !*tt.args.scrollOver { //nolint
			}
			gotErr := (*tt.args.errChn != nil) || (err != nil)
			require.Equal(t, tt.wantErr, gotErr, *tt.args.errChn)
			if gotErr {
				require.NotNil(t, tt.wantErrType, "wantErrType is nil")
				if tt.wantErrType != nil {
					require.Contains(t, (*tt.args.errChn).Error(), tt.wantErrType.Error())
				}
			}
		})
	}
}

func Test_toggleBoardViewSize(t *testing.T) {
	type args struct {
		boardSize *common.Size
	}
	tests := []struct {
		name     string
		args     args
		wantSize common.Size
	}{
		{
			name: "testSize0",
			args: args{
				boardSize: &common.Size{
					Width:  0,
					Height: 0,
				},
			},
			wantSize: common.Size{
				Width:  sizeIncrement,
				Height: sizeIncrement,
			},
		},
		{
			name: "testSize40",
			args: args{
				boardSize: &common.Size{
					Width:  defaultBoardSize,
					Height: defaultBoardSize,
				},
			},
			wantSize: common.Size{
				Width:  sizeIncrement,
				Height: sizeIncrement,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toggleBoardViewSize(tt.args.boardSize)
			require.Equal(t, tt.wantSize, *tt.args.boardSize)
		})
	}
}
