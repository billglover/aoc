package main

import (
	"fmt"
	"strings"
)

// Table represents the table on which pots can be placed.
type table struct {
	pots   []string
	offset int
}

func (t *table) padLeft(n int) {
	pad := strings.Repeat(".", n)
	t.pots = append(strings.Split(pad, ""), t.pots...)
	t.offset += n
}

func (t *table) padRight(n int) {
	pad := strings.Repeat(".", n)
	t.pots = append(t.pots, strings.Split(pad, "")...)
}

func newArrangement(s string, n int) table {
	t := table{
		pots:   strings.Split(s, ""),
		offset: n,
	}
	return t
}

func (t *table) neighbours(i int) string {
	i += t.offset
	ns := make([]string, 5)
	for n := -2; n <= 2; n++ {
		if i+n >= len(t.pots) {
			ns[2+n] = "."
		} else if i+n < 0 {
			ns[2+n] = "."
		} else {
			ns[2+n] = t.pots[i+n]
		}
	}

	return strings.Join(ns, "")
}

func (t *table) score() int {
	score := 0
	for p := range t.pots {
		if t.pots[p] == "#" {
			score += p - t.offset
		}
	}
	return score
}

func (t *table) trim() {
	for p := len(t.pots) - 1; p > 0; p-- {
		if t.pots[p] == "#" {
			t.pots = t.pots[:p+t.offset]
			return
		}
	}
}

func (t table) String() string {
	s := fmt.Sprintf("%s [%d] @ %d", strings.Join(t.pots, ""), t.offset, t.score())
	return s
}
