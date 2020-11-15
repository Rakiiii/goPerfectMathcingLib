package permatchalgh

import (
	"fmt"
	"math/rand"
	"testing"

	lslib "github.com/Rakiiii/goBipartitonLocalSearch"
	graphlib "github.com/Rakiiii/goGraph"
	gopair "github.com/Rakiiii/goPair"
)

var Dir string = "Testing/"

var testGraphPerfectMatching string = Dir + "GetPerfectMatchingGraph"
var testGraphGetHungryContractedGraphNI = Dir + "GetHungryContractedGraphNI"
var benchGraph = Dir + "graph_bench_1"
var testSeed int64 = 1238678900867786087

func TestGetPerfectMatchingByRandomAlgorithm(t *testing.T) {
	//t.Skip()
	fmt.Println("Start TestGetPerfectMatchingByRandomAlgorithm...")
	var parser = new(graphlib.Parser)
	var g, err = parser.ParseUnweightedUndirectedGraphFromFile(benchGraph)
	if err != nil {
		fmt.Println(err)
		return
	}
	g.Print()
	matcher := RandomMatcher{Rnd: rand.New(rand.NewSource(testSeed))}
	if !matcher.IsPerfectMatchingExist(g) {
		t.Error("Perfect matching does not exist")
	}

	matching, err := matcher.getPerfectMatchingByRandomAlgorithm(g)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(matching); i++ {
		fmt.Println("(", matching[i].First, ":", matching[i].Second, ")")
	}
	fmt.Println("TestGetPerfectMatchingByRandomAlgorithm=[ok]")
}

func TestGetPerfectMatchingByRandomAlgorithmWithFixedVertexes(t *testing.T) {
	//t.Skip()
	fmt.Println("Start TestGetPerfectMatchingByRandomAlgorithmWithFixedVertexes...")
	// var parser = new(graphlib.Parser)
	// var g, err = parser.ParseUnweightedUndirectedGraphFromFile(testGraphGetHungryContractedGraphNI)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	g := lslib.NewGraph()
	g.ParseGraph(benchGraph)
	g.HungryNumIndependent()
	g.Print()

	//matcher := RandomMatcher{rnd: rand.New(rand.NewSource(testSeed))  }
	//fixedVertexes := []gopair.IntPair{gopair.IntPair{First: 2, Second: 13}, gopair.IntPair{First: 11, Second: 24}}

	//  fixedVertexes := []gopair.IntPair{gopair.IntPair{First: 2, Second: 13}, gopair.IntPair{First: 6, Second: 29}}
	fixedVertexes := []gopair.IntPair{{First: 0, Second: 13}, {First: 4, Second: 28}, {First: 1, Second: 15}, {First: 2, Second: 25}, {First: 3, Second: 20}, {First: 5, Second: 16}, {First: 6, Second: 24}, {First: 7, Second: 23}, {First: 8, Second: 19}, {First: 9, Second: 21}, {First: 10, Second: 29}, {First: 11, Second: 17}, {First: 12, Second: 18}, {First: 14, Second: 26}, {First: 22, Second: 27}}
	matcher := NewRandomMathcerWithFixedVertexes(fixedVertexes)
	if !matcher.IsPerfectMatchingExist(g) {
		t.Error("Perfect matching does not exist")
		return
	}

	matcher.SetFixedVertexes([]gopair.IntPair{{First: 0, Second: 6}, {First: 3, Second: 4}})
	matching, err := matcher.GetPerfectMatching(g) //matcher.getPerfectMatchingByRandomAlgorithmWithFixedVertexes(g, fixedVertexes)
	if err == NoPerfectMatching {
		fmt.Println(err)
	} else {
		fmt.Println("Start result:")
		for i := 0; i < len(matching); i++ {
			fmt.Println("(", matching[i].First, ":", matching[i].Second, ")")
		}
		fmt.Println("TestGetPerfectMatchingByRandomAlgorithmWithFixedVertexes=[ok]")
	}

}

func TestInterfaceImplemenation(t *testing.T) {
	t.Skip()
	var matcher IPerfectMatcher
	matcher = NewRandomMatcher()
	var parser = new(graphlib.Parser)
	var g, _ = parser.ParseUnweightedUndirectedGraphFromFile(testGraphPerfectMatching)
	matcher.GetPerfectMatching(g)
}
