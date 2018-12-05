package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	partB := flag.Bool("partB", false, "run partB")
	flag.Parse()

	f, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("unable to open file:", err)
		os.Exit(1)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("unable to read file:", err)
		os.Exit(1)
	}

	if *partB == false {
		result := triggerReactions(data)
		fmt.Println("Result:", len(result))
		return
	}

	units := map[byte]bool{}
	for i := range data {
		b := data[i]
		if b > 90 {
			b = b - 32
		}
		units[b] = true
	}

	raw := string(data)
	minUnit, minLen := byte(0), len(data)
	for k := range units {
		filtered := strings.Map(func(r rune) rune {
			if r == rune(k) || r == rune(k+32) {
				return -1
			}
			return r
		}, raw)
		result := triggerReactions([]byte(filtered))
		if len(result) < minLen {
			minUnit, minLen = k, len(result)
		}

	}
	fmt.Println("Result:", string(minUnit), minLen)

}

func triggerReactions(d []byte) []byte {

	r := make([]byte, len(d))

	rIdx := 0
	dIdx := 1

	r[0] = d[0]

	for dIdx < len(d) {
		diff := 0
		if rIdx >= 0 {
			diff = int(d[dIdx]) - int(r[rIdx])
		}

		if diff == 32 || diff == -32 {
			rIdx--
			dIdx++
			continue
		}
		rIdx++
		r[rIdx] = d[dIdx]
		dIdx++
	}
	return r[:rIdx+1]
}
