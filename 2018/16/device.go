package device

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Opcode has a numeric ID, a name and a function, Do. Functions can operate on
// global registers so care should be taken in their use.
type Opcode struct {
	ID   int
	Name string
	Do   func(a, b, c int)
}

// Reg are globally accessible array of registers.
var Reg = [4]int{0, 0, 0, 0}

// Sample represents a sample gathered from the device.
type Sample struct {
	Before      [4]int
	Instruction [4]int
	After       [4]int
}

// LoadInstructions returns the instruction set.
func LoadInstructions() map[string]Opcode {
	ocs := []Opcode{
		// Addition
		{Name: "addr", Do: func(a, b, c int) { Reg[c] = Reg[a] + Reg[b] }},
		{Name: "addi", Do: func(a, b, c int) { Reg[c] = Reg[a] + b }},

		// Multiplication
		{Name: "mulr", Do: func(a, b, c int) { Reg[c] = Reg[a] * Reg[b] }},
		{Name: "muli", Do: func(a, b, c int) { Reg[c] = Reg[a] * b }},

		// Bitwise AND
		{Name: "banr", Do: func(a, b, c int) { Reg[c] = Reg[a] & Reg[b] }},
		{Name: "bani", Do: func(a, b, c int) { Reg[c] = Reg[a] & b }},

		// Bitwise OR
		{Name: "borr", Do: func(a, b, c int) { Reg[c] = Reg[a] | Reg[b] }},
		{Name: "bori", Do: func(a, b, c int) { Reg[c] = Reg[a] | b }},

		// Assignment
		{Name: "setr", Do: func(a, b, c int) { Reg[c] = Reg[a] }},
		{Name: "seti", Do: func(a, b, c int) { Reg[c] = a }},

		// Greater-than testing
		{Name: "gtir", Do: func(a, b, c int) {
			if a > Reg[b] {
				Reg[c] = 1
				return
			}
			Reg[c] = 0
		}},
		{Name: "gtri", Do: func(a, b, c int) {
			if Reg[a] > b {
				Reg[c] = 1
				return
			}
			Reg[c] = 0
		}},
		{Name: "gtrr", Do: func(a, b, c int) {
			if Reg[a] > Reg[b] {
				Reg[c] = 1
				return
			}
			Reg[c] = 0
		}},

		// Equality testing
		{Name: "eqir", Do: func(a, b, c int) {
			if a == Reg[b] {
				Reg[c] = 1
				return
			}
			Reg[c] = 0
		}},
		{Name: "eqri", Do: func(a, b, c int) {
			if Reg[a] == b {
				Reg[c] = 1
				return
			}
			Reg[c] = 0
		}},
		{Name: "eqrr", Do: func(a, b, c int) {
			if Reg[a] == Reg[b] {
				Reg[c] = 1
				return
			}
			Reg[c] = 0
		}},
	}

	var Instructions = map[string]Opcode{}
	for _, oc := range ocs {
		Instructions[oc.Name] = oc
	}

	return Instructions
}

// RunProgram takes a file name and exucutes the program.
func RunProgram(path string, ins map[int]Opcode) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	Reg = [4]int{0, 0, 0, 0}

	lines := strings.Split(string(data), "\n")
	for _, l := range lines {
		id, a, b, c := 0, 0, 0, 0
		fmt.Sscanf(l, "%d %d %d %d", &id, &a, &b, &c)

		s := Sample{
			Before:      Reg,
			Instruction: [4]int{id, a, b, c},
		}

		op := ins[id]
		op.Do(a, b, c)

		s.After = Reg
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

// ReadSamples takes a file name and returns a slice of samples.
func ReadSamples(path string) ([]Sample, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reBefore := regexp.MustCompile(`^Before:\s+\[(\d+),\s?(\d+),\s?(\d+),\s?(\d+)\]$`)
	reInstruction := regexp.MustCompile(`^(\d+)\s?(\d+)\s?(\d+)\s?(\d+)$`)
	reAfter := regexp.MustCompile(`^After:\s+\[(\d+),\s?(\d+),\s?(\d+),\s?(\d+)\]$`)

	samples := []Sample{}
	s := Sample{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		mB := reBefore.FindStringSubmatch(line)
		mI := reInstruction.FindStringSubmatch(line)
		mA := reAfter.FindStringSubmatch(line)

		switch {
		case len(mB) > 0:
			a, _ := strconv.Atoi(mB[1])
			b, _ := strconv.Atoi(mB[2])
			c, _ := strconv.Atoi(mB[3])
			d, _ := strconv.Atoi(mB[4])
			s.Before = [4]int{a, b, c, d}

		case len(mI) > 0:
			a, _ := strconv.Atoi(mI[1])
			b, _ := strconv.Atoi(mI[2])
			c, _ := strconv.Atoi(mI[3])
			d, _ := strconv.Atoi(mI[4])
			s.Instruction = [4]int{a, b, c, d}

		case len(mA) > 0:
			a, _ := strconv.Atoi(mA[1])
			b, _ := strconv.Atoi(mA[2])
			c, _ := strconv.Atoi(mA[3])
			d, _ := strconv.Atoi(mA[4])
			s.After = [4]int{a, b, c, d}
			samples = append(samples, s)
			s = Sample{}
		}
	}

	err = f.Close()
	if err != nil {
		return samples, err
	}

	if err := scanner.Err(); err != nil {
		return samples, err
	}

	return samples, nil
}
