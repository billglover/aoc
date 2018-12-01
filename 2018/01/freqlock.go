package freqlock

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

const maxIter = 1000

// LockFreq reads a series of line separated frequency adjustments. It
// returns the result of applying these frequency adjustments in turn
// to the initial frequency. An error is returned if it is unable to
// parse the adjustments.
func LockFreq(r io.Reader, f int) (int, error) {
	s := bufio.NewScanner(r)

	for s.Scan() {
		i, err := strconv.Atoi(s.Text())
		if err != nil {
			return f, err
		}
		f += i
	}

	if err := s.Err(); err != nil {
		return f, err
	}

	return f, nil
}

// RepeatedFreq reads a series of line separated frequency adjustments. It
// returns the first repeated frequency that results from applying the
// adjustments in turn to the initial frequency. The adjustments are applied
// in a loop (up to maxIter) until a repeating frequency is found. An error
// is returned if it is unable to parse the adjustments.
//
// Data from the io.Reader provided is read in its entirety, up to the EOF
// marker before processing begins.
func RepeatedFreq(r io.Reader, f int) (int, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return f, err
	}

	lines := strings.Split(string(data), "\n")

	seen := map[int]bool{f: true}

	for iter := 0; iter < maxIter; iter++ {
		for _, l := range lines {
			i, err := strconv.Atoi(l)
			if err != nil {
				return f, err
			}
			f += i

			if _, ok := seen[f]; ok == false {
				seen[f] = true
				continue
			}
			return f, nil
		}
	}

	return 0, fmt.Errorf("repeated frequency not found after %d iterations", maxIter)
}
