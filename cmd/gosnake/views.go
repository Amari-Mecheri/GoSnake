package main

import (
	"errors"
	"fmt"
	"gosnake/pkg/common"
	"gosnake/pkg/gameboard"
	"gosnake/pkg/gamestate"
	"gosnake/pkg/uimanager"
	"strconv"
)

const (
	boardViewTitle   = "boardView"
	scoreViewTitle   = "scoreView"
	messageViewTitle = "messageView"
	errorViewTitle   = "errorView"
	helpViewTitle    = "helpView"
	gameFrameTitle   = "frameView"
	panelViewTitle   = "panelView"
)

const (
	leftMost       = 0
	rightPanel     = 43
	maxX           = 63
	topMost        = 0
	topMessageView = 11
	topErrorView   = 14
	topHelpView    = 27
	maxY           = 41
)

func createViews(gameState gamestate.GameStater, userInterface uimanager.UIManagerer, boardSize common.Size) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	if err := createGameFrame(userInterface); err != nil {
		return err
	}

	if err := createPanelView(userInterface); err != nil {
		return err
	}

	if err := createErrorView(userInterface); err != nil {
		return err
	}

	if err := createHelpView(userInterface); err != nil {
		return err
	}

	if err := createScoreView(gameState, userInterface); err != nil {
		return err
	}

	if err := createMessageView(userInterface); err != nil {
		return err
	}

	if err := createBoardView(userInterface, boardSize); err != nil {
		return err
	}

	return nil
}

func clearView(userInterface uimanager.UIManagerer, viewName string) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	return userInterface.ClearView(viewName)
}

func createGameFrame(userInterface uimanager.UIManagerer) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	var gameFramePosition = common.ViewPosition{
		X1: leftMost,
		Y1: topMost,
		X2: maxX,
		Y2: maxY,
	}

	return userInterface.SetView(gameFrameTitle, gameFramePosition)
}

func createPanelView(userInterface uimanager.UIManagerer) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	var panelViewPosition = common.ViewPosition{
		X1: rightPanel,
		Y1: topMost,
		X2: maxX,
		Y2: maxY,
	}

	return userInterface.SetView(panelViewTitle, panelViewPosition)
}

func createErrorView(userInterface uimanager.UIManagerer) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	var errorViewPosition = common.ViewPosition{
		X1: rightPanel,
		Y1: topErrorView,
		X2: maxX,
		Y2: topHelpView - 1,
	}

	return userInterface.SetView(errorViewTitle, errorViewPosition)
}

func createBoardView(userInterface uimanager.UIManagerer, boardSize common.Size) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	var gameBoardPosition = common.ViewPosition{
		X1: leftMost,
		Y1: topMost,
		X2: boardSize.Width + 1,
		Y2: boardSize.Height + 1,
	}
	// takes the frame into account and avoids scrolling issues (!workaround)
	gameBoardPosition.X2++

	return userInterface.SetView(boardViewTitle, gameBoardPosition)
}

func createMessageView(userInterface uimanager.UIManagerer) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	var messageViewPosition = common.ViewPosition{
		X1: rightPanel,
		Y1: topMessageView,
		X2: maxX,
		Y2: topErrorView - 1,
	}

	return userInterface.SetView(messageViewTitle, messageViewPosition)
}

func createHelpView(userInterface uimanager.UIManagerer) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	var helpViewPosition = common.ViewPosition{
		X1: rightPanel,
		Y1: topHelpView,
		X2: maxX,
		Y2: maxY,
	}

	if err := userInterface.SetView(helpViewTitle, helpViewPosition); err != nil {
		return err
	}

	helpViewLayout := []string{
		"  The Snake Game",
		"GRAB the * CANDIES",
		"",
		"select board size",
		"   with ENTER",
		"",
		" SPACEBAR to start",
		"",
		"Keys:  BOTTOM, UP",
		"      LEFT, RIGHT",
		"",
		"",
		"  Ctrl+C to Quit",
	}

	return userInterface.SetViewLayout(helpViewTitle, helpViewLayout)
}

func createScoreView(gameState gamestate.GameStater, userInterface uimanager.UIManagerer) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	var scoreViewPosition = common.ViewPosition{
		X1: rightPanel,
		Y1: topMost,
		X2: maxX,
		Y2: topMessageView - 1,
	}

	if err = userInterface.SetView(scoreViewTitle, scoreViewPosition); err != nil {
		return err
	}

	var boardSize = strconv.Itoa(gameState.BoardSize().Width) +
		"x" + strconv.Itoa(gameState.BoardSize().Height)

	var (
		snkSize       string
		snakeSize     int
		snakePosition common.Position
		snkPosition   string
	)

	// if the snake has no body default values will be displayed
	snakeSize, err = gameState.SnakeSize()
	if err != nil {
		if errors.Is(err, gameboard.ErrInvalidSnakeReference) {
			snkPosition = ""
			snkSize = "0"
		} else {
			return err
		}
	} else {
		snakePosition, err = gameState.SnakePosition()
		snkPosition = fmt.Sprintf("%v", snakePosition)
		snkSize = strconv.Itoa(snakeSize)
	}

	scoreViewLayout := []string{
		" GAME BOARD " + boardSize,
		"",
		"ROUND: " + strconv.Itoa(gameState.Round()),
		"",
		"CANDIES: " + strconv.Itoa(gameState.Score()),
		"SNAKE SIZE:" + snkSize,
		"POSITION:" + snkPosition,
		"",
		"TOP SCORE: " + strconv.Itoa(gameState.HighScore()),
	}

	return userInterface.SetViewLayout(scoreViewTitle, scoreViewLayout)
}

func updateErrorView(errMsg error, userInterface uimanager.UIManagerer, title string) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	const (
		errorViewWidth = maxX - rightPanel - 2
		nbErrorLines   = 4
	)

	str := errMsg.Error()
	chunks, err := chunks(str, errorViewWidth, nbErrorLines)

	if err != nil {
		return err
	}

	return userInterface.DisplayRedLayout(errorViewTitle, []string{
		"",
		"",
		"   Program Error",
		"",
		title,
		"     crashed",
		"",
		chunks[0],
		chunks[1],
		chunks[2],
		chunks[3],
	})
}
