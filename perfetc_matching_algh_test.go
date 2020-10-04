package permatchalgh

import (
	"fmt"
	"math/rand"
	"testing"

	graphlib "github.com/Rakiiii/goGraph"
)

var testGraphPerfectMatching string = "Testing/GetPerfectMatchingGraph"
var testSeed int64 = 1238678900867786087

func TestGetPerfectMatchingByRandomAlgorithm(t *testing.T) {
	var parser = new(graphlib.Parser)
	var g, err = parser.ParseUnweightedUndirectedGraphFromFile(testGraphPerfectMatching)
	if err != nil {
		fmt.Println(err)
		return
	}
	matcher := RandomMatcher{rnd: rand.New(rand.NewSource(testSeed))}
	matching, err := matcher.GetPerfectMatchingByRandomAlgorithm(g)
	for i := 0; i < len(matching); i++ {
		fmt.Println("(", matching[i].First, ":", matching[i].Second, ")")
	}
}

func TestInterfaceImplemenation(t *testing.T) {
	var matcher IPerfectMatcher
	matcher = NewRandomMatcher()
	var parser = new(graphlib.Parser)
	var g, _ = parser.ParseUnweightedUndirectedGraphFromFile(testGraphPerfectMatching)
	matcher.GetPerfectMatching(g)
}
