package chocchart

import (
	"reflect"
	"testing"
)

var bakeCases = []struct {
	r1     int
	r2     int
	result []int
}{
	{r1: 7, r2: 3, result: []int{1, 0}},
	{r1: 2, r2: 3, result: []int{5}},
}

func TestBake(t *testing.T) {
	for _, tc := range bakeCases {
		result := bake(tc.r1, tc.r2)
		if reflect.DeepEqual(result, tc.result) == false {
			t.Errorf("got: %v, want %v", result, tc.result)
		}
	}
}

var bakeRoundCases = []struct {
	scores []int
	e1Ptr  int
	e2Ptr  int
	result []int
}{
	{scores: []int{3, 7}, e1Ptr: 0, e2Ptr: 1, result: []int{3, 7, 1, 0}},
	{scores: []int{2, 3}, e1Ptr: 0, e2Ptr: 1, result: []int{2, 3, 5}},
	{scores: []int{0, 0}, e1Ptr: 0, e2Ptr: 1, result: []int{0, 0, 0}},
}

func TestBakeRound(t *testing.T) {
	for _, tc := range bakeRoundCases {
		result := bakeRound(tc.scores, tc.e1Ptr, tc.e2Ptr)
		if reflect.DeepEqual(result, tc.result) == false {
			t.Errorf("got: %v, want %v", result, tc.result)
		}
	}
}

var nextRecipeCases = []struct {
	scores []int
	ptr    int
	result int
}{
	{scores: []int{3, 7}, ptr: 0, result: 0},
	{scores: []int{3, 7}, ptr: 1, result: 1},
	{scores: []int{3, 7, 1, 0}, ptr: 0, result: 0},
	{scores: []int{3, 7, 1, 0}, ptr: 1, result: 1},
	{scores: []int{3, 7, 1, 0, 1, 0}, ptr: 0, result: 4},
	{scores: []int{3, 7, 1, 0, 1, 0}, ptr: 1, result: 3},
	{scores: []int{3, 7, 1, 0, 1, 0, 1, 2}, ptr: 6, result: 0},
	{scores: []int{3, 7, 1, 0, 1, 0, 1, 2, 4, 5}, ptr: 8, result: 3},
}

func TestNextRecipe(t *testing.T) {
	for _, tc := range nextRecipeCases {
		result := nextRecipe(tc.scores, tc.ptr)
		if result != tc.result {
			t.Logf("scores: %v, ptr: %d", tc.scores, tc.ptr)
			t.Errorf("got: %v, want %v", result, tc.result)
		}
	}
}

var bakeRoundsCases = []struct {
	scores []int
	e1Ptr  int
	e2Ptr  int
	rounds int
	result []int
}{
	{scores: []int{3, 7}, e1Ptr: 0, e2Ptr: 1, rounds: 2, result: []int{3, 7}},
	{scores: []int{3, 7}, e1Ptr: 0, e2Ptr: 1, rounds: 4, result: []int{3, 7, 1, 0}},
	{scores: []int{3, 7}, e1Ptr: 0, e2Ptr: 1, rounds: 6, result: []int{3, 7, 1, 0, 1, 0}},
	{scores: []int{3, 7}, e1Ptr: 0, e2Ptr: 1, rounds: 9, result: []int{3, 7, 1, 0, 1, 0, 1, 2, 4}},
}

func TestBakeRounds(t *testing.T) {
	for _, tc := range bakeRoundsCases {
		result := bakeRounds(tc.scores, tc.e1Ptr, tc.e2Ptr, tc.rounds)
		if reflect.DeepEqual(result, tc.result) == false {
			t.Errorf("got: %v, want %v", result, tc.result)
		}
	}
}

var bakeNAfterCases = []struct {
	scores []int
	e1Ptr  int
	e2Ptr  int
	rounds int
	n      int
	result string
}{
	{scores: []int{3, 7}, e1Ptr: 0, e2Ptr: 1, rounds: 9, n: 10, result: "5158916779"},
	{scores: []int{3, 7}, e1Ptr: 0, e2Ptr: 1, rounds: 5, n: 10, result: "0124515891"},
	{scores: []int{3, 7}, e1Ptr: 0, e2Ptr: 1, rounds: 18, n: 10, result: "9251071085"},
	{scores: []int{3, 7}, e1Ptr: 0, e2Ptr: 1, rounds: 2018, n: 10, result: "5941429882"},
}

func TestBakeNAfter(t *testing.T) {
	for _, tc := range bakeNAfterCases {
		result := bakeNAfter(tc.scores, tc.e1Ptr, tc.e2Ptr, tc.rounds, tc.n)
		if result != tc.result {
			t.Errorf("got: %v, want %v", result, tc.result)
		}
	}
}

func TestPartOne(t *testing.T) {
	result := bakeNAfter([]int{3, 7}, 0, 1, 290431, 10)
	t.Logf("Part One: %s", result)
}

var matchScoreCases = []struct {
	scores []int
	e1Ptr  int
	e2Ptr  int
	match  string
	result int
}{
	{scores: []int{3, 7}, e1Ptr: 0, e2Ptr: 1, match: "51589", result: 9},
	{scores: []int{3, 7}, e1Ptr: 0, e2Ptr: 1, match: "01245", result: 5},
	{scores: []int{3, 7}, e1Ptr: 0, e2Ptr: 1, match: "92510", result: 18},
	{scores: []int{3, 7}, e1Ptr: 0, e2Ptr: 1, match: "59414", result: 2018},
}

func TestMatchScore(t *testing.T) {
	for _, tc := range matchScoreCases {
		result := matchScore(tc.scores, tc.e1Ptr, tc.e2Ptr, tc.match)
		if result != tc.result {
			t.Errorf("got: %v, want %v", result, tc.result)
		}
	}
}

func TestPartTwo(t *testing.T) {
	result := matchScore([]int{3, 7}, 0, 1, "290431")
	t.Logf("Part Two: %d", result)
}
