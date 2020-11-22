package permatchalgh

import (
	graphlib "github.com/Rakiiii/goGraph"
	gopair "github.com/Rakiiii/goPair"
	gonum "gonum.org/v1/gonum/mat"
)

type IPerfectMatcher interface {
	GetPerfectMatching(graphlib.IGraph) ([]gopair.IntPair, error)
}

type IPerfectMatcherChecker interface {
	IsPerfectMatchingExist(graph graphlib.IGraph) bool
}

type IElementMatcher interface {
	matchElements(float64, float64) float64
}

type IMatrixDeterminantChecker interface {
	isDetZero(matrix gonum.Matrix) bool
}
