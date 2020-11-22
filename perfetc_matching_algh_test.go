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
var crashSource int64 = 1605443878136428355

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
	matcher := RandomMatcher{Rnd: rand.New(rand.NewSource(testSeed)), detChecker: NewCondChecker(1)}
	if !matcher.IsPerfectMatchingExist(g) {
		fmt.Println("Perfect matching does not exist")
		fmt.Println("TestGetPerfectMatchingByRandomAlgorithm=[ok]")
		return
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
	g := lslib.NewGraph()
	g.ParseGraph(benchGraph)
	g.HungryNumIndependent()
	g.Print()

	//matcher := RandomMatcher{rnd: rand.New(rand.NewSource(testSeed))  }
	//fixedVertexes := []gopair.IntPair{gopair.IntPair{First: 2, Second: 13}, gopair.IntPair{First: 11, Second: 24}}

	//  fixedVertexes := []gopair.IntPair{gopair.IntPair{First: 2, Second: 13}, gopair.IntPair{First: 6, Second: 29}}
	fixedVertexes := []gopair.IntPair{{First: 3, Second: 22}, {First: 12, Second: 16}}
	//matcher := &RandomMathcerWithFixedVertexes{FixedVertexes: fixedVertexes, RandomMatcher: RandomMatcher{Rnd: rand.New(rand.NewSource(crashSource))}}
	matcher := NewRandomMathcerWithFixedVertexes(fixedVertexes)
	matcher.SetDetChecker(NewCondChecker(2))
	if !matcher.IsPerfectMatchingExist(g) {
		fmt.Println("Perfect matching does not exist")
		fmt.Println("TestGetPerfectMatchingByRandomAlgorithmWithFixedVertexes=[ok]")
		return
	}

	//matcher.SetFixedVertexes([]gopair.IntPair{{First: 0, Second: 6}, {First: 3, Second: 4}})
	matching, err := matcher.GetPerfectMatching(g) //matcher.getPerfectMatchingByRandomAlgorithmWithFixedVertexes(g, fixedVertexes)
	if err == NoPerfectMatching {
		fmt.Println(err)
		fmt.Println("TestGetPerfectMatchingByRandomAlgorithmWithFixedVertexes=[ok]")
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
