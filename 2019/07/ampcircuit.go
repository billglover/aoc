package ampcircuit

import (
	"github.com/billglover/aoc/2019/intcode"
)

// Perm generates all the permutations of a slice of ints. Alogirthm inspired
// by the solution shown here https://yourbasic.org/golang/generate-permutation-slice-string/
func Perm(a []int) [][]int {
	A := [][]int{}

	perm(0, a, func(a []int) {
		b := make([]int, len(a))
		copy(b, a)
		A = append(A, b)
	})
	return A
}

func perm(k int, a []int, f func([]int)) {
	if k > len(a) {
		f(a)
		return
	}

	perm(k+1, a, f)
	for j := k + 1; j < len(a); j++ {
		a[k], a[j] = a[j], a[k]
		perm(k+1, a, f)
		a[k], a[j] = a[j], a[k]
	}
}

func Amplify(p []int, phase []int, partB bool) (int, error) {
	inChans := make([]chan int, len(phase))
	outChans := make([]chan int, len(phase))
	doneChans := make([]chan bool, len(phase))
	out := make(chan int, 1)

	for s := 0; s < len(phase); s++ {
		m := make([]int, len(p))
		copy(m, p)
		if inChans[s] == nil {
			inChans[s] = make(chan int, 1)
		}
		ph := phase[s]
		in := inChans[s]
		outChans[s], doneChans[s] = intcode.Run(intcode.Mem(m), in)
		in <- ph
	}

	go func() {
		for v := range outChans[0] {
			inChans[1] <- v
		}
	}()

	go func() {
		for v := range outChans[1] {
			inChans[2] <- v
		}
	}()

	go func() {
		for v := range outChans[2] {
			inChans[3] <- v
		}
	}()

	go func() {
		for v := range outChans[3] {
			inChans[4] <- v
		}
	}()

	if partB {
		go func() {
			defer close(out)
			for v := range outChans[4] {
				out <- v
				inChans[0] <- v
			}
		}()
	}

	inChans[0] <- 0

	if !partB {
		signal := <-outChans[4]
		return signal, nil
	}

	signal := 0
	for v := range out {
		signal = v
	}

	return signal, nil
}

func MaxAmplify(p []int, phases [][]int, partB bool) (int, []int, error) {
	maxSignal := 0
	atPhase := make([]int, len(phases[0]))

	for _, phase := range phases {
		signal, err := Amplify(p, phase, partB)
		if err != nil {
			return maxSignal, atPhase, err
		}
		if signal > maxSignal {
			maxSignal = signal
			copy(atPhase, phase)
		}
	}

	return maxSignal, atPhase, nil
}
