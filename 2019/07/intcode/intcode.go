package intcode

import (
	"errors"
)

// Mem is the memory of the intcode computer.
type Mem []int

type mode int

const (
	position mode = iota
	immediate
)

func val(m Mem, p int, pm mode) int {
	switch pm {
	case position:
		return m[m[p]]
	case immediate:
		return m[p]
	}
	return 0
}

func opmodes(w int, op int) []mode {
	om := make([]mode, w)
	for p := 0; p < w; p++ {
		op = op / 10
		om[p] = mode(op % 10)
	}
	return om
}

// Run takes in intcode program and a slice of inputs. It executes the program
// and returns a slice of outputs. An error is returned if program execution
// fails.
func Run(p Mem, in []int) ([]int, error) {
	out := []int{}

	ip := 0

	for p[ip] != 99 {

		op := p[ip]   // operation at the Instruction Pointer
		i := op % 100 // instruction is the two right hand digits of op
		switch i {
		case 1: // add
			w := 4
			om := opmodes(w, op)
			r := val(p, ip+1, om[1]) + val(p, ip+2, om[2])
			loc := p[ip+3]
			p[loc] = r
			ip += w

		case 2: // multiply
			w := 4
			om := opmodes(w, op)
			r := val(p, ip+1, om[1]) * val(p, ip+2, om[2])
			loc := p[ip+3]
			p[loc] = r
			ip += w

		case 3: // in
			if len(in) == 0 {
				return out, errors.New("not enough input provided")
			}
			w := 2
			loc := p[ip+1]
			p[loc] = in[0]
			in = in[1:]
			ip += w

		case 4: // out
			w := 2
			om := opmodes(w, op)
			v := val(p, ip+1, om[1])
			out = append(out, v)
			ip += w

		case 5: // jit (Jump if True)
			w := 3
			om := opmodes(w, op)
			v := val(p, ip+1, om[1])
			addr := val(p, ip+2, om[2])
			if v != 0 {
				ip = addr
				continue
			}
			ip += w

		case 6: // jif (Jump if False)
			w := 3
			om := opmodes(w, op)
			v := val(p, ip+1, om[1])
			addr := val(p, ip+2, om[2])
			if v == 0 {
				ip = addr
				continue
			}
			ip += w

		case 7: // lt
			w := 4
			om := opmodes(w, op)
			v1 := val(p, ip+1, om[1])
			v2 := val(p, ip+2, om[2])
			loc := p[ip+3]
			v := 0
			if v1 < v2 {
				v = 1
			}
			p[loc] = v
			ip += w

		case 8: // eq
			w := 4
			om := opmodes(w, op)
			v1 := val(p, ip+1, om[1])
			v2 := val(p, ip+2, om[2])
			loc := p[ip+3]

			v := 0
			if v1 == v2 {
				v = 1
			}
			p[loc] = v
			ip += w

		case 99:
			return out, nil
		default:
			return out, errors.New("invalid instruction")
		}
	}

	return out, nil
}
