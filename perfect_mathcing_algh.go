package permatchalgh

import (
	"errors"
	"fmt"
	"math/rand"

	graphlib "github.com/Rakiiii/goGraph"
	gopair "github.com/Rakiiii/goPair"
	gonum "gonum.org/v1/gonum/mat"
)

var NoPerfectMatching error = errors.New("Perfect mathcing does not exist in graph")

func frinksPerfectMathcingAlgth(graph graphlib.IGraph) ([]gopair.IntPair, error) {
	perfectMatching := make([]gopair.IntPair, 0)

	return perfectMatching, nil
}

func IsPerfectMatchingExist(graph graphlib.IGraph) bool {
	matrix := constractRandomMatrix(graph)
	return isPerfectMatchingExist(matrix)
}

func isPerfectMatchingExist(matrix *gonum.Dense) bool {
	return gonum.Det(matrix) != 0.0
}

func GetPerfectMatchingByRandomAlgorithm(graph graphlib.IGraph) ([]gopair.IntPair, error) {
	n := graph.AmountOfVertex()
	if n%2 != 0 {
		return nil, NoPerfectMatching
	}
	perfectMatching := make([]gopair.IntPair, n/2)

	Binversed := gonum.NewDense(n, n, make([]float64, n*n))
	B := constractRandomMatrix(graph)
	if !isPerfectMatchingExist(B) {
		return nil, NoPerfectMatching
	}

	Binversed.Inverse(B)
	for k := 0; k < n/2; k++ {
		i, j := getFirstNonZeroElemntPosition(Binversed)
		perfectMatching[k] = gopair.IntPair{First: i + k*2, Second: j + k*2}
		if Binversed.RawMatrix().Rows > 2 {
			Binversed = getSubMatrix(i, j, Binversed)
		}
	}

	return perfectMatching, nil
}

func constractRandomMatrix(graph graphlib.IGraph) *gonum.Dense {
	vertexAmount := graph.AmountOfVertex()
	rawMatrix := make([]float64, vertexAmount*vertexAmount)
	for i := 0; i < vertexAmount; i++ {
		edges := graph.GetEdges(i)
		for _, e := range edges {
			rawMatrix[i*vertexAmount+e] = float64(rand.Int())
		}
	}
	return gonum.NewDense(vertexAmount, vertexAmount, rawMatrix)
}

func getSubMatrix(x, y int, oldMat *gonum.Dense) *gonum.Dense {
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

func getFirstNonZeroElemntPosition(mat *gonum.Dense) (int, int) {
	for i := 0; i < mat.RawMatrix().Rows; i++ {
		for j := 0; j < mat.RawMatrix().Cols; j++ {
			if mat.At(i, j) != 0 && i != j {
				return i, j
			}
		}
	}
	return -1, -1
}

func prtintMatrix(matrix *gonum.Dense) {
	for i := 0; i < matrix.RawMatrix().Rows; i++ {
		for j := 0; j < matrix.RawMatrix().Cols; j++ {
			fmt.Print(matrix.At(i, j), " ")
		}
		fmt.Println()
	}
}
