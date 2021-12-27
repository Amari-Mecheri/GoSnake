package snake

import (
	"gosnake/pkg/common"
	"gosnake/testdata"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnake_Size(t *testing.T) {
	type fields struct {
		body      []common.Position
		direction common.Direction
	}
	tests := []struct {
		name     string
		fields   fields
		wantSize int
		wantErr  bool
	}{
		{
			name: "TestEmptyBody",
			fields: fields{
				body:      nil,
				direction: testdata.Direction1_0,
			},
			wantSize: 0,
			wantErr:  false,
		},
		{
			name: "TestBodyOne",
			fields: fields{
				body: []common.Position{
					{
						X: 0,
						Y: 0,
					},
				},
				direction: testdata.Direction1_0,
			},
			wantSize: 1,
			wantErr:  false,
		},
		{
			name: "TestBodyTwo",
			fields: fields{
				body: []common.Position{
					{
						X: 0,
						Y: 0,
					},
					{
						X: 1,
						Y: 0,
					},
				},
				direction: testdata.Direction1_0,
			},
			wantSize: 2,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aSnake := &snake{
				body:      tt.fields.body,
				direction: tt.fields.direction,
			}
			gotSize, err := aSnake.Size()
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr)
			require.Equal(t, tt.wantSize, gotSize)
		})
	}
}

func TestSnake_SetDirection(t *testing.T) {
	type fields struct {
		body      []common.Position
		direction common.Direction
	}
	type args struct {
		direction common.Direction
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantDirection common.Direction
	}{
		{
			name: "TestDirection",
			fields: fields{
				body:      nil,
				direction: testdata.Direction1_0,
			},
			args: args{
				direction: common.Direction{
					DX: 0,
					DY: 1,
				},
			},
			wantDirection: common.Direction{
				DX: 0,
				DY: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aSnake := &snake{
				body:      tt.fields.body,
				direction: tt.fields.direction,
			}
			aSnake.SetDirection(tt.args.direction)
			gotDirection := aSnake.direction
			require.Equal(t, tt.wantDirection, gotDirection)
		})
	}
}

func TestSnake_Position(t *testing.T) {
	type fields struct {
		body      []common.Position
		direction common.Direction
	}
	tests := []struct {
		name         string
		fields       fields
		wantPosition common.Position
		wantErr      bool
	}{
		{
			name: "TestEmptyBody",
			fields: fields{
				body:      nil,
				direction: testdata.Direction1_0,
			},
			wantPosition: testdata.Position0_0,
			wantErr:      true,
		},
		{
			name: "TestBodyOne",
			fields: fields{
				body: []common.Position{
					{
						X: 1,
						Y: 2,
					},
				},
				direction: testdata.Direction1_0,
			},
			wantPosition: testdata.Position1_2,
			wantErr:      false,
		},
		{
			name: "TestBodyTwo",
			fields: fields{
				body: []common.Position{
					{
						X: 1,
						Y: 3,
					},
					{
						X: 2,
						Y: 3,
					},
				},
				direction: testdata.Direction1_0,
			},
			wantPosition: common.Position{
				X: 2,
				Y: 3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aSnake := &snake{
				body:      tt.fields.body,
				direction: tt.fields.direction,
			}
			gotPosition, err := aSnake.Position()
			gotErr := err != nil
			require.Equal(t, tt.wantErr, gotErr)
			require.Equal(t, tt.wantPosition, gotPosition)
		})
	}
}

func TestSnake_GetNextMove(t *testing.T) {
	type fields struct {
		body      []common.Position
		direction common.Direction
	}
	tests := []struct {
		name             string
		fields           fields
		wantNextPosition common.Position
		wantErr          bool
	}{
		{
			name: "TestEmptyBody",
			fields: fields{
				body:      nil,
				direction: testdata.Direction1_0,
			},
			wantNextPosition: testdata.Position0_0,
			wantErr:          true,
		},
		{
			name: "TestBodyOne",
			fields: fields{
				body: []common.Position{
					{
						X: 1,
						Y: 2,
					},
				},
				direction: testdata.Direction1_0,
			},
			wantNextPosition: testdata.Position2_2,
			wantErr:          false,
		},
		{
			name: "TestBodyTwo",
			fields: fields{
				body: []common.Position{
					{
						X: 1,
						Y: 3,
					},
					{
						X: 2,
						Y: 3,
					},
				},
				direction: testdata.Direction1_0,
			},
			wantNextPosition: common.Position{
				X: 3,
				Y: 3,
			},
			wantErr: false,
		},
		{
			name: "TestBodyTwoDown",
			fields: fields{
				body: []common.Position{
					{
						X: 1,
						Y: 3,
					},
					{
						X: 2,
						Y: 3,
					},
				},
				direction: common.Direction{
					DX: 0,
					DY: -1,
				},
			},
			wantNextPosition: testdata.Position2_2,
			wantErr:          false,
		},
		{
			name: "TestBodyTwoDownLeft",
			fields: fields{
				body: []common.Position{
					{
						X: 1,
						Y: 1,
					},
					{
						X: 0,
						Y: 0,
					},
				},
				direction: testdata.DirectionMinus1_Minus1,
			},
			wantNextPosition: testdata.PositionMinus1_Minus1,
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aSnake := &snake{
				body:      tt.fields.body,
				direction: tt.fields.direction,
			}
			gotNextPosition, err := aSnake.NextMove()
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr)
			require.Equal(t, tt.wantNextPosition, gotNextPosition)
		})
	}
}

func TestSnake_MoveTo(t *testing.T) {
	type fields struct {
		body      []common.Position
		direction common.Direction
	}
	type args struct {
		newPosition common.Position
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantTheTail  common.Position
		wantPosition common.Position
		wantErr      bool
	}{
		{
			name: "TestEmptyBody",
			fields: fields{
				body:      nil,
				direction: testdata.DirectionMinus1_Minus1,
			},
			args: args{
				newPosition: testdata.Position0_0,
			},
			wantTheTail:  testdata.Position0_0,
			wantPosition: testdata.Position0_0,
			wantErr:      true,
		},
		{
			name: "TestBodyOneDown",
			fields: fields{
				body: []common.Position{
					{
						X: 0,
						Y: 0,
					},
				},
				direction: common.Direction{
					DX: 0,
					DY: -1,
				},
			},
			args: args{
				newPosition: testdata.Position0_Minus1,
			},
			wantPosition: testdata.Position0_Minus1,
			wantTheTail:  testdata.Position0_0,
			wantErr:      false,
		},
		{
			name: "TestBodyTwoLeft",
			fields: fields{
				body: []common.Position{
					{
						X: 2,
						Y: 3,
					},
					{
						X: 1,
						Y: 3,
					},
				},
				direction: testdata.DirectionMinus1_0,
			},
			args: args{
				newPosition: common.Position{
					X: 0,
					Y: 3,
				},
			},
			wantTheTail: common.Position{
				X: 2,
				Y: 3,
			},
			wantPosition: common.Position{
				X: 0,
				Y: 3,
			},
			wantErr: false,
		},
		{
			name: "TestBodyTwoDownLeft",
			fields: fields{
				body: []common.Position{
					{
						X: 1,
						Y: 1,
					},
					{
						X: 0,
						Y: 0,
					},
				},
				direction: testdata.DirectionMinus1_Minus1,
			},
			args: args{
				newPosition: testdata.PositionMinus1_Minus1,
			},
			wantTheTail: common.Position{
				X: 1,
				Y: 1,
			},
			wantPosition: testdata.PositionMinus1_Minus1,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aSnake := &snake{
				body:      tt.fields.body,
				direction: tt.fields.direction,
			}
			gotTheTail, err := aSnake.MoveTo(tt.args.newPosition)
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr)
			require.Equal(t, tt.wantTheTail, gotTheTail)
			size, err := aSnake.Size()
			if err == nil {
				if size > 0 {
					position, err := aSnake.Position()
					if err == nil {
						require.Equal(t, tt.wantPosition, position)
					}
				}
			}
		})
	}
}

func TestSnake_GrowTo(t *testing.T) {
	type fields struct {
		body      []common.Position
		direction common.Direction
	}
	type args struct {
		newPosition common.Position
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantPosition common.Position
		wantSize     int
		wantErr      bool
	}{
		{
			name: "TestEmptyBody5,5",
			fields: fields{
				body:      nil,
				direction: testdata.DirectionMinus1_0,
			},
			args: args{
				newPosition: common.Position{
					X: 5,
					Y: 5,
				},
			},
			wantPosition: common.Position{
				X: 5,
				Y: 5,
			},
			wantSize: 1,
			wantErr:  false,
		},
		{
			name: "TestBodyOne12,13",
			fields: fields{
				body: []common.Position{
					{
						X: 5,
						Y: 7,
					},
				},
				direction: common.Direction{
					DX: 0,
					DY: -1,
				},
			},
			args: args{
				newPosition: common.Position{
					X: 12,
					Y: 13,
				},
			},
			wantPosition: common.Position{
				X: 12,
				Y: 13,
			},
			wantSize: 2,
			wantErr:  false,
		},
		{
			name: "TestBodyTwo1,2",
			fields: fields{
				body: []common.Position{
					{
						X: 2,
						Y: 3,
					},
					{
						X: 1,
						Y: 3,
					},
				},
				direction: testdata.DirectionMinus1_0,
			},
			args: args{
				newPosition: common.Position{
					X: 1,
					Y: 2,
				},
			},
			wantPosition: testdata.Position1_2,
			wantSize:     3,
			wantErr:      false,
		},
		{
			name: "TestBodyTwo-1,-1",
			fields: fields{
				body: []common.Position{
					{
						X: 1,
						Y: 1,
					},
					{
						X: 0,
						Y: 0,
					},
				},
				direction: testdata.DirectionMinus1_Minus1,
			},
			args: args{
				newPosition: testdata.PositionMinus1_Minus1,
			},
			wantPosition: testdata.PositionMinus1_Minus1,
			wantSize:     3,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aSnake := &snake{
				body:      tt.fields.body,
				direction: tt.fields.direction,
			}
			err := aSnake.GrowTo(tt.args.newPosition)
			gotErr := (err != nil)
			require.Equal(t, tt.wantErr, gotErr)
			size, err := aSnake.Size()
			if err == nil {
				require.Equal(t, tt.wantSize, size)
			}
			position, err := aSnake.Position()
			if err == nil {
				require.Equal(t, tt.wantPosition, position)
			}
		})
	}
}

func TestNew(t *testing.T) {
	var wantType *snake
	var got = New()
	require.IsType(t, wantType, got)
}
