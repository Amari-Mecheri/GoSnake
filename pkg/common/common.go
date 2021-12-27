package common

import (
	"errors"
	"fmt"
	"runtime"
)

// ViewPosition holds coordinates of a view rectangle
type ViewPosition struct {
	X1 int
	Y1 int
	X2 int
	Y2 int
}

// Direction specifies deltas in X and Y
type Direction struct {
	DX int
	DY int
}

// Position defines coordinates
type Position struct {
	X int
	Y int
}

// Size is Width and Height
type Size struct {
	Width  int
	Height int
}

// Sprite holds a rune and its position
type Sprite struct {
	Value    rune
	Position Position
}

// GetCurrentFuncName returns the caller's function name
func GetCurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

// ErrorWrapper is called by functions in defer. It catches panics and wraps err, if it is not nil, with the function name
func ErrorWrapper(funcName string, err *error) {
	if aPanic := recover(); aPanic != nil {
		// aPanic may be of type error or of type string
		// We may switch aPanic.(type)
		switch val := aPanic.(type) {
		case string:
			*err = errors.New(val)
		case error:
			*err = val
		default: // or simply convert aPanic to string then create a new error...
			strErr := fmt.Sprint(aPanic)
			*err = errors.New(strErr)
		}
	}
	if *err != nil {
		*err = fmt.Errorf(funcName+": %w", *err)
	}
}
