package intcode

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tc := []struct {
		program string
		result  []int
	}{
		{program: "1,0,0,0,99", result: []int{2, 0, 0, 0, 99}},
		{program: "2,3,0,3,99", result: []int{2, 3, 0, 6, 99}},
		{program: "2,4,4,5,99,0", result: []int{2, 4, 4, 5, 99, 9801}},
		{program: "1,1,1,4,99,5,6,0,99", result: []int{30, 1, 1, 4, 2, 5, 6, 0, 99}},
	}

	for _, c := range tc {
		got := Run(c.program)
		if reflect.DeepEqual(got, c.result) == false {
			t.Errorf("unexpected result:\n%v\n!=\n%v", got, c.result)
		}
	}
}

func TestRunPart1(t *testing.T) {
	p := "1,0,0,3,1,1,2,3,1,3,4,3,1,5,0,3,2,1,10,19,2,6,19,23,1,23,5,27,1,27,13,31,2,6,31,35,1,5,35,39,1,39,10,43,2,6,43,47,1,47,5,51,1,51,9,55,2,55,6,59,1,59,10,63,2,63,9,67,1,67,5,71,1,71,5,75,2,75,6,79,1,5,79,83,1,10,83,87,2,13,87,91,1,10,91,95,2,13,95,99,1,99,9,103,1,5,103,107,1,107,10,111,1,111,5,115,1,115,6,119,1,119,10,123,1,123,10,127,2,127,13,131,1,13,131,135,1,135,10,139,2,139,6,143,1,143,9,147,2,147,6,151,1,5,151,155,1,9,155,159,2,159,6,163,1,163,2,167,1,10,167,0,99,2,14,0,0"
	pa := strings.Split(p, ",")
	pa[1] = "12"
	pa[2] = "2"
	p = strings.Join(pa, ",")
	result := Run(p)
	fmt.Println("Part One:", result[0])
}

func TestRunPart2(t *testing.T) {
	p := "1,0,0,3,1,1,2,3,1,3,4,3,1,5,0,3,2,1,10,19,2,6,19,23,1,23,5,27,1,27,13,31,2,6,31,35,1,5,35,39,1,39,10,43,2,6,43,47,1,47,5,51,1,51,9,55,2,55,6,59,1,59,10,63,2,63,9,67,1,67,5,71,1,71,5,75,2,75,6,79,1,5,79,83,1,10,83,87,2,13,87,91,1,10,91,95,2,13,95,99,1,99,9,103,1,5,103,107,1,107,10,111,1,111,5,115,1,115,6,119,1,119,10,123,1,123,10,127,2,127,13,131,1,13,131,135,1,135,10,139,2,139,6,143,1,143,9,147,2,147,6,151,1,5,151,155,1,9,155,159,2,159,6,163,1,163,2,167,1,10,167,0,99,2,14,0,0"
	pa := strings.Split(p, ",")
	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			pa[1] = strconv.Itoa(noun)
			pa[2] = strconv.Itoa(verb)
			np := strings.Join(pa, ",")
			result := Run(np)
			if result[0] != 19690720 {
				continue
			}
			answer := 100*noun + verb
			fmt.Println("Part Two:", answer)
		}
	}
}
