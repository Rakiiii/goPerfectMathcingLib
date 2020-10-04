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
	perfectMatching := make([]gopair.IntPair, n/2)

	Binversed := gonum.NewDense(n, n, make([]float64, n*n))
	B := c.constractRandomMatrix(graph)
	if !c.isPerfectMatchingExist(B) {
		return nil, NoPerfectMatching
	}

	Binversed.Inverse(B)
	for k := 0; k < n/2; k++ {
		i, j := c.getFirstNonZeroElemntPosition(Binversed)
		perfectMatching[k] = gopair.IntPair{First: i + k*2, Second: j + k*2}
		if Binversed.RawMatrix().Rows > 2 {
			Binversed = c.getSubMatrix(i, j, Binversed)
		}
	}

	return perfectMatching, nil
}

func (c *RandomMatcher) constractRandomMatrix(graph graphlib.IGraph) *gonum.Dense {
	vertexAmount := graph.AmountOfVertex()
	rawMatrix := make([]float64, vertexAmount*vertexAmount)
	for i := 0; i < vertexAmount; i++ {
		edges := graph.GetEdges(i)
		for _, e := range edges {
			rawMatrix[i*vertexAmount+e] = float64(c.rnd.Int())
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

func (c *RandomMatcher) getFirstNonZeroElemntPosition(mat *gonum.Dense) (int, int) {
	for i := 0; i < mat.RawMatrix().Rows; i++ {
		for j := 0; j < mat.RawMatrix().Cols; j++ {
			if mat.At(i, j) != 0 && i != j {
				return i, j
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
}
