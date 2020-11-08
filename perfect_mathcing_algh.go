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
var FixedVertexesNotInited error = errors.New("Fixed vertexes at RandomMathcerWithFixedVertexes not initialized")

func frinksPerfectMathcingAlgth(graph graphlib.IGraph) ([]gopair.IntPair, error) {
	perfectMatching := make([]gopair.IntPair, 0)

	return perfectMatching, nil
}

type RandomMatcher struct {
	rnd     *rand.Rand
	matcher IElementMatcher
}

type RandomMathcerWithFixedVertexes struct {
	RandomMatcher
	fixedVertexes []gopair.IntPair
	matcher       IElementMatcher
}

func NewRandomMathcerWithNilFixedVertexes() *RandomMathcerWithFixedVertexes {
	return NewRandomMathcerWithFixedVertexes(nil)
}

func NewRandomMathcerWithFixedVertexes(fixedVertexes []gopair.IntPair) *RandomMathcerWithFixedVertexes {
	return &RandomMathcerWithFixedVertexes{RandomMatcher: *NewRandomMatcher(), fixedVertexes: fixedVertexes}
}

func (c *RandomMathcerWithFixedVertexes) SetFixedVertexes(fixedVertexes []gopair.IntPair) {
	c.fixedVertexes = fixedVertexes
}

func (c *RandomMathcerWithFixedVertexes) GetPerfectMatching(graph graphlib.IGraph) ([]gopair.IntPair, error) {
	if c.fixedVertexes == nil {
		return nil, FixedVertexesNotInited
	} else {
		return c.getPerfectMatchingByRandomAlgorithmWithFixedVertexes(graph, c.fixedVertexes)
	}

}

func NewRandomMatcher() *RandomMatcher {
	return &RandomMatcher{rnd: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (c *RandomMatcher) GetPerfectMatching(graph graphlib.IGraph) ([]gopair.IntPair, error) {
	return c.getPerfectMatchingByRandomAlgorithm(graph)
}

func (c *RandomMatcher) IsPerfectMatchingExist(graph graphlib.IGraph) bool {
	matrix := c.constractRandomMatrix(graph)
	return c.isPerfectMatchingExist(matrix)
}

func (c *RandomMatcher) isPerfectMatchingExist(matrix *gonum.Dense) bool {
	return gonum.Det(matrix) != 0.0
}

func (c *RandomMatcher) getPerfectMatchingByRandomAlgorithm(graph graphlib.IGraph) ([]gopair.IntPair, error) {
	n := graph.AmountOfVertex()
	if n%2 != 0 {
		return nil, NoPerfectMatching
	}
	perfectMatching := make([]gopair.IntPair, 0)

	B := c.constractRandomMatrix(graph)
	if !c.isPerfectMatchingExist(B) {
		return nil, NoPerfectMatching
	}
	Binversed := c.constructTatasMatrix(B, graph)
	var matrix matrixOfCorrectnes
	matrix.init(n)
	for k := 1; k <= n/2; k++ {
		// matrix.print()
		// fmt.Println()
		// printMatrix(Binversed)
		// fmt.Println()
		i, j := c.getFirstNonZeroElemntPosition(Binversed, graph, &matrix)
		x := matrix.getOriginalNumber(i, j)
		// fmt.Println("(i:", i, ";j:", j, ") x(", x.First, ":", x.Second, ")")
		perfectMatching = append(perfectMatching, x)
		matrix.updateMatrixOfCorrectnes(x.First, x.Second)
		if B.RawMatrix().Rows > 4 {
			// fmt.Println("b size row:", B.RawMatrix().Rows, " cols:", B.RawMatrix().Cols)
			B = c.getSubMatrix(i, j, B)
			Binversed = c.constructTatasMatrix(B, graph)
		}
	}

	return perfectMatching, nil
}

func (c *RandomMatcher) getPerfectMatchingByRandomAlgorithmWithFixedVertexes(graph graphlib.IGraph, fixedVertexes []gopair.IntPair) ([]gopair.IntPair, error) {
	n := graph.AmountOfVertex()
	if n%2 != 0 {
		return nil, NoPerfectMatching
	}

	B := c.constractRandomMatrix(graph)
	if !c.isPerfectMatchingExist(B) {
		return nil, NoPerfectMatching
	}

	Binversed := c.constructTatasMatrix(B, graph)
	// printMatrix(Binversed)
	var matrix matrixOfCorrectnes
	matrix.init(n)
	perfectMatching := make([]gopair.IntPair, 0)
	for _, vertexPair := range fixedVertexes {
		fixedPair := matrix.getFixedNumberFromPair(vertexPair)
		if Binversed.At(fixedPair.First, fixedPair.Second) == 0 || Binversed.At(fixedPair.First, fixedPair.Second) == -0 || Binversed.At(fixedPair.Second, fixedPair.First) == 0 || Binversed.At(fixedPair.Second, fixedPair.First) == -0 {
			return nil, NoPerfectMatching
		} else {
			perfectMatching = append(perfectMatching, vertexPair)
			matrix.updateMatrixOfCorrectnes(vertexPair.First, vertexPair.Second)
			if B.RawMatrix().Rows > 4 {
				B = c.getSubMatrix(fixedPair.First, fixedPair.Second, B)
				Binversed = c.constructTatasMatrix(B, graph)
			}
		}
	}

	for k := 0; k < (n/2)-len(fixedVertexes); k++ {
		if !c.isPerfectMatchingExist(B) {
			return nil, NoPerfectMatching
		}
		// matrix.print()
		// fmt.Println()
		// printMatrix(Binversed)
		// fmt.Println()
		i, j := c.getFirstNonZeroElemntPosition(Binversed, graph, &matrix)
		x := matrix.getOriginalNumber(i, j)
		// fmt.Println("(i:", i, ";j:", j, ") x(", x.First, ":", x.Second, ")")
		perfectMatching = append(perfectMatching, x)
		matrix.updateMatrixOfCorrectnes(x.First, x.Second)
		if B.RawMatrix().Rows > 4 {
			B = c.getSubMatrix(i, j, B)
			Binversed = c.constructTatasMatrix(B, graph)
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
			// fmt.Println("(i:", i, ";j:", j, ") origin (", originalPosition.First, ":", originalPosition.Second, ") mat:", mat.At(i, j), ";", mat.At(j, i), ")")
			if originalPosition.First != -1 && originalPosition.Second != -2 {
				// fmt.Println("contains:", contains(graph.GetEdges(originalPosition.First), originalPosition.Second))
				if mat.At(i, j) != 0 && mat.At(i, j) != -0 && i != j && contains(graph.GetEdges(originalPosition.First), originalPosition.Second) {
					if mat.At(j, i) != 0 && mat.At(j, i) != -0 && i != j {
						return i, j
					}
				}
			}
		}
	}
	return -1, -1
}

func (c *RandomMatcher) constructTatasMatrix(d *gonum.Dense, graph graphlib.IGraph) *gonum.Dense {
	resultMatrix := make([]float64, 0)
	for i := 0; i < d.RawMatrix().Rows; i++ {
		for j := 0; j < d.RawMatrix().Cols; j++ {
			// && contains(graph.GetEdges(i), j)
			if i != j {
				resultMatrix = append(resultMatrix, c.constructSecondLevelDet(d, i, j))
			} else {
				resultMatrix = append(resultMatrix, 0)
			}
		}
	}
	return gonum.NewDense(d.RawMatrix().Rows, d.RawMatrix().Cols, resultMatrix)
}

func (c *RandomMatcher) constructSecondLevelDet(matrix *gonum.Dense, rowPos int, colPos int) float64 {
	return gonum.Det(c.getSubMatrix(rowPos, colPos, matrix))
}

func (c *RandomMatcher) constructTatasMatrixElement(matrix *gonum.Dense, inversed *gonum.Dense, rowPos int, colPos int) float64 {
	if inversed.At(rowPos, colPos) == 0 {
		return 0
	}
	elem := c.constructSecondLevelElement(matrix, rowPos, colPos)
	if elem == 0 {
		return elem
	} else {
		return c.matcher.matchElements(inversed.At(rowPos, colPos), elem)
	}
}

func (c *RandomMatcher) constructSecondLevelElement(matrix *gonum.Dense, rowPos int, colPos int) float64 {
	return 0
}

func printMatrix(matrix *gonum.Dense) {
	for i := 0; i < matrix.RawMatrix().Rows; i++ {
		fmt.Print("New Line :")
		for j := 0; j < matrix.RawMatrix().Cols; j++ {
			fmt.Print(matrix.At(i, j), " ")
		}
		fmt.Println()
	}
	fmt.Println()
}

func printMatrixNonZero(matrix *gonum.Dense) {
	for i := 0; i < matrix.RawMatrix().Rows; i++ {
		for j := 0; j < matrix.RawMatrix().Cols; j++ {
			if matrix.At(i, j) != 0 {
				fmt.Print(matrix.At(i, j), " ")
			}
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
