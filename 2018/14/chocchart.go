package chocchart

import (
	"strconv"
)

func bake(r1, r2 int) []int {

	result := []int{}
	for r := r1 + r2; r > 0; r = r / 10 {
		result = append(result, r%10)
	}

	for i := len(result)/2 - 1; i >= 0; i-- {
		opp := len(result) - 1 - i
		result[i], result[opp] = result[opp], result[i]
	}

	if len(result) == 0 {
		result = []int{0}
	}

	return result
}

func bakeRound(scores []int, e1Ptr, e2Ptr int) []int {
	result := bake(scores[e1Ptr], scores[e2Ptr])
	scores = append(scores, result...)
	return scores
}

func nextRecipe(scores []int, ptr int) int {
	newPtr := (((ptr + 1 + scores[ptr]) % len(scores)) + len(scores)) % len(scores)
	//fmt.Println(len(scores), ptr, newPtr)
	return newPtr
}

func bakeRounds(scores []int, e1Ptr, e2Ptr, rounds int) []int {
	for len(scores) < rounds {
		scores = bakeRound(scores, e1Ptr, e2Ptr)
		e1Ptr, e2Ptr = nextRecipe(scores, e1Ptr), nextRecipe(scores, e2Ptr)
	}
	return scores[:rounds]
}

func bakeNAfter(scores []int, e1Ptr, e2Ptr, rounds, n int) string {
	result := bakeRounds(scores, e1Ptr, e2Ptr, rounds+n)
	lastN := result[rounds:]

	text := ""
	for i := range lastN {
		text += strconv.Itoa(lastN[i])
	}
	return text
}

func matchScore(scores []int, e1Ptr, e2Ptr int, match string) int {
	data := []byte(match)
	target := make([]int, len(data))
	for i := range data {
		target[i] = int(data[i]) - 48
	}
	index := -1

	for {
		scores = bakeRound(scores, e1Ptr, e2Ptr)
		e1Ptr, e2Ptr = nextRecipe(scores, e1Ptr), nextRecipe(scores, e2Ptr)

		if len(scores) < len(target)+3 {
			continue
		}

		if index = found(scores[len(scores)-len(target)-3:], target); index > 0 {
			index += (len(scores) - len(target) - 3)
			break
		}
	}

	return index
}

func found(a, b []int) int {
	if len(a) < len(b) {
		return -1
	}

	for i := 0; i < len(a)-len(b); i++ {
		match := true
		for j := 0; j < len(b); j++ {
			if a[i+j] != b[j] {
				match = false
			}
		}
		if match == true {
			return i
		}
	}
	return -1
}
