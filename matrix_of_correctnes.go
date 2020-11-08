package permatchalgh

import (
	"fmt"

	gopair "github.com/Rakiiii/goPair"
)

type matrixOfCorrectnes struct {
	matrix [][]gopair.IntPair
}

func (m *matrixOfCorrectnes) init(n int) {
	matrix := make([][]gopair.IntPair, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]gopair.IntPair, n)
		for j := 0; j < n; j++ {
			matrix[i][j] = gopair.IntPair{First: i, Second: j}
		}
	}
	m.matrix = matrix
}

func (m *matrixOfCorrectnes) updateMatrixOfCorrectnes(i int, j int) {
	for x := 0; x < len(m.matrix); x++ {
		for y := 0; y < len(m.matrix[x]); y++ {
			if x > i {
				m.matrix[x][y].First--
			}
			if x > j {
				m.matrix[x][y].First--
			}
			if y > i {
				m.matrix[x][y].Second--
			}
			if y > j {
				m.matrix[x][y].Second--
			}
			if x == i || y == j || y == i || x == j {
				m.matrix[x][y] = gopair.IntPair{First: -1, Second: -1}
			}
		}
	}
}

func (m *matrixOfCorrectnes) getOriginalNumber(i int, j int) gopair.IntPair {
	for posVer, x := range m.matrix {
		for posHor, y := range x {
			if y.First == i && y.Second == j {
				return gopair.IntPair{First: posVer, Second: posHor}
			}
		}

	}
	return gopair.IntPair{First: -1, Second: -1}
}

func (m *matrixOfCorrectnes) getFixedNumber(i int, j int) gopair.IntPair {
	return m.matrix[i][j]
}

func (m *matrixOfCorrectnes) getFixedNumberFromPair(pair gopair.IntPair) gopair.IntPair {
	return m.getFixedNumber(pair.First, pair.Second)
}

func (m *matrixOfCorrectnes) print() {
	for _, elem := range m.matrix {
		for _, x := range elem {
			if x.First >= 0 && x.Second >= 0 {
				fmt.Print("(", x.First, ":", x.Second, ") ")
			} else {
				fmt.Print("() ")
			}
		}
		fmt.Println()
	}
}
