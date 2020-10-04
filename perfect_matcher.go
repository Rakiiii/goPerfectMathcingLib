package permatchalgh

import (
	graphlib "github.com/Rakiiii/goGraph"
	gopair "github.com/Rakiiii/goPair"
)

type IPerfectMatcher interface {
	GetPerfectMatching(graphlib.IGraph) ([]gopair.IntPair, error)
}

type IPerfectMatcherChecker interface {
	IsPerfectMatchingExist(graph graphlib.IGraph) bool
}
