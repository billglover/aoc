package ampcircuit

import (
	"github.com/billglover/aoc/2019/07/intcode"
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

func Amplify(p []int, phase []int) (int, error) {
	signal, inPhase := 0, phase[0]

	for s := 0; s < len(phase); s++ {
		inPhase = phase[s]
		in := []int{inPhase, signal}

		out, err := intcode.Run(intcode.Mem(p), in)
		if err != nil {
			return signal, err
		}
		signal = out[0]
	}
	return signal, nil
}

func MaxAmplify(p []int, phases [][]int) (int, []int, error) {
	maxSignal := 0
	atPhase := make([]int, len(phases[0]))

	for _, phase := range phases {
		signal, err := Amplify(p, phase)
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
