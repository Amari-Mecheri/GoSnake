A version of the Snake Game in Go

Architecture of GoSnake:

By choice, the package "main" is made of functions, variables are local to the main function and are sent as arguments where needed.

uimanager, gamestate, gameboard, snake and candy packages export Interfaces, Methods implementing the corresponding interfaces and New() functions
They also define a struct with the needed properties. These structs are not exported.
The New() function of each package is used to get an instance of the corresponding object.
The New function also performs instantiation of the childs by calling the respective New() functions.
This allows the implementation of the mechanics of "classes" with private members, instantiated via a mandatory constructor/fabric (New function)

inheritance (is a) is implemented by using an anonymous field of type child
composition (has) is implemented by using a named field of type child

=> In the case of inheritance, exported methods of the child can be used and exported by the parent as if they were methods of the parent.
=> In the case of composition methods of the child are only accessible via the corresponding named field.

main function
| (var) ->	gamestate
			| (is a) -> gameboard
					| (has) -> snake
					| (has) -> candy
	| (var) ->	uimanager
				| (has) -> gocui

uimanager encapsulates functions from Gocui https:github.com/jroimartin/gocui

The main package controls the gamestate, creates the views layouts
and updates them to reflect the state via uimanager

The game logic is held by gamestate. The Play() methods plays a round.

The rounds are called in a loop controlled by tickers at intervals
The errors from the routines are channeled back to the main function

In order to manage the keys pressed, the bindings are all affected to the same eventHandler which selects the appropriate action
Since the eventHandler will only receive the key pressed as parameter, a closure is used to allow access to the main parameters

Most functions and methods start with a defer common.ErrorWrapper
Which handles unexpected panics and wraps any error with the function/method name
This allows identification of the function where the error occurred
as well as the chain of the functions called (since each parent function calls common.ErrorWrapper too).
Where needed custom errors are created and later identified with errors Is().
In the case of panics the content is checked, the panic can be an error or a string (cf common.ErrorWrapper test)
The errors and unexpected panics can be controlled and logged if needed
For example, if an error/panic occurs in the routines gameEngine or gameOverAnim, a red alert is displayed without terminating the program.

This version is the result of many iterations. The first version which was aesthetically almost identical to the final version, took a few hours to develop.
This test task was actually a learning task. As this was my first program in Go, the initial objective was to have a game up and running.
Next, improving code quality, refactoring, debugging, and unit testing; took most of the time.  Golang is not a complicated language but understanding some aspects and almost every other step took its share of time. All in all, I spent about 7 working days to complete the version published here.

Make commands:

make mock
make install-linters
make format
make lint
make test
make build
make run
make clean

