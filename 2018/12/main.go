package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const generations = 20

type rules map[string]string

func main() {
	// open file
	f, err := os.Open("test.txt")
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
	for i := 0; i < generations; i++ {
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
		if i%100000 == 0 {
			fmt.Println("iter:", i)
		}
	}
	fmt.Println("Table:", t.score())
}
