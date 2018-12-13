package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const generations = 50000000000

type rules map[string]string

func main() {
	// open file
	f, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("unable to open file:", err)
		os.Exit(1)
	}
	defer f.Close()

	// read all the data (we can assume small input files)
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("unable to read file:", err)
		os.Exit(1)
	}

	// split the file into lines
	lines := strings.Split(string(data), "\n")

	// read first line to get initial input
	initial := strings.Replace(lines[0], "initial state: ", "", 1)
	t := newArrangement(initial, 0)

	// read remaining lines to get rules
	r := make(rules)
	for _, l := range lines[2:] {
		parts := strings.Split(l, " ")
		in := parts[0]
		out := parts[2]
		r[in] = out
	}

	// match positions
	d0, d1, d2 := 0, -1, -2
	s0, s1 := 0, 0
	for i := 0; i < generations; i++ {
		if d0 == d1 && d1 == d2 {
			s0 += d0
			// fmt.Print("-")
			// fmt.Println("scores:", s0, s1)
			// fmt.Print("-")
			// fmt.Println("deltas:", d0, d1, d2)
			continue
		}
		s1 = s0
		d2, d1 = d1, d0

		pots := make([]string, len(t.pots)+11)
		for n := -5; n <= len(t.pots)+5; n++ {
			k := t.neighbours(n)
			if v, ok := r[k]; ok == true {
				pots[n+5] = v
			} else {
				pots[n+5] = "."
			}
		}

		t = newArrangement(strings.Join(pots, ""), 5)
		t.trim()

		s0 = t.score()
		d0 = s0 - s1
		// fmt.Println("scores:", s0, s1)
		// fmt.Println("deltas:", d0, d1, d2)
	}
	fmt.Println("---")
	fmt.Println("Total:", s0)
}
