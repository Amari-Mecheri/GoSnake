package common

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCurrentFuncName(t *testing.T) {
	tests := []struct {
		name         string
		wantFuncName string
	}{
		{
			name:         "TestGetCurrentFuncName",
			wantFuncName: "common.TestGetCurrentFuncName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFuncName := GetCurrentFuncName()
			require.Contains(t, gotFuncName, tt.wantFuncName)
		})
	}
}

func TestErrorWrapper(t *testing.T) {
	type args struct {
		funcName string
		err      *error
	}
	tests := []struct {
		name         string
		args         args
		theFunc      func()
		wantErrType  string // We use a string to check the content of the panic message;
		wantFuncName string
	}{
		{
			name: "TestPanicDivideByZero",
			theFunc: func() {
				var a, b int
				a = a / b //nolint
			},
			args: args{
				funcName: "TestPanicDivideByZero",
				err:      new(error),
			},
			wantFuncName: "TestPanicDivideByZero",
			wantErrType:  "runtime error",
		},
		{
			name: "TestPanicRandomZero",
			theFunc: func() {
				aBig, err := rand.Int(rand.Reader, big.NewInt(0))
				if err != nil {
					t.Fail()
				}
				a := aBig.Int64()
				a++ //nolint
			},
			args: args{
				funcName: "TestPanicRandomZero",
				err:      new(error),
			},
			wantFuncName: "TestPanicRandomZero",
			wantErrType:  "argument to Int is <= 0",
		},
		{
			name: "TestPanicANumber",
			theFunc: func() {
				panic(5)
			},
			args: args{
				funcName: "TestPanicANumber",
				err:      new(error),
			},
			wantFuncName: "TestPanicANumber",
			wantErrType:  "5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			func() {
				defer ErrorWrapper(tt.args.funcName, tt.args.err)
				tt.theFunc()
			}()
			gotErrMsg := fmt.Sprint(*tt.args.err)
			require.Contains(t, gotErrMsg, tt.wantErrType)
			require.Contains(t, gotErrMsg, tt.wantFuncName)
		})
	}
}
