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
var notCheckedVertex int = 0
var connectedVertex int = 1
var checkedVertex int = 2

func frinksPerfectMathcingAlgth(graph graphlib.IGraph) ([]gopair.IntPair, error) {
	perfectMatching := make([]gopair.IntPair, 0)

	return perfectMatching, nil
}

type DetChecker struct{}

func (d *DetChecker) isDetZero(matrix gonum.Matrix) bool {
	return gonum.Det(matrix) != 0.0
}

func NewDetChecker() *DetChecker {
	return &DetChecker{}
}

type CondChecker struct {
	cond float64
}

func NewCondChecker(cond float64) *CondChecker {
	return &CondChecker{cond: cond}
}

func (d *CondChecker) isDetZero(matrix gonum.Matrix) bool {
	return gonum.Cond(matrix, d.cond) < 100
}

type RandomMatcher struct {
	Rnd        *rand.Rand
	matcher    IElementMatcher
	detChecker IMatrixDeterminantChecker
}

type RandomMathcerWithFixedVertexes struct {
	RandomMatcher
	FixedVertexes []gopair.IntPair
	matcher       IElementMatcher
}

func NewRandomMathcerWithNilFixedVertexes() *RandomMathcerWithFixedVertexes {
	return NewRandomMathcerWithFixedVertexes(nil)
}

func NewRandomMathcerWithFixedVertexes(fixedVertexes []gopair.IntPair) *RandomMathcerWithFixedVertexes {
	return &RandomMathcerWithFixedVertexes{RandomMatcher: *NewRandomMatcher(), FixedVertexes: fixedVertexes}
}

func (c *RandomMathcerWithFixedVertexes) SetFixedVertexes(fixedVertexes []gopair.IntPair) {
	c.FixedVertexes = fixedVertexes
}

func (c *RandomMathcerWithFixedVertexes) GetPerfectMatching(graph graphlib.IGraph) ([]gopair.IntPair, error) {
	if c.FixedVertexes == nil {
		return nil, FixedVertexesNotInited
	} else {
		return c.getPerfectMatchingByRandomAlgorithmWithFixedVertexes(graph, c.FixedVertexes)
	}

}

func (c *RandomMathcerWithFixedVertexes) IsPerfectMatchingExist(graph graphlib.IGraph) bool {
	if c.FixedVertexes == nil {
		return false
	} else {
		return c.isPerfectMatchingExistWithFixedVertexes(graph, c.FixedVertexes)
	}
}

func NewRandomMatcher() *RandomMatcher {
	time := time.Now().UnixNano()
	fmt.Println("time:", time)
	return &RandomMatcher{Rnd: rand.New(rand.NewSource(time)), detChecker: NewDetChecker()}
}

func (c *RandomMatcher) GetPerfectMatching(graph graphlib.IGraph) ([]gopair.IntPair, error) {
	return c.getPerfectMatchingByRandomAlgorithm(graph)
}

func (c *RandomMatcher) SetDetChecker(checker IMatrixDeterminantChecker) {
	c.detChecker = checker
}

func (c *RandomMatcher) IsPerfectMatchingExist(graph graphlib.IGraph) bool {
	matrix := c.constractRandomMatrix(graph)
	return c.isPerfectMatchingExist(matrix)
}

func (c *RandomMatcher) isPerfectMatchingExistWithFixedVertexes(graph graphlib.IGraph, fixedVertexes []gopair.IntPair) bool {
	n := graph.AmountOfVertex()
	if n%2 != 0 {
		return false
	}

	B := c.constractRandomMatrix(graph)
	if !c.isPerfectMatchingExist(B) {
		return false
	}
	Binversed := c.constructTatasMatrix(B, graph)
	var matrix matrixOfCorrectnes
	matrix.init(n)
	perfectMatching := make([]gopair.IntPair, 0)
	for _, vertexPair := range fixedVertexes {
		fixedPair := matrix.getFixedNumberFromPair(vertexPair)
		if Binversed.At(fixedPair.First, fixedPair.Second) == 0 || Binversed.At(fixedPair.First, fixedPair.Second) == -0 || Binversed.At(fixedPair.Second, fixedPair.First) == 0 || Binversed.At(fixedPair.Second, fixedPair.First) == -0 {
			return false
		} else {
			perfectMatching = append(perfectMatching, vertexPair)
			matrix.updateMatrixOfCorrectnes(vertexPair.First, vertexPair.Second)
			if B.RawMatrix().Rows > 4 {
				B = c.getSubMatrix(fixedPair.First, fixedPair.Second, B)
				Binversed = c.constructTatasMatrix(B, graph)
			}
		}
	}
	return c.isPerfectMatchingExist(B)
}

func (c *RandomMatcher) isPerfectMatchingExist(matrix *gonum.Dense) bool {
	return c.detChecker.isDetZero(matrix)
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
		i, j := c.getFirstNonZeroElemntPosition(Binversed, graph, &matrix)
		x := matrix.getOriginalNumber(i, j)
		perfectMatching = append(perfectMatching, x)
		matrix.updateMatrixOfCorrectnes(x.First, x.Second)
		if B.RawMatrix().Rows > 4 {
			B = c.getSubMatrix(i, j, B)
			Binversed = c.constructTatasMatrix(B, graph)
		} else {
			if Binversed.RawMatrix().Rows > 2 {
				Binversed = c.getSubMatrix(i, j, Binversed)
			}
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
		if !c.isPerfectMatchingExist(B) || !isGraphSingleConnected(B) || !isMatrixContainsNonZeroElements(Binversed) {
			return nil, NoPerfectMatching
		}
		i, j := c.getFirstNonZeroElemntPosition(Binversed, graph, &matrix)
		x := matrix.getOriginalNumber(i, j)
		perfectMatching = append(perfectMatching, x)
		matrix.updateMatrixOfCorrectnes(x.First, x.Second)
		if B.RawMatrix().Rows > 4 {
			B = c.getSubMatrix(i, j, B)
			Binversed = c.constructTatasMatrix(B, graph)
		} else {
			if Binversed.RawMatrix().Rows > 2 {
				Binversed = c.getSubMatrix(i, j, Binversed)
			}
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
			value := float64(c.Rnd.Int31n(1000))
			if rawMatrix[i*vertexAmount+e] == 0 {
				rawMatrix[i*vertexAmount+e] = value
			}
			if rawMatrix[e*vertexAmount+i] == 0 {
				rawMatrix[e*vertexAmount+i] = -value
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
			if i != j && d.At(i, j) != 0 {
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

func isGraphSingleConnected(matrix *gonum.Dense) bool {
	vertexMarkers := make([]int, matrix.RawMatrix().Cols)
	vertexMarkers[0] = connectedVertex
	for contains(vertexMarkers, connectedVertex) {
		vertex := find(vertexMarkers, connectedVertex)
		//fmt.Println("markers:", vertexMarkers, " vertex:", vertex)
		for i := 0; i < matrix.RawMatrix().Rows; i++ {
			if matrix.At(vertex, i) != 0 && vertexMarkers[i] == notCheckedVertex {
				vertexMarkers[i] = connectedVertex
			}
		}
		vertexMarkers[vertex] = checkedVertex
	}
	return !contains(vertexMarkers, notCheckedVertex)
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

func isMatrixContainsNonZeroElements(matrix *gonum.Dense) bool {
	for i := 0; i < matrix.RawMatrix().Rows; i++ {
		for j := 0; j < matrix.RawMatrix().Cols; j++ {
			if matrix.At(i, j) != 0 {
				return true
			}
		}
	}
	return false
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func find(s []int, e int) int {
	for pos, a := range s {
		if a == e {
			return pos
		}
	}
	return -1
}

func containsInPair(s []gopair.IntPair, e int) bool {
	for _, a := range s {
		if a.First == e || a.Second == e {
			return true
		}
	}
	return false
}
