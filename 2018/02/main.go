package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	// use a flag for code specific to part B, credit: @lizthegrey
	b := flag.Bool("b", false, "return the result for part b")
	flag.Parse()

	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(data), "\n")

	if *b == false {
		partA(lines)
	} else {
		partB(lines)
	}

}

// PartA counts the number of lines in the input that have exactly two or
// three repeated characters. It returns the multiple of the two counts.
func partA(lines []string) {
	twos := 0
	threes := 0

	for _, l := range lines {
		chars := map[byte]int{}
		for i := range l {
			chars[l[i]] = chars[l[i]] + 1
		}

		seen2 := false
		seen3 := false

		for _, v := range chars {
			switch v {
			case 2:
				seen2 = true
			case 3:
				seen3 = true
			default:
				continue
			}
		}

		if seen2 {
			twos++
		}

		if seen3 {
			threes++
		}
	}

	fmt.Println(twos * threes)
}

// PartB looks for pairs of lines which differ by exactly one character at
// the same position.
func partB(lines []string) {
	sort.Strings(lines)

	for n := range lines {
		if n < len(lines)-1 {

			l1 := lines[n]
			l2 := lines[n+1]

			diffCount := 0
			diffLoc := 0
			for i := range l1 {
				if l1[i] != l2[i] {
					diffCount++
					diffLoc = i
				}
			}

			if diffCount == 1 {
				box1 := l1[:diffLoc] + l1[diffLoc+1:]
				box2 := l2[:diffLoc] + l2[diffLoc+1:]
				fmt.Println(box1)
				fmt.Println(box2)
			}

		}
	}
}
