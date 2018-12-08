package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func main() {
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

	asciiEntries := strings.Split(string(data), " ")
	entries := make([]int, len(asciiEntries))
	for i := range asciiEntries {
		entries[i], err = strconv.Atoi(asciiEntries[i])
		if err != nil {
			fmt.Println("unable to parse input data:", err)
			os.Exit(1)
		}
	}

	partB := true
	_, sum, err := readNodeValues(entries, partB)
	if err != nil {
		fmt.Println("unable to parse nodes data:", err)
		os.Exit(1)
	}

	fmt.Println("Sum of Metadata:", sum)

}

// ReadNodeValues reads nodes from the provided slice of ints. It returns the
// number of values read, along with the sum of the node values.
func readNodeValues(data []int, partB bool) (int, int, error) {
	if len(data) == 0 {
		return 0, 0, nil
	}

	ptr := 0
	numChildren := data[ptr]
	ptr++
	numMetadata := data[ptr]
	ptr++

	children := make(map[int]int)
	for childIndex := 1; childIndex <= numChildren; childIndex++ {
		ptrAdv, childValue, err := readNodeValues(data[ptr:], partB)
		if err != nil {
			return 0, 0, errors.Wrap(err, "unable to read children")
		}

		ptr += ptrAdv
		children[childIndex] = childValue
	}

	metadata := data[ptr : ptr+numMetadata]
	ptr += numMetadata

	if partB == false || numChildren == 0 {
		value := 0
		for i := range metadata {
			value += metadata[i]
		}

		for i := range children {
			value += children[i]
		}
		return ptr, value, nil
	}

	value := 0
	for i := range metadata {
		cValue := children[metadata[i]]
		value += cValue
	}
	return ptr, value, nil
}
