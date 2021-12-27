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
// The errors from the routines are channeled back to the main function
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
	"fmt"
	"gosnake/pkg/common"
	"gosnake/pkg/gamestate"
	"gosnake/pkg/uimanager"
	"strings"
	"time"
)

const (
	defaultBoardSize = 40
	sizeIncrement    = 10
	refreshInterval  = 100 * time.Millisecond // Defines the animations refresh rate
)

func main() {
	var (
		gameState     = gamestate.New()
		userInterface = uimanager.New()
		boardSize     = common.Size{
			Width:  defaultBoardSize,
			Height: defaultBoardSize,
		}
		scrollOver = true
		err        error // main function errors
		errChn     error // errors channeled from routines are written in errChn
	)

	// When terminating if any, errChn or err are displayed
	defer reportError(common.GetCurrentFuncName(), &err, &errChn)

	// Inits the user interface library
	if err = openUI(userInterface); err != nil {
		return
	}
	defer closeUI(userInterface)

	// Inits the state and creates the gameBoard
	if err = initGame(gameState, common.Size{
		Width:  defaultBoardSize,
		Height: defaultBoardSize,
	}); err != nil {
		return
	}

	// Creates the UI layout
	if err = createViews(gameState, userInterface, boardSize); err != nil {
		return
	}

	// Clears the board view
	if err = clearView(userInterface, boardViewTitle); err != nil {
		return
	}

	// Creates and displays the snake and the candy
	if err = displayPlayers(gameState, userInterface); err != nil {
		return
	}

	// Attaches the event handler
	if err = setEventHandler(gameState,
		userInterface, &scrollOver, &boardSize, &errChn); err != nil {
		return
	}

	// Enters the user interface main loop
	// which will quit when it receives uimanager.ErrQuit from the event handler
	err = eventLoop(userInterface)
}

func reportError(funcName string, err, errChn *error) {
	// Called in defer to report errors before quitting

	var errPanic = errors.New("runtime error")

	// Is there a aPanic going on?
	if aPanic := recover(); aPanic != nil {
		// Test the type of aPanic
		switch val := aPanic.(type) {
		case string:
			*err = errors.New(val)
		case error:
			*err = val
		default: // or simply convert aPanic to string then create a new error...
			strErr := fmt.Sprint(aPanic)
			*err = errors.New(strErr)
		}

		*err = fmt.Errorf(funcName+": %w", *err)
	}
	// Is there an error coming from routines?
	// Then errChn will be reported first
	if *errChn != nil {
		*err = *errChn
	}

	if *err != nil {
		// If there was a panic (runtime error) somewhere we emphasize the report
		if strings.Contains((*err).Error(), errPanic.Error()) {
			fmt.Println("Panic occurred")
		}
		// uimanager uses an error to tell the mainloop to quit
		// So it is not an error
		if !errors.Is(*err, uimanager.ErrQuit) {
			fmt.Println(*err)
		}
	}
}

func initGame(gameState gamestate.GameStater, boardSize common.Size) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	if err := gameState.InitBoard(boardSize); err != nil {
		return err
	}

	return nil
}

func displayPlayers(gameState gamestate.GameStater, userInterface uimanager.UIManagerer) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	// Creates the objects
	var listSprite []common.Sprite

	if listSprite, err = createObjects(gameState); err != nil {
		return err
	}

	// Displays the objects to the board view
	return updateView(userInterface, boardViewTitle, listSprite)
}

func openUI(userInterface uimanager.UIManagerer) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	return userInterface.OpenUIManager()
}

func closeUI(userInterface uimanager.UIManagerer) {
	userInterface.Close()
}

func setEventHandler(gameState gamestate.GameStater, userInterface uimanager.UIManagerer,
	scollOver *bool, boardSize *common.Size, errChan *error) (err error) {

	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	// We use a closure
	theHandler := func(key uimanager.Key) error {
		// which will have access to the surrounding parameters
		return handleKeyPress(gameState, userInterface, key,
			scollOver, boardSize, errChan)
	}

	return userInterface.OnKeyPress(theHandler)
}

func eventLoop(userInterface uimanager.UIManagerer) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	return userInterface.MainLoop()
}

func createObjects(gameState gamestate.GameStater) (listSprite []common.Sprite, err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	// Creates the snake and candy
	return gameState.CreateObjects()
}

func updateView(userInterface uimanager.UIManagerer, viewName string, spriteList []common.Sprite) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	return userInterface.Update(viewName, spriteList)
}

func handleKeyPress(gameState gamestate.GameStater, userInterface uimanager.UIManagerer,
	key uimanager.Key, scrollOver *bool, boardSize *common.Size, errChan *error) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	switch key {
	case uimanager.KeyCtrlC:
		return userInterface.Quit()
	case uimanager.KeyArrowUp:
		gameState.MoveUp()
		return nil
	case uimanager.KeyArrowDown:
		gameState.MoveDown()
		return nil
	case uimanager.KeyArrowLeft:
		gameState.MoveLeft()
		return nil
	case uimanager.KeyArrowRight:
		gameState.MoveRight()
		return nil
	case uimanager.KeySpace:
		if !gameState.GameInProgress() && *scrollOver {
			if gameState.Dirty() {
				if err := prepareGame(gameState, userInterface, boardSize); err != nil {
					return err
				}
			}

			if err := startGame(gameState, userInterface, scrollOver, errChan); err != nil {
				return err
			}
		}

		return nil
	case uimanager.KeyEnter:
		if !gameState.GameInProgress() && *scrollOver {
			toggleBoardViewSize(boardSize)

			if err := prepareGame(gameState, userInterface, boardSize); err != nil {
				return err
			}
		}

		return nil
	}

	return nil
}

func toggleBoardViewSize(boardSize *common.Size) {
	// Change the size of the board by cycling threw 10;20;30;40 (default)
	if boardSize.Width < defaultBoardSize {
		boardSize.Width += sizeIncrement
		boardSize.Height += sizeIncrement
	} else {
		boardSize.Width = sizeIncrement
		boardSize.Height = sizeIncrement
	}
}

func prepareGame(gameState gamestate.GameStater, userInterface uimanager.UIManagerer, boardSize *common.Size) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	if err := initGame(gameState, *boardSize); err != nil {
		return err
	}

	if err := createBoardView(userInterface, *boardSize); err != nil {
		return err
	}

	if err := clearView(userInterface, boardViewTitle); err != nil {
		return err
	}

	if err := displayPlayers(gameState, userInterface); err != nil {
		return err
	}

	if err := createScoreView(gameState, userInterface); err != nil {
		return err
	}

	return nil
}

func startGame(gameState gamestate.GameStater, userInterface uimanager.UIManagerer,
	scrollOver *bool, errChn *error) (err error) {

	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	gameState.Start()

	go func() {
		// If for whatever reason a panic occures we handle, display and chanel it
		defer handlePanic(userInterface, errChn)

		// launch the gameEngine
		errChan := make(chan error)

		go gameEngine(gameState, userInterface, errChan)
		*errChn = <-errChan

		// Launch the game over animation
		if *errChn == nil && !gameState.GameInProgress() {
			errChan := make(chan error)
			go gameOverAnim(userInterface, scrollOver, errChan)
			*errChn = <-errChan
		}
	}()

	return err
}

func handlePanic(userInterface uimanager.UIManagerer, errChn *error) {
	if aPanic := recover(); aPanic != nil {
		var err error
		switch val := aPanic.(type) {
		case string:
			err = errors.New(val)
		case error:
			err = val
		default: // or simply convert aPanic to string then create a new error...
			strErr := fmt.Sprint(aPanic)
			err = errors.New(strErr)
		}
		err = fmt.Errorf(common.GetCurrentFuncName()+": %w", err)
		errChan := make(chan error)
		go handleRoutineError(userInterface, errChan, &err, "   Start Game")
		*errChn = <-errChan
	}
}

func gameEngine(gameState gamestate.GameStater, userInterface uimanager.UIManagerer, errChan chan error) {
	var err error

	defer gameState.SetGameInProgress(false)
	defer handleRoutineError(userInterface, errChan, &err, "   Game Engine")
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	// The game loop
	ticker := time.NewTicker(refreshInterval)
	for range ticker.C {
		var spriteList []common.Sprite

		if spriteList, err = gameState.Play(); err != nil {
			break
		}

		if err = createScoreView(gameState, userInterface); err != nil {
			break
		}

		if err = updateView(userInterface, boardViewTitle, spriteList); err != nil {
			break
		}

		if !gameState.GameInProgress() {
			break
		}
	}
}

func gameOverAnim(userInterface uimanager.UIManagerer, scrollOver *bool, errChan chan error) {
	var err error

	defer func() { *scrollOver = true }()
	defer handleRoutineError(userInterface, errChan, &err, "  Game Over Anim")
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	// These constants could be calculated with the view's available width
	const (
		posMax        = 16
		scrollMessage = "                    GAME OVER!!!              "
		chunkLength   = 18
	)

	scrollPosition := 0

	*scrollOver = false

	// The scroll loop
	ticker := time.NewTicker(refreshInterval)
	for range ticker.C {

		chunk := scrollMessage[scrollPosition : scrollPosition+chunkLength]

		if err = userInterface.UpdateLn(messageViewTitle, chunk); err != nil {
			break
		}

		if scrollPosition++; scrollPosition > posMax {
			break
		}
	}
}

func handleRoutineError(userInterface uimanager.UIManagerer, errChan chan error, err *error, title string) {
	if *err != nil {
		// For demo purpose only since the error message is truncated
		// It might not be readable until the user presses CTRL+C
		err2 := updateErrorView(*err, userInterface, title)
		// if there is an error from the UpdateErrorView, we report it first
		if err2 != nil {
			err = &err2
		}
	}
	// The error is channeled out via errChan
	errChan <- *err
}

func chunks(str string, lenChunk, nbChunks int) (chunks []string, err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)
	// The string str is split in nbChunks strings of size lenChunk
	chunks = append(chunks, str[0:])
	i := 0

	for lenChunk >= 0 && len(chunks[i]) > lenChunk {
		chunks = append(chunks, chunks[i][lenChunk:])
		chunks[i] = chunks[i][0:lenChunk]
		i++
	}

	// To meet the number of chunks requested, empty strings are added
	for i := 0; i < nbChunks; i++ {
		if len(chunks) == i {
			chunks = append(chunks, "")
		}
	}

	return chunks, nil
}
