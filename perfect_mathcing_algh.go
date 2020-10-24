package permatchalgh

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	graphlib "github.com/Rakiiii/goGraph"
	gopair "github.com/Rakiiii/goPair"
	gonum "gonum.org/v1/gonum/mat"
)

var NoPerfectMatching error = errors.New("Perfect mathcing does not exist in graph")

func frinksPerfectMathcingAlgth(graph graphlib.IGraph) ([]gopair.IntPair, error) {
	perfectMatching := make([]gopair.IntPair, 0)

	return perfectMatching, nil
}

type RandomMatcher struct {
	rnd *rand.Rand
}

func NewRandomMatcher() *RandomMatcher {
	return &RandomMatcher{rnd: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (c *RandomMatcher) GetPerfectMatching(graph graphlib.IGraph) ([]gopair.IntPair, error) {
	return c.GetPerfectMatchingByRandomAlgorithm(graph)
}

func (c *RandomMatcher) GetPerfectMatchingWithFixedVertexes(graph graphlib.IGraph, fixedVertexes []gopair.IntPair) ([]gopair.IntPair, error) {
	return c.GetPerfectMatchingByRandomAlgorithmWithFixedVertexes(graph, fixedVertexes)
}

func (c *RandomMatcher) IsPerfectMatchingExist(graph graphlib.IGraph) bool {
	matrix := c.constractRandomMatrix(graph)
	return c.isPerfectMatchingExist(matrix)
}

func (c *RandomMatcher) isPerfectMatchingExist(matrix *gonum.Dense) bool {
	return gonum.Det(matrix) != 0.0
}

func (c *RandomMatcher) GetPerfectMatchingByRandomAlgorithm(graph graphlib.IGraph) ([]gopair.IntPair, error) {
	n := graph.AmountOfVertex()
	if n%2 != 0 {
		return nil, NoPerfectMatching
	}
	perfectMatching := make([]gopair.IntPair, 0)

	Binversed := gonum.NewDense(n, n, make([]float64, n*n))
	B := c.constractRandomMatrix(graph)
	if !c.isPerfectMatchingExist(B) {
		return nil, NoPerfectMatching
	}

	Binversed.Inverse(B)
	var matrix matrixOfCorrectnes
	matrix.init(n)
	//printMatrix(B)

	for k := 0; k < n/2; k++ {
		//printMatrix(Binversed)
		//matrix.print()
		i, j := c.getFirstNonZeroElemntPosition(Binversed, graph, &matrix)
		x := matrix.getOriginalNumber(i, j)
		//fmt.Println()
		//fmt.Println("(", i, ":", j, ") original position", "(", x.First, ":", x.Second, ") ")
		//fmt.Println()
		perfectMatching = append(perfectMatching, x)
		matrix.updateMatrixOfCorrectnes(x.First, x.Second)
		if Binversed.RawMatrix().Rows > 2 {
			Binversed = c.getSubMatrix(i, j, Binversed)
		}
	}

	return perfectMatching, nil
}

func (c *RandomMatcher) GetPerfectMatchingByRandomAlgorithmWithFixedVertexes(graph graphlib.IGraph, fixedVertexes []gopair.IntPair) ([]gopair.IntPair, error) {
	n := graph.AmountOfVertex()
	if n%2 != 0 {
		return nil, NoPerfectMatching
	}

	Binversed := gonum.NewDense(n, n, make([]float64, n*n))
	B := c.constractRandomMatrix(graph)
	if !c.isPerfectMatchingExist(B) {
		return nil, NoPerfectMatching
	}

	Binversed.Inverse(B)
	var matrix matrixOfCorrectnes
	matrix.init(n)
	perfectMatching := make([]gopair.IntPair, 0)
	for _, vertexPair := range fixedVertexes {
		if Binversed.At(vertexPair.First, vertexPair.Second) == 0 || Binversed.At(vertexPair.First, vertexPair.Second) == -0 || Binversed.At(vertexPair.Second, vertexPair.First) == 0 || Binversed.At(vertexPair.Second, vertexPair.First) == -0 {
			return nil, NoPerfectMatching
		} else {
			perfectMatching = append(perfectMatching, vertexPair)
			matrix.updateMatrixOfCorrectnes(vertexPair.First, vertexPair.Second)
		}
	}

	Binversed = c.getSubMatrixFromSlice(fixedVertexes, Binversed)
	for k := 0; k < (n/2)-len(fixedVertexes); k++ {
		//printMatrix(Binversed)
		//matrix.print()
		i, j := c.getFirstNonZeroElemntPosition(Binversed, graph, &matrix)
		x := matrix.getOriginalNumber(i, j)
		//fmt.Println()
		//fmt.Println("(", i, ":", j, ") original position", "(", x.First, ":", x.Second, ") ")
		//fmt.Println()
		perfectMatching = append(perfectMatching, x)
		matrix.updateMatrixOfCorrectnes(x.First, x.Second)
		if Binversed.RawMatrix().Rows > 2 {
			Binversed = c.getSubMatrix(i, j, Binversed)
		}
	}

	return perfectMatching, nil
}

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

func (m *matrixOfCorrectnes) print() {
	for _, elem := range m.matrix {
		for _, x := range elem {
			fmt.Print("(", x.First, ":", x.Second, ") ")
		}
		fmt.Println()
	}
}

func (c *RandomMatcher) constractRandomMatrix(graph graphlib.IGraph) *gonum.Dense {
	vertexAmount := graph.AmountOfVertex()
	rawMatrix := make([]float64, vertexAmount*vertexAmount)
	for i := 0; i < vertexAmount; i++ {
		edges := graph.GetEdges(i)
		for _, e := range edges {
			value := float64(c.rnd.Int31n(1000))
			if rawMatrix[i*vertexAmount+e] == 0 {
				rawMatrix[i*vertexAmount+e] = value
			}
			if rawMatrix[e*vertexAmount+i] == 0 {
				rawMatrix[e*vertexAmount+i] = value
			}
		}
	}
	return gonum.NewDense(vertexAmount, vertexAmount, rawMatrix)
}

func (c *RandomMatcher) getSubMatrix(x, y int, oldMat *gonum.Dense) *gonum.Dense {
	newMatrixSize := oldMat.RawMatrix().Rows - 2
	newRawMatrix := make([]float64, newMatrixSize*newMatrixSize)
	position := 0
	for i := 0; i < oldMat.RawMatrix().Rows; i++ {
		for j := 0; j < oldMat.RawMatrix().Cols; j++ {
			if i != x && i != y && j != x && j != y {
				newRawMatrix[position] = oldMat.At(i, j)
				position++
			}
		}
	}
	return gonum.NewDense(newMatrixSize, newMatrixSize, newRawMatrix)
}

func (c *RandomMatcher) getSubMatrixFromSlice(vertexPairs []gopair.IntPair, oldMat *gonum.Dense) *gonum.Dense {
	newMatrixSize := oldMat.RawMatrix().Rows - 2*len(vertexPairs)
	newRawMatrix := make([]float64, newMatrixSize*newMatrixSize)
	position := 0
	for i := 0; i < oldMat.RawMatrix().Rows; i++ {
		for j := 0; j < oldMat.RawMatrix().Cols; j++ {
			if !containsInPair(vertexPairs, i) && !containsInPair(vertexPairs, j) {
				newRawMatrix[position] = oldMat.At(i, j)
				position++
			}
		}
	}
	return gonum.NewDense(newMatrixSize, newMatrixSize, newRawMatrix)
}

func (c *RandomMatcher) getFirstNonZeroElemntPosition(mat *gonum.Dense, graph graphlib.IGraph, matrixOfCorrectnes *matrixOfCorrectnes) (int, int) {
	for i := 0; i < mat.RawMatrix().Rows; i++ {
		for j := 0; j < mat.RawMatrix().Cols; j++ {
			originalPosition := matrixOfCorrectnes.getOriginalNumber(i, j)
			if mat.At(i, j) != 0 && mat.At(i, j) != -0 && i != j && contains(graph.GetEdges(originalPosition.First), originalPosition.Second) {
				if mat.At(j, i) != 0 && mat.At(j, i) != -0 && i != j {
					//fmt.Println("value is ", mat.At(i, j))
					return i, j
				}
			}
		}
	}
	return -1, -1
}

func printMatrix(matrix *gonum.Dense) {
	for i := 0; i < matrix.RawMatrix().Rows; i++ {
		for j := 0; j < matrix.RawMatrix().Cols; j++ {
			fmt.Print(matrix.At(i, j), " ")
		}
		fmt.Println()
	}
	fmt.Println()
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsInPair(s []gopair.IntPair, e int) bool {
	for _, a := range s {
		if a.First == e || a.Second == e {
			return true
		}
	}
	return false
}
