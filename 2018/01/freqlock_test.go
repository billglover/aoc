package freqlock_test

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/billglover/aoc/2018/01/freqlock"
)

var lockCases = []struct {
	input   string
	initial int
	result  int
}{
	{input: "1 1 1", initial: 0, result: 3},
	{input: "+1 +1 -2", initial: 0, result: 0},
	{input: "-1 -2 -3", initial: 0, result: -6},
	{input: "1 3 -2\n", initial: 0, result: 2},
}

func TestLockFreq(t *testing.T) {
	for _, tc := range lockCases {
		r := strings.NewReader(strings.Replace(tc.input, " ", "\n", -1))
		f, err := freqlock.LockFreq(r, tc.initial)
		if err != nil {
			t.Error("unexpected error returned:", err)
		}

		if f != tc.result {
			t.Errorf("got: %d, expected: %d", f, tc.result)
		}
	}
}

var repeatCases = []struct {
	input   string
	initial int
	result  int
	err     error
}{
	{input: "1 -2 3 1 1 -2 1 3 -4", initial: 0, result: 2, err: nil},
	{input: "+1 -1", initial: 0, result: 0, err: nil},
	{input: "+3 +3 +4 -2 -4", initial: 0, result: 10, err: nil},
	{input: "-6 +3 +8 +5 -6", initial: 0, result: 5, err: nil},
	{input: "+7 +7 -2 -7 -4", initial: 0, result: 14, err: nil},
}

func TestRepeatedFreq(t *testing.T) {
	for _, tc := range repeatCases {
		r := strings.NewReader(strings.Replace(tc.input, " ", "\n", -1))
		f, err := freqlock.RepeatedFreq(r, tc.initial)
		if err != nil {
			t.Error("unexpected error returned:", err)
		}

		if f != tc.result {
			t.Errorf("got: %d, expected: %d", f, tc.result)
		}
	}
}

func TestAoCSolution(t *testing.T) {

	if aoc == false {
		t.SkipNow()
	}

	f, err := os.Open("input.txt")
	if err != nil {
		t.Fatal("unable to open file:", err)
	}

	result, err := freqlock.LockFreq(f, 0)
	if err != nil {
		t.Error("unable to determine final frequency:", err)
	}
	fmt.Println("final frequency:", result)

	f.Seek(0, 0)

	result, err = freqlock.RepeatedFreq(f, 0)
	if err != nil {
		t.Error("unable to determine repeated frequency:", err)
	}
	fmt.Println("first repeated frequency:", result)
	fmt.Println()

	err = f.Close()
	if err != nil {
		t.Fatal("unable to close file:", err)
	}
}

var aoc bool

func init() {
	flag.BoolVar(&aoc, "aoc", false, "run tests against Accent on Code input")
	flag.Parse()
}
