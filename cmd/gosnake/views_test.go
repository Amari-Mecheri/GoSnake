package main

import (
	"gosnake/pkg/uimanager"
	"testing"
)

func Test_updateErrorView(t *testing.T) {
	type args struct {
		errMsg        error
		userInterface uimanager.UIManagerer
		origin        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := updateErrorView(tt.args.errMsg, tt.args.userInterface, tt.args.origin); (err != nil) != tt.wantErr {
				t.Errorf("updateErrorView() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
