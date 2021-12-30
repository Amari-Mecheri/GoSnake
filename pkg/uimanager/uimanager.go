package uimanager

import (
	"gosnake/pkg/common"

	"github.com/jroimartin/gocui"
)

// UIManagerer is the interface for uiManager
type UIManagerer interface {
	OpenUIManager() (err error)
	Close()
	MainLoop() (err error)
	Update(viewName string, spriteList []common.Sprite) (err error)
	UpdateLn(viewName, msg string) (err error)
	SetView(viewName string, position common.ViewPosition) (err error)
	ClearView(viewName string) (err error)
	DisplayRedLayout(viewName string, layout []string) (err error)
	SetViewLayout(viewName string, layout []string) (err error)
	OnKeyPress(fn func(Key) error) (err error)
	Quit() (err error)
}

// Key is an alias
type Key gocui.Key

// Aliases to active keys
const (
	KeyCtrlC      Key = Key(gocui.KeyCtrlC)
	KeyArrowUp        = Key(gocui.KeyArrowUp)
	KeyArrowDown      = Key(gocui.KeyArrowDown)
	KeyArrowLeft      = Key(gocui.KeyArrowLeft)
	KeyArrowRight     = Key(gocui.KeyArrowRight)
	KeySpace          = Key(gocui.KeySpace)
	KeyEnter          = Key(gocui.KeyEnter)
)

// Aliases to gocui constants
var (
	ErrQuit    error = gocui.ErrQuit
	ActiveKeys       = []Key{
		KeyCtrlC,
		KeyArrowUp,
		KeyArrowDown,
		KeyArrowLeft,
		KeyArrowRight,
		KeySpace,
		KeyEnter,
	}
)

// uiManager encapsulates gocui library
type uiManager struct {
	gui *gocui.Gui
}

// New returns an instance of uiManager
func New() UIManagerer {
	return new(uiManager)
}

// OpenUIManager inits and gets a pointer to the user interface library
func (uim *uiManager) OpenUIManager() (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}

	uim.gui = gui

	return nil
}

// Close the UI library
func (uim *uiManager) Close() {
	uim.gui.Close()
}

// MainLoop updates the UI and manages events
func (uim *uiManager) MainLoop() (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	return uim.gui.MainLoop()
}

// Update a view with a list of sprites
func (uim *uiManager) Update(viewName string, spriteList []common.Sprite) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	var view *gocui.View

	if view, err = uim.gui.View(viewName); err != nil {
		return err
	}

	uim.gui.Update(
		func(g *gocui.Gui) error {
			return uim.displaySprites(view, spriteList)
		})

	return nil
}

// UpdateLn prints a line to the view
func (uim *uiManager) UpdateLn(viewName, msg string) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	var view *gocui.View

	if view, err = uim.gui.View(viewName); err != nil {
		return err
	}

	uim.gui.Update(
		func(g *gocui.Gui) error {
			return uim.writeLn(view, msg)
		})

	return nil
}

func (uim *uiManager) displaySprites(view *gocui.View, spriteList []common.Sprite) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	for i := range spriteList {
		if err := view.SetCursor(spriteList[i].Position.X,
			spriteList[i].Position.Y); err != nil {
			return err
		}

		view.EditWrite(spriteList[i].Value)
	}

	return nil
}

// SetView adds the view to the display manager
func (uim *uiManager) SetView(viewName string, position common.ViewPosition) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	if view, err := uim.gui.SetView(viewName, position.X1, position.Y1,
		position.X2, position.Y2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		view.Overwrite = true
		view.Autoscroll = false
		view.Wrap = true
	}

	return nil
}

// ClearView clears a view
func (uim *uiManager) ClearView(viewName string) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	view, err := uim.gui.View(viewName)
	if err != nil {
		return err
	}

	view.Clear()

	return nil
}

// DisplayRedLayout displays the view layout with red background
func (uim *uiManager) DisplayRedLayout(viewName string, layout []string) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	view, err := uim.gui.View(viewName)
	if err != nil {
		return err
	}

	view.Clear()
	view.BgColor = gocui.ColorRed

	for i := range layout {
		if err := view.SetCursor(0, i); err != nil {
			return err
		}

		if err := writeLn(view, layout[i]); err != nil {
			return err
		}
	}

	uim.gui.Update(func(g *gocui.Gui) error { return nil })

	return nil
}

// SetViewLayout defines the view layout
func (uim *uiManager) SetViewLayout(viewName string, layout []string) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	view, err := uim.gui.View(viewName)
	if err != nil {
		return err
	}

	view.Clear()

	for i := range layout {
		if err := view.SetCursor(0, i); err != nil {
			return err
		}

		if err := writeLn(view, layout[i]); err != nil {
			return err
		}
	}

	return nil
}

func writeLn(view *gocui.View, str string) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	for i := range str {
		view.EditWrite(rune(str[i]))
	}

	return nil
}

func (uim *uiManager) writeLn(view *gocui.View, str string) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	if err := view.SetCursor(0, 0); err != nil {
		return err
	}

	for i := range str {
		view.EditWrite(rune(str[i]))
	}

	return nil
}

// OnKeyPress attaches all active keys to an eventHandler
func (uim *uiManager) OnKeyPress(fn func(Key) error) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	for i := range ActiveKeys {
		if err := uim.setKeybinding(gocui.Key(ActiveKeys[i]), fn); err != nil {
			return err
		}
	}

	return nil
}

//
func (uim *uiManager) setKeybinding(key gocui.Key, fn func(Key) error) (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	return uim.gui.SetKeybinding("", key, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return fn(Key(key))
		})
}

// Quit stops the mainLoop
func (uim *uiManager) Quit() (err error) {
	defer common.ErrorWrapper(common.GetCurrentFuncName(), &err)

	return ErrQuit
}
