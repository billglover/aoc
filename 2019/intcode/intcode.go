package intcode

// Mem is the memory of the intcode computer.
type Mem []int

type mode int

const (
	position mode = iota
	immediate
	relative
)

func val(m Mem, p int, pm mode, base int) int {
	switch pm {
	case position:
		return m[m[p]]
	case immediate:
		return m[p]
	case relative:
		return m[base+m[p]]
	}
	return 0
}

func target(m Mem, p int, pm mode, base int) int {
	switch pm {
	case position:
		return m[p]
	case relative:
		return base + m[p]
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

// Run takes in intcode program and a channel for inputs. It returns an output
// channel and a channel to idnicate when the program terminates.
func Run(p Mem, in <-chan int) (chan int, chan bool) {
	ip := 0
	out := make(chan int, 1)
	done := make(chan bool)
	m := make([]int, 10000)
	copy(m, p)
	p = m

	go func() {
		base := 0
		for {
			op := p[ip]   // operation at the Instruction Pointer
			i := op % 100 // instruction is the two right hand digits of op

			switch i {
			case 1: // add
				w := 4
				om := opmodes(w, op)
				r := val(p, ip+1, om[1], base) + val(p, ip+2, om[2], base)
				loc := target(p, ip+3, om[3], base)
				p[loc] = r
				ip += w

			case 2: // multiply
				w := 4
				om := opmodes(w, op)
				r := val(p, ip+1, om[1], base) * val(p, ip+2, om[2], base)
				loc := target(p, ip+3, om[3], base)
				p[loc] = r
				ip += w

			case 3: // in
				w := 2
				om := opmodes(w, op)
				loc := target(p, ip+1, om[1], base)
				p[loc] = <-in
				ip += w

			case 4: // out
				w := 2
				om := opmodes(w, op)
				v := val(p, ip+1, om[1], base)
				out <- v
				ip += w

			case 5: // jit (Jump if True)
				w := 3
				om := opmodes(w, op)
				v := val(p, ip+1, om[1], base)
				addr := val(p, ip+2, om[2], base)
				if v != 0 {
					ip = addr
					continue
				}
				ip += w

			case 6: // jif (Jump if False)
				w := 3
				om := opmodes(w, op)
				v := val(p, ip+1, om[1], base)
				addr := val(p, ip+2, om[2], base)
				if v == 0 {
					ip = addr
					continue
				}
				ip += w

			case 7: // lt
				w := 4
				om := opmodes(w, op)
				v1 := val(p, ip+1, om[1], base)
				v2 := val(p, ip+2, om[2], base)
				loc := target(p, ip+3, om[3], base)
				v := 0
				if v1 < v2 {
					v = 1
				}
				p[loc] = v
				ip += w

			case 8: // eq
				w := 4
				om := opmodes(w, op)
				v1 := val(p, ip+1, om[1], base)
				v2 := val(p, ip+2, om[2], base)
				loc := target(p, ip+3, om[3], base)

				v := 0
				if v1 == v2 {
					v = 1
				}
				p[loc] = v
				ip += w

			case 9:
				w := 2
				om := opmodes(w, op)
				v := val(p, ip+1, om[1], base)
				base += v
				ip += w

			case 99:
				close(out)
				done <- true
				return
			default:
				close(out)
				done <- false
				return
			}
		}
	}()
	return out, done
}
