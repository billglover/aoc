package intcode

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Run(p string) []int {
	pa := strings.Split(p, ",")
	m := make([]int, len(pa))
	for i, v := range pa {
		opcode, err := strconv.Atoi(v)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		m[i] = opcode
	}

	for pc := 0; pc < len(m); pc = pc + 4 {
		op := m[pc]
		switch op {
		case 1:
			op1, op2, store := m[m[pc+1]], m[m[pc+2]], m[pc+3]
			m[store] = op1 + op2
		case 2:
			op1, op2, store := m[m[pc+1]], m[m[pc+2]], m[pc+3]
			m[store] = op1 * op2
		case 99:
			return m
		default:
			fmt.Println(op)
			fmt.Fprintln(os.Stderr, "invalid opcode:", op)
			os.Exit(2)
		}
	}

	return m
}
