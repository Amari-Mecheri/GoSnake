package candy

import (
	"gosnake/pkg/common"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCandy_CandyAlive(t *testing.T) {
	type fields struct {
		alive    bool
		position common.Position
	}
	tests := []struct {
		name      string
		fields    fields
		wantAlive bool
	}{
		{
			name: "TestCandyAlive",
			fields: fields{
				alive: true,
				position: common.Position{
					X: 5,
					Y: 5,
				},
			},
			wantAlive: true,
		},
		{
			name: "TestCandyDead",
			fields: fields{
				alive: false,
				position: common.Position{
					X: 5,
					Y: 5,
				},
			},
			wantAlive: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			candy := &candy{
				alive:    tt.fields.alive,
				position: tt.fields.position,
			}
			gotAlive := candy.Alive()
			require.Equal(t, tt.wantAlive, gotAlive)
		})
	}
}
func TestCandy_RemoveCandy(t *testing.T) {
	type fields struct {
		alive    bool
		position common.Position
	}
	tests := []struct {
		name      string
		fields    fields
		wantAlive bool
	}{
		{
			name: "TestCandyAlive",
			fields: fields{
				alive: true,
				position: common.Position{
					X: 5,
					Y: 5,
				},
			},
			wantAlive: false,
		},
		{
			name: "TestCandyDead",
			fields: fields{
				alive: false,
				position: common.Position{
					X: 5,
					Y: 5,
				},
			},
			wantAlive: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			candy := &candy{
				alive:    tt.fields.alive,
				position: tt.fields.position,
			}
			candy.Remove()
			gotAlive := candy.Alive()
			require.Equal(t, tt.wantAlive, gotAlive)
		})
	}
}
func TestCandy_GetPosition(t *testing.T) {
	type fields struct {
		alive    bool
		position common.Position
	}
	tests := []struct {
		name         string
		fields       fields
		wantPosition common.Position
	}{
		{
			name: "TestCandyAlive",
			fields: fields{
				alive: true,
				position: common.Position{
					X: 5,
					Y: 5,
				},
			},
			wantPosition: common.Position{
				X: 5,
				Y: 5,
			},
		},
		{
			name: "TestCandyDead",
			fields: fields{
				alive: false,
				position: common.Position{
					X: 15,
					Y: 7,
				},
			},
			wantPosition: common.Position{
				X: 15,
				Y: 7,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			candy := &candy{
				alive:    tt.fields.alive,
				position: tt.fields.position,
			}
			gotPosition := candy.Position()
			require.Equal(t, tt.wantPosition, gotPosition)
		})
	}
}

func TestCandy_Init(t *testing.T) {
	type fields struct {
		alive    bool
		position common.Position
	}
	type args struct {
		newPosition common.Position
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantPosition common.Position
		wantAlive    bool
	}{
		{
			name: "TestCandyCreateFromDead",
			fields: fields{
				alive: false,
				position: common.Position{
					X: 15,
					Y: 7,
				},
			},
			args: args{
				newPosition: common.Position{
					X: 20,
					Y: 20,
				},
			},
			wantPosition: common.Position{
				X: 20,
				Y: 20,
			},
			wantAlive: true,
		},
		{
			name: "TestCandyCreateFromAlive",
			fields: fields{
				alive: true,
				position: common.Position{
					X: 11,
					Y: 17,
				},
			},
			args: args{
				newPosition: common.Position{
					X: 10,
					Y: 9,
				},
			},
			wantPosition: common.Position{
				X: 10,
				Y: 9,
			},
			wantAlive: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			candy := &candy{
				alive:    tt.fields.alive,
				position: tt.fields.position,
			}
			candy.Init(tt.args.newPosition)
			gotPosition := candy.Position()
			gotAlive := candy.Alive()

			require.Equal(t, tt.wantAlive, gotAlive)
			require.Equal(t, tt.wantPosition, gotPosition)
		})
	}
}

func TestNew(t *testing.T) {
	var wantType *candy
	var got = New()
	require.IsType(t, wantType, got)
}
