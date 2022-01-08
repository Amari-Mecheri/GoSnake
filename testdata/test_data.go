package testdata

import (
	"gosnake/pkg/common"
)

func Duplicate(source [][]rune) (duplicate [][]rune) {
	duplicate = make([][]rune, len(source))
	for i := range source {
		duplicate[i] = make([]rune, len(source[i]))
		copy(duplicate[i], source[i])
	}

	return duplicate
}

var (
	Direction0_0 = common.Direction{
		DX: 0,
		DY: 0,
	}
	Direction1_0 = common.Direction{
		DX: 1,
		DY: 0,
	}
	DirectionMinus1_Minus1 = common.Direction{
		DX: -1,
		DY: -1,
	}
	DirectionMinus1_0 = common.Direction{
		DX: -1,
		DY: 0,
	}
	Position0_0 = common.Position{
		X: 0,
		Y: 0,
	}
	Position0_1 = common.Position{
		X: 0,
		Y: 1,
	}
	Position1_1 = common.Position{
		X: 1,
		Y: 1,
	}
	Position0_Minus1 = common.Position{
		X: 0,
		Y: -1,
	}
	Position1_2 = common.Position{
		X: 1,
		Y: 2,
	}
	Position2_2 = common.Position{
		X: 2,
		Y: 2,
	}
	Position4_3 = common.Position{
		X: 4,
		Y: 3,
	}
	Position0_3 = common.Position{
		X: 0,
		Y: 3,
	}
	Position3_4 = common.Position{
		X: 3,
		Y: 4,
	}
	Position3_0 = common.Position{
		X: 3,
		Y: 0,
	}
	PositionMinus1_Minus1 = common.Position{
		X: -1,
		Y: -1,
	}
	Size0_0 = common.Size{
		Width:  0,
		Height: 0,
	}
	Size3_3 = common.Size{
		Width:  3,
		Height: 3,
	}
	Size4_4 = common.Size{
		Width:  4,
		Height: 4,
	}
	SizeMinus1_Minus1 = common.Size{
		Width:  -1,
		Height: -1,
	}
	Board3_3_OneFreeSpotPos1_1 = [][]rune{
		{
			'a', 'b', 'c',
		},
		{
			'a', ' ', 'c',
		},
		{
			'a', 'b', 'c',
		},
	}
	Board3_3 = [][]rune{
		{
			'a', 'b', 'c',
		},
		{
			'a', 'b', 'c',
		},
		{
			'a', 'b', 'c',
		},
	}
	Board3_3Candy1_1 = [][]rune{
		{
			'a', 'b', 'c',
		},
		{
			'a', '*', 'c',
		},
		{
			'a', 'b', 'c',
		},
	}
	Board3_3Snake_1 = [][]rune{
		{
			'a', 'S', 'c',
		},
		{
			'a', 'S', 'd',
		},
		{
			'a', 'S', 'c',
		},
	}
	Board3_3WrongAllocation = [][]rune{
		{
			'a', 'b', 'c', 'd',
		},
		{
			'a', 'b', 'c',
		},
		{
			'a', 'b', 'c',
		},
	}
)
