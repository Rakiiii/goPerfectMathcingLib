package permatchalgh

import (
	"fmt"
	"testing"

	graphlib "github.com/Rakiiii/goGraph"
)

var testGraphPerfectMatching string = "Testing/GetPerfectMatchingGraph"

func TestGetPerfectMatchingByRandomAlgorithm(t *testing.T) {
	var parser = new(graphlib.Parser)
	var g, err = parser.ParseUnweightedUndirectedGraphFromFile(testGraphPerfectMatching)
	if err != nil {
		fmt.Println(err)
		return
	}
	matching, err := GetPerfectMatchingByRandomAlgorithm(g)
	for i := 0; i < len(matching); i++ {
		fmt.Println("(", matching[i].First, ":", matching[i].Second, ")")
	}
}
